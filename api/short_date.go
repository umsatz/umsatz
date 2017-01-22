package main

import (
	"encoding/json"
	"time"

	"database/sql/driver"
)

type shortDate time.Time

func (d shortDate) MarshalJSON() ([]byte, error) {
	if time.Time(d).Format("2006-01-02") == "0001-01-01" {
		return []byte{34, 34}, nil
	}
	return json.Marshal(time.Time(d).Format("2006-01-02"))
}

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

// Scan implements the database/sql Scanner interface
func (d *shortDate) Scan(value interface{}) error {
	if date, ok := value.(time.Time); ok {
		*d = shortDate(date)
	}
	return nil
}

// Value implements the database/sql Valuer interface
func (d shortDate) Value() (driver.Value, error) {
	return driver.Value(time.Time(d)), nil
}
