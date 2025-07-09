package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"

	"github.com/arkeonetwork/arkeo/directory/types"
	"github.com/arkeonetwork/arkeo/sentinel"
	"github.com/arkeonetwork/arkeo/sentinel/conf"
	arkeotypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func TestInsertProvider(t *testing.T) {
	m, db := getMockDirectoryDBForTest(t)
	defer m.Close()
	testTime := time.Now()

	entity, err := db.InsertProvider(context.Background(), nil)
	assert.NotNil(t, err)
	assert.Nil(t, entity)

	testPubKey := arkeotypes.GetRandomPubKey()
	p := &ArkeoProvider{
		Entity: Entity{
			ID:      1,
			Created: testTime,
			Updated: testTime,
		},
		Pubkey:              testPubKey.String(),
		Service:             "mock",
		Bond:                "1000",
		MetadataURI:         "http://localhost",
		MetadataNonce:       1,
		Status:              "ONLINE",
		MinContractDuration: 10,
		MaxContractDuration: 1000,
		SettlementDuration:  10,
		SubscriptionRate: []cosmostypes.Coin{
			cosmostypes.NewCoin("uarkeo", math.NewInt(10)),
		},
		PayAsYouGoRate: []cosmostypes.Coin{
			cosmostypes.NewCoin("uarkeo", math.NewInt(10)),
		},
	}
	p.Bond = "whatever"
	entity, err = db.InsertProvider(context.Background(), p)
	assert.NotNil(t, err)
	assert.Nil(t, entity)

	p.Bond = "1000"
	m.ExpectQuery("insert into providers.*").WithArgs(p.Pubkey, p.Service, int64(1000)).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "created", "updated"}).
				AddRow(int64(1), testTime, testTime),
		)
	entity, err = db.InsertProvider(context.Background(), p)
	assert.Nil(t, err)
	assert.NotNil(t, entity)
	assert.Equal(t, int64(1), entity.ID)
	assert.Equal(t, testTime, entity.Created)
	assert.Equal(t, testTime, entity.Updated)
	assert.Nil(t, m.ExpectationsWereMet())
}

