package db

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type Block struct {
	Entity
	Height    int64     `db:"height"`
	Hash      string    `db:"hash"`
	BlockTime time.Time `db:"block_time"`
}

func (d *DirectoryDB) InsertBlock(ctx context.Context, b *Block) (*Entity, error) {
	if b == nil {
		return nil, fmt.Errorf("nil block")
	}
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()
	return insert(ctx, conn, sqlUpsertBlock, b.Height, b.Hash, b.BlockTime)
}

func (d *DirectoryDB) FindLatestBlock(ctx context.Context) (*Block, error) {
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()

	block := &Block{}
	if err = selectOne(ctx, conn, sqlFindLatestBlock, block); err != nil {
		return nil, errors.Wrapf(err, "error selecting")
	}
	return block, nil
}

type BlockGap struct {
	Start int64 `db:"gap_start"`
	End   int64 `db:"gap_end"`
}

func (g BlockGap) String() string {
	return fmt.Sprintf("[%d-%d]", g.Start, g.End)
}
