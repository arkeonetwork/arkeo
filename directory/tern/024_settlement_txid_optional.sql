alter table contract_settlement_events drop constraint contract_settlement_events_txid_key;
alter table contract_settlement_events add constraint contract_settlement_contract_nonce_key unique (contract_id, nonce);
alter table contract_settlement_events drop constraint contract_settlement_events_txid_check;

---- create above / drop below ----
truncate contract_settlement_events;
alter table contract_settlement_events add constraint contract_settlement_events_txid_check check ( txid != '' ); 
alter table contract_settlement_events drop constraint contract_settlement_contract_nonce_key;
alter table contract_settlement_events add constraint contract_settlement_events_txid_key unique (txid);