func TestUpdateProvider(t *testing.T) {
	m, db := getMockDirectoryDBForTest(t)
	defer m.Close()
	testTime := time.Now()
	entity, err := db.UpdateProvider(context.Background(), nil)
	assert.NotNil(t, err)
	assert.Nil(t, entity)

	testPubKey := arkeotypes.GetRandomPubKey()
	p := &ArkeoProvider{
		Entity: Entity{
			ID:      1,
			Created: testTime,
			Updated: testTime,
		},
		Pubkey:              testPubKey.String(),
		Service:             "mock",
		Bond:                "1000",
		MetadataURI:         "http://localhost",
		MetadataNonce:       1,
		Status:              "ONLINE",
		MinContractDuration: 10,
		MaxContractDuration: 1000,
		SettlementDuration:  10,
		SubscriptionRate: []cosmostypes.Coin{
			cosmostypes.NewCoin("uarkeo", math.NewInt(10)),
		},
		PayAsYouGoRate: []cosmostypes.Coin{
			cosmostypes.NewCoin("uarkeo", math.NewInt(10)),
		},
	}
	m.ExpectBegin().WillReturnError(fmt.Errorf("fail to begin tx"))
	entity, err = db.UpdateProvider(context.Background(), p)
	assert.NotNil(t, err)
	assert.Nil(t, entity)
	assert.Nil(t, m.ExpectationsWereMet())
	// happy path
	m1, db1 := getMockDirectoryDBForTest(t)
	defer m1.Close()
	m1.ExpectBegin()
	m1.ExpectQuery("update providers.*").
		WithArgs(p.Pubkey, p.Service, p.Bond, p.MetadataURI, p.MetadataNonce, p.Status, p.MinContractDuration, p.MaxContractDuration, p.SettlementDuration).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "created", "updated"}).
				AddRow(int64(1), testTime, testTime),
		)
	m1.ExpectExec("DELETE FROM provider_subscription_rates.*").
		WithArgs(p.Pubkey, p.Service).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	m1.ExpectExec("DELETE FROM provider_pay_as_you_go_rates.*").
		WithArgs(p.Pubkey, p.Service).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	m1.ExpectExec("INSERT INTO provider_subscription_rates.*").
		WithArgs(int64(1), p.SubscriptionRate[0].Denom, p.SubscriptionRate[0].Amount.Int64()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	m1.ExpectExec("INSERT INTO provider_pay_as_you_go_rates.*").
		WithArgs(int64(1), p.PayAsYouGoRate[0].Denom, p.PayAsYouGoRate[0].Amount.Int64()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	m1.ExpectCommit()
	entity, err = db1.UpdateProvider(context.Background(), p)
	assert.Nil(t, err)
	assert.NotNil(t, entity)
	assert.Nil(t, m1.ExpectationsWereMet())
}

func TestFindProvider(t *testing.T) {
	m, db := getMockDirectoryDBForTest(t)
	defer m.Close()
	testTime := time.Now()
	testPubKey := arkeotypes.GetRandomPubKey()
	m.ExpectQuery("select.*from providers.*").
		WithArgs(testPubKey.String(), "mock").
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "created", "updated", "pubkey", "service", "bond", "metadata_uri", "metadata_nonce", "status", "min_contract_duration", "max_contract_duration", "settlement_duration",
		}).
			AddRow(int64(1), testTime, testTime, testPubKey.String(), "mock", "1200", "http://localhost", uint64(1), "ONLINE", int64(10), int64(1000), int64(10)))
	m.ExpectQuery("SELECT.*FROM provider_subscription_rates.*").
		WithArgs(int64(1)).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "provider_id", "token_name", "token_amount"}).
				AddRow(int64(1), int64(1), "uarkeo", int64(100)),
		)
	m.ExpectQuery("SELECT.*FROM provider_pay_as_you_go_rates.*").
		WithArgs(int64(1)).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "provider_id", "token_name", "token_amount"}).
				AddRow(int64(1), int64(1), "uarkeo", int64(100)),
		)
	p, err := db.FindProvider(context.Background(), testPubKey.String(), "mock")
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.Nil(t, m.ExpectationsWereMet())
}

func TestUpdateValidatorPayoutEvent(t *testing.T) {
	m, db := getMockDirectoryDBForTest(t)
	defer m.Close()
	testTime := time.Now()
	testAddr := arkeotypes.GetRandomBech32Addr()
	evt := arkeotypes.EventValidatorPayout{
		Validator: testAddr,
		Reward:    math.NewInt(1024),
	}
	m.ExpectQuery("insert into validator_payout_events.*").
		WithArgs(testAddr.String(), int64(1), int64(1024)).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "created", "updated"}).
				AddRow(int64(1), testTime, testTime))

	entity, err := db.UpsertValidatorPayoutEvent(context.Background(), evt, 1)
	assert.Nil(t, err)
	assert.NotNil(t, entity)
	assert.Equal(t, int64(1), entity.ID)
	assert.Equal(t, testTime, entity.Created)
	assert.Equal(t, testTime, entity.Updated)
	assert.Nil(t, m.ExpectationsWereMet())
}

