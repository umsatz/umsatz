# umsatz-api

[![Build Status](https://travis-ci.org/umsatz/api.svg)](https://travis-ci.org/umsatz/api)

fiscalPeriod api for umsatz

## Tests

To run the testsuite you need to have a PostgreSQL server running & deployed.
Umsatz uses [sql-migrate][1] for schema management. Thus you need to create the
required database upfront:

``` bash
$ createdb umsatz
```

## Executing

``` bash
$ DATABASE=umsatz go run umsatz.go -http.addr=:8080
or
$ REV_DSN="user=postgres database=foo" go run umsatz.go -http.addr=:8080
```

## Remove build artifacts

go clean -i -r

[1]:https://github.com/nicolai86/sql-migrate