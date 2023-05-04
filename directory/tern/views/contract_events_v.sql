create or replace view contract_events_v as
(
with evts as (
    select id, contract_id, txid, created, height, 'open_contract' as evt_name
    from open_contract_events
    union all
    select id, contract_id, txid, created, height, 'close_contract'
    from close_contract_events
    union all
    select id, contract_id, txid, created, height, 'contract_settlement'
    from contract_settlement_events)
select evts.created,
       evts.height,
       evts.evt_name,
       evts.txid,
       p.service,
       p.id     as provider_id,
       c.id     as contract_id,
       evts.id as event_id,
       p.pubkey as provider_pubkey,
       c.client_pubkey,
       c.delegate_pubkey
from providers p
         join contracts c on p.id = c.provider_id
         join evts on c.id = evts.contract_id
);
