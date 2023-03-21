package db

import (
	"context"
	"fmt"
	"time"

	"github.com/arkeonetwork/arkeo/common/logging"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type DBConfig struct {
	Host         string
	Port         uint
	User         string
	Pass         string
	DBName       string
	PoolMaxConns int
	PoolMinConns int
	SSLMode      string
}

type DirectoryDB struct {
	pool *pgxpool.Pool
}

// base entity for db types
type Entity struct {
	ID      int64     `db:"id"`
	Created time.Time `db:"created"`
	Updated time.Time `db:"updated"`
}

var log = logging.WithoutFields()

// obtain a db connection, callers must call conn.Release() when finished to return the conn to the pool
func (d *DirectoryDB) getConnection() (*pgxpool.Conn, error) {
	return d.pool.Acquire(context.Background())
}

func New(config DBConfig) (*DirectoryDB, error) {
	connStrTemplate := "postgres://%s:%s@%s:%d/%s?pool_max_conns=%d&pool_min_conns=%d&sslmode=%s"
	url := fmt.Sprintf(connStrTemplate, config.User, config.Pass, config.Host, config.Port, config.DBName, config.PoolMaxConns, config.PoolMinConns, config.SSLMode)
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing url to config from: \"%s\"", url)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "error connecting to db")
	}

	log.Infof("connected pool for db %s on %s:%d", config.DBName, config.Host, config.Port)
	return &DirectoryDB{pool}, nil
}
