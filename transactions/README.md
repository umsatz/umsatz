# transactions

golang json api to manage &amp; retrieve bank transactions

## TODO

- listTransactions should use the database table to return transactions
- add refresh binary:
  - move server (current main.go) into api package
  - refresh binary into worker package
- add tests
- replace gorilla.mux /w http.ServerMux