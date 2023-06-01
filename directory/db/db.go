package db

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"github.com/arkeonetwork/arkeo/common/logging"
	"github.com/arkeonetwork/arkeo/sentinel"
	atypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

type DBConfig struct {
	Host                    string `mapstructure:"host" json:"host"`
	Port                    uint   `mapstructure:"port" json:"port"`
	User                    string `mapstructure:"user" json:"user"`
	Pass                    string `mapstructure:"pass" json:"pass"`
	DBName                  string `mapstructure:"name" json:"name"`
	PoolMaxConns            int    `mapstructure:"pool_max_conns" json:"pool_max_conns"`
	PoolMinConns            int    `mapstructure:"pool_min_conns" json:"pool_min_conns"`
	SSLMode                 string `mapstructure:"ssl_mode" json:"ssl_mode"`
	ConnectionTimeoutSecond int    `mapstructure:"connection_timeout" json:"connection_timeout"`
}

type IDataStorage interface {
	FindLatestBlock() (*Block, error)
	InsertBlock(b *Block) (*Entity, error)
	UpsertValidatorPayoutEvent(evt atypes.EventValidatorPayout, height int64) (*Entity, error)
	FindProvider(pubkey, service string) (*ArkeoProvider, error)
	UpsertContract(providerID int64, evt atypes.EventOpenContract) (*Entity, error)
	FindContract(contractId uint64) (*ArkeoContract, error)
	CloseContract(contractID uint64, height int64) (*Entity, error)
	UpdateProvider(provider *ArkeoProvider) (*Entity, error)
	UpsertContractSettlementEvent(evt atypes.EventSettleContract) (*Entity, error)
	UpsertProviderMetadata(providerID, nonce int64, data sentinel.Metadata) (*Entity, error)
	InsertBondProviderEvent(providerID int64, evt atypes.EventBondProvider, height int64, txID string) (*Entity, error)
	InsertProvider(provider *ArkeoProvider) (*Entity, error)
}

var _ IDataStorage = &DirectoryDB{}

type Acquireable interface {
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
}

var _ Acquireable = &pgxpool.Pool{}

type IConnection interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	pgxscan.Querier
	Release()
	Begin(ctx context.Context) (pgx.Tx, error)
}

var _ IConnection = &pgxpool.Conn{}

// ErrNotFound indicate the record doesn't exist in DB
var ErrNotFound = pgx.ErrNoRows

type (
	connectionHijacker func() (IConnection, error)
	DirectoryDB        struct {
		pool     Acquireable
		hijacker connectionHijacker // this is only used for test purpose
	}
)

// base entity for db types
type Entity struct {
	ID      int64     `db:"id"`
	Created time.Time `db:"created"`
	Updated time.Time `db:"updated"`
}

var log = logging.WithoutFields()

// obtain a db connection, callers must call conn.Release() when finished to return the conn to the pool
func (d *DirectoryDB) getConnection() (IConnection, error) {
	if d.hijacker != nil {
		return d.hijacker()
	}
	return d.pool.Acquire(context.Background())
}

func New(config DBConfig) (*DirectoryDB, error) {
	connStrTemplate := "postgres://%s:%s@%s:%d/%s?pool_max_conns=%d&pool_min_conns=%d&sslmode=%s"
	url := fmt.Sprintf(connStrTemplate, config.User, config.Pass, config.Host, config.Port, config.DBName, config.PoolMaxConns, config.PoolMinConns, config.SSLMode)
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing url to config from: \"%s\"", url)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "error connecting to db")
	}

	log.Infof("connected pool for db %s on %s:%d", config.DBName, config.Host, config.Port)
	return &DirectoryDB{
		pool: pool,
	}, nil
}
