package sentinel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

var ModuleBasics = module.NewBasicManager()

// TODO: this should receive events from arceo chain to update its database
// TODO: clean up contracts from memory after they expire
type MemStore struct {
	storeLock   *sync.Mutex
	db          map[string]types.Contract
	client      http.Client
	baseURL     string
	blockHeight int64
	logger      log.Logger
}

func NewMemStore(baseURL string, logger log.Logger) *MemStore {
	return &MemStore{
		storeLock: &sync.Mutex{},
		db:        make(map[string]types.Contract),
		client: http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
		logger:  logger,
	}
}

func (k *MemStore) Key(pubkey, service, spender string) string {
	return fmt.Sprintf("%s/%s/%s", pubkey, service, spender)
}

func (k *MemStore) GetHeight() int64 {
	return k.blockHeight
}

func (k *MemStore) SetHeight(height int64) {
	k.blockHeight = height
}

func (k *MemStore) Get(key string) (types.Contract, error) {
	k.storeLock.Lock()
	defer k.storeLock.Unlock()
	contract, ok := k.db[key]
	// contract is not in cache or contract expired , fetch it
	if !ok || contract.IsExpired(k.blockHeight) {
		crtUpStream, err := k.fetchContract(key)
		if err != nil {
			return crtUpStream, err
		}
		if !crtUpStream.IsExpired(k.blockHeight) {
			k.db[key] = crtUpStream
		}
		return crtUpStream, nil
	}
	// contract still valid
	return contract, nil
}

func (k *MemStore) Put(contract types.Contract) {
	k.storeLock.Lock()
	defer k.storeLock.Unlock()
	key := contract.Key()
	if contract.IsExpired(k.blockHeight) {
		delete(k.db, key)
		return
	}
	k.db[key] = contract
}

func (k *MemStore) GetActiveContract(provider common.PubKey, service common.Service, spender common.PubKey) (types.Contract, error) {
	k.storeLock.Lock()
	defer k.storeLock.Unlock()
	// iterate through the map to find the contract
	for _, contract := range k.db {
		if !contract.IsExpired(k.GetHeight()) && contract.Provider.Equals(provider) && contract.Service == service && contract.GetSpender().Equals(spender) {
			return contract, nil
		}
	}
	// we should also probably call arkeo if we don't find the contract as we do below.
	return types.Contract{}, fmt.Errorf("contract not found")
}

func (k *MemStore) fetchContract(key string) (types.Contract, error) {
	// TODO: this should cache a "miss" for 5 seconds, to stop DoS/thrashing
	type fetch struct {
		Contract types.Contract `json:"contract"`
	}

	var data fetch
	requestURL := fmt.Sprintf("%s/arkeo/contract/%s", k.baseURL, key)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		k.logger.Error("fail to create http request", "error", err)
		return types.Contract{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		k.logger.Error("fail to send http request", "error", err)
		return types.Contract{}, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		k.logger.Error("fail to read from response body", "error", err)
		return types.Contract{}, err
	}

	err = json.Unmarshal(resBody, &data)
	if err != nil {
		k.logger.Error("fail to unmarshal response", "error", err)
		return types.Contract{}, err
	}
	return data.Contract, nil
}
