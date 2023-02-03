package sentinel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"github.com/cosmos/cosmos-sdk/types/module"
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
}

func NewMemStore(baseURL string) *MemStore {
	return &MemStore{
		storeLock: &sync.Mutex{},
		db:        make(map[string]types.Contract),
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
	k.storeLock.Lock()
	defer k.storeLock.Unlock()
	contract, ok := k.db[key]
	// contract is not in cache or contract expired , fetch it
	if !ok || contract.IsClose(k.blockHeight) {
		crtUpStream, err := k.fetchContract(key)
		if err != nil {
			return crtUpStream, err
		}
		if !crtUpStream.IsClose(k.blockHeight) {
			k.db[key] = crtUpStream
		}
	}
	// contract still valid
	return contract, nil
}

func (k *MemStore) Put(key string, value types.Contract) {
	k.storeLock.Lock()
	defer k.storeLock.Unlock()
	if value.IsClose(k.blockHeight) {
		delete(k.db, key)
		return
	}
	k.db[key] = value
}

func (k *MemStore) fetchContract(key string) (types.Contract, error) {
	// TODO: this should cache a "miss" for 5 seconds, to stop DoS/thrashing
	var contract types.Contract

	type fetchContract struct {
		ProviderPubKey common.PubKey      `protobuf:"bytes,1,opt,name=provider_pub_key,json=providerPubKey,proto3,casttype=github.com/arkeonetwork/arkeo/common.PubKey" json:"provider_pub_key,omitempty"`
		Chain          common.Chain       `protobuf:"varint,2,opt,name=chain,proto3,casttype=github.com/arkeonetwork/arkeo/common.Chain" json:"chain,omitempty"`
		Client         common.PubKey      `protobuf:"bytes,3,opt,name=client,proto3,casttype=github.com/arkeonetwork/arkeo/common.PubKey" json:"client,omitempty"`
		Delegate       common.PubKey      `protobuf:"bytes,4,opt,name=delegate,proto3,casttype=github.com/arkeonetwork/arkeo/common.PubKey" json:"delegate,omitempty"`
		Type           types.ContractType `protobuf:"varint,5,opt,name=type,proto3,enum=arkeo.arkeo.ContractType" json:"type,omitempty"`
		Height         string             `protobuf:"varint,6,opt,name=height,proto3" json:"height,omitempty"`
		Duration       string             `protobuf:"varint,7,opt,name=duration,proto3" json:"duration,omitempty"`
		Rate           string             `protobuf:"varint,8,opt,name=rate,proto3" json:"rate,omitempty"`
		Deposit        string             `protobuf:"varint,9,opt,name=deposit,proto3" json:"deposit,omitempty"`
		Paid           string             `protobuf:"varint,10,opt,name=paid,proto3" json:"paid,omitempty"`
		Nonce          string             `protobuf:"varint,11,opt,name=nonce,proto3" json:"nonce,omitempty"`
		ClosedHeight   string             `protobuf:"varint,12,opt,name=closed_height,json=closedHeight,proto3" json:"closed_height,omitempty"`
	}

	type fetch struct {
		Contract fetchContract `json:"contract"`
	}

	var data fetch
	requestURL := fmt.Sprintf("%s/arkeo/contract/%s", k.baseURL, key)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println(err)
		return contract, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return contract, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return contract, err
	}

	err = json.Unmarshal(resBody, &data)
	if err != nil {
		fmt.Println(err)
		return contract, err
	}

	contract.ProviderPubKey = data.Contract.ProviderPubKey
	contract.Chain = data.Contract.Chain
	contract.Client = data.Contract.Client
	contract.Delegate = data.Contract.Delegate
	contract.Type = data.Contract.Type
	contract.Height, _ = strconv.ParseInt(data.Contract.Height, 10, 64)
	contract.Duration, _ = strconv.ParseInt(data.Contract.Duration, 10, 64)
	contract.Rate, _ = strconv.ParseInt(data.Contract.Rate, 10, 64)
	contract.Deposit, _ = cosmos.NewIntFromString(data.Contract.Deposit)
	contract.Paid, _ = cosmos.NewIntFromString(data.Contract.Paid)
	contract.Nonce, _ = strconv.ParseInt(data.Contract.Nonce, 10, 64)
	contract.ClosedHeight, _ = strconv.ParseInt(data.Contract.ClosedHeight, 10, 64)

	return contract, nil
}
