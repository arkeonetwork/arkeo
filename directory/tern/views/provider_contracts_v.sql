create or replace view provider_contracts_v as
(
with indexed_height as (select height
                        from indexer_status
                        limit 1)
select p.id as provider_id,
       c.id as contract_id,
       p.pubkey,
       p.service,
       c.delegate_pubkey,
       c.client_pubkey,
       c.height,
       c.contract_type,
       c.duration,
       c.duration-((select height from indexed_height)-c.height) as remaining,
--        c.closed_height,
       c.rate_asset,
       c.rate_amount,
       c.open_cost,
       c.updated,
       c.created
from providers p join contracts c on p.id = c.provider_id
    );
