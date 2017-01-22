package main

import (
	"database/sql"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/umsatz/api/Godeps/_workspace/src/github.com/lib/pq"
	"github.com/umsatz/api/Godeps/_workspace/src/github.com/nicolai86/sql-migrate"
)

type testManager struct {
	db     *sql.DB
	server *httptest.Server
}

func (t *testManager) Clear() {
	t.db.Exec(`TRUNCATE fiscal_periods, positions, accounts`)
}

var TestManager testManager

func TestMain(m *testing.M) {
	os.Setenv("DATABASE", "umsatz_test")
	db := setupDb()
	TestManager = testManager{
		db:     db,
		server: httptest.NewServer(newUmsatzServer(db, "http://127.0.0.1:9000")),
	}
	migrations := runMigrations(db)

	TestManager.Clear()
	ret := m.Run()

	if _, err := migrate.Exec(db, "postgres", migrations, migrate.Down); err != nil {
		panic(err)
	}
	TestManager.db.Close()
	os.Exit(ret)
}
