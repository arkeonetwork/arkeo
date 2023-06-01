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
	sqlInsertBlock     = `insert into blocks(height,hash,block_time) values($1,$2,$3) returning id, created, updated`
	sqlFindLatestBlock = `
		select ` + blockCols + `
		from blocks b
		where b.height = (select max(height) from blocks)
	`
)
