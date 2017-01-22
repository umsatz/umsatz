package main

import "time"

type metaData struct {
	Company string `json:"company"`
	TaxID   string `json:"tax_id"`
}
type fiscalPeriod struct {
	Name      string
	StartsAt  time.Time
	EndsAt    time.Time
	Accounts  []account
	Positions []position
	Expenses  *positionGroup
	Incomes   *positionGroup
}
type account struct {
	Code  string
	Label string
}
type position struct {
	ID                       int
	InvoiceDate              time.Time
	InvoiceNumber            string
	FromAccountCode          string
	ToAccountCode            string
	Currency                 string
	TotalAmountCents         int
	OriginalTotalAmountCents int
	Tax                      int
	Type                     string
	Description              string
}
type positionGroup struct {
	positions []position
}

func (pg *positionGroup) TotalWithTax() int {
	var total = 0
	for _, p := range pg.positions {
		total = total + p.TotalAmountCents
	}
	return total
}

func (pg *positionGroup) TotalWithoutTax() int {
	var total = 0
	for _, p := range pg.positions {
		total = total + p.TotalAmountCents - int(float32(p.TotalAmountCents)/float32(float32(p.Tax)/float32(100)+1.0))
	}
	return pg.TotalWithTax() - total
}

func (p position) TotalWithoutTax() int {
	return int(float32(p.TotalAmountCents) / (100.0 + float32(p.Tax)) * 100.0)
}

func (p position) Sign() string {
	if p.Type == "expense" {
		return "-"
	}
	return ""
}

func (f *fiscalPeriod) TotalWithTax() int {
	return f.Incomes.TotalWithTax() - f.Expenses.TotalWithTax()
}
func (f *fiscalPeriod) TotalWithoutTax() int {
	return f.Incomes.TotalWithoutTax() - f.Expenses.TotalWithoutTax()
}
func (f *fiscalPeriod) Sign(accountCode string) string {
	ps := f.PositionsWithAccount(accountCode)
	allExpense := true
	for _, p := range ps {
		allExpense = allExpense && p.Type == "expense"
	}
	if allExpense {
		return "-"
	}
	return ""
}

func (f *fiscalPeriod) UsedAccounts() []account {
	var acs []account
	for i := range f.Accounts {
		if f.CountPositionsWithAccount(f.Accounts[i].Code) > 0 {
			acs = append(acs, f.Accounts[i])
		}
	}
	return acs
}

func expenses(ps []position) *positionGroup {
	var pg = new(positionGroup)
	for _, p := range ps {
		if p.Type == "expense" {
			pg.positions = append(pg.positions, p)
		}
	}
	return pg
}
func incomes(ps []position) *positionGroup {
	var pg = new(positionGroup)
	for _, p := range ps {
		if p.Type == "income" {
			pg.positions = append(pg.positions, p)
		}
	}
	return pg
}

func (f *fiscalPeriod) CountPositionsWithAccount(accountCode string) int {
	return len(f.PositionsWithAccount(accountCode))
}

func (f *fiscalPeriod) PositionsWithAccount(accountCode string) []position {
	var ps []position
	for _, p := range f.Positions {
		if p.ToAccountCode == accountCode {
			ps = append(ps, p)
		}
	}
	return ps
}

func (f *fiscalPeriod) TotalAmountFromAccount(accountCode string) int {
	total := 0
	for _, position := range f.PositionsWithAccount(accountCode) {
		total += position.TotalAmountCents
	}
	return total
}

func (f *fiscalPeriod) KnownAccount(accountCode string) bool {
	for _, a := range f.Accounts {
		if a.Code == accountCode {
			return true
		}
	}
	return false
}

func (f *fiscalPeriod) AccountByCode(accountCode string) account {
	for _, a := range f.Accounts {
		if a.Code == accountCode {
			return a
		}
	}
	return account{}
}
