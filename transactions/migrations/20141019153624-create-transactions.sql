-- +migrate Up
CREATE TABLE public.transactions (
  id                    SERIAL PRIMARY KEY,
  purpose               TEXT,
  description           TEXT,
  status                CHARACTER VARYING(32)  NOT NULL,
  created_at            TIMESTAMPTZ NOT NULL,
  valuta_date           TIMESTAMPTZ,
  mandate_reference     CHARACTER VARYING(64) NOT NULL,
  customer_reference    CHARACTER VARYING(64) NOT NULL,
  total_amount_cents    INT NOT NULL,
  total_amount_currency CHARACTER VARYING(3) NOT NULL,
  fee_amount_cents      INT,
  fee_amount_currency   CHARACTER VARYING(3),

  local_bank_id         INT NOT NULL,
  remote_bank_id        INT NOT NULL,

  CONSTRAINT localbankfk
    FOREIGN KEY (local_bank_id)
    REFERENCES bank_accounts (id) MATCH FULL,

  CONSTRAINT remotebankfk
    FOREIGN KEY (remote_bank_id)
    REFERENCES bank_accounts (id) MATCH FULL
);

-- +migrate Down
DROP TABLE public.transactions;
