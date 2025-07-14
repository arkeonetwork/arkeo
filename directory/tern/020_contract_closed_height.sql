alter table contracts add column settlement_height bigint not null default 0;
---- create above / drop below ----
drop view if exists open_contracts_v;
alter table contracts drop column settlement_height;
