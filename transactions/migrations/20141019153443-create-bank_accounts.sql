-- +migrate Up
CREATE TABLE public.bank_accounts (
  id                    SERIAL PRIMARY KEY,

  -- aqbanking attribute of transactions
  bank_code             CHARACTER VARYING(32) NOT NULL,
  account_number        CHARACTER VARYING(32) NOT NULL,
  iban                  CHARACTER VARYING(64),
  bic                   CHARACTER VARYING(20),
  name                  CHARACTER VARYING(64)
);

-- +migrate Down
DROP TABLE public.bank_accounts;
