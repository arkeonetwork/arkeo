alter table provider_bond_events
    add constraint provider_bond_events_txid_unq
        unique (txid);


alter table provider_mod_events
    add constraint provider_mod_events_txid_unq
        unique (txid);

---- create above / drop below ----
alter table provider_mod_events drop constraint provider_mod_events_txid_unq;
alter table provider_bond_events drop constraint provider_bond_events_txid_unq;
