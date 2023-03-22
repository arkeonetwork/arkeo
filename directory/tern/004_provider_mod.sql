create table provider_mod_events
(
    id                    bigserial                 not null
        constraint provider_mod_events_pk
            primary key,
    created               timestamptz default now() not null,
    updated               timestamptz default now() not null,
    provider_id           bigint                    not null references providers (id),
    txid                  text                      not null check ( txid != '' ),
    metadata_uri          text check ( metadata_uri != '' ),
    metadata_nonce        numeric check ( metadata_nonce >= 0 ),
    status                text references provider_status (status),
    min_contract_duration numeric,
    max_contract_duration numeric,
    subscription_rate     numeric,
    paygo_rate            numeric
);

create index prov_mod_evts_prov_id_idx on provider_mod_events (provider_id);

---- create above / drop below ----

drop table provider_mod_events;
