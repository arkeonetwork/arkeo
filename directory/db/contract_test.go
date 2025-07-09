package db

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/math"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"

	"github.com/arkeonetwork/arkeo/common/cosmos"

	arkeotypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func TestFindContract(t *testing.T) {
	m, db := getMockDirectoryDBForTest(t)
	defer m.Close()
	testTime := time.Now()
	testPubKey := arkeotypes.GetRandomPubKey()
	m.ExpectQuery(`select .* from contracts c*`).
		WithArgs(uint64(1)).
		WillReturnRows(
			pgxmock.NewRows([]string{
				"id", "created", "updated", "provider", "service", "delegate_pubkey", "client_pubkey", "height", "contract_type", "duration", "rate_asset",
				"rate_amount", "open_cost", "deposit", "auth", "queries_per_minute", "settlement_duration", "paid", "reserve_contrib_asset",
				"reserve_contrib_usd", "settlement_height", "provider_id",
			}).AddRow(int64(1), testTime, testTime, testPubKey.String(), "mock", testPubKey.String(), testPubKey.String(), int64(1024), "PayAsYouGo",
				int64(10), "uarkeo", int64(10), int64(10), int64(100000), "STRICT", int64(10), int64(10), int64(1000), int64(100), int64(100), int64(2048), int64(1)),
		)
	ctx := context.Background()
	contract, err := db.GetContract(ctx, 1)
	assert.Nil(t, err)
	assert.NotNil(t, contract)
	assert.Equal(t, int64(1), contract.ContractID)
	assert.Equal(t, testPubKey.String(), contract.Provider)
	assert.Equal(t, "mock", contract.Service)
	assert.Equal(t, testPubKey.String(), contract.DelegatePubkey)
	assert.Equal(t, testPubKey.String(), contract.ClientPubkey)
	assert.Equal(t, int64(1024), contract.Height)
	assert.Equal(t, "PayAsYouGo", contract.ContractType)
	assert.Equal(t, int64(10), contract.Duration)
	assert.Equal(t, "uarkeo", contract.RateAsset)
	assert.Equal(t, int64(10), contract.RateAmount)
	assert.Equal(t, int64(10), contract.OpenCost)
	assert.Equal(t, int64(100000), contract.Deposit)
	assert.Equal(t, "STRICT", contract.Authorization)
	assert.Equal(t, int64(2048), contract.SettlementHeight)
	assert.Equal(t, int64(1), contract.ProviderID)
	assert.Equal(t, int64(10), contract.QueriesPerMinute)
	assert.Equal(t, int64(1000), contract.Paid)
	assert.Equal(t, int64(10), contract.SettlementDurtion)
	assert.Equal(t, int64(100), contract.ReserveContribAsset)
	assert.Equal(t, int64(100), contract.ReserveContribUSD)
	assert.Nil(t, m.ExpectationsWereMet())
}

func TestUpdateContract(t *testing.T) {
	m, db := getMockDirectoryDBForTest(t)
	defer m.Close()
	testTime := time.Now()
	testPubKey := arkeotypes.GetRandomPubKey()
	evt := arkeotypes.EventOpenContract{
		Provider:           testPubKey,
		ContractId:         1,
		Service:            "mock",
		Client:             testPubKey,
		Delegate:           testPubKey,
		Type:               arkeotypes.ContractType_PAY_AS_YOU_GO,
		Height:             1024,
		Duration:           10,
		Rate:               cosmostypes.NewCoin("uarkeo", cosmos.NewInt(10)),
		OpenCost:           1000,
		Deposit:            math.NewInt(10000),
		SettlementDuration: 10,
		Authorization:      arkeotypes.ContractAuthorization_STRICT,
		QueriesPerMinute:   10,
	}
	m.ExpectQuery("insert into contracts.*").
		WithArgs(int64(1), evt.Delegate, evt.Client,
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
			evt.ContractId).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "created", "updated"}).
				AddRow(int64(1), testTime, testTime),
		)
	entity, err := db.UpsertContract(context.Background(), 1, evt)
	assert.Nil(t, err)
	assert.NotNil(t, entity)
	assert.Equal(t, int64(1), entity.ID)
	assert.Equal(t, testTime, entity.Created)
	assert.Equal(t, testTime, entity.Updated)
	assert.Nil(t, m.ExpectationsWereMet())
}

func TestCloseContract(t *testing.T) {
	m, db := getMockDirectoryDBForTest(t)
	defer m.Close()
	testTime := time.Now()
	m.ExpectQuery("update contracts.*").
		WithArgs(int64(1024), uint64(1)).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "created", "updated"}).
				AddRow(int64(1), testTime, testTime),
		)
	entity, err := db.CloseContract(context.Background(), 1, 1024)
	assert.Nil(t, err)
	assert.NotNil(t, entity)
	assert.Equal(t, int64(1), entity.ID)
	assert.Equal(t, testTime, entity.Created)
	assert.Equal(t, testTime, entity.Updated)
	assert.Nil(t, m.ExpectationsWereMet())
}

func TestUpsertContractSettltementEvent(t *testing.T) {
	m, db := getMockDirectoryDBForTest(t)
	defer m.Close()
	testTime := time.Now()

	testPubKey := arkeotypes.GetRandomPubKey()
	evt := arkeotypes.EventSettleContract{
		Provider:   testPubKey,
		ContractId: 1,
		Service:    "mock",
		Client:     testPubKey,
		Delegate:   testPubKey,
		Type:       arkeotypes.ContractType_PAY_AS_YOU_GO,
		Nonce:      1,
		Height:     1024,
		Paid:       math.NewInt(1000),
		Reserve:    math.NewInt(1000),
	}
	m.ExpectQuery("UPDATE contracts.*").
		WithArgs(evt.Nonce, evt.Paid.Int64(), evt.Reserve.Int64(), evt.ContractId).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "created", "updated"}).
				AddRow(int64(1), testTime, testTime),
		)
	entity, err := db.UpsertContractSettlementEvent(context.Background(), evt)
	assert.Nil(t, err)
	assert.NotNil(t, entity)
	assert.Equal(t, int64(1), entity.ID)
	assert.Equal(t, testTime, entity.Created)
	assert.Equal(t, testTime, entity.Updated)
	assert.Nil(t, m.ExpectationsWereMet())
}
