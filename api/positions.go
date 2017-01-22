package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type position struct {
	ID                    int                 `json:"id,omitempty"`
	AccountCodeFrom       string              `json:"accountCodeFrom"`
	AccountCodeTo         string              `json:"accountCodeTo"`
	PositionType          string              `json:"type"`
	InvoiceDate           shortDate           `json:"invoiceDate"`
	BookingDate           shortDate           `json:"bookingDate"`
	InvoiceNumber         string              `json:"invoiceNumber"`
	TotalAmountCents      int                 `json:"totalAmountCents"`
	TotalAmountCentsInEur int                 `json:"totalAmountCentsEur"`
	Currency              string              `json:"currency"`
	Tax                   int                 `json:"tax"`
	FiscalPeriodID        int                 `json:"fiscalPeriodId"`
	Description           string              `json:"description"`
	CreatedAt             time.Time           `json:"createdAt"`
	UpdatedAt             time.Time           `json:"updatedAt"`
	AttachmentPath        string              `json:"attachmentPath"`
	Errors                map[string][]string `json:"errors",omitempty`
}

type positionsApp struct {
	Db           *sql.DB
	CurrencyHost string
}

func newPositionsApp(db *sql.DB, currencyHost string) http.Handler {
	return &positionsApp{
		Db:           db,
		CurrencyHost: currencyHost,
	}
}

func (app *positionsApp) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		app.IndexHandler(rw, req)
	} else if req.Method == "POST" {
		app.createHandler(rw, req)
	} else if req.Method == "PUT" {
		app.updateHandler(rw, req)
	} else if req.Method == "DELETE" {
		app.deleteHandler(rw, req)
	}
}

type positionTypeValidator struct {
}

func (p *position) IsValid() (map[string][]string, bool) {
	var errors = make(map[string][]string)
	var valid = true

	if p.Currency == "" {
		valid = false
		errors["currency"] = []string{"must be present"}
	}

	if p.AccountCodeTo == "" {
		valid = false
		errors["accountCodeTo"] = []string{"must be present"}
	}

	if p.AccountCodeFrom == "" {
		valid = false
		errors["accountCodeFrom"] = []string{"must be present"}
	}

	if p.InvoiceDate == shortDate(time.Time{}) {
		valid = false
		errors["invoiceDate"] = []string{"must be present"}
	}

	if p.InvoiceNumber == "" {
		valid = false
		errors["invoiceNumber"] = []string{"must be present"}
	}

	if p.PositionType == "" {
		valid = false
		errors["positionType"] = []string{"must be present"}
	} else {
		if p.PositionType != "income" && p.PositionType != "expense" {
			valid = false
			errors["positionType"] = []string{"must be either 'income' or 'expense'"}
		}
	}

	return errors, valid
}

type exchangeInfo struct {
	Date  string             `json:"date"`
	Rates map[string]float32 `json:"rates"`
}

// TODO(rr) rewrite as event based approach:
// create/ update writes change event into channel, worker consumes channel and
// updates positions as needed, after commit
func (app *positionsApp) setTotalAmountCentsInEur(p *position) error {
	if p.Currency == "EUR" {
		p.TotalAmountCentsInEur = p.TotalAmountCents
	} else {
		url := fmt.Sprintf(`http://%v/%v`, app.CurrencyHost, time.Time(p.InvoiceDate).Format("2006-01-02"))

		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)

		var v exchangeInfo
		decoder.Decode(&v)
		// TODO handle the case if the currency is missing or unknown
		p.TotalAmountCentsInEur = int(float32(p.TotalAmountCents) / v.Rates[p.Currency])
	}
	return nil
}

const listPositionsSQL = `SELECT
			id,
			account_code_from,
			account_code_to,
			type,
			invoice_date,
			booking_date,
			invoice_number,
			total_amount_cents,
			total_amount_cents_in_eur,
			currency,
			tax,
			fiscal_period_id,
			attachment_path,
			description,
			created_at,
			updated_at
		FROM positions
		WHERE fiscal_period_id = $1
		ORDER BY invoice_date ASC`

type sqlScanner interface {
	Scan(dest ...interface{}) error
}

const selectPositionSQL = `SELECT
		id,
		account_code_from,
		account_code_to,
		type,
		invoice_date,
		booking_date,
		invoice_number,
		total_amount_cents,
		total_amount_cents_in_eur,
		currency,
		tax,
		fiscal_period_id,
		attachment_path,
		description,
		created_at,
		updated_at
	FROM positions
	WHERE id = $1`

func (app *positionsApp) scanPosition(rows sqlScanner) (position, error) {
	p := position{}
	if err := rows.Scan(&p.ID,
		&p.AccountCodeFrom,
		&p.AccountCodeTo,
		&p.PositionType,
		&p.InvoiceDate,
		&p.BookingDate,
		&p.InvoiceNumber,
		&p.TotalAmountCents,
		&p.TotalAmountCentsInEur,
		&p.Currency,
		&p.Tax,
		&p.FiscalPeriodID,
		&p.AttachmentPath,
		&p.Description,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		return position{}, err
	}
	return p, nil
}

