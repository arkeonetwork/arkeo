CREATE TABLE provider_subscription_rates (
    id SERIAL PRIMARY KEY,
    provider_id BIGINT NOT NULL REFERENCES providers (id),
    token_name TEXT NOT NULL,
    token_amount NUMERIC NOT NULL,
    UNIQUE (provider_id, token_name)
);

ALTER TABLE providers DROP COLUMN subscription_rate;

CREATE TABLE provider_pay_as_you_go_rates (
    id SERIAL PRIMARY KEY,
    provider_id BIGINT NOT NULL REFERENCES providers (id),
    token_name TEXT NOT NULL,
    token_amount NUMERIC NOT NULL,
    UNIQUE (provider_id, token_name)
);

ALTER TABLE providers DROP COLUMN paygo_rate;

-- 1) Generic events table
CREATE TABLE IF NOT EXISTS public.generic_events (
     id          BIGSERIAL     PRIMARY KEY,
     height      BIGINT        NOT NULL,            -- block height
     txid        TEXT          NULL,                -- tx hash (hex) or NULL for BeginBlock/EndBlock
     event_type  TEXT          NOT NULL,            -- e.g. "burn", "transfer", etc.
     attributes  JSONB         NOT NULL,            -- raw list of {key,value,index} objects
     created_at  TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS generic_events_id_idx
    ON public.generic_events(id);

CREATE INDEX IF NOT EXISTS generic_events_event_type_idx
    ON public.generic_events(event_type);

CREATE INDEX IF NOT EXISTS generic_events_height_idx
    ON public.generic_events(height);

-- Rebuild provider views now that subscription-rate table exists
DROP VIEW IF EXISTS public.providers_v;
DROP VIEW IF EXISTS public.providers_base_v;
{{ template "views/providers_base_v_v2.sql" . }}
{{ template "views/providers_v_v1.sql" . }}

---- create above / drop below ----

drop table provider_subscription_rates;
drop table provider_pay_as_you_go_rates;
