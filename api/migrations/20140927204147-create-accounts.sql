-- +migrate Up

SET client_min_messages = 'warning';

CREATE TABLE IF NOT EXISTS public.accounts (
  id         SERIAL                PRIMARY KEY,
  code       character varying(5),
  label      character varying(16) NOT NULL,

  CONSTRAINT uniq_account_code UNIQUE(code)
);

-- +migrate Down
DROP TABLE public.accounts CASCADE;