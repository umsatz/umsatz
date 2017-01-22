// Umsatz provides a tiny JSON API which handles accounting related informations.
// This includes fiscal periods, positions within fiscal periods and accounts.
package main

//go:generate go-bindata -pkg main -o bindata.go migrations/

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"
	"time"

	"database/sql"
	_ "expvar"

	_ "github.com/umsatz/api/Godeps/_workspace/src/github.com/lib/pq"
	"github.com/umsatz/api/Godeps/_workspace/src/github.com/nicolai86/sql-migrate"
)

type link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

func baseURI(h *http.Header) string {
	baseURI := h.Get("X-Requested-Uri")
	if strings.HasSuffix(baseURI, "/") {
		baseURI = baseURI[:len(baseURI)-1]
	}
	return baseURI
}

func newLink(h *http.Header, rel string, href string) link {
	absoluteHref := fmt.Sprintf("%v%v", baseURI(h), href)
	return link{rel, absoluteHref}
}

func setupDb() *sql.DB {
	database := os.Getenv("DATABASE")
	if database == "" {
		database = "umsatz"
	}

	revDsn := os.Getenv("REV_DSN")
	if revDsn == "" {
		user, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		revDsn = "user=" + user.Username + " dbname=" + database + " sslmode=disable"
	}

	newDb, err := sql.Open("postgres", revDsn)
	if err != nil {
		log.Fatal("failed to connect to postgres", err)
	}
	newDb.SetMaxIdleConns(100)

	return newDb
}

func logHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"%s\t%s\t%s\t%s\t%s",
			r.RemoteAddr,
			time.Now().Format("2006-01-02T15:04:05 -0700"),
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	}
}

func jsonHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	}
}

func routingHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	// TODO(rr): links should be registede by every application, e.g. accounts, periods, positions
	routes := []link{
		newLink(&req.Header, "fiscalPeriods", "/fiscalPeriods"),
		newLink(&req.Header, "positions", "/positions/?fiscal_period_id={fiscalPeriodID}"),
		newLink(&req.Header, "position", "/position/{id}"),
	}

	bytes, err := json.Marshal(routes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	io.WriteString(w, string(bytes))
}

func authHandler(next http.Handler, authAddress string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// nginx adds the X-Requires-Authentication header.
		var authenticationRequired = r.Header.Get("X-Requires-Authentication") != ""
		if !authenticationRequired {
			next.ServeHTTP(w, r)
			return
		}

		var sessionKey = r.Header.Get("X-UMSATZ-SESSION")

		var resp, err = http.Get("http://" + authAddress + "/validate/" + sessionKey)
		if err != nil {
			log.Printf("err: %v", err)
		}
		if err == nil && resp.StatusCode == http.StatusOK {
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
	}
}

func newUmsatzServer(db *sql.DB, currencyAddress string) http.Handler {
	r := http.DefaultServeMux

	r.Handle("/fiscalPeriods/", http.StripPrefix("/fiscalPeriods/", newFiscalPeriodApp(db)))
	r.Handle("/positions/", http.StripPrefix("/positions/", newPositionsApp(db, currencyAddress)))
	r.Handle("/accounts/", http.StripPrefix("/accounts/", newAccountingApp(db)))

	r.Handle("/", http.HandlerFunc(routingHandler))

	return http.Handler(logHandler(jsonHandler(r)))
}

func runMigrations(db *sql.DB) migrate.MigrationSource {
	migrate.SetTable("gorp_migrations")

	migrations := &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "migrations",
	}

	if _, err := migrate.Exec(db, "postgres", migrations, migrate.Up); err != nil {
		log.Printf("unable to migrate: %v", err)
	}

	return migrations
}

// Set by make file on build
var (
	Version string
	Commit  string
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	var (
		httpAddress         = flag.String("http.addr", ":8080", "listen address")
		currencyHTTPAddress = flag.String("currency.addr", "127.0.0.1:8081", "currency listen address")
		authHTTPAddress     = flag.String("auth.addr", "127.0.0.1:8082", "auth listen address")
		printVersion        = flag.Bool("version", false, "print version and exit")
	)
	flag.Parse()

	if *printVersion {
		fmt.Printf("%s", Version)
		os.Exit(0)
	}

	db := setupDb()
	runMigrations(db)

	log.Printf("listening on %s", *httpAddress)
	log.Fatal(http.ListenAndServe(*httpAddress, authHandler(newUmsatzServer(db, *currencyHTTPAddress), *authHTTPAddress)))
}
