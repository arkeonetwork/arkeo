CREATE OR REPLACE VIEW public.providers_v AS
SELECT
    b.id,
    b.pubkey,
    b.service,
    b.bond,
    b.metadata_uri,
    b.metadata_nonce,
    b.status,
    b.min_contract_duration,
    b.max_contract_duration,
    b.created,
    b.updated,
    b.metadata_nonce_value,
    b.metadata_version,
    b.metadata_moniker,
    b.metadata_website,
    b.metadata_description,
    b.metadata_location,
    b.metadata_free_rate_limit,
    b.metadata_free_rate_limit_duration,
    b.contract_count,
    b.birth_height,
    b.cur_height,
    b.total_paid,
    (b.cur_height - (b.birth_height)::numeric) AS age
FROM
    providers_base_v b;