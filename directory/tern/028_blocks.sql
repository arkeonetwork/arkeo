create table blocks
(
    id              bigserial                 not null
        constraint blocks_pk
            primary key,
    created         timestamptz default now() not null,
    updated         timestamptz default now() not null,
    height          numeric                   not null check ( height > 0 ) unique,
    hash            text                      not null check ( hash != '' ) unique,
    block_time  timestamptz not null
);

---- create above / drop below ----
drop table blocks;
