package db

import (
	"context"

	"github.com/arkeonetwork/arkeo/directory/types"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
)

// TODO: missing provider pubkey and service, need to expand ProviderID
type ArkeoContract struct {
	Entity
	ContractID       int64              `json:"contract_id" db:"contract_id"`
	DelegatePubkey   string             `json:"delegate_pubkey" db:"delegate_pubkey"`
	ClientPubkey     string             `json:"client_pubkey" db:"client_pubkey"`
	Height           int64              `json:"height" db:"height"`
	ContractType     types.ContractType `json:"contract_type" db:"contract_type"`
	Duration         int64              `json:"duration" db:"duration"`
	Rate             int64              `json:"rate" db:"rate"`
	OpenCost         int64              `json:"open_cost" db:"open_cost"`
	ClosedHeight     int64              `json:"closed_height" db:"closed_height"`
	ProviderID       int64              `json:"-" db:"provider_id"`
	Deposit          int64              `json:"deposit" db:"deposit"`
	Authorization    types.AuthType     `json:"authorization" db:"auth"`
	QueriesPerMinute int64              `json:"queries_per_minute" db:"queries_per_minute"`
}

func (d *DirectoryDB) FindContract(contractId uint64) (*ArkeoContract, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	contract := ArkeoContract{}
	if err = selectOne(conn, sqlFindContract, &contract, contractId); err != nil {
		return nil, errors.Wrapf(err, "error selecting")
	}

	// not found
	if contract.ClientPubkey == "" {
		return nil, nil
	}
	return &contract, nil
}

func (d *DirectoryDB) FindContractsByPubKeys(service, providerPubkey, delegatePubkey string) ([]*ArkeoContract, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	results := make([]*ArkeoContract, 0, 128)
	if err = pgxscan.Select(context.Background(), conn, &results, sqlFindContractsByPubKeys, service, providerPubkey, delegatePubkey); err != nil {
		return nil, errors.Wrapf(err, "error scanning")
	}

	return results, nil
}

func (d *DirectoryDB) FindContractByPubKeys(service, providerPubkey, delegatePubkey string, height int64) (*ArkeoContract, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	contract := ArkeoContract{}
	if err = selectOne(conn, sqlFindContractByPubKeys, &contract, service, providerPubkey, delegatePubkey, height); err != nil {
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
		evt.Duration, evt.Rate, evt.OpenCost, evt.Height, evt.Deposit, evt.SettlementDuration, evt.Authorization, evt.QueriesPerMinute, evt.ContractId)
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
		evt.Duration, evt.Rate, evt.OpenCost, evt.Deposit, evt.SettlementDuration, evt.Authorization, evt.QueriesPerMinute)
}

func (d *DirectoryDB) UpsertCloseContractEvent(contractID int64, evt types.CloseContractEvent) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return upsert(conn, sqlUpsertCloseContractEvent, contractID, evt.ClientPubkey, evt.GetDelegatePubkey(), evt.EventHeight, evt.TxID)
}
