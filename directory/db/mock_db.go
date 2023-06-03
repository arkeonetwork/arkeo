package db

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/arkeonetwork/arkeo/sentinel"
	atypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

var _ IDataStorage = &MockDataStorage{}

type MockDataStorage struct {
	mock.Mock
}

func (s *MockDataStorage) FindLatestBlock(ctx context.Context) (*Block, error) {
	args := s.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Block), args.Error(1)
}

func (s *MockDataStorage) InsertBlock(ctx context.Context, b *Block) (*Entity, error) {
	args := s.Called(ctx, b)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) UpsertValidatorPayoutEvent(ctx context.Context, evt atypes.EventValidatorPayout, height int64) (*Entity, error) {
	args := s.Called(ctx, evt, height)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) FindProvider(ctx context.Context, pubkey, service string) (*ArkeoProvider, error) {
	args := s.Called(ctx, pubkey, service)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*ArkeoProvider), args.Error(1)
}

func (s *MockDataStorage) UpsertContract(ctx context.Context, providerID int64, evt atypes.EventOpenContract) (*Entity, error) {
	args := s.Called(ctx, providerID, evt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) GetContract(ctx context.Context, contractId uint64) (*ArkeoContract, error) {
	args := s.Called(ctx, contractId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*ArkeoContract), args.Error(1)
}

func (s *MockDataStorage) CloseContract(ctx context.Context, contractID uint64, height int64) (*Entity, error) {
	args := s.Called(ctx, contractID, height)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) UpdateProvider(ctx context.Context, provider *ArkeoProvider) (*Entity, error) {
	args := s.Called(ctx, provider)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) UpsertContractSettlementEvent(ctx context.Context, evt atypes.EventSettleContract) (*Entity, error) {
	args := s.Called(ctx, evt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) UpsertProviderMetadata(ctx context.Context, providerID, nonce int64, data sentinel.Metadata) (*Entity, error) {
	args := s.Called(ctx, providerID, nonce, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) InsertBondProviderEvent(ctx context.Context, providerID int64, evt atypes.EventBondProvider, height int64, txID string) (*Entity, error) {
	args := s.Called(ctx, providerID, evt, height, txID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) InsertProvider(ctx context.Context, provider *ArkeoProvider) (*Entity, error) {
	args := s.Called(ctx, provider)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	//nolint:forcetypeassert
	return args.Get(0).(*Entity), args.Error(1)
}
