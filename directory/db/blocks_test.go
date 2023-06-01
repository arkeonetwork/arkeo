package db

import (
	"testing"
	"time"

	arkeotypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
	pgxmock "github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
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
	result, err := db.InsertBlock(nil)
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

	b, err := db.InsertBlock(&Block{Height: 1, Hash: hash, BlockTime: blockTime})
	assert.Nil(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, b.ID, int64(1))
	assert.Nil(t, m.ExpectationsWereMet())
}

// func TestFindLatestBlock(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	db, err := New(config)
// 	if err != nil {
// 		t.Errorf("error getting db: %+v", err)
// 	}
// 	b, err := db.FindLatestBlock()
// 	if err != nil {
// 		t.Fatalf("error: %+v", err)
// 	}
// 	log.Infof("found block b %v", b)
// }

// func TestFindBlockGaps(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	db, err := New(config)
// 	if err != nil {
// 		t.Errorf("error getting db: %+v", err)
// 	}
// 	b, err := db.FindBlockGaps()
// 	if err != nil {
// 		log.Fatalf("error finding gaps: %+v", err)
// 	}
// 	log.Infof("gaps: %v", b)
// }
