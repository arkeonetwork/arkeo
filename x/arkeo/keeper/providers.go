package keeper

import (
	"errors"

	"github.com/ArkeoNetwork/arkeo-protocol/common"
	"github.com/ArkeoNetwork/arkeo-protocol/common/cosmos"
	"github.com/ArkeoNetwork/arkeo-protocol/x/arkeo/types"
)

func (k KVStore) setProvider(ctx cosmos.Context, key string, record types.Provider) {
	store := ctx.KVStore(k.storeKey)
	buf := k.cdc.MustMarshal(&record)
	if buf == nil {
		store.Delete([]byte(key))
	} else {
		store.Set([]byte(key), buf)
	}
}

func (k KVStore) getProvider(ctx cosmos.Context, key string, record *types.Provider) (bool, error) {
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

// GetProviderIterator iterate providers
func (k KVStore) GetProviderIterator(ctx cosmos.Context) cosmos.Iterator {
	return k.getIterator(ctx, prefixProvider)
}

// GetProvider get the entire Provider metadata struct based on given asset
func (k KVStore) GetProvider(ctx cosmos.Context, pubkey common.PubKey, chain common.Chain) (types.Provider, error) {
	record := types.NewProvider(pubkey, chain)
	_, err := k.getProvider(ctx, k.GetKey(ctx, prefixProvider, record.Key()), &record)

	return record, err
}

// SetProvider save the entire Provider metadata struct to key value store
func (k KVStore) SetProvider(ctx cosmos.Context, provider types.Provider) error {
	if provider.PubKey.IsEmpty() || provider.Chain.IsEmpty() {
		return errors.New("cannot save a provider with an empty pubkey or chain")
	}
	k.setProvider(ctx, k.GetKey(ctx, prefixProvider, provider.Key()), provider)
	return nil
}

// ProviderExists check whether the given provider exist in the data store
func (k KVStore) ProviderExists(ctx cosmos.Context, pubkey common.PubKey, chain common.Chain) bool {
	record := types.NewProvider(pubkey, chain)
	return k.has(ctx, k.GetKey(ctx, prefixProvider, record.Key()))
}

func (k KVStore) RemoveProvider(ctx cosmos.Context, pubkey common.PubKey, chain common.Chain) {
	record := types.NewProvider(pubkey, chain)
	k.del(ctx, k.GetKey(ctx, prefixProvider, record.Key()))
}
