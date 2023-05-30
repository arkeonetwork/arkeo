package indexer

import (
	"fmt"
	"sync"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/arkeonetwork/arkeo/common/logging"
	"github.com/arkeonetwork/arkeo/directory/db"
	"github.com/arkeonetwork/arkeo/directory/types"
	arkeotypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func TestCreateProvider(t *testing.T) {
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

	// fail to insert provider should result an error
	failCreateProvider := mockDb.On("InsertProvider", mock.Anything).Return(nil, fmt.Errorf("fail to add provider"))
	result, err := s.createProvider(types.BondProviderEvent{
		Pubkey:       arkeotypes.GetRandomPubKey().String(),
		Service:      "mock",
		Height:       0,
		TxID:         arkeotypes.GetRandomTxID(),
		BondRelative: "1",
		BondAbsolute: "1",
	})
	assert.NotNil(t, err)
	assert.Nil(t, result)
	failCreateProvider.Unset()

	// when insert provider fails which result in a nil entity, it should return an error
	failCreateProvider = mockDb.On("InsertProvider", mock.Anything).Return(nil, nil)
	result, err = s.createProvider(types.BondProviderEvent{
		Pubkey:       arkeotypes.GetRandomPubKey().String(),
		Service:      "mock",
		Height:       0,
		TxID:         arkeotypes.GetRandomTxID(),
		BondRelative: "1",
		BondAbsolute: "1",
	})
	assert.NotNil(t, err)
	assert.Nil(t, result)

	// happy path
	failCreateProvider.Unset()

	mockDb.On("InsertProvider", mock.Anything).Return(&db.Entity{
		ID:      0,
		Created: time.Now(),
		Updated: time.Now(),
	}, nil)
	result, err = s.createProvider(types.BondProviderEvent{
		Pubkey:       arkeotypes.GetRandomPubKey().String(),
		Service:      "mock",
		Height:       0,
		TxID:         arkeotypes.GetRandomTxID(),
		BondRelative: "1",
		BondAbsolute: "1",
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

func TestHandleBondProviderEvent(t *testing.T) {
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
	testPubKey := arkeotypes.GetRandomPubKey().String()

	// fail to find provider should result in an error
	mockFindProvider := mockDb.On("FindProvider", testPubKey, "mock").Return(nil, fmt.Errorf("fail to find provider"))
	err := s.handleBondProviderEvent(types.BondProviderEvent{
		Pubkey:       testPubKey,
		Service:      "mock",
		Height:       0,
		TxID:         arkeotypes.GetRandomTxID(),
		BondRelative: "1",
		BondAbsolute: "1",
	})
	assert.NotNil(t, err)
	mockFindProvider.Unset()

	// when provider doesn't exist , it should try to create one
	// if it fail to create , then it should return an error
	mockFindProvider = mockDb.On("FindProvider", testPubKey, "mock").Return(nil, nil)
	mockInsertProvider := mockDb.On("InsertProvider", mock.Anything).Return(nil, fmt.Errorf("fail to create provider"))
	err = s.handleBondProviderEvent(types.BondProviderEvent{
		Pubkey:       testPubKey,
		Service:      "mock",
		Height:       0,
		TxID:         arkeotypes.GetRandomTxID(),
		BondRelative: "1",
		BondAbsolute: "1",
	})
	assert.NotNil(t, err)

	mockInsertProvider.Unset()

	// fail to insert bond provider event should result in an error
	mockInsertProvider = mockDb.On("InsertProvider", mock.Anything).Return(&db.Entity{
		ID:      0,
		Created: time.Now(),
		Updated: time.Now(),
	}, nil)

	mockInsertBondProviderEvent := mockDb.On("InsertBondProviderEvent", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("fail to insert bond provider"))
	err = s.handleBondProviderEvent(types.BondProviderEvent{
		Pubkey:       testPubKey,
		Service:      "mock",
		Height:       0,
		TxID:         arkeotypes.GetRandomTxID(),
		BondRelative: "1",
		BondAbsolute: "1",
	})
	assert.NotNil(t, err)
	mockInsertBondProviderEvent.Unset()

	// happy path
	mockInsertBondProviderEvent = mockDb.On("InsertBondProviderEvent", mock.Anything, mock.Anything).Return(&db.Entity{
		ID:      0,
		Created: time.Now(),
		Updated: time.Now(),
	}, nil)
	err = s.handleBondProviderEvent(types.BondProviderEvent{
		Pubkey:       testPubKey,
		Service:      "mock",
		Height:       0,
		TxID:         arkeotypes.GetRandomTxID(),
		BondRelative: "1",
		BondAbsolute: "1",
	})
	assert.Nil(t, err)

	// when a bond provider already exists, it should update the provider
	// when update provider fails, it should return an error
	mockFindProvider.Unset()
	mockFindProvider = mockDb.On("FindProvider", testPubKey, "mock").Return(&db.ArkeoProvider{}, nil)
	mockUpdateProvider := mockDb.On("UpdateProvider", mock.Anything).Return(nil, fmt.Errorf("fail to update provider"))
	err = s.handleBondProviderEvent(types.BondProviderEvent{
		Pubkey:       testPubKey,
		Service:      "mock",
		Height:       0,
		TxID:         arkeotypes.GetRandomTxID(),
		BondRelative: "1",
		BondAbsolute: "1",
	})
	assert.NotNil(t, err)

	// happy path with update provider
	mockUpdateProvider.Unset()
	mockUpdateProvider = mockDb.On("UpdateProvider", mock.Anything).Return(&db.Entity{
		ID:      0,
		Created: time.Now(),
		Updated: time.Now(),
	}, nil)
	err = s.handleBondProviderEvent(types.BondProviderEvent{
		Pubkey:       testPubKey,
		Service:      "mock",
		Height:       0,
		TxID:         arkeotypes.GetRandomTxID(),
		BondRelative: "1",
		BondAbsolute: "1",
	})
	assert.Nil(t, err)
}

func TestHandleModProviderEvent(t *testing.T) {
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

	// fail to find provider should result in an error
	mockFindProvider := mockDb.On("FindProvider", testPubKey.String(), "mock").Return(nil, fmt.Errorf("fail to find provider"))
	err := s.handleModProviderEvent(arkeotypes.EventModProvider{
		Creator:             arkeotypes.GetRandomBech32Addr(),
		Provider:            testPubKey,
		Service:             "mock",
		MetadataUri:         "",
		MetadataNonce:       0,
		Status:              0,
		MinContractDuration: 0,
		MaxContractDuration: 0,
		SubscriptionRate:    nil,
		PayAsYouGoRate:      nil,
		Bond:                sdk.NewInt(100),
		SettlementDuration:  0,
	})
	assert.NotNil(t, err)
	mockFindProvider.Unset()

	// fail to update provider cause an error
	mockFindProvider = mockDb.On("FindProvider", testPubKey.String(), "mock").Return(&db.ArkeoProvider{}, nil)
	mockUpdateProvider := mockDb.On("UpdateProvider", mock.Anything).Return(nil, fmt.Errorf("fail to update provider"))
	err = s.handleModProviderEvent(arkeotypes.EventModProvider{
		Creator:             arkeotypes.GetRandomBech32Addr(),
		Provider:            testPubKey,
		Service:             "mock",
		MetadataUri:         "",
		MetadataNonce:       0,
		Status:              0,
		MinContractDuration: 0,
		MaxContractDuration: 0,
		SubscriptionRate:    nil,
		PayAsYouGoRate:      nil,
		Bond:                sdk.NewInt(100),
		SettlementDuration:  0,
	})
	assert.NotNil(t, err)
	mockUpdateProvider.Unset()
	mockUpdateProvider = mockDb.On("UpdateProvider", mock.Anything).Return(&db.Entity{
		ID:      0,
		Created: time.Now(),
		Updated: time.Now(),
	}, nil)
	err = s.handleModProviderEvent(arkeotypes.EventModProvider{
		Creator:             arkeotypes.GetRandomBech32Addr(),
		Provider:            testPubKey,
		Service:             "mock",
		MetadataUri:         "",
		MetadataNonce:       0,
		Status:              0,
		MinContractDuration: 0,
		MaxContractDuration: 0,
		SubscriptionRate:    nil,
		PayAsYouGoRate:      nil,
		Bond:                sdk.NewInt(100),
		SettlementDuration:  0,
	})
	assert.Nil(t, err)
}
