package db

import (
	"context"
	"testing"
)

var config = DBConfig{
	Host:         "localhost",
	Port:         5432,
	User:         "arkeo",
	Pass:         "arkeo123",
	DBName:       "arkeo_directory",
	PoolMaxConns: 2,
	PoolMinConns: 1,
	SSLMode:      "prefer",
}

func TestNew(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := New(config)
	if err != nil {
		t.Errorf("error: %+v", err)
	}
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		t.Errorf("error acquiring connection: %+v", err)
	}
	defer conn.Release()
	log.Infof("got connection %s", conn)
}
