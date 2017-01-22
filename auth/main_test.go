package main

import (
	"encoding/json"
	"fmt"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

type testManager struct {
	fakeRegistration *Registration
	server           *httptest.Server
}

type testTOTPValidator bool

func (t testTOTPValidator) Validate(string, string) bool {
	return bool(t)
}

var TestManager testManager

func TestMain(m *testing.M) {
	storage.EncryptionKey = "asdasdasdasdasdasdasdasd"
	var path, _ = os.Getwd()
	storage.ConfigPath = path + "/fixtures/secrets.json"
	if err := storage.Load(); err != nil {
		panic(err)
	}

	TestManager = testManager{
		fakeRegistration: storage.ActiveRegistration,
		server:           httptest.NewServer(authRouter()),
	}
	ret := m.Run()
	os.Exit(ret)
}

func TestSignupSuccess(t *testing.T) {
	storage.ActiveRegistration = nil

	var resp, err = http.Post(TestManager.server.URL+"/signup/", "application/json", strings.NewReader(`{
			"firstname": "Max",
			"lastname":  "Mustermann",
			"email":     "max@mustermann.com",
			"password":  "password",
			"company":   "Musterfirma"
		}`))
	if err != nil {
		t.Fatalf("Expected response to succeed: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected http status code %v to be %v", resp.StatusCode, http.StatusCreated)
	}
}

// TODO convert into table driven test - all tests that require a registration
func TestSignupFails(t *testing.T) {
	storage.ActiveRegistration = TestManager.fakeRegistration

	var resp, err = http.Post(TestManager.server.URL+"/signup/", "application/json", strings.NewReader(`{
			"firstname": "Max",
			"lastname":  "Mustermann",
			"email":     "max@mustermann.com",
			"password":  "password",
			"company":   "Musterfirma"
		}`))
	if err != nil {
		t.Fatalf("Expected response to succeed: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected http status code %v to be %v", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestShowExistingRegistration(t *testing.T) {
	storage.ActiveRegistration = TestManager.fakeRegistration

	var resp, err = http.Get(TestManager.server.URL + "/registration/")
	if err != nil {
		t.Fatalf("Expected request to succeed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected http status code %v to be %v", resp.StatusCode, http.StatusOK)
	}
}

func TestUpdateExistingRegistration(t *testing.T) {
	storage.ActiveRegistration = TestManager.fakeRegistration
	var previousTOTPKey = storage.ActiveRegistration.TOTPKey
	var previousPassword = storage.ActiveRegistration.Password

	req, _ := http.NewRequest("PUT", fmt.Sprintf("%v/registration/", TestManager.server.URL), strings.NewReader(`{
		"firstname":"Marx",
		"lastname": "Murx",
		"company": "Mini Corb",
		"tax_id": "ASD123",
		"totp": "1337hax0r",
		"password": "bla"
	}`))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Expected request to succeed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected http status code %v to be %v", resp.StatusCode, http.StatusOK)
	}

	if storage.ActiveRegistration.Firstname != "Marx" ||
		storage.ActiveRegistration.Lastname != "Murx" ||
		storage.ActiveRegistration.CompanyName != "Mini Corb" ||
		storage.ActiveRegistration.TaxID != "ASD123" {
		t.Fatalf("Expected update to succeed")
	}

	if storage.ActiveRegistration.TOTPKey != previousTOTPKey ||
		storage.ActiveRegistration.Password != previousPassword {
		t.Fatalf("Should not overwrite sensitive information")
	}
}

func TestShowMissingRegistration(t *testing.T) {
	storage.ActiveRegistration = nil

	var resp, err = http.Get(TestManager.server.URL + "/registration/")
	if err != nil {
		t.Fatalf("Expected request to succeed: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected http status code %v to be %v", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestShowQRWithRegistration(t *testing.T) {
	storage.ActiveRegistration = TestManager.fakeRegistration
	var realCreatedAt = storage.ActiveRegistration.CreatedAt
	storage.ActiveRegistration.CreatedAt = time.Now()
	defer func() {
		storage.ActiveRegistration.CreatedAt = realCreatedAt
	}()

	var resp, err = http.Get(TestManager.server.URL + "/qr/")
	if err != nil {
		t.Fatalf("Expected response to succeed. Got %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected response %v to be %v", resp.StatusCode, http.StatusOK)
	}

	if _, err := png.Decode(resp.Body); err != nil {
		t.Fatalf("Expected /qr/ to render a valid image. Failed %v", err)
	}
}

func TestShowQRWithoutRegistration(t *testing.T) {
	storage.ActiveRegistration = nil

	var resp, err = http.Get(TestManager.server.URL + "/qr/")
	if err != nil {
		t.Fatalf("Expected response to succeed. Got %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected response %v to be %v", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestCreateSessionSuccessAndValidate(t *testing.T) {
	storage.ActiveRegistration = TestManager.fakeRegistration
	storage.Validator = testTOTPValidator(true)

	var resp, err = http.Post(TestManager.server.URL+"/signin/", "application/json", strings.NewReader(`{
		"password": "password",
		"otp": "dontcare"
 	}`))
	if err != nil {
		t.Fatalf("Expected request to succeed: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected http status code %v to be %v", resp.StatusCode, http.StatusCreated)
	}

	var s Session
	var dec = json.NewDecoder(resp.Body)
	if err := dec.Decode(&s); err != nil {
		t.Fatalf("Failed to decode session: %v", err)
	}

	if s.Key == "" {
		t.Fatalf("Expected session with valid key")
	}

	resp, err = http.Post(TestManager.server.URL+"/validate/"+s.Key, "application/json", nil)
	if err != nil {
		t.Fatalf("Expected request to succeed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected http status %v to be %v", resp.StatusCode, http.StatusOK)
	}
}

func TestCreateSessionFails(t *testing.T) {
	storage.ActiveRegistration = nil

	var resp, err = http.Post(TestManager.server.URL+"/signin/", "application/json", strings.NewReader(`{
		"password": "password",
		"otp": "dontcare"
 	}`))
	if err != nil {
		t.Fatalf("Expected request to succeed: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected http status code %v to be %v", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestValidateSessionFails(t *testing.T) {
	// TODO test validate session
}