func (app *positionsApp) IndexHandler(w http.ResponseWriter, req *http.Request) {
	params, _ := url.ParseQuery(req.URL.RawQuery)

	fiscalPeriodID, ok := params["fiscal_period_id"]
	if !ok {
		log.Println("missing parameter fiscal_period_id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rows, err := app.Db.Query(listPositionsSQL, fiscalPeriodID[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer rows.Close()

	var positions []position
	for rows.Next() {
		p, err := app.scanPosition(rows)
		if err != nil {
			log.Println("database error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		positions = append(positions, p)
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(positions); err != nil {
		fmt.Println("ERRRRRORRR %v, %v", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "{}")
	}
}

const deletePositionSQL = `DELETE FROM positions WHERE id = $1 RETURNING attachment_path`

func (app *positionsApp) deleteHandler(w http.ResponseWriter, req *http.Request) {
	parts := strings.Split(req.URL.Path, "/")
	positionID, _ := strconv.Atoi(parts[len(parts)-1])

	var attachmentPath string
	if err := app.Db.QueryRow(deletePositionSQL, positionID).Scan(&attachmentPath); err != nil {
		log.Println("database error", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf(`{ "errors": "%v" }`, err))
		return
	}
	os.Remove(attachmentPath)

	io.WriteString(w, "")
}

const updatePositionSQL = `UPDATE
		positions
	SET
    account_code_from = $1,
    account_code_to = $2,
    type = $3,
    invoice_date = $4,
    booking_date = $5,
    invoice_number = $6,
    total_amount_cents = $7,
    total_amount_cents_in_eur = $8,
    currency = $9,
    tax = $10,
    fiscal_period_id = $11,
    description = $12,
    attachment_path = $13,
    updated_at = now()::timestamptz
  WHERE ID = $14`

func (app *positionsApp) updateHandler(w http.ResponseWriter, req *http.Request) {
	parts := strings.Split(req.URL.Path, "/")
	positionID, _ := strconv.Atoi(parts[len(parts)-1])

	var p, err = app.scanPosition(app.Db.QueryRow(selectPositionSQL, positionID))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&p); err != nil && err != io.EOF {
		log.Fatal("decode error", err)
	}

	if errors, ok := p.IsValid(); !ok {
		log.Println("INFO: unable to update position due to validation errors: %v", errors)
		w.WriteHeader(http.StatusBadRequest)

		p.Errors = errors
		if b, err := json.Marshal(p); err == nil {
			io.WriteString(w, string(b))
		}
		return
	}

	// TODO(rr) rewrite with evented approach. we might need to add a flag as well to persist if the lookup was fine
	if err := app.setTotalAmountCentsInEur(&p); err != nil {
		fmt.Println("currency lookup error %v", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "{}")
		return
	}

	_, updateError := app.Db.Exec(updatePositionSQL,
		p.AccountCodeFrom,
		p.AccountCodeTo,
		p.PositionType,
		time.Time(p.InvoiceDate),
		time.Time(p.BookingDate),
		p.InvoiceNumber,
		p.TotalAmountCents,
		p.TotalAmountCentsInEur,
		p.Currency,
		p.Tax,
		p.FiscalPeriodID,
		p.Description,
		p.AttachmentPath,
		positionID)

	enc := json.NewEncoder(w)
	if err := enc.Encode(p); err != nil || updateError != nil {
		fmt.Printf(`Error updating position: %v, %v\n`, err, updateError)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "{}")
	}
}

const insertPositionSQL = `INSERT INTO
	positions
  	(account_code_from, account_code_to, type, invoice_date, booking_date, invoice_number, total_amount_cents, total_amount_cents_in_eur, currency, tax, fiscal_period_id, description, attachment_path)
  VALUES
  	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
  RETURNING
  	id, account_code_from, account_code_to, type, invoice_date, booking_date, invoice_number, total_amount_cents, total_amount_cents_in_eur, currency, tax, fiscal_period_id, description, attachment_path`

func (app *positionsApp) insertPosition(p *position) error {
	return app.Db.QueryRow(insertPositionSQL,
		p.AccountCodeFrom,
		p.AccountCodeTo,
		p.PositionType,
		p.InvoiceDate,
		p.BookingDate,
		p.InvoiceNumber,
		p.TotalAmountCents,
		p.TotalAmountCentsInEur,
		p.Currency,
		p.Tax,
		p.FiscalPeriodID,
		p.Description,
		p.AttachmentPath,
	).Scan(
		&p.ID,
		&p.AccountCodeFrom,
		&p.AccountCodeTo,
		&p.PositionType,
		&p.InvoiceDate,
		&p.BookingDate,
		&p.InvoiceNumber,
		&p.TotalAmountCents,
		&p.TotalAmountCentsInEur,
		&p.Currency,
		&p.Tax,
		&p.FiscalPeriodID,
		&p.Description,
		&p.AttachmentPath,
	)
}

func (app *positionsApp) createHandler(w http.ResponseWriter, req *http.Request) {
	dec := json.NewDecoder(req.Body)
	var p position
	if err := dec.Decode(&p); err != nil && err != io.EOF {
		log.Println("decode error", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf(`{ "errors": "%v" }`, err))
		return
	}

	if errors, ok := p.IsValid(); !ok {
		log.Println("INFO: unable to insert position due to validation errors: %+v", errors)
		w.WriteHeader(http.StatusBadRequest)

		p.Errors = errors
		if b, err := json.Marshal(p); err == nil {
			io.WriteString(w, string(b))
		}
		return
	}

	if err := app.setTotalAmountCentsInEur(&p); err != nil {
		fmt.Println("currency lookup error %v", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "{}")
		return
	}

	insertError := app.insertPosition(&p)

	enc := json.NewEncoder(w)
	// fmt.Println(string(b))
	if err := enc.Encode(p); err != nil || insertError != nil {
		fmt.Println("INSERT ERRR %v, %v", err, insertError)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "{}")
	}
}
