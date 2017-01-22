package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestAccountValidation(t *testing.T) {
	account := account{}

	if account.IsValid() {
		t.Fatalf("expected empty account to be invalid")
	}

	account.Label = "ASD"
	if account.IsValid() {
		t.Fatalf("Still missing a code")
	}

	account.Code = "2000"
	if !account.IsValid() {
		t.Fatalf("Should be valid")
	}
}

func TestCreateAccount(t *testing.T) {
	defer TestManager.Clear()

	resp, _ := http.Post(TestManager.server.URL+"/accounts/", "application/json", strings.NewReader(`{"code": "2000","label": "Income"}`))

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create account: %v", resp.Body)
	}

	var a account
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&a)

	if a.Code != "2000" {
		t.Fatalf("Created account /w wrong code")
	}

	if a.Label != "Income" {
		t.Fatalf("Created account wo/ wrong label")
	}
}

func insertFakeAccount() account {
	var app = accountingApp{
		Db: TestManager.db,
	}
	var a = account{
		Code:  "2000",
		Label: "Income",
	}
	app.insertAccount(&a)
	return a
}

func TestAccountIndex(t *testing.T) {
	defer TestManager.Clear()
	insertFakeAccount()

	resp, _ := http.Get(TestManager.server.URL + "/accounts/")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %+v", "200", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var accounts []account

	_ = decoder.Decode(&accounts)
	if len(accounts) != 1 {
		t.Fatalf("Received wrong number of accounts: %v - '%v'", accounts, resp.Body)
	}
}

func TestUpdateAccount(t *testing.T) {
	defer TestManager.Clear()
	var a = insertFakeAccount()

	req, _ := http.NewRequest("PUT", fmt.Sprintf("%v/accounts/%d", TestManager.server.URL, a.ID), strings.NewReader(`{"code":"2001","label":"Expense"}`))
	resp, _ := http.DefaultClient.Do(req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected response to be success, got: %v", resp.StatusCode)
	}
}
