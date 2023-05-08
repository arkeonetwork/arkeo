create or replace view providers_base_v as
(
with indexed_height as (select height
                        from indexer_status
                        limit 1)

select p.id,
       p.pubkey,
       p.service,
       p.bond,
       p.metadata_uri,
       p.metadata_nonce,
       p.status,
       p.min_contract_duration,
       p.max_contract_duration,
       p.created,
       p.updated,
       (select count(1) from contracts oc where oc.provider_id = p.id)        as contract_count,
    --    (select count(1) from open_contracts_v oc where oc.provider_id = p.id) as open_contract_count,
       (select min(bond_evts.height)
        from provider_bond_events bond_evts
        where bond_evts.provider_id = p.id)                                   as birth_height,
       (select indexed_height.height from indexed_height)                        cur_height,
        (
            select sum(settle_events.paid)
            from contracts c
                join contract_settlement_events settle_events on c.id = settle_events.contract_id
            where c.provider_id = p.id
        ) as total_paid
from providers p
);
