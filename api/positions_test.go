package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestPositionValidations(t *testing.T) {
	position := position{}
	if errs, ok := position.IsValid(); ok {
		t.Fatalf("expected empty position to be invalid: %v", errs)
	}

	position.Currency = "EUR"
	if errs, ok := position.IsValid(); ok {
		t.Fatalf("expected empty position to be invalid")
	} else {
		if _, ok := errs["currency"]; ok {
			t.Fatalf("expected error on currency: %v", errs)
		}
	}

	position.PositionType = "income"
	if _, ok := position.IsValid(); ok {
		t.Fatalf("expected empty position to be invalid")
	}

	position.AccountCodeFrom = "5900"
	if _, ok := position.IsValid(); ok {
		t.Fatalf("expected empty position to be invalid")
	}

	position.AccountCodeTo = "5900"
	if _, ok := position.IsValid(); ok {
		t.Fatalf("expected empty position to be invalid")
	}

	position.InvoiceDate = shortDate(time.Now())
	if _, ok := position.IsValid(); ok {
		t.Fatalf("expected empty position to be invalid")
	}

	position.InvoiceNumber = "20140101"
	if _, ok := position.IsValid(); !ok {
		t.Fatalf("expect position to be valid")
	}
}

func TestPositionsIndex(t *testing.T) {
	defer TestManager.Clear()

	var period fiscalPeriod
	TestManager.db.QueryRow(`INSERT INTO fiscal_periods (name) VALUES ('2014') RETURNING id, name`).Scan(&period.ID, &period.Name)
	TestManager.db.Exec(`INSERT INTO positions
        (
          account_code_from,
          account_code_to,
          type,

          invoice_date,
          booking_date,
          invoice_number,

          total_amount_cents,
          currency,

          tax,
          fiscal_period_id,
          description,
          attachment_path
        )
      VALUES
        (
          '5900',
          '1100',
          'expense',

          NOW(),
          NOW(),
          '2001312',

          0,
          'EUR',

          0,
          $1,
          '',
          ''
        )`, period.ID)

	request, _ := http.NewRequest("GET", fmt.Sprintf("/positions?fiscal_period_id=%d", period.ID), strings.NewReader(""))
	response := httptest.NewRecorder()

	(&positionsApp{Db: TestManager.db}).IndexHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %+v", "200", response.Code)
	}

	decoder := json.NewDecoder(response.Body)
	var positions []position
	_ = decoder.Decode(&positions)
	if len(positions) != 1 {
		t.Fatalf("Received wrong number of positions: %v - %v", positions, response.Body)
	}
}

func TestPositionCreation(t *testing.T) {
	defer TestManager.Clear()

	p := insertFakeFiscalPeriod()
	payload := fmt.Sprintf(`{
      "fiscalPeriodId":   %d,
      "accountCodeFrom":  "5900",
      "accountCodeTo":    "1100",
      "type":             "income",
      "invoiceDate":      "2014-02-02",
      "invoiceNumber":    "20140201",
      "totalAmountCents": 2099,
      "currency":         "EUR",
      "tax":              700,
      "description":      "Kunde A Februar"
    }`, p.ID)
	resp, _ := http.Post(TestManager.server.URL+"/positions/", "application/json", strings.NewReader(payload))

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %+v", "200", resp.StatusCode, resp.Body)
	}

	// read returned payload
	decoder := json.NewDecoder(resp.Body)
	var pos position
	decoder.Decode(&pos)

	jsb, _ := json.Marshal(pos)
	var realAttributes map[string]interface{}
	json.Unmarshal(jsb, &realAttributes)

	// read given payload
	var expectedAttributes map[string]interface{}
	decoder = json.NewDecoder(strings.NewReader(payload))
	decoder.Decode(&expectedAttributes)

	// compare values
	for key, value := range expectedAttributes {
		if value != realAttributes[key] {
			t.Fatalf("did not properly unmarshall: %v. %v -> %v", key, value, realAttributes[key])
		}
	}
}

func insertFakePosition(p fiscalPeriod) position {
	var app = positionsApp{
		Db: TestManager.db,
	}
	var pos = position{
		FiscalPeriodID:   p.ID,
		AccountCodeFrom:  "5900",
		AccountCodeTo:    "1100",
		PositionType:     "income",
		InvoiceDate:      shortDate(time.Date(2014, 2, 2, 0, 0, 0, 0, time.UTC)),
		InvoiceNumber:    "20140201",
		TotalAmountCents: 2099,
		Currency:         "EUR",
		Tax:              700,
		Description:      "Example Description",
	}
	app.insertPosition(&pos)
	return pos
}

func TestPositionUpdate(t *testing.T) {
	defer TestManager.Clear()
	p := insertFakeFiscalPeriod()
	pos := insertFakePosition(p)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("%v/positions/%d", TestManager.server.URL, pos.ID), strings.NewReader(`{"type":"expense"}`))
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Update status code%v:\n\tbody: %+v", "200", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var updatedPos position
	decoder.Decode(&updatedPos)

	if updatedPos.PositionType != "expense" {
		t.Fatalf("position should have been expense now, got '%v'", updatedPos.PositionType)
	}
}

func TestPositionCreationWithMissingAttributes(t *testing.T) {
	defer TestManager.Clear()
	p := insertFakeFiscalPeriod()
	payload := fmt.Sprintf(`{
      "fiscalPeriodId":   %d,
      "accountCodeFrom":  "5900",
      "accountCodeTo":    "1100",
      "invoiceDate":      "2014-02-02",
      "invoiceNumber":    "20140201",
      "totalAmountCents": 2099,
      "currency":         "EUR",
      "tax":              700,
      "description":      "Kunde A Februar"
    }`, p.ID)
	resp, _ := http.Post(TestManager.server.URL+"/positions/", "application/json", strings.NewReader(payload))

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Non-expected status code%v:\n\tbody: %+v", "200", resp.StatusCode, resp.Body)
	}
}
