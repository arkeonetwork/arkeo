package db

import (
	"github.com/pkg/errors"

	"github.com/arkeonetwork/arkeo/directory/types"
)

func (d *DirectoryDB) GetArkeoNetworkStats() (*types.ArkeoStats, error) {
	conn, err := d.getConnection()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()
	stats := types.ArkeoStats{}
	if err = selectOne(conn, sqlGetNetworkStats, &stats); err != nil {
		return nil, errors.Wrapf(err, "error getting stats")
	}

	return &stats, nil
}
