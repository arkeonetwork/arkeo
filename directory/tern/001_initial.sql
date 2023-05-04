create table provider_status
(
    id     serial not null
        constraint provider_status_pk
            primary key,
    status text   not null unique
);

insert into provider_status(status) values ('ONLINE');
insert into provider_status(status) values ('OFFLINE');

create table contract_types
(
    id  serial not null
        constraint contract_types_pk
            primary key,
    val text   not null unique
);

insert into contract_types(val)
values ('PAY_AS_YOU_GO');
insert into contract_types(val)
values ('SUBSCRIPTION');

create table auth_types
(
    id  serial not null
        constraint auth_types_pk
            primary key,
    val text   not null unique
);

insert into auth_types(val)
values ('STRICT');
insert into auth_types(val)
values ('OPEN');

---- create above / drop below ----
-- undo --
drop table contract_types;
drop table provider_status;
