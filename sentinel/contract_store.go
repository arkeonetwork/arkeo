package sentinel

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type ContractConfigurationStore struct {
	logger zerolog.Logger
	db     *leveldb.DB
}

type CORs struct {
	AllowOrigins []string `json:"allow_origins"`
	AllowMethods []string `json:"allow_methods"`
	AllowHeaders []string `json:"allow_headers"`
}

func NewCORs() CORs {
	return CORs{
		AllowOrigins: make([]string, 0),
		AllowMethods: make([]string, 0),
		AllowHeaders: make([]string, 0),
	}
}

type ContractConfiguration struct {
	ContractId           uint64   `json:"contract_id"`
	LastTimeStamp        int64    `json:"last_timestamp"`
	PerUserRateLimit     int      `json:"per_user_rate_limit"`
	CORs                 CORs     `json:"cors"`
	WhitelistIPAddresses []string `json:"white_listed_ip_addresses"`
}

func (c ContractConfiguration) Key() string {
	return strconv.FormatUint(c.ContractId, 10)
}

type ContractConfigurations []ContractConfiguration

func NewContractConfiguration(contractId uint64, cors CORs, ips []string, rateLimit int) ContractConfiguration {
	return ContractConfiguration{
		ContractId:           contractId,
		LastTimeStamp:        0,
		CORs:                 cors,
		WhitelistIPAddresses: ips,
		PerUserRateLimit:     rateLimit,
	}
}

func NewContractConfigurationStore(levelDbFolder string) (*ContractConfigurationStore, error) {
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
	return &ContractConfigurationStore{
		logger: log.With().Str("module", "contract-config-storage").Logger(),
		db:     db,
	}, nil
}

func (s *ContractConfigurationStore) Set(item ContractConfiguration) error {
	key := item.Key()
	buf, err := json.Marshal(item)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to marshal to claim store item")
		return err
	}
	if err := s.db.Put([]byte(key), buf, nil); err != nil {
		s.logger.Error().Err(err).Msg("fail to set claim item")
		return err
	}
	return nil
}

func (s *ContractConfigurationStore) Batch(items ContractConfigurations) error {
	batch := new(leveldb.Batch)
	for _, item := range items {
		key := item.Key()
		buf, err := json.Marshal(item)
		if err != nil {
			s.logger.Error().Err(err).Msg("fail to marshal to contract configuration store item")
			return err
		}
		batch.Put([]byte(key), buf)
	}
	return s.db.Write(batch, nil)
}

func (s *ContractConfigurationStore) Get(id uint64) (item ContractConfiguration, err error) {
	item = NewContractConfiguration(id, NewCORs(), make([]string, 0), 0)
	key := item.Key()
	ok, err := s.db.Has([]byte(key), nil)
	if !ok || err != nil {
		return
	}
	buf, err := s.db.Get([]byte(key), nil)
	if err := json.Unmarshal(buf, &item); err != nil {
		s.logger.Error().Err(err).Msg("fail to unmarshal to contract configuration store item")
		return item, err
	}

	return
}

// Has check whether the given key exist in key value store
func (s *ContractConfigurationStore) Has(id uint64) (ok bool) {
	key := strconv.FormatUint(id, 10)
	ok, _ = s.db.Has([]byte(key), nil)
	return
}

// Remove remove the given item from key values store
func (s *ContractConfigurationStore) Remove(id uint64) error {
	key := strconv.FormatUint(id, 10)
	return s.db.Delete([]byte(key), nil)
}

// List send back tx out to retry depending on arg failed only
func (s *ContractConfigurationStore) List() ContractConfigurations {
	iterator := s.db.NewIterator(util.BytesPrefix([]byte(nil)), nil)
	defer iterator.Release()
	var results ContractConfigurations
	for iterator.Next() {
		buf := iterator.Value()
		if len(buf) == 0 {
			continue
		}

		var item ContractConfiguration
		if err := json.Unmarshal(buf, &item); err != nil {
			s.logger.Error().Err(err).Msg("fail to unmarshal to contract configuration store item")
			continue
		}

		results = append(results, item)
	}

	return results
}

// Close underlying db
func (s *ContractConfigurationStore) Close() error {
	return s.db.Close()
}

func (s *ContractConfigurationStore) GetInternalDb() *leveldb.DB {
	return s.db
}
