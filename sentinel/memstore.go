package sentinel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/sentinel/types"
	arkeoTypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

var ModuleBasics = module.NewBasicManager()

// TODO: this should receive events from arkeo chain to update its database
// TODO: clean up contracts from memory after they expire
type MemStore struct {
	storeLock   *sync.Mutex
	db          map[string]types.SentinelContract
	client      http.Client
	baseURL     string
	blockHeight int64
	logger      log.Logger
}

func NewMemStore(baseURL string, logger log.Logger) *MemStore {
	return &MemStore{
		storeLock: &sync.Mutex{},
		db:        make(map[string]types.SentinelContract),
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

func (k *MemStore) Get(key string) (types.SentinelContract, error) {
	k.storeLock.Lock()
	defer k.storeLock.Unlock()
	sentinelContract, ok := k.db[key]
	contract := sentinelContract.ArkeoContract
	// contract is not in cache or contract expired , fetch it
	if !ok || contract.IsExpired(k.blockHeight) {
		crtUpStream, err := k.fetchContract(key)
		if err != nil {
			return crtUpStream, err
		}
		if !crtUpStream.ArkeoContract.IsExpired(k.blockHeight) {
			k.db[key] = crtUpStream
		}
		return crtUpStream, nil
	}
	// contract still valid
	return sentinelContract, nil
}

func (k *MemStore) Put(contract types.SentinelContract) {
	k.storeLock.Lock()
	defer k.storeLock.Unlock()
	key := contract.ArkeoContract.Key()
	if contract.ArkeoContract.IsExpired(k.blockHeight) {
		delete(k.db, key)
		return
	}
	k.db[key] = contract
}

func (k *MemStore) GetActiveContract(provider common.PubKey, service common.Service, spender common.PubKey) (types.SentinelContract, error) {
	k.storeLock.Lock()
	defer k.storeLock.Unlock()
	// iterate through the map to find the contract
	for _, contract := range k.db {
		arkeoContract := contract.ArkeoContract
		if !arkeoContract.IsExpired(k.GetHeight()) && arkeoContract.Provider.Equals(provider) &&
			arkeoContract.Service == service && arkeoContract.GetSpender().Equals(spender) {
			return contract, nil
		}
	}
	// we should also probably call arkeo if we don't find the contract as we do below.
	return types.SentinelContract{}, fmt.Errorf("contract not found")
}

func (k *MemStore) fetchContract(key string) (types.SentinelContract, error) {
	// TODO: this should cache a "miss" for 5 seconds, to stop DoS/thrashing
	var sentinelContract types.SentinelContract
	var contract arkeoTypes.Contract
	sentinelContract.ArkeoContract = contract

	type fetchContract struct {
		ProviderPubKey   common.PubKey        `protobuf:"bytes,1,opt,name=provider_pub_key,json=providerPubKey,proto3,casttype=github.com/arkeonetwork/arkeo/common.PubKey" json:"provider_pub_key,omitempty"`
		Service          common.Service       `protobuf:"varint,2,opt,name=service,proto3,casttype=github.com/arkeonetwork/arkeo/common.Service" json:"service,omitempty"`
		Client           common.PubKey        `protobuf:"bytes,3,opt,name=client,proto3,casttype=github.com/arkeonetwork/arkeo/common.PubKey" json:"client,omitempty"`
		Delegate         common.PubKey        `protobuf:"bytes,4,opt,name=delegate,proto3,casttype=github.com/arkeonetwork/arkeo/common.PubKey" json:"delegate,omitempty"`
		MeterType        arkeoTypes.MeterType `protobuf:"varint,5,opt,name=meter_type,proto3,enum=arkeo.arkeo.MeterType" json:"meter_type,omitempty"`
		UserType         arkeoTypes.UserType  `protobuf:"varint,6,opt,name=user_type,proto3,enum=arkeo.arkeo.UserType" json:"user_type,omitempty"`
		Height           string               `protobuf:"varint,7,opt,name=height,proto3" json:"height,omitempty"`
		Duration         string               `protobuf:"varint,8,opt,name=duration,proto3" json:"duration,omitempty"`
		Rate             string               `protobuf:"varint,9,opt,name=rate,proto3" json:"rate,omitempty"`
		Deposit          string               `protobuf:"varint,10,opt,name=deposit,proto3" json:"deposit,omitempty"`
		Paid             string               `protobuf:"varint,11,opt,name=paid,proto3" json:"paid,omitempty"`
		Nonces           map[string]int64     `protobuf:"varint,12,opt,name=nonces,proto3,castkey=github.com/arkeonetwork/arkeo/common.PubKey" json:"nonce,omitempty"`
		SettlementHeight string               `protobuf:"varint,13,opt,name=settlement_height,json=settlementHeight,proto3" json:"settlement_height,omitempty"`
	}

	type fetch struct {
		Contract fetchContract `json:"contract"`
	}

	var data fetch
	requestURL := fmt.Sprintf("%s/arkeo/contract/%s", k.baseURL, key)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		k.logger.Error("fail to create http request", "error", err)
		return sentinelContract, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		k.logger.Error("fail to send http request", "error", err)
		return sentinelContract, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		k.logger.Error("fail to read from response body", "error", err)
		return sentinelContract, err
	}

	err = json.Unmarshal(resBody, &data)
	if err != nil {
		k.logger.Error("fail to unmarshal response", "error", err)
		return sentinelContract, err
	}

	contract.Provider = data.Contract.ProviderPubKey
	contract.Service = data.Contract.Service
	contract.Client = data.Contract.Client
	contract.Delegate = data.Contract.Delegate
	contract.MeterType = data.Contract.MeterType
	contract.UserType = data.Contract.UserType
	contract.Height, _ = strconv.ParseInt(data.Contract.Height, 10, 64)
	contract.Duration, _ = strconv.ParseInt(data.Contract.Duration, 10, 64)
	contract.Rate, _ = strconv.ParseInt(data.Contract.Rate, 10, 64)
	contract.Deposit, _ = cosmos.NewIntFromString(data.Contract.Deposit)
	contract.Paid, _ = cosmos.NewIntFromString(data.Contract.Paid)
	contract.SettlementHeight, _ = strconv.ParseInt(data.Contract.SettlementHeight, 10, 64)

	// todo second call to get nonces is required. this should be done when a user calls
	// for the first time if we don't have a nonce for them, we check on chain

	return sentinelContract, nil
}
