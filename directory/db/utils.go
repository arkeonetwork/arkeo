package db

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"
)

func insert(ctx context.Context, conn IConnection, sql string, params ...interface{}) (*Entity, error) {
	var (
		id      int64
		created time.Time
		updated time.Time
		err     error
	)
	log.Debugf("sql: %s\nparams: %v", sql, params)
	row := conn.QueryRow(ctx, sql, params...)
	if err = row.Scan(&id, &created, &updated); err != nil {
		return nil, errors.Wrap(err, "fail to insert")
	}

	return &Entity{ID: id, Created: created, Updated: updated}, nil
}

func update(ctx context.Context, conn IConnection, sql string, params ...interface{}) (*Entity, error) {
	var (
		id      int64
		created time.Time
		updated time.Time
		err     error
	)
	log.Debugf("sql: %s", sql)
	row := conn.QueryRow(ctx, sql, params...)
	if err = row.Scan(&id, &created, &updated); err != nil {
		return nil, errors.Wrap(err, "error inserting")
	}

	return &Entity{ID: id, Created: created, Updated: updated}, nil
}

// if the query returns no rows, the passed target remains unchanged. target must be a pointer
func selectOne(ctx context.Context, conn IConnection, query string, target interface{}, params ...interface{}) error {
	log.Debugf("sql: %s\nparams: %v", query, params)
	if err := pgxscan.Get(ctx, conn, target, query, params...); err != nil {
		return errors.Wrapf(err, "error selecting with params: %v", params)
	}
	return nil
}

// nolint
//func selectMany(ctx context.Context, conn IConnection, sql string, params ...interface{}) ([]map[string]interface{}, error) {
//	results := make([]map[string]interface{}, 0, 512)
//	if err := pgxscan.Select(ctx, conn, &results, sql, params...); err != nil {
//		return nil, errors.Wrapf(err, "error selecting many")
//	}
//	return results, nil
//}

func upsert(ctx context.Context, conn IConnection, sql string, params ...interface{}) (*Entity, error) {
	row := conn.QueryRow(ctx, sql, params...)

	var (
		id      int64
		created time.Time
		updated time.Time
		err     error
	)

	if err = row.Scan(&id, &created, &updated); err != nil {
		return nil, fmt.Errorf("error upserting: %+v", err)
	}

	entity := &Entity{
		ID:      id,
		Created: created,
		Updated: updated,
	}

	return entity, nil
}

func getFlavor() sqlbuilder.Flavor {
	return sqlbuilder.PostgreSQL
}

// Upsert Generic Event
func (d *DirectoryDB) InsertGenericEvent(ctx context.Context, eventType, txID string, height int64, attrJSON []byte) (*Entity, error) {
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()

	return insert(ctx, conn, sqlInsertGenericEvent, eventType, txID, height, attrJSON)
}
