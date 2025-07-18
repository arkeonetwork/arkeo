package db

const (
	blockCols = `
		b.id,
		b.created,
		b.updated,
		b.height,
		b.hash,
		b.block_time
	`

	sqlUpsertBlock = `
		INSERT INTO blocks(id, height, hash, block_time)
		VALUES (1, $1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		  SET height     = EXCLUDED.height,
			  hash       = EXCLUDED.hash,
			  block_time = EXCLUDED.block_time
		RETURNING id, created, updated
		`

	sqlFindLatestBlock = `
		select ` + blockCols + `
		from blocks b
		where b.height = (select max(height) from blocks)
	`
)
