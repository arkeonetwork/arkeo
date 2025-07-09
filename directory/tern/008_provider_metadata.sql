BEGIN;

CREATE TABLE IF NOT EXISTS public.provider_metadata (
    id bigserial PRIMARY KEY,
    created timestamptz DEFAULT now() NOT NULL,
    updated timestamptz DEFAULT now() NOT NULL,
    provider_id bigint NOT NULL,
    nonce numeric NOT NULL, -- (Remove if you don't need nonce)
    "version" text,
    moniker text,
    website text,
    description text,
    "location" point,
    port text,
    source_chain text,
    event_stream_host text,
    claim_store_location text,
    free_rate_limit bigint,
    free_rate_limit_duration bigint,
    subscribe_rate_limit bigint,
    subscribe_rate_limit_duration bigint,
    paygo_rate_limit bigint,
    paygo_rate_limit_duration bigint
    );

ALTER TABLE IF EXISTS public.provider_metadata
    ADD CONSTRAINT provider_metadata_provider_id_fkey
    FOREIGN KEY (provider_id)
    REFERENCES public.providers (id);

ALTER TABLE IF EXISTS public.provider_metadata
    ADD CONSTRAINT provider_metadata_provider_id_uniq
    UNIQUE (provider_id);

---- create above / drop below ----

DROP TABLE IF EXISTS provider_metadata;