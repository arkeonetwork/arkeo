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

-- Rebuild provider views now that subscription-rate table exists
DROP VIEW IF EXISTS public.providers_v;
DROP VIEW IF EXISTS public.providers_base_v;
{{ template "views/providers_base_v_v2.sql" . }}
{{ template "views/providers_v_v1.sql" . }}

---- create above / drop below ----

drop table provider_subscription_rates;
drop table provider_pay_as_you_go_rates;
