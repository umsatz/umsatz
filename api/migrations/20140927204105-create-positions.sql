-- +migrate Up
SET client_min_messages = 'warning';

-- +migrate StatementBegin
DO $$
BEGIN

IF NOT EXISTS (
  select 1 from pg_type where typname = 'position_type'
) THEN
  CREATE TYPE  position_type AS ENUM ('income', 'expense');
END IF;

END
$$;
-- +migrate StatementEnd

-- +migrate StatementBegin
DO $$
BEGIN

IF NOT EXISTS (
  select 1 from pg_type where typname = 'position_currency'
) THEN
  CREATE TYPE  position_currency AS ENUM ('EUR', 'USD', 'GBP');
END IF;

END
$$;
-- +migrate StatementEnd

CREATE TABLE IF NOT EXISTS public.positions (
  id                        SERIAL                 PRIMARY KEY,
  account_code_from         character varying(5)   NOT NULL,
  account_code_to           character varying(5)   NOT NULL,
  type                      position_type          NOT NULL,
  invoice_date              TIMESTAMPTZ            NOT NULL,
  booking_date              TIMESTAMPTZ            DEFAULT NULL,
  invoice_number            character varying(32)  NOT NULL,
  total_amount_cents        int                    NOT NULL DEFAULT 0,
  total_amount_cents_in_eur int                    NOT NULL DEFAULT 0,
  currency                  position_currency      NOT NULL,
  tax                       int                    NOT NULL,
  fiscal_period_id          int                    NOT NULL,
  attachment_path           character varying(256),
  description               text,
  created_at                TIMESTAMPTZ            NOT NULL DEFAULT NOW(),
  updated_at                TIMESTAMPTZ            NOT NULL DEFAULT NOW(),

  CONSTRAINT fiscalPeriodfk FOREIGN KEY (fiscal_period_id) REFERENCES fiscal_periods (id) MATCH FULL
);

-- +migrate Down
DROP TABLE public.positions;

DROP TYPE IF EXISTS position_type;
DROP TYPE IF EXISTS position_currency;
