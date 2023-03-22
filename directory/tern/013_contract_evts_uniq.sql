alter table open_contract_events
    add constraint open_contract_events_txid_unq unique (txid);
---- create above / drop below ----

alter table open_contract_events drop constraint open_contract_events_txid_unq;