func TestInsertBondProviderEvent(t *testing.T) {
	m, db := getMockDirectoryDBForTest(t)
	defer m.Close()
	testTime := time.Now()
	txHash := arkeotypes.GetRandomTxID()
	testPubKey := arkeotypes.GetRandomPubKey()
	m.ExpectQuery("insert into provider_bond_events.*").
		WithArgs(int64(1), int64(1024), txHash, "1000", "2000").
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "created", "updated"}).
				AddRow(int64(1), testTime, testTime))
	entity, err := db.InsertBondProviderEvent(context.Background(), int64(1), arkeotypes.EventBondProvider{
		Provider: testPubKey,
		Service:  "mock",
		BondRel:  math.NewInt(1000),
		BondAbs:  math.NewInt(2000),
	}, 1024, txHash)
	assert.Nil(t, err)
	assert.NotNil(t, entity)
	assert.Equal(t, int64(1), entity.ID)
	assert.Equal(t, testTime, entity.Created)
	assert.Equal(t, testTime, entity.Updated)
	assert.Nil(t, m.ExpectationsWereMet())
}

func TestInsertModProviderEvent(t *testing.T) {
	m, db := getMockDirectoryDBForTest(t)
	defer m.Close()
	testTime := time.Now()
	testPubKey := arkeotypes.GetRandomPubKey()
	evt := types.ModProviderEvent{
		Pubkey:              testPubKey.String(),
		Service:             "mock",
		Height:              1024,
		TxID:                arkeotypes.GetRandomTxID(),
		MetadataURI:         "http://localhost",
		MetadataNonce:       1,
		Status:              "ONLINE",
		MinContractDuration: 10,
		MaxContractDuration: 1000,
		SettlementDuration:  10,
		SubscriptionRate: []cosmostypes.Coin{
			cosmostypes.NewCoin("uarkeo", math.NewInt(10)),
		},
		PayAsYouGoRate: []cosmostypes.Coin{
			cosmostypes.NewCoin("uarkeo", math.NewInt(10)),
		},
	}
	m.ExpectQuery("insert into provider_mod_events.*").
		WithArgs(int64(1), evt.Height, evt.TxID, evt.MetadataURI, evt.MetadataNonce, evt.Status,
			evt.MinContractDuration, evt.MaxContractDuration).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "created", "updated"}).
				AddRow(int64(1), testTime, testTime))

	entity, err := db.InsertModProviderEvent(context.Background(), int64(1), evt, evt.TxID, evt.Height)
	assert.Nil(t, err)
	assert.NotNil(t, entity)
	assert.Equal(t, int64(1), entity.ID)
	assert.Equal(t, testTime, entity.Created)
	assert.Equal(t, testTime, entity.Updated)
	assert.Nil(t, m.ExpectationsWereMet())
}

func TestUpsertProviderMetadata(t *testing.T) {
	m, db := getMockDirectoryDBForTest(t)
	defer m.Close()
	testTime := time.Now()
	testPubKey := arkeotypes.GetRandomPubKey()
	metadata := sentinel.Metadata{
		Configuration: conf.Configuration{
			Moniker:                     "whatever",
			Website:                     "http://www.whatever.com",
			Description:                 "aha",
			Location:                    "",
			Port:                        "1234",
			SourceChain:                 "arkeo",
			EventStreamHost:             "arkeo",
			ClaimStoreLocation:          "arkeo",
			ContractConfigStoreLocation: "arkeo",
			ProviderPubKey:              testPubKey,
			FreeTierRateLimit:           0,
		},
		Version: "1",
	}
	m.ExpectQuery("insert into provider_metadata.*").
		WithArgs(int64(1), int64(1), metadata.Configuration.Moniker,
			metadata.Configuration.Website,
			metadata.Configuration.Description,
			sql.NullString{Valid: false},
			metadata.Configuration.FreeTierRateLimit).
		WillReturnRows(
			pgxmock.NewRows([]string{"id", "created", "updated"}).
				AddRow(int64(1), testTime, testTime))
	entity, err := db.UpsertProviderMetadata(context.Background(), 1, 1, metadata)
	assert.Nil(t, err)
	assert.NotNil(t, entity)
	assert.Equal(t, int64(1), entity.ID)
	assert.Equal(t, testTime, entity.Created)
	assert.Equal(t, testTime, entity.Updated)
	assert.Nil(t, m.ExpectationsWereMet())
}
