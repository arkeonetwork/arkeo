create table contract_settlement_events
(
    id            bigserial                 not null
        constraint contract_settlement_events_pk
            primary key,
    created       timestamptz default now() not null,
    updated       timestamptz default now() not null,
    contract_id   bigint                    not null references contracts (id),
    txid          text                      not null check ( txid != '' ) unique,
    client_pubkey text                      not null check ( client_pubkey != '' ),
    height        numeric                   not null check ( height > 0 ),
    nonce         bigint                    not null,
    paid          bigint                    not null,
    reserve       bigint                    not null
);

create index contract_settle_evts_contract_id_idx on contract_settlement_events (contract_id);

drop view providers_v;
drop view providers_base_v;
{{ template "views/providers_base_v_v2.sql" . }}
{{ template "views/providers_v_v1.sql" . }}

---- create above / drop below ----

-- {{ template "views/providers_base_v_v1.sql" . }}
drop view providers_v;
drop view providers_base_v;
{{ template "views/providers_base_v_v2.sql" . }}
{{ template "views/providers_v_v1.sql" . }}
drop table contract_settlement_events;
