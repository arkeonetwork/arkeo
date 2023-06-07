package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"

	arkeotypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func TestInsertBlock(t *testing.T) {
	m, err := pgxmock.NewPool()
	assert.Nil(t, err)
	defer m.Close()

	db := DirectoryDB{
		hijacker: func() (IConnection, error) {
			return &MockDB{
				pool: m,
			}, nil
		},
	}
	result, err := db.InsertBlock(context.Background(), nil)
	assert.NotNil(t, err)
	assert.Nil(t, result)
	hash := arkeotypes.GetRandomTxID()
	blockTime := time.Now()
	returnTime := time.Now()

	m.ExpectQuery(`insert into blocks(height,hash,block_time)*`).
		WithArgs(int64(1), hash, AnyTime{}).WillReturnRows(
		pgxmock.NewRows([]string{"id", "created", "updated"}).
			AddRow(int64(1), returnTime, returnTime),
	)

	b, err := db.InsertBlock(context.Background(), &Block{Height: 1, Hash: hash, BlockTime: blockTime})
	assert.Nil(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, b.ID, int64(1))
	assert.Nil(t, m.ExpectationsWereMet())
}

func TestFindLatestBlock(t *testing.T) {
	m, err := pgxmock.NewPool()
	assert.Nil(t, err)
	defer m.Close()
	mockDb := &MockDB{
		pool: m,
	}
	db := DirectoryDB{
		hijacker: func() (IConnection, error) {
			return mockDb, nil
		},
	}

	testTime := time.Now()
	hash := arkeotypes.GetRandomTxID()
	m.ExpectQuery(`select b.id, b.created, b.updated, b.height, b.hash, b.block_time from blocks b where*`).WillReturnRows(
		pgxmock.NewRows([]string{
			"id", "created", "updated", "height", "hash", "block_time",
		}).AddRow(int64(1), testTime, testTime, int64(1024), hash, testTime),
	)
	b, err := db.FindLatestBlock(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, b)
	assert.Nil(t, m.ExpectationsWereMet())
	assert.Equal(t, b.ID, int64(1))
	assert.Equal(t, b.Height, int64(1024))
	assert.Equal(t, b.Hash, hash)
	assert.Equal(t, b.BlockTime, testTime)
	assert.Equal(t, b.Created, testTime)
	assert.Equal(t, b.Updated, testTime)

	// when query fail , it should return nil block and err
	m, err = pgxmock.NewPool()
	assert.Nil(t, err)
	defer m.Close()
	mockDb = &MockDB{
		pool: m,
	}

	m.ExpectQuery(`select.*from blocks b where*`).
		WillReturnError(fmt.Errorf("fail to query latest block"))
	b, err = db.FindLatestBlock(context.Background())
	assert.NotNil(t, err)
	assert.Nil(t, b)
	assert.Nil(t, m.ExpectationsWereMet())
}
