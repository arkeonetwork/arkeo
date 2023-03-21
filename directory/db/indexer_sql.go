package db

var (
	sqlUpsertIndexerStatus = `
		insert into indexer_status(id,height) values ($1,$2)
		on conflict on constraint indexer_status_pk do update 
		set height = $2, updated = now()
		where indexer_status.id = $1
		returning id, created, updated
	`
	sqlUpdateIndexerStatus = `update indexer_status set height = $2, updated = now() where id = $1 returning id, created, updated`
	sqlFindIndexerStatus   = `select id,height from indexer_status where id = $1`
)
