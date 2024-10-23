package sentinel

import (
	"encoding/json"
	"fmt"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

type ProviderConfigurationStore struct {
	logger zerolog.Logger
	db     *leveldb.DB
}

func NewProviderConfigurationStore(levelDbFolder string) (*ProviderConfigurationStore, error) {

	var db *leveldb.DB
	var err error
	if len(levelDbFolder) == 0 {
		log.Warn().Msg("level db folder is empty, create in memory storage")
		// no directory given, use in memory store
		storage := storage.NewMemStorage()
		db, err = leveldb.Open(storage, nil)
		if err != nil {
			return nil, fmt.Errorf("fail to in memory open level db: %w", err)
		}
	} else {
		db, err = leveldb.OpenFile(levelDbFolder, nil)
		if err != nil {
			return nil, fmt.Errorf("fail to open level db %s: %w", levelDbFolder, err)
		}
	}

	return &ProviderConfigurationStore{
		logger: log.With().Str("module", "provider-config-store").Logger(),
		db:     db,
	}, nil
}

type ProviderConfiguration struct {
	PubKey              common.PubKey        `json:"pubkey,omitempty"`
	Service             common.Service       `json:"service,omitempty"`
	Bond                cosmos.Int           `json:"bond,omitempty"`
	BondRelative        cosmos.Int           `json:"bond_relative,omitempty"`
	MetadataUri         string               `json:"metadata_uri,omitempty"`
	MetadataNonce       uint64               `json:"metadata_nonce,omitempty"`
	Status              types.ProviderStatus `json:"status,omitempty"`
	MinContractDuration int64                `json:"min_contract_duration,omitempty"`
	MaxContractDuration int64                `json:"max_contract_duration,omitempty"`
	SubscriptionRate    cosmos.Coins         `json:"subscription_rate"`
	PayAsYouGoRate      cosmos.Coins         `json:"pay_as_you_go_rate"`
	SettlementDuration  int64                `json:"settlement_duration,omitempty"`
}

// GetProviderModOrBondConfig retrieves a ProviderConfiguration by its PubKey
func (ps *ProviderConfigurationStore) Get(pubKey common.PubKey) (ProviderConfiguration, error) {
	data, err := ps.db.Get([]byte(pubKey.String()), nil)
	if err != nil {
		return ProviderConfiguration{}, err
	}

	var config ProviderConfiguration
	err = json.Unmarshal(data, &config)
	if err != nil {
		return ProviderConfiguration{}, err
	}
	return config, nil
}

// SetProviderModOrBondConfig saves or updates a ProviderConfiguration in the database
func (ps *ProviderConfigurationStore) Set(config ProviderConfiguration) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = ps.db.Put([]byte(config.PubKey.String()), data, nil)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProviderConfigurationStore) Remove(key common.PubKey) error {
	return p.db.Delete([]byte(key.String()), nil)
}
