package db

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	atypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

// TODO: missing provider pubkey and service, need to expand ProviderID
// TODO: add paid
// TODO: add settlement duration
// TODO: add nonce
// TODO: reserve contributions, both in asset and USD terms
type ArkeoContract struct {
	Entity
	ContractID          int64       `json:"contract_id" db:"id"`
	Provider            string      `json:"provider" db:"provider"`
	Service             string      `json:"service" db:"service"`
	DelegatePubkey      string      `json:"delegate_pubkey" db:"delegate_pubkey"`
	ClientPubkey        string      `json:"client_pubkey" db:"client_pubkey"`
	Height              int64       `json:"height" db:"height"`
	ContractType        string      `json:"contract_type" db:"contract_type"`
	Duration            int64       `json:"duration" db:"duration"`
	RateAsset           string      `json:"-" db:"rate_asset"`
	RateAmount          int64       `json:"-" db:"rate_amount"`
	Rate                cosmos.Coin `json:"rate" db:"-"`
	OpenCost            int64       `json:"open_cost" db:"open_cost"`
	ClosedHeight        int64       `json:"closed_height" db:"closed_height"`
	ProviderID          int64       `json:"-" db:"provider_id"`
	Deposit             int64       `json:"deposit" db:"deposit"`
	Authorization       string      `json:"authorization" db:"auth"`
	QueriesPerMinute    int64       `json:"queries_per_minute" db:"queries_per_minute"`
	Nonce               int64       `json:"nonce" db:"nonce"`
	Paid                int64       `json:"paid" db:"paid"`
	SettlementDurtion   int64       `json:"settlement_duration" db:"settlement_duration"`
	ReserveContribAsset int64       `json:"reserve_contrib_asset" db:"reserve_contrib_asset"`
	ReserveContribUSD   int64       `json:"reserve_contrib_usd" db:"reserve_contrib_usd"`
}

// GetContract query db to find contract that match the given contract id
func (d *DirectoryDB) GetContract(ctx context.Context, contractId uint64) (*ArkeoContract, error) {
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to obtain db connection,err: %w", err)
	}
	defer conn.Release()

	contract := ArkeoContract{}
	if err = selectOne(ctx, conn, sqlGetContractByID, &contract, contractId); err != nil {
		return nil, errors.Wrapf(err, "error selecting")
	}

	if len(contract.RateAsset) > 0 {
		contract.Rate = cosmos.NewInt64Coin(contract.RateAsset, contract.RateAmount)
	}

	return &contract, nil
}

// UpsertContract update database with the given open contract event, if the contract doesn't exist , it will create a new one
func (d *DirectoryDB) UpsertContract(ctx context.Context, providerID int64, evt atypes.EventOpenContract) (*Entity, error) {
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()

	if evt.Delegate.String() == "" {
		evt.Delegate = evt.Client
	}

	return upsert(ctx,
		conn,
		sqlUpsertContract,
		providerID,
		evt.Delegate,
		evt.Client,
		evt.Type,
		evt.Duration,
		evt.Rate.Denom,
		evt.Rate.Amount.Int64(),
		evt.OpenCost,
		evt.Height,
		evt.Deposit.Int64(),
		evt.SettlementDuration,
		evt.Authorization,
		evt.QueriesPerMinute,
		evt.ContractId,
	)
}

func (d *DirectoryDB) CloseContract(ctx context.Context, contractID uint64, height int64) (*Entity, error) {
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()

	return update(ctx, conn, sqlCloseContract, height, contractID)
}

func (d *DirectoryDB) UpsertContractSettlementEvent(ctx context.Context, evt atypes.EventSettleContract) (*Entity, error) {
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()

	return upsert(ctx, conn, sqlUpsertContractSettlementEvent, evt.Nonce, evt.Paid.Int64(), evt.Reserve.Int64(), evt.ContractId)
}

func (d *DirectoryDB) UpsertOpenContractEvent(ctx context.Context, contractID int64, evt atypes.EventOpenContract) (*Entity, error) {
	conn, err := d.getConnection(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error obtaining db connection")
	}
	defer conn.Release()

	/*
		return upsert(conn, sqlUpsertOpenContractEvent, evt.ContractId, evt.Client, evt.Type, evt.Height, evt.TxID,
			evt.Duration, evt.Rate, evt.OpenCost, evt.Deposit, evt.SettlementDuration, evt.Authorization, evt.QueriesPerMinute)
	*/
	return nil, nil
}
