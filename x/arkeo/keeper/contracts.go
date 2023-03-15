package keeper

import (
	"errors"
	"strconv"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	gogotypes "github.com/gogo/protobuf/types"
)

func (k KVStore) setContract(ctx cosmos.Context, contract types.Contract) {
	store := ctx.KVStore(k.storeKey)
	key := k.GetContractKey(ctx, contract.Id)
	buf := k.cdc.MustMarshal(&contract)
	if buf == nil {
		store.Delete([]byte(key))
	} else {
		store.Set([]byte(key), buf)
	}
}

func (k KVStore) getContract(ctx cosmos.Context, id uint64, contract *types.Contract) (bool, error) {
	store := ctx.KVStore(k.storeKey)
	key := k.GetContractKey(ctx, id)
	if !store.Has([]byte(key)) {
		return false, nil
	}

	bz := store.Get([]byte(key))
	if err := k.cdc.Unmarshal(bz, contract); err != nil {
		return true, err
	}
	return true, nil
}

// GetContractIterator iterate contract
func (k KVStore) GetContractIterator(ctx cosmos.Context) cosmos.Iterator {
	return k.getIterator(ctx, prefixContract)
}

// GetContract get the entire Contract metadata struct based on given asset
func (k KVStore) GetContract(ctx cosmos.Context, id uint64) (types.Contract, error) {
	contract := types.Contract{}
	_, err := k.getContract(ctx, id, &contract)
	return contract, err
}

// SetContract save the entire Contract metadata struct to key value store
func (k KVStore) SetContract(ctx cosmos.Context, contract types.Contract) error {
	if contract.Provider.IsEmpty() || contract.Service.IsEmpty() || contract.Client.IsEmpty() {
		return errors.New("cannot save a contract with an empty provider pubkey, service, or client address")
	}
	k.setContract(ctx, contract)
	return nil
}

// ContractExists check whether the given contract exist in the data store
func (k KVStore) ContractExists(ctx cosmos.Context, id uint64) bool {
	return k.has(ctx, k.GetContractKey(ctx, id))
}

func (k KVStore) RemoveContract(ctx cosmos.Context, id uint64) {
	k.del(ctx, k.GetContractKey(ctx, id))
}

func (k KVStore) setContractExpirationSet(ctx cosmos.Context, key string, record types.ContractExpirationSet) {
	store := ctx.KVStore(k.storeKey)
	buf := k.cdc.MustMarshal(&record)
	if len(record.ContractSet.ContractIds) == 0 {
		store.Delete([]byte(key))
	} else {
		store.Set([]byte(key), buf)
	}
}

func (k KVStore) getContractExpirationSet(ctx cosmos.Context, key string, record *types.ContractExpirationSet) (bool, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(key)) {
		return false, nil
	}

	bz := store.Get([]byte(key))
	if err := k.cdc.Unmarshal(bz, record); err != nil {
		return true, err
	}
	return true, nil
}

func (k KVStore) getContractExpirationSetKey(ctx cosmos.Context, height int64) string {
	return k.GetKey(ctx, prefixContractExpirationSet, strconv.FormatInt(height, 10))
}

func (k KVStore) getUserContractSetKey(ctx cosmos.Context, userPubKey common.PubKey) string {
	return k.GetKey(ctx, prefixUserContractSet, userPubKey.String())
}

func (k KVStore) GetContractKey(ctx cosmos.Context, id uint64) string {
	return k.GetKey(ctx, prefixContract, strconv.FormatUint(id, 10))
}

// GetContractExpirationSetIterator iterate contract expiration sets
func (k KVStore) GetContractExpirationSetIterator(ctx cosmos.Context) cosmos.Iterator {
	return k.getIterator(ctx, prefixContractExpirationSet)
}

// GetContractExpirationSet get a series of contract expirations
func (k KVStore) GetContractExpirationSet(ctx cosmos.Context, height int64) (types.ContractExpirationSet, error) {
	record := types.ContractExpirationSet{
		Height: height,
	}
	_, err := k.getContractExpirationSet(ctx, k.getContractExpirationSetKey(ctx, height), &record)
	return record, err
}

// SetContractExpirationSet save the series of Contract Expirations
func (k KVStore) SetContractExpirationSet(ctx cosmos.Context, record types.ContractExpirationSet) error {
	if record.Height <= 0 {
		return errors.New("cannot save a contract expiration set with an invalid height (less than or equal to zero)")
	}
	k.setContractExpirationSet(ctx, k.getContractExpirationSetKey(ctx, record.Height), record)
	return nil
}

