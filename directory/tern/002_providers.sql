create table providers
(
    id                    bigserial                 not null
        constraint providers_pk
            primary key,
    created               timestamptz default now() not null,
    updated               timestamptz default now() not null,
    pubkey                text                      not null,
    chain                 text                      not null,
    bond                  numeric                   not null,
    metadata_uri          text,
    metadata_nonce        numeric,
    status                text references provider_status (status),
    min_contract_duration numeric,
    max_contract_duration numeric,
    subscription_rate     numeric,
    paygo_rate            numeric
);

alter table providers
    add constraint pubkey_chain_uniq unique (pubkey, chain);

---- create above / drop below ----

drop table providers;
