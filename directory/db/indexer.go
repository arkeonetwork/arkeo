package db

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
)

type IndexerStatus struct {
	ID     int64  `db:"id"`
	Height uint64 `db:"height"`
}

func (d *DirectoryDB) UpsertIndexerStatus(indexerStatus *IndexerStatus) (*Entity, error) {
	if indexerStatus == nil {
		return nil, fmt.Errorf("nil IndexerStatus")
	}
	conn, err := d.getConnection()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()

	return insert(conn, sqlUpsertIndexerStatus, indexerStatus.ID, indexerStatus.Height)
}

func (d *DirectoryDB) UpdateIndexerStatus(indexerStatus *IndexerStatus) (*Entity, error) {
	if indexerStatus == nil {
		return nil, fmt.Errorf("nil IndexerStatus")
	}
	conn, err := d.getConnection()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()

	return update(conn,
		sqlUpdateIndexerStatus,
		indexerStatus.ID,
		indexerStatus.Height,
	)
}

func (d *DirectoryDB) FindIndexerStatus(id int64) (*IndexerStatus, error) {
	conn, err := d.getConnection()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()
	indexerStatus := IndexerStatus{Height: math.MaxUint64} // used to designate not found... need a better way!
	if err = selectOne(conn, sqlFindIndexerStatus, &indexerStatus, id); err != nil {
		return nil, errors.Wrapf(err, "error selecting")
	}
	// not found
	if indexerStatus.Height == math.MaxUint64 {
		return nil, nil
	}
	return &indexerStatus, nil
}
