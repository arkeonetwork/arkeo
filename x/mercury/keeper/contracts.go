package keeper

import (
	"errors"
	"mercury/common"
	"mercury/common/cosmos"
	"mercury/x/mercury/types"
)

func (k KVStore) setContract(ctx cosmos.Context, key string, record types.Contract) {
	store := ctx.KVStore(k.storeKey)
	buf := k.cdc.MustMarshal(&record)
	if buf == nil {
		store.Delete([]byte(key))
	} else {
		store.Set([]byte(key), buf)
	}
}

func (k KVStore) getContract(ctx cosmos.Context, key string, record *types.Contract) (bool, error) {
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

// GetContractIterator iterate contract
func (k KVStore) GetContractIterator(ctx cosmos.Context) cosmos.Iterator {
	return k.getIterator(ctx, prefixContract)
}

// GetContract get the entire Contract metadata struct based on given asset
func (k KVStore) GetContract(ctx cosmos.Context, pubkey common.PubKey, chain common.Chain, client cosmos.AccAddress) (types.Contract, error) {
	record := types.NewContract(pubkey, chain, client)
	_, err := k.getContract(ctx, k.GetKey(ctx, prefixContract, record.Key()), &record)

	return record, err
}

// SetContract save the entire Contract metadata struct to key value store
func (k KVStore) SetContract(ctx cosmos.Context, contract types.Contract) error {
	if contract.ProviderPubKey.IsEmpty() || contract.Chain.IsEmpty() || contract.ClientAddress.Empty() {
		return errors.New("cannot save a contract with an empty provider pubkey, chain, or client address")
	}
	k.setContract(ctx, k.GetKey(ctx, prefixContract, contract.Key()), contract)
	return nil
}

// ContractExists check whether the given contract exist in the data store
func (k KVStore) ContractExists(ctx cosmos.Context, pubkey common.PubKey, chain common.Chain, client cosmos.AccAddress) bool {
	record := types.NewContract(pubkey, chain, client)
	return k.has(ctx, k.GetKey(ctx, prefixContract, record.Key()))
}

func (k KVStore) RemoveContract(ctx cosmos.Context, pubkey common.PubKey, chain common.Chain, client cosmos.AccAddress) {
	record := types.NewContract(pubkey, chain, client)
	k.del(ctx, k.GetKey(ctx, prefixContract, record.Key()))
}
