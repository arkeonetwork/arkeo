package sentinel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mercury/x/mercury/types"
	"net/http"
	"time"
)

// TODO: this should receive events from arceo chain to update its database
// TODO: clean up contracts from memory after they expire
type MemStore struct {
	db          map[string]types.Contract
	client      http.Client
	baseURL     string
	blockHeight int64
}

func NewMemStore(baseURL string) *MemStore {
	return &MemStore{
		db: make(map[string]types.Contract),
		client: http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (k *MemStore) Key(pubkey, chain, spender string) string {
	return fmt.Sprintf("%s/%s/%s", pubkey, chain, spender)
}

func (k *MemStore) GetHeight() int64 {
	return k.blockHeight
}

func (k *MemStore) SetHeight(height int64) {
	k.blockHeight = height
}

func (k *MemStore) Get(key string) (types.Contract, error) {
	contract := k.db[key]
	if contract.IsClose(k.blockHeight) {
		return k.fetchContract(key)
	}
	return contract, nil
}

func (k *MemStore) Put(key string, value types.Contract) {
	k.db[key] = value
}

func (k *MemStore) fetchContract(key string) (types.Contract, error) {
	// TODO: this should cache a "miss" for 5 seconds, to stop DoS/thrashing

	var contract types.Contract
	requestURL := fmt.Sprintf("%s/%s", k.baseURL, key)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println(err)
		return contract, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return contract, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return contract, err
	}

	err = json.Unmarshal(resBody, &contract)
	if err != nil {
		return contract, err
	}

	if contract.IsClose(k.blockHeight) {
		delete(k.db, key) // clean up
		return types.Contract{}, nil
	}

	k.Put(key, contract)

	return contract, nil
}