func (k KVStore) RemoveContractExpirationSet(ctx cosmos.Context, height int64) {
	k.del(ctx, k.GetKey(ctx, prefixContractExpirationSet, strconv.FormatInt(height, 10)))
}

func (kvStore KVStore) GetAndIncrementNextContractId(ctx cosmos.Context) uint64 {
	contractId := kvStore.GetNextContractId(ctx)
	kvStore.SetNextContractId(ctx, contractId+1) // increment and set
	return contractId
}

func (kvStore KVStore) GetNextContractId(ctx cosmos.Context) uint64 {
	var contractId uint64
	store := ctx.KVStore(kvStore.storeKey)

	bz := store.Get([]byte(prefixContractNextId))
	if bz == nil {
		// initialize the id number, this assignment is not necessary, but it's here for clarity
		contractId = 0
	} else {
		val := gogotypes.UInt64Value{}
		kvStore.cdc.MustUnmarshal(bz, &val)
		contractId = val.GetValue()
	}
	return contractId
}

func (k KVStore) SetNextContractId(ctx cosmos.Context, contractId uint64) {
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: contractId})
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(prefixContractNextId), bz)
}

func (k KVStore) SetUserContractSet(ctx cosmos.Context, contractSet types.UserContractSet) error {
	if len(contractSet.User) == 0 {
		return errors.New("cannot save a user contract set with a blank user")
	}
	k.setUserContractSet(ctx, k.getUserContractSetKey(ctx, contractSet.User), contractSet)
	return nil
}

func (k KVStore) setUserContractSet(ctx cosmos.Context, key string, set types.UserContractSet) {
	store := ctx.KVStore(k.storeKey)
	buf := k.cdc.MustMarshal(&set)
	if len(set.ContractSet.ContractIds) == 0 {
		store.Delete([]byte(key))
	} else {
		store.Set([]byte(key), buf)
	}
}

func (k KVStore) GetUserContractSet(ctx cosmos.Context, user common.PubKey) (types.UserContractSet, error) {
	record := types.UserContractSet{
		User: user,
	}
	_, err := k.getUserContractSet(ctx, k.getUserContractSetKey(ctx, user), &record)
	return record, err
}

func (k KVStore) getUserContractSet(ctx cosmos.Context, key string, contractSet *types.UserContractSet) (bool, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(key)) {
		return false, nil
	}

	bz := store.Get([]byte(key))
	if err := k.cdc.Unmarshal(bz, contractSet); err != nil {
		return true, err
	}
	return true, nil
}

func (k KVStore) GetActiveContractForUser(ctx cosmos.Context, user common.PubKey, provider common.PubKey, service common.Service) (types.Contract, error) {
	contractSet, err := k.GetUserContractSet(ctx, user)
	if err != nil {
		return types.Contract{}, err
	}

	if contractSet.ContractSet == nil || len(contractSet.ContractSet.ContractIds) == 0 {
		return types.Contract{}, nil
	}

	for _, contractId := range contractSet.ContractSet.ContractIds {
		contract, err := k.GetContract(ctx, contractId)
		if err != nil {
			return types.Contract{}, err
		}
		if contract.Provider.Equals(provider) && contract.Service.Equals(service) && contract.IsOpen(ctx.BlockHeight()) {
			return contract, nil
		}
	}

	return types.Contract{}, nil
}

// RemoveFromUserContractSet remove a contract from a user's contract set and saves the updated set to the store
func (k KVStore) RemoveFromUserContractSet(ctx cosmos.Context, user common.PubKey, contractId uint64) error {
	contractSet, err := k.GetUserContractSet(ctx, user)
	if err != nil {
		return err
	}

	err = contractSet.RemoveContractFromSet(contractId)
	if err != nil {
		return err
	}

	if len(contractSet.ContractSet.ContractIds) == 0 {
		// set is empty, remove key.
		k.del(ctx, k.getUserContractSetKey(ctx, user))
		return nil
	}

	return k.SetUserContractSet(ctx, contractSet)
}

func (k KVStore) GetUserContractSetIterator(ctx cosmos.Context) cosmos.Iterator {
	return k.getIterator(ctx, prefixUserContractSet)
}
