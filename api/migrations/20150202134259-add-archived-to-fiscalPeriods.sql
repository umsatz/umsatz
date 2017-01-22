-- +migrate Up
SET client_min_messages = 'warning';

ALTER TABLE public.fiscal_periods
  ADD COLUMN archived boolean NOT NULL default 'f';

-- +migrate Down
ALTER TABLE public.fiscal_periods
  DROP COLUMN archived;