package sentinel

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

type NonceStore struct {
	logger zerolog.Logger
	db     *leveldb.DB
}

type NonceRecord struct {
	ContractId uint64 `json:"contract_id"`
	Nonce      int64  `json:"nonce"`
	UpdatedAt  int64  `json:"updated_at"`
}

func NewNonceStore(levelDbFolder string) (*NonceStore, error) {
	var db *leveldb.DB
	var err error
	if len(levelDbFolder) == 0 {
		log.Warn().Msg("nonce store folder is empty, create in memory storage")
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
	return &NonceStore{
		logger: log.With().Str("module", "nonce-storage").Logger(),
		db:     db,
	}, nil
}

func (s *NonceStore) Get(contractId uint64) (int64, error) {
	key := strconv.FormatUint(contractId, 10)
	exist, err := s.db.Has([]byte(key), nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to check nonce existence")
		return 0, err
	}
	if !exist {
		return 0, nil // Start from 0 if not found
	}

	value, err := s.db.Get([]byte(key), nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to get nonce record")
		return 0, err
	}
	
	var record NonceRecord
	if err := json.Unmarshal(value, &record); err != nil {
		s.logger.Error().Err(err).Msg("fail to unmarshal nonce record")
		return 0, err
	}
	
	return record.Nonce, nil
}

func (s *NonceStore) Set(contractId uint64, nonce int64) error {
	record := NonceRecord{
		ContractId: contractId,
		Nonce:      nonce,
		UpdatedAt:  time.Now().Unix(),
	}
	
	key := strconv.FormatUint(contractId, 10)
	buf, err := json.Marshal(record)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to marshal nonce record")
		return err
	}
	
	if err := s.db.Put([]byte(key), buf, nil); err != nil {
		s.logger.Error().Err(err).Msg("fail to set nonce record")
		return err
	}
	
	return nil
}

func (s *NonceStore) Close() error {
	return s.db.Close()
}