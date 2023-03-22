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
	sqlInsertBlock = `
		insert into blocks(height,hash,block_time)
		values($1,$2,$3)
		returning id, created, updated
	`
	sqlFindLatestBlock = `
		select ` + blockCols + `
		from blocks b
		where b.height = (select max(height) from blocks)
	`
	sqlFindBlockGaps = `
		select previousHeight + 1 as gap_start, height - 1 as gap_end
		from (select lag(b.height, 1) over (partition by 1 order by b.height) as previousHeight,
								b.height
					from blocks b
					order by b.height) as x
		where x.height - x.previousHeight > 1
		order by gap_start
	`
)
