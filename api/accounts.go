package main

//go:generate go-bindata -pkg main -o bindata.go migrations/

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"database/sql"
	_ "expvar"
)

type accountingApp struct {
	Db *sql.DB
}

type account struct {
	ID     int      `json:"id,omitempty"`
	Code   string   `json:"code"`
	Label  string   `json:"label"`
	Errors []string `json:"errors,omitempty"`
}

func newAccountingApp(db *sql.DB) http.Handler {
	return &accountingApp{Db: db}
}

func (app *accountingApp) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		app.indexHandler(rw, req)
	} else if req.Method == "POST" {
		app.createHandler(rw, req)
	} else if req.Method == "PUT" {
		app.updateHandler(rw, req)
	}
}

func (app *accountingApp) loadAccounts() ([]account, error) {
	rows, err := app.Db.Query(`SELECT id, code, label FROM accounts ORDER BY code ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []account
	for rows.Next() {
		account := account{}
		if err := rows.Scan(&account.ID, &account.Code, &account.Label); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (app *accountingApp) loadAccount(id int) (account, error) {
	a := account{}
	err := app.Db.
		QueryRow(`SELECT id, code, label FROM accounts WHERE id = $1`, id).
		Scan(&a.ID, &a.Code, &a.Label)

	if err != nil {
		return account{}, err
	}

	return a, nil
}

func (app *accountingApp) updateAccount(a account) error {
	_, err := app.Db.
		Exec(`UPDATE accounts SET code = $1, label = $2 WHERE ID = $3`, a.Code, a.Label, a.ID)
	if err != nil {
		return err
	}
	return nil
}

func (app *accountingApp) insertAccount(a *account) error {
	err := app.Db.
		QueryRow(`INSERT INTO accounts (code, label) VALUES ($1, $2) RETURNING id, code, label`, a.Code, a.Label).
		Scan(&a.ID, &a.Code, &a.Label)

	if err != nil {
		return err
	}
	return nil
}

func (a *account) AddError(attr string, errorMsg string) {
	a.Errors = append(a.Errors, attr+":"+errorMsg)
}

func (a *account) IsValid() bool {
	a.Errors = make([]string, 0)

	if a.Code == "" {
		a.AddError("code", "must be present")
	}
	if a.Label == "" {
		a.AddError("label", "must be present")
	}

	return len(a.Errors) == 0
}

func (app *accountingApp) indexHandler(w http.ResponseWriter, req *http.Request) {
	accounts, err := app.loadAccounts()
	if err != nil {
		log.Println("database error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(accounts); err != nil {
		log.Println("json marshal error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (app *accountingApp) updateHandler(w http.ResponseWriter, req *http.Request) {
	var account account

	parts := strings.Split(req.URL.Path, "/")

	if id, err := strconv.Atoi(parts[len(parts)-1]); err != nil {
		// TODO
	} else {
		if a, err := app.loadAccount(id); err != nil {
			log.Fatalf("unknown account: %v\n%v\n", parts[len(parts)-1], err)
		} else {
			account = a
		}
	}

	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&account); err != nil && err != io.EOF {
		log.Fatal("decode error", err)
	}

	enc := json.NewEncoder(w)

	if !account.IsValid() {
		log.Println("INFO: unable to update account due to validation errors: %v", account.Errors)
		w.WriteHeader(http.StatusBadRequest)

		if err := enc.Encode(account); err != nil {
			log.Println("INFO: unable to encode account: %v", err)
		}
		return
	}

	if err := app.updateAccount(account); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]string{"psql": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	enc.Encode(account)
}

func (app *accountingApp) createHandler(w http.ResponseWriter, req *http.Request) {
	dec := json.NewDecoder(req.Body)
	var account account
	if err := dec.Decode(&account); err != nil && err != io.EOF {
		log.Println("decode error", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf(`{ "errors": "%v" }`, err))
		return
	}

	enc := json.NewEncoder(w)

	if !account.IsValid() {
		log.Println("INFO: unable to insert account due to validation errors: %+v", account.Errors)
		w.WriteHeader(http.StatusBadRequest)

		enc.Encode(account)
		return
	}

	if err := app.insertAccount(&account); err != nil {
		log.Println("INFO: unable to insert account due to sql errors: %+v", err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	enc.Encode(account)
}
