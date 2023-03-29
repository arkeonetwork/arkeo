package db

import (
	"context"

	"github.com/arkeonetwork/arkeo/directory/types"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
)

type ArkeoContract struct {
	Entity
	ProviderID     int64              `db:"provider_id"`
	DelegatePubkey string             `db:"delegate_pubkey"`
	ClientPubkey   string             `db:"client_pubkey"`
	Height         int64              `db:"height"`
	ContractType   types.ContractType `db:"contract_type"`
	Duration       int64              `db:"duration"`
	Rate           int64              `db:"rate"`
	OpenCost       int64              `db:"open_cost"`
	ClosedHeight   int64              `db:"closed_height"`
}

func (d *DirectoryDB) FindContract(providerID int64, delegatePubkey string, height int64) (*ArkeoContract, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	contract := ArkeoContract{}
	if err = selectOne(conn, sqlFindContract, &contract, providerID, delegatePubkey, height); err != nil {
		return nil, errors.Wrapf(err, "error selecting")
	}

	// not found
	if contract.ID == 0 {
		return nil, nil
	}
	return &contract, nil
}

func (d *DirectoryDB) FindContractsByPubKeys(chain, providerPubkey, delegatePubkey string) ([]*ArkeoContract, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	results := make([]*ArkeoContract, 0, 128)
	if err = pgxscan.Select(context.Background(), conn, &results, sqlFindContractsByPubKeys, chain, providerPubkey, delegatePubkey); err != nil {
		return nil, errors.Wrapf(err, "error scanning")
	}

	return results, nil
}

func (d *DirectoryDB) FindContractByPubKeys(chain, providerPubkey, delegatePubkey string, height int64) (*ArkeoContract, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	contract := ArkeoContract{}
	if err = selectOne(conn, sqlFindContractByPubKeys, &contract, chain, providerPubkey, delegatePubkey, height); err != nil {
		return nil, errors.Wrapf(err, "error selecting")
	}

	// not found
	if contract.ID == 0 {
		return nil, nil
	}
	return &contract, nil
}

func (d *DirectoryDB) UpsertContract(providerID int64, evt types.OpenContractEvent) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return upsert(conn, sqlUpsertContract, providerID, evt.GetDelegatePubkey(), evt.ClientPubkey, evt.ContractType,
		evt.Duration, evt.Rate, evt.OpenCost, evt.Height)
}

func (d *DirectoryDB) CloseContract(contractID, height int64) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return update(conn, sqlCloseContract, height, contractID)
}

func (d *DirectoryDB) UpsertContractSettlementEvent(contractID int64, evt types.ContractSettlementEvent) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return upsert(conn, sqlUpsertContractSettlementEvent, contractID, evt.TxID, evt.ClientPubkey, evt.EventHeight,
		evt.Nonce, evt.Paid, evt.Reserve)
}

func (d *DirectoryDB) UpsertOpenContractEvent(contractID int64, evt types.OpenContractEvent) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return upsert(conn, sqlUpsertOpenContractEvent, contractID, evt.ClientPubkey, evt.ContractType, evt.EventHeight, evt.TxID,
		evt.Duration, evt.Rate, evt.OpenCost)
}

func (d *DirectoryDB) UpsertCloseContractEvent(contractID int64, evt types.CloseContractEvent) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return upsert(conn, sqlUpsertCloseContractEvent, contractID, evt.ClientPubkey, evt.GetDelegatePubkey(), evt.EventHeight, evt.TxID)
}
