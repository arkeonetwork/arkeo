package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

var config = DBConfig{
	Host:                    "localhost",
	Port:                    5432,
	User:                    "arkeo",
	Pass:                    "arkeo123",
	DBName:                  "arkeo_directory",
	PoolMaxConns:            2,
	PoolMinConns:            1,
	SSLMode:                 "prefer",
	ConnectionTimeoutSecond: 1,
}

func TestNew(t *testing.T) {
	// when config fail to parse should result in an error
	db, err := New(DBConfig{
		PoolMaxConns: -100,
		SSLMode:      "",
	})
	assert.NotNil(t, err)
	assert.Nil(t, db)
	db, err = New(config)
	assert.Nil(t, err)
	assert.NotNil(t, db)
}

type MockDB struct {
	pool pgxmock.PgxPoolIface
}

func (m *MockDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return m.pool.QueryRow(ctx, sql, args...)
}

func (m *MockDB) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return m.pool.Query(ctx, query, args...)
}
func (m *MockDB) Release() {}
func (m *MockDB) Begin(ctx context.Context) (pgx.Tx, error) {
	return m.pool.Begin(ctx)
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v interface{}) bool {
	_, ok := v.(time.Time)
	return ok
}

func getMockDirectoryDBForTest(t *testing.T) (pgxmock.PgxPoolIface, *DirectoryDB) {
	m, err := pgxmock.NewPool()
	assert.Nil(t, err)

	mockDb := &MockDB{
		pool: m,
	}
	db := &DirectoryDB{
		hijacker: func() (IConnection, error) {
			return mockDb, nil
		},
	}
	return m, db
}
