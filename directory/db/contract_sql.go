package db

const (
	sqlUpsertIndexerStatus = `
		INSERT INTO indexer_status (id, height)
		VALUES (1, $1)
		ON CONFLICT (id) DO UPDATE SET
			height = EXCLUDED.height,
			updated = now()
		RETURNING id, created, updated
	`

	sqlGetContractByID = ` select c.id,
	c.created,
	c.updated,
	p.pubkey as provider,
	p.service,
	c.delegate_pubkey,
	c.client_pubkey,
	c.height,
	c.contract_type,
	c.duration,
	c.rate_asset,
	c.rate_amount,
	c.open_cost,
	c.deposit,
	c.auth,
	c.queries_per_minute,
	c.settlement_duration,
	c.paid,
	c.reserve_contrib_asset,
	c.reserve_contrib_usd,
	c.settlement_height,
	c.provider_id
	from contracts c 
	left outer join providers p on p.id = c.provider_id
		where c.id = $1`

	sqlUpsertContract = `
		insert into contracts(provider_id,delegate_pubkey,client_pubkey,contract_type,duration,rate_asset,rate_amount,open_cost,height,deposit,settlement_duration,settlement_height,auth,queries_per_minute,id)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		on conflict on constraint contracts_pk
		do update set provider_id=$1,delegate_pubkey=$2,client_pubkey=$3,contract_type = $4, duration = $5, rate_asset = $6, rate_amount = $7, open_cost = $8, height=$9, deposit = $10, settlement_duration = $11, settlement_height = $12, auth = $13, queries_per_minute = $14, updated = now() 
		where contracts.id = $15
		returning id, created, updated
	`

	sqlCloseContract = `
		update contracts
		set settlement_height = $1
		where id = $2
		returning id, created, updated
	`

	sqlUpsertCloseContractEventRecord = `
		INSERT INTO close_contract_events (
			contract_id,
			txid,
			client_pubkey,
			delegate_pubkey,
			height
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			contract_id = EXCLUDED.contract_id,
			txid = EXCLUDED.txid,
			client_pubkey = EXCLUDED.client_pubkey,
			delegate_pubkey = EXCLUDED.delegate_pubkey,
			height = EXCLUDED.height,
			updated = now()
		returning id, created, updated
	`

	sqlUpsertContractSettlementEvent = `
		UPDATE contracts
		SET nonce = $1, paid = paid + $2, reserve_contrib_asset = reserve_contrib_asset + $3, settlement_height = $5
		WHERE id = $4
		returning id, created, updated
	`

	sqlUpsertContractSettlementEventRecord = `
		INSERT INTO contract_settlement_events (
			contract_id,
			txid,
			client_pubkey,
			height,
			nonce,
			paid,
			reserve
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			contract_id = EXCLUDED.contract_id,
			txid = EXCLUDED.txid,
			client_pubkey = EXCLUDED.client_pubkey,
			height = EXCLUDED.height,
			nonce = EXCLUDED.nonce,
			paid = EXCLUDED.paid,
			reserve = EXCLUDED.reserve,
			updated = now()
		returning id, created, updated
	`

	sqlInsertOpenContractEventRecord = `
		INSERT INTO open_contract_events (
			contract_id,
			txid,
			client_pubkey,
			contract_type,
			duration,
			rate,
			open_cost,
			height,
			deposit,
			settlement_duration,
			auth,
			queries_per_minute
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		returning id, created, updated
	`
)
