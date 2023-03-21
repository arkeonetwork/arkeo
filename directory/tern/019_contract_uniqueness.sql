alter table contracts drop constraint pubkey_prov_dlgt_uniq;

alter table contracts
    add constraint contracts_provider_delegate_height_key unique (provider_id, delegate_pubkey, height);

---- create above / drop below ----
alter table contracts drop constraint contracts_provider_delegate_height_key;
