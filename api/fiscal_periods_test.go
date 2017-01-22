package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFiscalPeriodsIndex(t *testing.T) {
	app := &fiscalPeriodApp{TestManager.db}
	defer TestManager.Clear()

	app.Db.Exec(`INSERT INTO fiscal_periods (name) VALUES ('2014')`)

	request, _ := http.NewRequest("GET", "/fiscalPeriods", strings.NewReader(""))
	response := httptest.NewRecorder()

	app.indexHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code%v:\n\tbody: %+v", "200", response.Code)
	}

	decoder := json.NewDecoder(response.Body)
	var fiscalPeriods []fiscalPeriod
	_ = decoder.Decode(&fiscalPeriods)
	if len(fiscalPeriods) != 1 {
		t.Fatalf("Received wrong number of fiscalPeriods: %v - %v", fiscalPeriods, response.Body)
	}
}

func insertFakeFiscalPeriod() fiscalPeriod {
	var app = fiscalPeriodApp{
		Db: TestManager.db,
	}
	var p = fiscalPeriod{
		Name: "2110",
	}
	app.insertFiscalPeriod(&p)
	return p
}

func TestFiscalPeriodCreation(t *testing.T) {
	defer TestManager.Clear()

	resp, _ := http.Post(TestManager.server.URL+"/fiscalPeriods/", "application/json", strings.NewReader(`{"name": "2020"}`))

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create a new fiscalPeriod: %v\n%v", resp.StatusCode, resp.Body)
	}

	dec := json.NewDecoder(resp.Body)
	resPeriod := fiscalPeriod{}
	dec.Decode(&resPeriod)

	if resPeriod.Name != "2020" {
		t.Fatalf("failed to update the period")
	}
}

func TestFiscalPeriodUpdate(t *testing.T) {
	defer TestManager.Clear()

	p := insertFakeFiscalPeriod()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%v/fiscalPeriods/%d", TestManager.server.URL, p.ID), strings.NewReader(`{"name":"2011"}`))
	resp, _ := http.DefaultClient.Do(req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected response to be success, got: %v", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	resPeriod := fiscalPeriod{}
	dec.Decode(&resPeriod)

	if resPeriod.Name != "2011" {
		t.Fatalf("failed to update the period")
	}
}

func TestFiscalPeriodDeletion(t *testing.T) {
	defer TestManager.Clear()

	p := insertFakeFiscalPeriod()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%v/fiscalPeriods/%d", TestManager.server.URL, p.ID), nil)
	resp, _ := http.DefaultClient.Do(req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to deleted the period: %v\n%v", resp.StatusCode, resp.Body)
	}
}
