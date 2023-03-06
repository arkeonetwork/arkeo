package sentinel

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/arkeonetwork/arkeo/common"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type ClaimStore struct {
	logger zerolog.Logger
	db     *leveldb.DB
}

type Claim struct {
	Provider   common.PubKey
	ContractId uint64
	Spender    common.PubKey
	Nonce      int64
	Height     int64
	Signature  string
	Claimed    bool
}

func NewClaim(contractId uint64, spender common.PubKey, nonce, height int64, signature string) Claim {
	return Claim{
		ContractId: contractId,
		Spender:    spender,
		Nonce:      nonce,
		Height:     height,
		Signature:  signature,
		Claimed:    false,
	}
}

func NewClaimStore(levelDbFolder string) (*ClaimStore, error) {
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
	return &ClaimStore{
		logger: log.With().Str("module", "claim-storage").Logger(),
		db:     db,
	}, nil
}

func (s *ClaimStore) Set(item Claim) error {
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

func (s *ClaimStore) Batch(items []Claim) error {
	batch := new(leveldb.Batch)
	for _, item := range items {
		key := item.Key()
		buf, err := json.Marshal(item)
		if err != nil {
			s.logger.Error().Err(err).Msg("fail to marshal to claim store item")
			return err
		}
		batch.Put([]byte(key), buf)
	}
	return s.db.Write(batch, nil)
}

func (s *ClaimStore) Get(key string) (item Claim, err error) {
	ok, err := s.db.Has([]byte(key), nil)
	if !ok || err != nil {
		return
	}
	buf, err := s.db.Get([]byte(key), nil)
	if err := json.Unmarshal(buf, &item); err != nil {
		s.logger.Error().Err(err).Msg("fail to unmarshal to claim store item")
		return item, err
	}

	return
}

// Has check whether the given key exist in key value store
func (s *ClaimStore) Has(key string) (ok bool) {
	ok, _ = s.db.Has([]byte(key), nil)
	return
}

// Remove remove the given item from key values store
func (s *ClaimStore) Remove(key string) error {
	return s.db.Delete([]byte(key), nil)
}

// List send back tx out to retry depending on arg failed only
func (s *ClaimStore) List() []Claim {
	iterator := s.db.NewIterator(util.BytesPrefix([]byte(nil)), nil)
	defer iterator.Release()
	var results []Claim
	for iterator.Next() {
		buf := iterator.Value()
		if len(buf) == 0 {
			continue
		}

		var item Claim
		if err := json.Unmarshal(buf, &item); err != nil {
			s.logger.Error().Err(err).Msg("fail to unmarshal to claim store item")
			continue
		}

		results = append(results, item)
	}

	return results
}

// Close underlying db
func (s *ClaimStore) Close() error {
	return s.db.Close()
}

func (s *ClaimStore) GetInternalDb() *leveldb.DB {
	return s.db
}

func (c Claim) Key() string {
	return strconv.FormatUint(c.ContractId, 10)
}
