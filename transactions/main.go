package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rubenv/sql-migrate"
	"github.com/splicers/jet"
	"github.com/umsatz/go-aqbanking"
)

type pin struct {
	Blz string `json:"blz"`
	UID string `json:"uid"`
	PIN string `json:"pin"`
}

func (p *pin) BankCode() string {
	return p.Blz
}

func (p *pin) UserID() string {
	return p.UID
}

func (p *pin) Pin() string {
	return p.PIN
}

func loadPins(filename string) []aqbanking.Pin {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("%v", err)
		return nil
	}

	var _pins []pin
	if err = json.NewDecoder(f).Decode(&_pins); err != nil {
		log.Fatal("%v", err)
		return nil
	}

	var pins = make([]aqbanking.Pin, len(_pins))
	for i, pin := range _pins {
		pins[i] = aqbanking.Pin(&pin)
	}

	return pins
}

func logHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func renderError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)

	bytes, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})

	io.WriteString(w, string(bytes))
}

type _user struct {
	UserID     string `json:"user_id"`
	CustomerID string `json:"customer_id"`
	BankCode   string `json:"bank_code"`
	Name       string `json:"name"`
}

type account struct {
	Name          string `json:"name"`
	Owner         string `json:"owner"`
	AccountNumber string `json:"account_number"`
	BankCode      string `json:"bank_code"`
	BIC           string `json:"bic"`
	IBAN          string `json:"iban"`
	Currency      string `json:"currency"`
	Country       string `json:"country"`
}

type transaction struct {
	Purpose             string    `json:"purpose"`
	Text                string    `json:"text"`
	Status              string    `json:"status"`
	Date                time.Time `json:"date"`
	ValutaDate          time.Time `json:"valuta_date"`
	MandateReference    string    `json:"mandate_reference"`
	CustomerReference   string    `json:"customer_reference"`
	Total               float32   `json:"total"`
	TotalCurrency       string    `json:"total_currency"`
	Fee                 float32   `json:"fee"`
	FeeCurrency         string    `json:"fee_currency"`
	LocalBankCode       string    `json:"local_bank_code"`
	LocalAccountNumber  string    `json:"local_account_number"`
	LocalIBAN           string    `json:"local_iban"`
	LocalBIC            string    `json:"local_bic"`
	LocalName           string    `json:"local_name"`
	RemoteBankCode      string    `json:"remote_bank_code"`
	RemoteAccountNumber string    `json:"remote_account_number"`
	RemoteIBAN          string    `json:"remote_iban"`
	RemoteBIC           string    `json:"remote_bic"`
	RemoteName          string    `json:"remote_name"`
}

func listTransactions(aq *aqbanking.AQBanking) http.Handler {
	handler := func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		iban := vars["iban"]

		accounts, err := aq.Accounts()
		if err != nil {
			renderError(w, err)
			return
		}

		var account *aqbanking.Account
		for i := range accounts.Accounts {
			if accounts.Accounts[i].IBAN == iban {
				account = &accounts.Accounts[i]
			}
		}
		if account == nil {
			renderError(w, errors.New("account not found"))
			return
		}

		transactions, listError := aq.Transactions(account, nil, nil)
		if listError != nil {
			renderError(w, listError)
			return
		}

		trx := make([]transaction, len(transactions))
		for i, t := range transactions {
			trx[i] = transaction{
				Purpose:             t.Purpose,
				Text:                t.Text,
				Status:              t.Status,
				Date:                t.Date,
				ValutaDate:          t.ValutaDate,
				CustomerReference:   t.CustomerReference,
				Total:               t.Total,
				TotalCurrency:       t.TotalCurrency,
				Fee:                 t.Fee,
				FeeCurrency:         t.FeeCurrency,
				LocalBankCode:       t.LocalBankCode,
				LocalAccountNumber:  t.LocalAccountNumber,
				LocalIBAN:           t.LocalIBAN,
				LocalBIC:            t.LocalBIC,
				LocalName:           t.LocalName,
				RemoteBankCode:      t.RemoteBankCode,
				RemoteAccountNumber: t.RemoteAccountNumber,
				RemoteIBAN:          t.RemoteIBAN,
				RemoteBIC:           t.RemoteBIC,
				RemoteName:          t.RemoteName,
			}
		}

		var bytes []byte
		bytes, err = json.MarshalIndent(trx, "", "  ")
		if err != nil {
			renderError(w, err)
			return
		}

		io.WriteString(w, string(bytes))

	}
	return http.HandlerFunc(handler)
}

func listAccounts(aq *aqbanking.AQBanking) http.Handler {
	handler := func(w http.ResponseWriter, req *http.Request) {
		accounts, err := aq.Accounts()

		if err != nil {
			renderError(w, err)
			return
		}

		as := make([]account, len(accounts.Accounts))
		for i, a := range accounts.Accounts {
			as[i] = account{a.Name, a.Owner, a.AccountNumber, a.Bank.BankCode, a.BIC, a.IBAN, a.Currency, a.Country}
		}

		var bytes []byte
		bytes, err = json.MarshalIndent(as, "", "  ")
		if err != nil {
			renderError(w, err)
			return
		}

		io.WriteString(w, string(bytes))
	}
	return http.HandlerFunc(handler)
}

