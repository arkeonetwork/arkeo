create table provider_bond_events
(
    id          bigserial                 not null
        constraint provider_bond_events_pk
            primary key,
    created     timestamptz default now() not null,
    updated     timestamptz default now() not null,
    provider_id bigint                    not null references providers (id),
    txid        text                      not null check ( txid != '' ),
    bond_rel    numeric                   not null,
    bond_abs    numeric                   not null
);

create index prov_bond_evts_prov_id_idx on provider_bond_events (provider_id);

---- create above / drop below ----

drop table if exists provider_bond_events;
