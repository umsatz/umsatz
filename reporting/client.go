package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type endpoint struct {
	URI     string
	AuthURI string
}

func (e *endpoint) getMetaData() metaData {
	var resp, err = http.Get(e.AuthURI + "/registration/")
	if err != nil {
		// TODO logging
		return metaData{}
	}

	var data map[string]interface{}
	var dec = json.NewDecoder(resp.Body)
	dec.Decode(&data)

	log.Printf("%v", data)

	var meta = metaData{
		Company: data["company"].(string),
		TaxID:   data["tax_id"].(string),
	}
	return meta
}

// TODO request data from umsatz api, wrapping them in models

func newEndpoint(URI, authURI string) endpoint {
	return endpoint{
		URI:     URI,
		AuthURI: authURI,
	}
}

type shortDate time.Time

func (d *shortDate) UnmarshalJSON(data []byte) (err error) {
	strDate := string(data)
	time, err := time.Parse("2006-01-02", strDate[1:len(strDate)-1])
	if err != nil {
		d = &shortDate{}
		return nil
	}
	*d = shortDate(time)
	return nil
}

type umsatzFiscalPeriod struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	StartsAt shortDate `json:"startsAt"`
	EndsAt   shortDate `json:"endsAt"`
}
type umsatzAccount struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}
type umsatzPosition struct {
	ID                    int       `json:"id,omitempty"`
	AccountCodeFrom       string    `json:"accountCodeFrom"`
	AccountCodeTo         string    `json:"accountCodeTo"`
	PositionType          string    `json:"type"`
	InvoiceDate           shortDate `json:"invoiceDate"`
	BookingDate           shortDate `json:"bookingDate"`
	InvoiceNumber         string    `json:"invoiceNumber"`
	TotalAmountCents      int       `json:"totalAmountCents"`
	TotalAmountCentsInEur int       `json:"totalAmountCentsEur"`
	Currency              string    `json:"currency"`
	Tax                   int       `json:"tax"`
	Description           string    `json:"description"`
}

func (e *endpoint) FetchAccounts() ([]account, error) {
	resp, err := http.Get(fmt.Sprintf("%v/accounts/", e.URI))
	if err != nil {
		return nil, err
	}

	var dec = json.NewDecoder(resp.Body)
	var uacs []umsatzAccount

	if err := dec.Decode(&uacs); err != nil {
		return nil, err
	}

	var acs []account
	for _, acc := range uacs {
		acs = append(acs, account{
			Code:  acc.Code,
			Label: acc.Label,
		})
	}
	return acs, nil
}

func (e *endpoint) FetchPositions(id string) ([]position, error) {
	resp, err := http.Get(fmt.Sprintf("%v/positions/?fiscal_period_id=%v", e.URI, id))
	if err != nil {
		return nil, err
	}

	var dec = json.NewDecoder(resp.Body)
	var ups []umsatzPosition
	if err := dec.Decode(&ups); err != nil {
		return nil, err
	}

	var ps []position
	for i, p := range ups {
		var amount = p.TotalAmountCents
		if p.Currency != "EUR" {
			amount = p.TotalAmountCentsInEur
		}

		ps = append(ps, position{
			ID:                       i + 1,
			FromAccountCode:          p.AccountCodeFrom,
			ToAccountCode:            p.AccountCodeTo,
			Currency:                 p.Currency,
			TotalAmountCents:         amount,
			OriginalTotalAmountCents: p.TotalAmountCents,
			Tax:           p.Tax,
			InvoiceDate:   time.Time(p.InvoiceDate),
			Description:   p.Description,
			Type:          p.PositionType,
			InvoiceNumber: p.InvoiceNumber,
		})
	}
	return ps, nil
}

func (e *endpoint) FetchFiscalPeriod(id string) (fiscalPeriod, error) {
	resp, err := http.Get(fmt.Sprintf("%v/fiscalPeriods/", e.URI))
	if err != nil {
		return fiscalPeriod{}, err
	}

	var dec = json.NewDecoder(resp.Body)
	var data []umsatzFiscalPeriod
	if err := dec.Decode(&data); err != nil {
		return fiscalPeriod{}, err
	}

	var requestedID, _ = strconv.Atoi(id)
	for _, period := range data {
		if period.ID == requestedID {
			var acs, _ = e.FetchAccounts()
			var ps, _ = e.FetchPositions(id)
			return fiscalPeriod{
				Name:      period.Name,
				StartsAt:  time.Time(period.StartsAt),
				EndsAt:    time.Time(period.EndsAt),
				Accounts:  acs,
				Positions: ps,
				Expenses:  expenses(ps),
				Incomes:   incomes(ps),
			}, nil
		}
	}
	return fiscalPeriod{}, fmt.Errorf("Unknown fiscal period")
}
