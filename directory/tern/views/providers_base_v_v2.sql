CREATE OR REPLACE VIEW public.providers_base_v AS
 WITH indexed_height AS (
         SELECT indexer_status.height
           FROM indexer_status
         LIMIT 1
        )
SELECT p.id,
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
       m.nonce AS metadata_nonce_value,
       m.version AS metadata_version,
       m.moniker AS metadata_moniker,
       m.website AS metadata_website,
       m.description AS metadata_description,
       m.location AS metadata_location,
       m.free_rate_limit AS metadata_free_rate_limit,
       m.free_rate_limit_duration AS metadata_free_rate_limit_duration,
       ( SELECT count(1) AS count
FROM contracts oc
WHERE (oc.provider_id = p.id)) AS contract_count,
    ( SELECT min(bond_evts.height) AS min
FROM provider_bond_events bond_evts
WHERE (bond_evts.provider_id = p.id)) AS birth_height,
    ( SELECT indexed_height.height
FROM indexed_height) AS cur_height,
    COALESCE(( SELECT sum(settle_events.paid) AS sum
    FROM (contracts c
    JOIN contract_settlement_events settle_events ON ((c.id = settle_events.contract_id)))
    WHERE (c.provider_id = p.id)), (0)::numeric) AS total_paid
FROM (providers p
    LEFT JOIN provider_metadata m ON ((m.provider_id = p.id)));