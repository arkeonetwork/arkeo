package db

import (
	"github.com/stretchr/testify/mock"

	"github.com/arkeonetwork/arkeo/directory/types"
	"github.com/arkeonetwork/arkeo/sentinel"
	atypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

var _ IDataStorage = &MockDataStorage{}

type MockDataStorage struct {
	mock.Mock
}

func (s *MockDataStorage) FindLatestBlock() (*Block, error) {
	args := s.Called()
	return args.Get(0).(*Block), args.Error(1)
}

func (s *MockDataStorage) InsertBlock(b *Block) (*Entity, error) {
	args := s.Called(b)
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) UpsertValidatorPayoutEvent(evt types.ValidatorPayoutEvent) (*Entity, error) {
	args := s.Called(evt)
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) FindProvider(pubkey, service string) (*ArkeoProvider, error) {
	args := s.Called(pubkey, service)
	return args.Get(0).(*ArkeoProvider), args.Error(1)
}

func (s *MockDataStorage) UpsertContract(providerID int64, evt atypes.EventOpenContract) (*Entity, error) {
	args := s.Called(providerID, evt)
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) FindContract(contractId uint64) (*ArkeoContract, error) {
	args := s.Called(contractId)
	return args.Get(0).(*ArkeoContract), args.Error(1)
}

func (s *MockDataStorage) CloseContract(contractID uint64, height int64) (*Entity, error) {
	args := s.Called(contractID, height)
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) UpdateProvider(provider *ArkeoProvider) (*Entity, error) {
	args := s.Called(provider)
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) UpsertContractSettlementEvent(evt types.ContractSettlementEvent) (*Entity, error) {
	args := s.Called(evt)
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) UpsertProviderMetadata(providerID, nonce int64, data sentinel.Metadata) (*Entity, error) {
	args := s.Called(providerID, nonce, data)
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) InsertBondProviderEvent(providerID int64, evt types.BondProviderEvent) (*Entity, error) {
	args := s.Called(providerID, evt)
	return args.Get(0).(*Entity), args.Error(1)
}

func (s *MockDataStorage) InsertProvider(provider *ArkeoProvider) (*Entity, error) {
	args := s.Called(provider)
	return args.Get(0).(*Entity), args.Error(1)
}