func listUsers(aq *aqbanking.AQBanking) http.Handler {
	handler := func(w http.ResponseWriter, req *http.Request) {
		users, err := aq.Users()

		if err != nil {
			renderError(w, err)
			return
		}

		us := make([]_user, len(users.Users))
		for i, u := range users.Users {
			us[i] = _user{
				UserID:     u.UserID,
				CustomerID: u.CustomerID,
				BankCode:   u.BankCode,
				Name:       u.Name,
			}
		}

		var bytes []byte
		bytes, err = json.MarshalIndent(us, "", "  ")
		if err != nil {
			renderError(w, err)
			return
		}

		io.WriteString(w, string(bytes))
	}
	return http.HandlerFunc(handler)
}

type newUser struct {
	_user
	Pin         string `json:"pin"`
	HbciVersion int    `json:"hbci_version"`
	ServerURI   string `json:"server_uri"`
}

func addUser(aq *aqbanking.AQBanking) http.Handler {
	handler := func(w http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		var u newUser

		if err := decoder.Decode(&u); err != nil {
			renderError(w, err)
			return
		}

		if u.BankCode == "" {
			renderError(w, errors.New("bank_code must be given"))
			return
		}

		if u.UserID == "" {
			renderError(w, errors.New("user_id must be given"))
			return
		}

		if u.ServerURI == "" {
			renderError(w, errors.New("server_uri must be given"))
			return
		}

		if u.HbciVersion == 0 {
			u.HbciVersion = 300
		}

		user := aqbanking.User{
			ID:          -1,
			UserID:      u.UserID,
			CustomerID:  u.CustomerID,
			BankCode:    u.BankCode,
			Name:        u.Name,
			ServerURI:   u.ServerURI,
			HbciVersion: u.HbciVersion,
		}

		if err := aq.AddPinTanUser(&user); err != nil {
			renderError(w, fmt.Errorf("unable to add user. %v\n", err))
			return
		}
		user.FetchAccounts(aq)

		var bytes []byte
		var err error
		bytes, err = json.MarshalIndent(u._user, "", "  ")
		if err != nil {
			renderError(w, err)
			return
		}

		io.WriteString(w, string(bytes))
	}
	return http.HandlerFunc(handler)
}

func runMigrations(db *jet.Db) {
	migrations := &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "migrations",
	}

	if n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up); err != nil {
		log.Printf("unable to migrate: %v", err)
	} else {
		log.Printf("Applied %d migrations!\n", n)
	}
}

func setupDb() *jet.Db {
	database := os.Getenv("DATABASE")
	if database == "" {
		database = "aqbanking"
	}

	revDsn := os.Getenv("REV_DSN")
	if revDsn == "" {
		user, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		revDsn = "user=" + user.Username + " dbname=" + database + " sslmode=disable"
	}

	newDb, err := jet.Open("postgres", revDsn)
	if err != nil {
		log.Fatal("failed to connect to postgres", err)
	}
	newDb.SetMaxIdleConns(100)
	newDb.ColumnConverter = (&umsatzSnakeConv{})

	return newDb
}

type umsatzSnakeConv struct{}

func (conv *umsatzSnakeConv) ColumnToFieldName(col string) string {
	name := ""
	if l := len(col); l > 0 {
		chunks := strings.Split(col, "_")
		for i, v := range chunks {
			chunks[i] = strings.Title(v)
		}
		name = strings.Replace(strings.Join(chunks, ""), "Id", "ID", -1)
	}
	return name
}

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	log.SetPrefix(fmt.Sprintf("pid:%d ", syscall.Getpid()))
}

func newTransactionServer(db *jet.Db, aq *aqbanking.AQBanking) http.Handler {
	r := mux.NewRouter()

	r.Handle("/users", logHandler(listUsers(aq))).Methods("GET")
	r.Handle("/users", logHandler(addUser(aq))).Methods("POST")
	r.Handle("/accounts", logHandler(listAccounts(aq))).Methods("GET")
	r.Handle("/accounts/{iban}/transactions", logHandler(listTransactions(aq))).Methods("GET")

	return r
}

// Set by make file on build
var (
	Version string
	Commit  string
)

func main() {
	var (
		aqbankingRoot = flag.String("aqbanking.conf", "./aq", "conf directory for aqbanking")
		httpAddress   = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	db := setupDb()
	runMigrations(db)

	aq, err := aqbanking.NewAQBanking("transactions", *aqbankingRoot)
	if err != nil {
		log.Fatal("unable to init aqbanking: %v", err)
	}

	// important: properly shut down aqbanking
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			aq.Free()
			os.Exit(1)
		}
	}()

	log.Printf("using aqbanking %d.%d.%d\n",
		aq.Version.Major,
		aq.Version.Minor,
		aq.Version.Patchlevel,
	)

	for _, pin := range loadPins("pins.json") {
		aq.RegisterPin(pin)
	}

	var us *aqbanking.UserCollection
	us, err = aq.Users()
	fmt.Printf("%d users\n", len(us.Users))
	us.Free()

	log.Printf("listening on %s", *httpAddress)
	log.Fatal(http.ListenAndServe(*httpAddress, newTransactionServer(db, aq)))
}
