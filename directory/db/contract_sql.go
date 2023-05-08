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
	c.rate,
	c.open_cost,
	c.deposit,
	c.auth,
	c.queries_per_minute,
	c.closed_height
	`
)

const (
	sqlFindContract = ` select ` + contractCols + `
	-- id,
	-- created,
	-- updated,
	-- provider_id,
	-- delegate_pubkey,
	-- client_pubkey,
	-- height,
	-- contract_type,
	-- duration,
	-- rate,
	-- open_cost,
	-- deposit,
	-- auth,
	-- queries_per_minute,
	-- closed_height
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
	-- c.rate,
	-- c.open_cost
	-- c.deposit
	-- c.auth
	-- c.queries_per_minute
	from providers p join contracts c on p.id = c.provider_id
	where p.service = $1 and p.pubkey = $2 and c.delegate_pubkey = $3
	order by c.id desc
	`
	sqlFindContractByPubKeys = `select ` + contractCols + `
	from providers p join contracts c on p.id = c.provider_id
	where p.service = $1 and p.pubkey = $2 and c.delegate_pubkey = $3 and c.height = $4
	`
	sqlUpsertContract = `
		insert into contracts(provider_id,delegate_pubkey,client_pubkey,contract_type,duration,rate,open_cost,height,deposit,settlement_duration,auth,queries_per_minute,id)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		on conflict on constraint contracts_provider_delegate_height_key
		do update set contract_type = $4, duration = $5, rate = $6, open_cost = $7, updated = now() 
		where contracts.provider_id = $1
		  and contracts.delegate_pubkey = $2
			and contracts.height = $8
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
	insert into contract_settlement_events(contract_id,txid,client_pubkey,height,nonce,paid,reserve)
	values ($1,$2,$3,$4,$5,$6,$7)
	on conflict on constraint contract_settlement_contract_nonce_key
	do update set updated = now()
	where contract_settlement_events.contract_id = $1
	  and contract_settlement_events.nonce = $5
	returning id, created, updated
`
)
