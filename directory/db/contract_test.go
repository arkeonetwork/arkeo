package db

import (
	"testing"
	"time"

	arkeotypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
	pgxmock "github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

func TestFindContract(t *testing.T) {
	m, err := pgxmock.NewPool()
	assert.Nil(t, err)
	defer m.Close()
	mockDb := &MockDB{
		pool: m,
	}
	db := DirectoryDB{
		hijacker: func() (IConnection, error) {
			return mockDb, nil
		},
	}
	testTime := time.Now()
	testPubKey := arkeotypes.GetRandomPubKey()
	m.ExpectQuery(`select .* from contracts c*`).
		WithArgs(uint64(1)).
		WillReturnRows(
			pgxmock.NewRows([]string{
				"id", "created", "updated", "provider_id", "delegate_pubkey", "client_pubkey", "height", "contract_type", "duration", "rate_asset",
				"rate_amount", "open_cost", "deposit", "auth", "queries_per_minute", "settlement_duration", "paid", "reserve_contrib_asset",
				"reserve_contrib_usd", "closed_height",
			}).AddRow(int64(1), testTime, testTime, int64(1), testPubKey.String(), testPubKey.String(), int64(1024), "PayAsYouGo",
				int64(10), "uarkeo", int64(10), int64(10), int64(100000), "STRICT", int64(10), int64(10), int64(1000), int64(100), int64(100), int64(2048)),
		)
	m.ExpectQuery(`SELECT .* FROM providers WHERE id*`).
		WithArgs(int64(1)).
		WillReturnRows(
			pgxmock.NewRows([]string{"pubkey", "service"}).
				AddRow(testPubKey.String(), "mock"))
	contract, err := db.FindContract(1)
	assert.Nil(t, err)
	assert.NotNil(t, contract)
	assert.Equal(t, contract.ContractID, int64(1))
	assert.Equal(t, contract.Provider, testPubKey.String())
	assert.Equal(t, contract.Service, "mock")
	assert.Equal(t, contract.DelegatePubkey, testPubKey.String())
	assert.Equal(t, contract.ClientPubkey, testPubKey.String())
	assert.Equal(t, contract.Height, int64(1024))
	assert.Equal(t, contract.ContractType, "PayAsYouGo")
	assert.Equal(t, contract.Duration, int64(10))
	assert.Equal(t, contract.RateAsset, "uarkeo")
	assert.Equal(t, contract.RateAmount, int64(10))
	assert.Equal(t, contract.OpenCost, int64(10))
	assert.Equal(t, contract.Deposit, int64(100000))
	assert.Equal(t, contract.Authorization, "STRICT")
	assert.Equal(t, contract.ClosedHeight, int64(2048))
	assert.Equal(t, contract.ProviderID, int64(1))
	assert.Equal(t, contract.QueriesPerMinute, int64(10))
	assert.Equal(t, contract.Paid, int64(1000))
	assert.Equal(t, contract.SettlementDurtion, int64(10))
	assert.Equal(t, contract.ReserveContribAsset, int64(100))
	assert.Equal(t, contract.ReserveContribUSD, int64(100))

}
