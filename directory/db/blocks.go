package db

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
)

type Block struct {
	Entity
	Height    int64     `db:"height"`
	Hash      string    `db:"hash"`
	BlockTime time.Time `db:"block_time"`
}

func (d *DirectoryDB) InsertBlock(b *Block) (*Entity, error) {
	if b == nil {
		return nil, fmt.Errorf("nil block")
	}
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	return insert(conn, sqlInsertBlock, b.Height, b.Hash, b.BlockTime)
}

func (d *DirectoryDB) FindLatestBlock() (*Block, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	block := &Block{} // used to designate not found... need a better way!
	if err = selectOne(conn, sqlFindLatestBlock, block); err != nil {
		return nil, errors.Wrapf(err, "error selecting")
	}
	// not found
	if block.Height == 0 {
		return nil, nil
	}
	return block, nil
}

type BlockGap struct {
	Start int64 `db:"gap_start"`
	End   int64 `db:"gap_end"`
}

func (g BlockGap) String() string {
	return fmt.Sprintf("%d-%d", g.Start, g.End)
}

func (d *DirectoryDB) FindBlockGaps() ([]*BlockGap, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	results := make([]*BlockGap, 0, 128)
	if err = pgxscan.Select(context.Background(), conn, &results, sqlFindBlockGaps); err != nil {
		return nil, errors.Wrapf(err, "error scanning")
	}

	return results, nil
}
