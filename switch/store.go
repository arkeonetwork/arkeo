package switchd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mercury/x/mercury/types"
	"net/http"
	"time"
)

type MemStore struct {
	db          map[string]types.Contract
	client      http.Client
	baseURL     string
	blockHeight int64
}

func NewStore(baseURL string) *MemStore {
	return &MemStore{
		db: make(map[string]types.Contract),
		client: http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (k MemStore) SetHeight(height int64) {
	k.blockHeight = height
}

func (k MemStore) Get(key string) (types.Contract, error) {
	contract := k.db[key]
	if contract.IsClose(k.blockHeight) {
		return k.fetchContract(key)
	}
	return contract, nil
}

func (k MemStore) Put(key string, value types.Contract) {
	k.db[key] = value
}

func (k MemStore) fetchContract(key string) (types.Contract, error) {
	var contract types.Contract
	requestURL := fmt.Sprintf("%s/%s", k.baseURL, key)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
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
