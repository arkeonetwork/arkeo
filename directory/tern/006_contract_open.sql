create table open_contract_events
(
    id            bigserial                 not null
        constraint open_contract_events_pk
            primary key,
    created       timestamptz default now() not null,
    updated       timestamptz default now() not null,
    contract_id   bigint                    not null references contracts (id),
    txid          text                      not null check ( txid != '' ),
    client_pubkey text                      not null check ( client_pubkey != '' ),
    contract_type text                      not null references contract_types (val),
    duration      bigint                    not null,
    rate          bigint                    not null,
    open_cost     bigint                    not null,
    height        numeric                   not null check ( height > 0 )
);

create index open_cntrc_evts_cntrc_id_idx on open_contract_events (contract_id);

---- create above / drop below ----
drop table open_contract_events;
