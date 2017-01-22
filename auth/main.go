package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/json"
	_ "expvar"
	"flag"
	"fmt"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/umsatz/auth/Godeps/_workspace/src/github.com/pquerna/otp"
	"github.com/umsatz/auth/Godeps/_workspace/src/github.com/pquerna/otp/totp"
)

type Signin struct {
	Password string `json:"password"`
	OTPKey   string `json:"otp"`
}

// Registration is stored in an encrypted file. Auth ensures that only one registration exists
type Registration struct {
	Firstname   string    `json:"firstname"`
	Lastname    string    `json:"lastname"`
	TaxID       string    `json:"tax_id"`
	CompanyName string    `json:"company"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	TOTPKey     string    `json:"totp"`
	CreatedAt   time.Time `json:"created_at"`
}

// UpdateRegistration is used to hide secret attributes from change requests
type UpdateRegistration struct {
	*Registration
	Password  string    `json:"password,omitempty"`
	TOTPKey   string    `json:"totp,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type registrationResponse struct {
	*Registration
	Password  string `json:"password,omitempty"`
	TOTPKey   string `json:"totp,omitempty"`
	QRCodeURL string `json:"qrcode_url"`
}

func (r *Registration) IsValid() (bool, []string) {
	var errs = make([]string, 0)
	var valid = true

	const cutset = " "
	if strings.Trim(r.Firstname, cutset) == "" {
		valid = false
		errs = append(errs, "Firstname must not be blank")
	}
	if strings.Trim(r.Lastname, cutset) == "" {
		valid = false
		errs = append(errs, "Lastname must not be blank")
	}
	// TODO(rr) valid email?
	if strings.Trim(r.Email, cutset) == "" {
		valid = false
		errs = append(errs, "Email must not be blank")
	}
	if strings.Trim(r.Password, cutset) == "" {
		valid = false
		errs = append(errs, "Password must not be blank")
	}
	return valid, errs
}

// Session of a given user. Can be as many as you like
type Session struct {
	Key        string    `json:"key"`
	ValidUntil time.Time `json:"valid_until"`
}

type TOTPValidator interface {
	Validate(string, string) bool
}

type totpValidator struct{}

func (v *totpValidator) Validate(otpKey, secret string) bool {
	valid, _ := totp.ValidateCustom(
		otpKey,
		secret,
		time.Now().UTC(),
		totp.ValidateOpts{
			Period:    30,
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		},
	)
	return valid
}

// RegistrationStorage stores the active registration in an encrypted file
type RegistrationStorage struct {
	ConfigPath    string
	EncryptionKey string
	Validator     TOTPValidator

	mu                 *sync.Mutex
	ActiveRegistration *Registration
	Sessions           []*Session
}

// Load unencryptes the encrypted config
func (s *RegistrationStorage) Load() error {
	f, err := os.OpenFile(s.ConfigPath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher([]byte(s.EncryptionKey))

	if err != nil {
		panic(err)
	}

	encrypted, _ := ioutil.ReadAll(f)
	iv := encrypted[:aes.BlockSize]

	encrypted = encrypted[aes.BlockSize:]

	decrypter := cipher.NewCFBDecrypter(block, iv)

	decrypted := make([]byte, len(encrypted))
	decrypter.XORKeyStream(decrypted, encrypted)

	var r Registration
	json.NewDecoder(bytes.NewBuffer(decrypted)).Decode(&r)
	if ok, _ := r.IsValid(); ok {
		storage.mu.Lock()
		storage.ActiveRegistration = &r
		storage.mu.Unlock()
	} else {
		return fmt.Errorf("decryption failed")
	}
	return nil
}

// Sync writes the active registration to an encrypted file
func (s *RegistrationStorage) Sync() {
	// either 16, 24, or 32 bytes
	block, err := aes.NewCipher([]byte(s.EncryptionKey))

	if err != nil {
		panic(err)
	}

	var buf = &bytes.Buffer{}
	json.NewEncoder(buf).Encode(*storage.ActiveRegistration)

	encrypted := make([]byte, aes.BlockSize+buf.Len())
	iv := encrypted[:aes.BlockSize]

	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypter.XORKeyStream(encrypted[aes.BlockSize:], buf.Bytes())

	f, err := os.OpenFile(s.ConfigPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	f.Write(encrypted)
	f.Close()
}

// TODO(rr) move into main function to avoid globals
var (
	storage = NewRegistrationStorage()
)

func NewRegistrationStorage() RegistrationStorage {
	return RegistrationStorage{
		Sessions: make([]*Session, 0),
		mu:       &sync.Mutex{},
	}
}

func logHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"%s\t%s\t%s\t%s\t%s",
			r.RemoteAddr,
			time.Now().Format("2006-01-02T15:04:05 -0700"),
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	}
}

