-- +migrate Up
SET client_min_messages = 'warning';

CREATE TABLE IF NOT EXISTS public.fiscal_periods (
  id          SERIAL         PRIMARY KEY,
  name        character varying(64) NOT NULL,
  created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
  starts_at   date           NOT NULL DEFAULT NOW(),
  ends_at     date           NOT NULL DEFAULT NOW()
);

-- +migrate Down
DROP TABLE public.fiscal_periods;
