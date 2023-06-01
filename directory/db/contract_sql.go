package db

const (
	contractCols = `
	c.id,
	c.created,
	c.updated,
	c.provider_id,
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
	c.closed_height
	`
)

const (
	sqlFindContract = ` select ` + contractCols + `
	from contracts c
		where c.id = $1
	`

	sqlFindContractsByPubKeys = `select ` + contractCols + `
	-- c.id,
	-- c.created,
	-- c.updated,
	-- c.provider_id,
	-- c.delegate_pubkey,
	-- c.client_pubkey,
	-- c.contract_type,
	-- c.duration,
	-- c.rate_asset,
	-- c.rate_amount,
	-- c.open_cost,
	-- c.deposit,
	-- c.auth,
	-- c.queries_per_minute,
	-- c.settlement_duration,
	-- c.paid,
	-- c.reserve_contrib_asset,
	-- c.reserve_contrib_usd
	from providers p join contracts c on p.id = c.provider_id
	where p.service = $1 and p.pubkey = $2 and c.delegate_pubkey = $3
	order by c.id desc
	`
	sqlFindContractByPubKeys = `select ` + contractCols + `
	from providers p join contracts c on p.id = c.provider_id
	where p.service = $1 and p.pubkey = $2 and c.delegate_pubkey = $3 and c.height = $4
	`

	sqlUpsertContract = `
		insert into contracts(provider_id,delegate_pubkey,client_pubkey,contract_type,duration,rate_asset,rate_amount,open_cost,height,deposit,settlement_duration,auth,queries_per_minute,id)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		on conflict on constraint contracts_provider_delegate_height_key
		do update set contract_type = $4, duration = $5, rate_asset = $6, rate_amount = $7, open_cost = $8, deposit = $10, settlement_duration = $11, auth = $12, queries_per_minute = $13, id = $14, updated = now() 
		where contracts.provider_id = $1
		  and contracts.delegate_pubkey = $2
			and contracts.height = $9
		returning id, created, updated
	`
	sqlCloseContract = `
	update contracts
	set closed_height = $1
	where id = $2
	returning id, created, updated
	`

	/*
		sqlUpsertOpenContractEvent = `
		insert into open_contract_events(contract_id,client_pubkey,contract_type,height,txid,duration,rate,open_cost,deposit,settlement_duration,auth,queries_per_minute)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		on conflict on constraint open_contract_events_txid_unq
		do update set updated = now()
		where open_contract_events.txid = $5
		returning id, created, updated
		`
	*/

	sqlUpsertCloseContractEvent = `
	insert into close_contract_events(contract_id,client_pubkey,delegate_pubkey,height,txid)
	values ($1,$2,$3,$4,$5)
	on conflict on constraint close_contract_events_txid_key
	do update set updated = now()
	where close_contract_events.txid = $5
	returning id, created, updated
	`

	sqlUpsertContractSettlementEvent = `
		UPDATE contracts
		SET nonce = $1, paid = $2, reserve_contrib_asset = $3
		WHERE id = $4
	returning id, created, updated
`
)
