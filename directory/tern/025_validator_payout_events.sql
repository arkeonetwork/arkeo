create table validator_payout_events
(
    id        bigserial                 not null
        constraint validator_payout_events_pk
            primary key,
    created   timestamptz default now() not null,
    updated   timestamptz default now() not null,
    validator text                      not null,
    height    bigint                    not null,
    paid      numeric
);

alter table validator_payout_events add constraint validator_payout_evts_validator_height_key unique (validator, height);
---- create above / drop below ----
drop table validator_payout_events;
