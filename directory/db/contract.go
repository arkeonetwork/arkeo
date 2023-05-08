package db

import (
	"context"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/directory/types"
	atypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
)

// TODO: missing provider pubkey and service, need to expand ProviderID
// TODO: add paid
// TODO: add settlement duration
// TODO: add nonce
// TODO: reserve contributions, both in asset and USD terms
type ArkeoContract struct {
	Entity
	ContractID          int64              `json:"contract_id" db:"id"`
	Provider            string             `json:"provider" db:"-"`
	Client              string             `json:"service" db:"-"`
	DelegatePubkey      string             `json:"delegate_pubkey" db:"delegate_pubkey"`
	ClientPubkey        string             `json:"client_pubkey" db:"client_pubkey"`
	Height              int64              `json:"height" db:"height"`
	ContractType        types.ContractType `json:"contract_type" db:"contract_type"`
	Duration            int64              `json:"duration" db:"duration"`
	RateAsset           string             `json:"-" db:"rate_asset"`
	RateAmount          int64              `json:"-" db:"rate_amount"`
	Rate                cosmos.Coin        `json:"rate" db:"-"`
	OpenCost            int64              `json:"open_cost" db:"open_cost"`
	ClosedHeight        int64              `json:"closed_height" db:"closed_height"`
	ProviderID          int64              `json:"-" db:"provider_id"`
	Deposit             int64              `json:"deposit" db:"deposit"`
	Authorization       types.AuthType     `json:"authorization" db:"auth"`
	QueriesPerMinute    int64              `json:"queries_per_minute" db:"queries_per_minute"`
	Nonce               int64              `json:"nonce" db:"nonce"`
	Paid                int64              `json:"paid" db:"paid"`
	SettlementDurtion   int64              `json:"settlement_duration" db:"settlement_duration"`
	ReserveContribAsset int64              `json:"reserve_contrib_asset" db:"reserve_contrib_asset"`
	ReserveContribUSD   int64              `json:"reserve_contrib_usd" db:"reserve_contrib_usd"`
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

func (d *DirectoryDB) UpsertContract(providerID int64, evt atypes.EventOpenContract) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return upsert(conn, sqlUpsertContract, providerID, evt.Delegate, evt.Client, evt.Type,
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

func (d *DirectoryDB) UpsertOpenContractEvent(contractID int64, evt atypes.EventOpenContract) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	/*
		return upsert(conn, sqlUpsertOpenContractEvent, evt.ContractId, evt.Client, evt.Type, evt.Height, evt.TxID,
			evt.Duration, evt.Rate, evt.OpenCost, evt.Deposit, evt.SettlementDuration, evt.Authorization, evt.QueriesPerMinute)
	*/
	return nil, nil
}

func (d *DirectoryDB) UpsertCloseContractEvent(contractID int64, evt types.CloseContractEvent) (*Entity, error) {
	conn, err := d.getConnection()
	defer conn.Release()
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}

	return upsert(conn, sqlUpsertCloseContractEvent, contractID, evt.ClientPubkey, evt.GetDelegatePubkey(), evt.EventHeight, evt.TxID)
}
