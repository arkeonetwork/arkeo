package db

import (
	"context"

	"github.com/pkg/errors"

	"github.com/arkeonetwork/arkeo/directory/types"
)

func (d *DirectoryDB) GetArkeoNetworkStats(ctx context.Context) (*types.ArkeoStats, error) {
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()
	stats := types.ArkeoStats{}
	if err = selectOne(ctx, conn, sqlGetNetworkStats, &stats); err != nil {
		return nil, errors.Wrapf(err, "error getting stats")
	}

	return &stats, nil
}

func (d *DirectoryDB) GetArkeoNetworkStatsByService(ctx context.Context, service string) (*types.ArkeoStats, error) {
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	defer conn.Release()
	stats := types.ArkeoStats{}
	if err = selectOne(ctx, conn, sqlGetNetworkStatsByService, &stats, service); err != nil {
		return nil, errors.Wrapf(err, "error getting stats")
	}

	return &stats, nil
}
