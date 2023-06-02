package db

import (
	"github.com/stretchr/testify/mock"

	"github.com/arkeonetwork/arkeo/sentinel"
	atypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

var _ IDataStorage = &MockDataStorage{}

type MockDataStorage struct {
	mock.Mock
}

func (s *MockDataStorage) FindLatestBlock() (*Block, error) {
	args := s.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Block), args.Error(1)
}

func (s *MockDataStorage) InsertBlock(b *Block) (*Entity, error) {
	args := s.Called(b)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) UpsertValidatorPayoutEvent(evt atypes.EventValidatorPayout, height int64) (*Entity, error) {
	args := s.Called(evt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) FindProvider(pubkey, service string) (*ArkeoProvider, error) {
	args := s.Called(pubkey, service)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*ArkeoProvider), args.Error(1)
}

func (s *MockDataStorage) UpsertContract(providerID int64, evt atypes.EventOpenContract) (*Entity, error) {
	args := s.Called(providerID, evt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) GetContract(contractId uint64) (*ArkeoContract, error) {
	args := s.Called(contractId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*ArkeoContract), args.Error(1)
}

func (s *MockDataStorage) CloseContract(contractID uint64, height int64) (*Entity, error) {
	args := s.Called(contractID, height)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) UpdateProvider(provider *ArkeoProvider) (*Entity, error) {
	args := s.Called(provider)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) UpsertContractSettlementEvent(evt atypes.EventSettleContract) (*Entity, error) {
	args := s.Called(evt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) UpsertProviderMetadata(providerID, nonce int64, data sentinel.Metadata) (*Entity, error) {
	args := s.Called(providerID, nonce, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) InsertBondProviderEvent(providerID int64, evt atypes.EventBondProvider, height int64, txID string) (*Entity, error) {
	args := s.Called(providerID, evt, height, txID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) InsertProvider(provider *ArkeoProvider) (*Entity, error) {
	args := s.Called(provider)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}
