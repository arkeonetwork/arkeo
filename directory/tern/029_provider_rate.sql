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

---- create above / drop below ----

drop table provider_subscription_rates;
drop table provider_pay_as_you_go_rates;
