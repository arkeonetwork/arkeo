package db

var (
	sqlInsertProvider = `
		insert into providers(pubkey,service,bond) values ($1,$2,$3) returning id, created, updated
	`

	sqlUpdateProvider = `
		update providers
		set bond = $3,
		    metadata_uri = $4,
			metadata_nonce = $5,
			status = $6,
			min_contract_duration = $7,
			max_contract_duration = $8,
			settlement_duration = $9,
			updated = now()
		where pubkey = $1
		  and service = $2
		returning id, created, updated
	`

	sqlFindProvider = `
		select 
			id,
			created,
			updated,
			pubkey,
			service,
			coalesce(bond,0) as bond,
			coalesce(metadata_uri,'') as metadata_uri,
			coalesce(metadata_nonce,0) as metadata_nonce,
			coalesce(status,'OFFLINE') as status,
			coalesce(min_contract_duration,-1) as min_contract_duration,
			coalesce(max_contract_duration,-1) as max_contract_duration,
			coalesce(settlement_duration,-1) as settlement_duration
		from providers p
		where p.pubkey = $1
		  and p.service = $2
	`
	sqlInsertBondProviderEvent = `insert into provider_bond_events(provider_id,height,txid,bond_rel,bond_abs) values ($1,$2,$3,$4,$5)
		on conflict on constraint provider_bond_events_txid_unq
		do update set updated = now()
		where provider_bond_events.txid = $3
		returning id, created, updated
	`
	sqlInsertModProviderEvent = `insert into provider_mod_events(provider_id,height,txid,metadata_uri,metadata_nonce,status,min_contract_duration,max_contract_duration)
		values ($1,$2,$3,$4,$5,$6,$7,$8)
		on conflict on constraint provider_mod_events_txid_unq
		do update set updated = now()
		where provider_mod_events.txid = $3
		returning id, created, updated
	`

	sqlUpsertProviderMetadata = `INSERT INTO provider_metadata (
			provider_id, nonce, moniker, website, description, location, free_rate_limit
		) VALUES (
			$1, $2, $3, $4, $5, CAST(NULLIF($6, '') AS point), $7
		)
		ON CONFLICT (provider_id)
		DO UPDATE SET
			nonce = EXCLUDED.nonce,
			moniker = EXCLUDED.moniker,
			website = EXCLUDED.website,
			description = EXCLUDED.description,
			location = EXCLUDED.location,
			free_rate_limit = EXCLUDED.free_rate_limit,
			updated = now()
		RETURNING id, created, updated;
	`

	sqlUpsertValidatorPayoutEvent = `insert into validator_payout_events(validator,height,paid)
	values ($1,$2,$3)
	on conflict on constraint validator_payout_evts_validator_height_key
	do update set updated = now()
	where validator_payout_events.validator = $1
	  and validator_payout_events.height = $2
	returning id, created, updated
	`

	sqlDeleteSubscriptionRates = `
		DELETE FROM provider_subscription_rates
		WHERE provider_id IN (
			SELECT id
			FROM providers
			WHERE pubkey = $1 AND service = $2
		);
	`

	sqlInsertSubscriptionRates = `
		INSERT INTO provider_subscription_rates (provider_id, token_name, token_amount) VALUES
	`

	sqlFindProviderSubscriptionRates = `
		SELECT * FROM provider_subscription_rates
        WHERE provider_id = $1
	`

	sqlDeletePayAsYouGoRates = `
		DELETE FROM provider_pay_as_you_go_rates
		WHERE provider_id IN (
			SELECT id
			FROM providers
			WHERE pubkey = $1 AND service = $2
		);
	`

	sqlInsertPayAsYouGoRates = `
		INSERT INTO provider_pay_as_you_go_rates (provider_id, token_name, token_amount) VALUES
	`

	sqlFindProviderPayAsYouGoRates = `
		SELECT * FROM provider_pay_as_you_go_rates
        WHERE provider_id = $1
	`

	sqlFindSubscriberContractsByService = `
		SELECT
		  COALESCE(ocv.id, 0) as contract_id,
		  ocv.created AS contract_created,
		  ocv.updated AS contract_updated,
		  COALESCE(ocv.provider_id, 0) as provider_id,
		  COALESCE(ocv.delegate_pubkey, '') as delegate_pubkey,
		  COALESCE(ocv.client_pubkey, '') as client_pubkey,
		  COALESCE(ocv.height, 0) as height,
		  COALESCE(ocv.contract_type, '') as contract_type,
		  COALESCE(ocv.duration, 0) as duration,
		  COALESCE(ocv.rate_asset, '')as rate_asset,
		  COALESCE(ocv.rate_amount, 0) as rate_amount,
		  COALESCE(ocv.auth, '') as auth,
		  COALESCE(ocv.open_cost, 0) as open_cost,
		  COALESCE(ocv.deposit, 0) as deposit,
		  COALESCE(ocv.queries_per_minute, 0) as queries_per_minute,
		  COALESCE(ocv.settlement_duration, 0) as settlement_duration,
		  COALESCE(ocv.nonce, 0) as nonce,
		  COALESCE(ocv.paid, 0) as paid,
		  COALESCE(ocv.reserve_contrib_asset, 0) as reserve_contrib_asset,
		  COALESCE(ocv.reserve_contrib_usd, 0) as reserve_contrib_usd,
		  COALESCE(ocv.settlement_height, 0) as settlement_height,
		  COALESCE(ocv.start_height, 0) as start_height,
		  COALESCE(ocv.current_height, 0) as current_height,
		  COALESCE(ocv.remaining, 0) as remaining,
		  COALESCE(pv.pubkey, '') as pubkey,
		  COALESCE(pv.service, '') as service,
		  COALESCE(pv.bond, 0) as bond,
		  COALESCE(pv.metadata_uri, '') as metadata_uri,
		  COALESCE(pv.metadata_nonce, 0) as metadata_nonce,
		  COALESCE(pv.status, '') AS provider_status,
		  COALESCE(pv.min_contract_duration, 0) as min_contract_duration,
		  COALESCE(pv.max_contract_duration, 0) as max_contract_duration,
		  pv.created AS provider_created,
		  pv.updated AS provider_updated,
		  COALESCE(pv.metadata_nonce_value, 0) as metadata_nonce_value,
		  COALESCE(pv.metadata_version, '') as metadata_version,
		  COALESCE(pv.metadata_moniker, '') as metadata_moniker,
		  COALESCE(pv.metadata_website, '') as metadata_website,
		  COALESCE(pv.metadata_description, '') as metadata_description,
		  COALESCE(pv.metadata_location::text, '(0,0)') as metadata_location,
		  COALESCE(pv.metadata_free_rate_limit, 0) as metadata_free_rate_limit,
		  COALESCE(pv.metadata_free_rate_limit_duration, 0) as metadata_free_rate_limit_duration,
		  COALESCE(pv.contract_count, 0) as contract_count,
		  COALESCE(pv.birth_height, 0) as birth_height,
		  COALESCE(pv.cur_height, 0) as cur_height,
		  COALESCE(pv.total_paid, 0) as total_paid,
		  COALESCE(pv.age, 0) as age
		FROM open_contracts_v ocv
		INNER JOIN providers_v pv ON ocv.provider_id = pv.id
		WHERE ocv.client_pubkey = $1
		AND pv.service = $2
	`
)
