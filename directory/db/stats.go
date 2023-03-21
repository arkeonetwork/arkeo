package db

import (
	"github.com/arkeonetwork/arkeo/directory/types"
	"github.com/pkg/errors"
)

func (d *DirectoryDB) GetArkeoNetworkStats() (*types.ArkeoStats, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	stats := types.ArkeoStats{}
	if err = selectOne(conn, sqlGetNetworkStats, &stats); err != nil {
		return nil, errors.Wrapf(err, "error getting stats")
	}

	return &stats, nil
}
