package main

import (
	"fmt"
	"html/template"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var rxp = regexp.MustCompile(`(\d{3})*(\d{2})$`)

func formatCurrency(amount int, currency string) string {
	formatted := string(rxp.ReplaceAll([]byte(strconv.Itoa(amount)), []byte(".$1,$2")))
	if formatted[0] == '.' {
		formatted = formatted[1:]
	}
	formatted = strings.Replace(formatted, ".,", ",", -1)
	return fmt.Sprintf("%s %v", formatted, currency)
}

func formatDateShort(date time.Time) string {
	return date.Format("02/01/06")
}

func formatDate(date time.Time) string {
	return date.Format("02.01.2006")
}

func now() time.Time {
	return time.Now()
}

var templateHelperFuncs = template.FuncMap{
	"currency":        formatCurrency,
	"shortDateFormat": formatDateShort,
	"dateFormat":      formatDate,
	"now":             now,
}
