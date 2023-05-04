create table contracts
(
    id              bigserial                 not null
        constraint contracts_pk
            primary key,
    created         timestamptz default now() not null,
    updated         timestamptz default now() not null,
    provider_id     bigint                    not null references providers (id),
    delegate_pubkey text                      not null check ( delegate_pubkey != '' ),
    client_pubkey   text                      not null check ( client_pubkey != '' ),
    height          bigint                    not null check ( height > 0 ),
    contract_type   text                      not null references contract_types (val),
    duration        bigint                    not null,
    rate            bigint                    not null,
    auth   text                      not null references auth_types (val),
    open_cost       bigint                    not null,
    deposit         bigint                    not null,
    settlement_duration         bigint                    not null,
    queries_per_minute         bigint                    not null
);

alter table contracts
    add constraint pubkey_prov_dlgt_uniq unique (provider_id, delegate_pubkey);

---- create above / drop below ----

drop table contracts;
