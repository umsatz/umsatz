package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type fiscalPeriod struct {
	ID       int       `json:"id,omitempty"`
	Name     string    `json:"name"`
	StartsAt shortDate `json:"startsAt"`
	EndsAt   shortDate `json:"endsAt"`
	Archived bool      `json:"archived"`
}

type fiscalPeriodApp struct {
	Db *sql.DB
}

func newFiscalPeriodApp(db *sql.DB) http.Handler {
	return &fiscalPeriodApp{Db: db}
}

func (app *fiscalPeriodApp) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		app.indexHandler(rw, req)
	} else if req.Method == "POST" {
		app.createHandler(rw, req)
	} else if req.Method == "PUT" {
		app.updateHandler(rw, req)
	} else if req.Method == "DELETE" {
		app.deleteHandler(rw, req)
	}
}

type fiscalPeriodResponse struct {
	fiscalPeriod
	PositionsCount    int    `json:"positionsCount"`
	TotalIncomeCents  int    `json:"totalIncomeCents"`
	TotalExpenseCents int    `json:"totalExpenseCents"`
	PositionsURL      string `json:"positionsUrl"`
}

const scanPeriodSQL = `
	SELECT
		id,
		name,
		starts_at,
		ends_at,
		(SELECT count(*)
			FROM positions
			WHERE fiscal_period_id = fiscal_periods.id
		) AS positions_count,
		COALESCE( (SELECT sum(CASE currency WHEN 'EUR' THEN total_amount_cents ELSE total_amount_cents_in_eur END)
			FROM positions WHERE fiscal_period_id = fiscal_periods.id
			AND type = 'income'
		), 0) AS total_income_cents,
		COALESCE( (SELECT sum(CASE currency WHEN 'EUR' THEN total_amount_cents ELSE total_amount_cents_in_eur END)
			FROM positions WHERE fiscal_period_id = fiscal_periods.id
			AND type = 'expense'
		), 0) AS total_expense_cents,
		archived
	FROM fiscal_periods
	ORDER BY starts_at ASC`

func (app *fiscalPeriodApp) scanPeriod(rows *sql.Rows) (fiscalPeriodResponse, error) {
	var period fiscalPeriodResponse
	if err := rows.Scan(
		&period.ID,
		&period.Name,
		&period.StartsAt,
		&period.EndsAt,
		&period.PositionsCount,
		&period.TotalIncomeCents,
		&period.TotalExpenseCents,
		&period.Archived,
	); err != nil {
		return fiscalPeriodResponse{}, err
	}
	return period, nil
}

func (app *fiscalPeriodApp) loadFiscalPeriods() ([]fiscalPeriodResponse, error) {
	rows, err := app.Db.Query(scanPeriodSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var periods []fiscalPeriodResponse
	for rows.Next() {
		period, err := app.scanPeriod(rows)
		if err != nil {
			return nil, err
		}
		periods = append(periods, period)
	}
	return periods, nil
}

func (app *fiscalPeriodApp) indexHandler(w http.ResponseWriter, req *http.Request) {
	fiscalPeriods, err := app.loadFiscalPeriods()

	if err != nil {
		log.Println("database error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(w)

	baseURI := baseURI(&req.Header)
	if strings.HasSuffix(baseURI, "/fiscalPeriods") {
		baseURI = baseURI[:len(baseURI)-13]
	}

	for i, fiscalPeriod := range fiscalPeriods {
		fiscalPeriods[i].PositionsURL = fmt.Sprintf(`%vpositions/?fiscal_period_id=%d`, baseURI, fiscalPeriod.ID)
	}

	if err := enc.Encode(fiscalPeriods); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (app *fiscalPeriodApp) insertFiscalPeriod(p *fiscalPeriod) error {
	const insertSQL = `
		INSERT INTO fiscal_periods
			(name, starts_at, ends_at)
		VALUES
			($1, $2, $3)
		RETURNING id
	`
	return app.Db.QueryRow(insertSQL, p.Name, p.StartsAt, p.EndsAt).Scan(&p.ID)
}

func (app *fiscalPeriodApp) createHandler(w http.ResponseWriter, req *http.Request) {
	dec := json.NewDecoder(req.Body)
	enc := json.NewEncoder(w)

	newPeriod := fiscalPeriodResponse{}
	if err := dec.Decode(&newPeriod); err != nil {
		enc.Encode(map[string]string{"decoding": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := app.insertFiscalPeriod(&newPeriod.fiscalPeriod); err != nil {
		enc.Encode(map[string]string{"inserting": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	baseURI := baseURI(&req.Header)
	if strings.HasSuffix(baseURI, "/fiscalPeriods") {
		baseURI = baseURI[:len(baseURI)-13]
	}
	newPeriod.PositionsURL = fmt.Sprintf(`%vpositions/?fiscal_period_id=%d`, baseURI, newPeriod.ID)

	w.WriteHeader(http.StatusCreated)
	enc.Encode(newPeriod)
}

func (app *fiscalPeriodApp) deleteHandler(w http.ResponseWriter, req *http.Request) {
	enc := json.NewEncoder(w)

	parts := strings.Split(req.URL.Path, "/")
	fiscalPeriodID, _ := strconv.Atoi(parts[len(parts)-1])
	var tx *sql.Tx
	var err error
	if tx, err = app.Db.Begin(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]string{"sql error": err.Error()})
		return
	}

	app.Db.Exec(`DELETE FROM positions WHERE fiscal_period_id = $1`, fiscalPeriodID)
	app.Db.Exec(`DELETE FROM fiscal_periods WHERE id = $1`, fiscalPeriodID)
	if err := tx.Commit(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]string{"sql error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *fiscalPeriodApp) updateHandler(w http.ResponseWriter, req *http.Request) {
	enc := json.NewEncoder(w)

	fiscalPeriodID, _ := strconv.Atoi(req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:])

	dec := json.NewDecoder(req.Body)

	existingPeriod := fiscalPeriod{}
	if err := dec.Decode(&existingPeriod); err != nil {
		enc.Encode(map[string]string{"decoding": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := app.Db.Exec(`UPDATE fiscal_periods SET name = $1, starts_at = $2, ends_at = $3, archived = $4 WHERE id = $5`, existingPeriod.Name, time.Time(existingPeriod.StartsAt), time.Time(existingPeriod.EndsAt), existingPeriod.Archived, fiscalPeriodID); err != nil {
		enc.Encode(map[string]string{"inserting": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	enc.Encode(existingPeriod)
}
