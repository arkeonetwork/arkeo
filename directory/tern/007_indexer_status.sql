create table indexer_status
(
    id          numeric not null
        constraint indexer_status_pk
            primary key,
    created     timestamptz default now() not null,
    updated     timestamptz default now() not null,
    height      numeric not null check ( height >= 0 )
);

---- create above / drop below ----
drop table if exists indexer_status;