func authRouter() http.Handler {
	var r = http.DefaultServeMux

	// print out current registration informations, except for sensitive informations
	r.Handle("/registration/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if storage.ActiveRegistration == nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "no registration exists",
			})
			return
		}

		if req.Method == "GET" {
			json.NewEncoder(w).Encode(registrationResponse{
				Registration: storage.ActiveRegistration,
				QRCodeURL:    "/qr",
			})
		} else if req.Method == "PUT" {
			var r = UpdateRegistration{
				Registration: storage.ActiveRegistration,
			}

			storage.mu.Lock()
			json.NewDecoder(req.Body).Decode(&r)
			storage.mu.Unlock()

			storage.Sync()

			json.NewEncoder(w).Encode(r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

	}))

	// print qr code for setup /w time based stuff
	r.Handle("/qr/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if storage.ActiveRegistration == nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "no registration exists",
			})
			return
		}

		// QR is no longer available after 30 minutes
		if storage.ActiveRegistration.CreatedAt.Add(time.Minute * 30).Before(time.Now()) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "qr timeframe passed",
			})
		}

		key, err := otp.NewKeyFromURL(storage.ActiveRegistration.TOTPKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		var buf bytes.Buffer
		img, err := key.Image(400, 400)
		if err != nil {
			panic(err)
		}
		png.Encode(&buf, img)

		io.Copy(w, (&buf))
	}))

	// register with service
	r.Handle("/signup/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if storage.ActiveRegistration != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "active registration exists",
			})
			return
		}

		if req.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "only POST requests allowed",
			})
			return
		}

		var dec = json.NewDecoder(req.Body)

		var r Registration
		dec.Decode(&r)

		if ok, errs := r.IsValid(); !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errs)
			return
		}

		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "umsatz.app",
			AccountName: r.Email,
			Period:      30,
			Digits:      otp.DigitsSix,
			Algorithm:   otp.AlgorithmSHA1,
		})
		if err != nil {
			// TODO
		}
		r.TOTPKey = key.String()
		r.CreatedAt = time.Now()

		storage.mu.Lock()
		storage.ActiveRegistration = &r
		storage.mu.Unlock()
		storage.Sync()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(registrationResponse{
			Registration: &r,
			QRCodeURL:    "/qr/",
		})
	}))

	// validate a session
	r.Handle("/validate/", http.StripPrefix("/validate/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// TODO check if session exists
		if storage.ActiveRegistration == nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "no registration exists",
			})
			return
		}

		var known = false
		for _, s := range storage.Sessions {
			if s.Key == req.URL.Path {
				known = true
				storage.mu.Lock()
				s.ValidUntil = time.Now().Add(time.Minute * 30)
				storage.mu.Unlock()
				break
			}
		}

		if !known {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
	})))

	// create a new session
	r.Handle("/signin/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if storage.ActiveRegistration == nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "no registration exists",
			})
			return
		}

		var signin Signin
		var dec = json.NewDecoder(req.Body)
		dec.Decode(&signin)

		if signin.Password != storage.ActiveRegistration.Password {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("wrong password: %q", signin.Password),
			})
			return
		}

		key, err := otp.NewKeyFromURL(storage.ActiveRegistration.TOTPKey)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to restore TOTP. Registration invalid",
			})
		}

		var valid = storage.Validator.Validate(signin.OTPKey, key.Secret())

		if !valid {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("wrong totp: %v", signin.OTPKey),
			})
			return
		}

		var s = Session{
			Key:        fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%v-%v", storage.EncryptionKey, time.Now().Format("2006010203040507"))))),
			ValidUntil: time.Now().Add(time.Minute * 30),
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "X-UMSATZ-SESSION",
			Value: s.Key,
		})

		storage.mu.Lock()
		storage.Sessions = append(storage.Sessions, &s)
		storage.mu.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(s)
	}))

	return logHandler(r)
}

// Runs inside go routine
func (s *RegistrationStorage) cleanSessons() {
	for {
		select {
		case <-time.After(time.Minute * 1):
			if s.ActiveRegistration != nil {
				s.mu.Lock()
				for i := range storage.Sessions {
					var s = storage.Sessions[i]
					if s.ValidUntil.Before(time.Now()) {
						log.Printf("removing stale session %q", s.Key)
						storage.Sessions[i] = nil
						storage.Sessions = append(storage.Sessions[:i], storage.Sessions[i+1:]...)
					}
				}
				s.mu.Unlock()
			}
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	var (
		listen = flag.String("http.addr", ":8080", "http listen address")
		config = flag.String("config", "secrets.json", "encrypted configuration file")
	)
	flag.Parse()

	storage.ConfigPath = *config

	storage.EncryptionKey = os.Getenv("ENCRYPTION_KEY")
	storage.Validator = &totpValidator{}

	go storage.cleanSessons()

	if len(storage.EncryptionKey) != 16 && len(storage.EncryptionKey) != 24 && len(storage.EncryptionKey) != 32 {
		panic("key must be 16, 24 or 32 bytes length")
	}

	if err := storage.Load(); err != nil {
		log.Printf("Failed to load config: %v", err)
	}

	log.Printf("listening on %v", *listen)
	log.Printf("%v", http.ListenAndServe(*listen, authRouter()))
}
