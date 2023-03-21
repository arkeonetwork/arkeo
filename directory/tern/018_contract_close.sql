create table close_contract_events
(
    id              bigserial                 not null
        constraint close_contract_events_pk
            primary key,
    created         timestamptz default now() not null,
    updated         timestamptz default now() not null,
    contract_id     bigint                    not null references contracts (id),
    txid            text                      not null check ( txid != '' ) unique,
    client_pubkey   text                      not null check ( client_pubkey != '' ),
    delegate_pubkey text                      not null, -- can be ''
    height          numeric                   not null check ( height > 0 )
);

create index close_contract_evts_contract_id_idx on close_contract_events (contract_id);

---- create above / drop below ----
drop table close_contract_events;
