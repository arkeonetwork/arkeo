package indexer

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"cosmossdk.io/math"
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/arkeonetwork/arkeo/common/logging"
	"github.com/arkeonetwork/arkeo/directory/db"
	arkeotypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func TestHandleOpenContractEvent(t *testing.T) {
	mockDb := new(db.MockDataStorage)
	s := Service{
		params:         ServiceParams{},
		db:             mockDb,
		done:           make(chan struct{}),
		wg:             &sync.WaitGroup{},
		logger:         logging.WithoutFields(),
		tmClient:       nil,
		blockFillQueue: make(chan db.BlockGap),
	}
	testPubKey := arkeotypes.GetRandomPubKey()
	mockFindProvider := mockDb.On("FindProvider", mock.Anything, testPubKey.String(), "mock").
		Return(nil, fmt.Errorf("fail to find provider"))
	eventOpenContract := arkeotypes.EventOpenContract{
		Provider:           testPubKey,
		ContractId:         2,
		Service:            "mock",
		Client:             arkeotypes.GetRandomPubKey(),
		Delegate:           arkeotypes.GetRandomPubKey(),
		Type:               arkeotypes.ContractType_PAY_AS_YOU_GO,
		Height:             1024,
		Duration:           10,
		Rate:               cosmostypes.NewCoin("uarkeo", math.NewInt(2)),
		OpenCost:           0,
		Deposit:            math.NewInt(100000),
		SettlementDuration: 10,
		Authorization:      arkeotypes.ContractAuthorization_STRICT,
		QueriesPerMinute:   10,
	}
	err := s.handleOpenContractEvent(context.Background(), eventOpenContract)
	assert.NotNil(t, err)
	mockFindProvider.Unset()
	mockDb.On("FindProvider", mock.Anything, testPubKey.String(), "mock").Return(&db.ArkeoProvider{
		Entity: db.Entity{
			ID:      1,
			Created: time.Now(),
			Updated: time.Now(),
		},
	}, nil)
	mockUpdateContract := mockDb.On("UpsertContract", mock.Anything, int64(1), mock.Anything).Return(nil, fmt.Errorf("fail to update contract"))
	err = s.handleOpenContractEvent(context.Background(), eventOpenContract)
	assert.NotNil(t, err)
	mockUpdateContract.Unset()

	mockDb.On("UpsertContract", mock.Anything, int64(1), mock.Anything).Return(&db.Entity{
		ID:      2,
		Created: time.Now(),
		Updated: time.Now(),
	}, nil)
	err = s.handleOpenContractEvent(context.Background(), eventOpenContract)
	assert.Nil(t, err)
}

func TestHandleCloseContractEvent(t *testing.T) {
	mockDb := new(db.MockDataStorage)
	s := Service{
		params:         ServiceParams{},
		db:             mockDb,
		done:           make(chan struct{}),
		wg:             &sync.WaitGroup{},
		logger:         logging.WithoutFields(),
		tmClient:       nil,
		blockFillQueue: make(chan db.BlockGap),
	}
	testPubKey := arkeotypes.GetRandomPubKey()
	eventCloseContract := arkeotypes.EventCloseContract{
		ContractId: 1,
		Provider:   testPubKey,
		Service:    "mock",
		Client:     arkeotypes.GetRandomPubKey(),
		Delegate:   arkeotypes.GetRandomPubKey(),
	}
	mockCloseContract := mockDb.On("CloseContract", mock.Anything, uint64(1), int64(1)).Return(nil, fmt.Errorf("fail to close contract"))
	err := s.handleCloseContractEvent(context.Background(), eventCloseContract, 1)
	assert.NotNil(t, err)
	mockCloseContract.Unset()

	mockDb.On("CloseContract", mock.Anything, uint64(1), int64(1)).Return(&db.Entity{
		ID:      2,
		Created: time.Now(),
		Updated: time.Now(),
	}, nil)
	err = s.handleCloseContractEvent(context.Background(), eventCloseContract, 1)
	assert.Nil(t, err)
}

func TestHandleContractSettlementEvent(t *testing.T) {
	mockDb := new(db.MockDataStorage)
	s := Service{
		params:         ServiceParams{},
		db:             mockDb,
		done:           make(chan struct{}),
		wg:             &sync.WaitGroup{},
		logger:         logging.WithoutFields(),
		tmClient:       nil,
		blockFillQueue: make(chan db.BlockGap),
	}
	testPubKey := arkeotypes.GetRandomPubKey()
	mockSettlement := mockDb.On("UpsertContractSettlementEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("fail to upsert contract settlement event"))
	eventSettleContract := arkeotypes.EventSettleContract{
		Provider:   testPubKey,
		ContractId: 1,
		Service:    "mock",
		Client:     arkeotypes.GetRandomPubKey(),
		Delegate:   arkeotypes.GetRandomPubKey(),
		Type:       arkeotypes.ContractType_PAY_AS_YOU_GO,
		Nonce:      0,
		Height:     1024,
		Paid:       math.NewInt(1024),
		Reserve:    math.NewInt(100000),
	}
	err := s.handleContractSettlementEvent(context.Background(), eventSettleContract)
	assert.NotNil(t, err)
	mockSettlement.Unset()
	mockDb.On("UpsertContractSettlementEvent", mock.Anything, mock.Anything).Return(&db.Entity{
		ID:      1,
		Created: time.Now(),
		Updated: time.Now(),
	}, nil)
	err = s.handleContractSettlementEvent(context.Background(), eventSettleContract)
	assert.Nil(t, err)
}
