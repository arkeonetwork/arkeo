alter table provider_bond_events
    add height bigint not null check ( height > 0 );
alter table provider_mod_events
    add height bigint not null check ( height > 0 );

---- create above / drop below ----
alter table provider_mod_events drop column height;
alter table provider_bond_events drop column height;
