package keeper

import (
	"errors"
	"fmt"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func (k KVStore) setNonce(ctx cosmos.Context, key string, nonce int64) {
	store := ctx.KVStore(k.storeKey)
	protoNonce := types.ProtoInt64{Value: nonce}
	buf := k.cdc.MustMarshal(&protoNonce)
	if buf == nil {
		store.Delete([]byte(key))
	} else {
		store.Set([]byte(key), buf)
	}
}

func (k KVStore) getNonce(ctx cosmos.Context, key string) (int64, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(key)) {
		return 0, nil
	}

	bz := store.Get([]byte(key))
	nonce := types.ProtoInt64{}
	if err := k.cdc.Unmarshal(bz, &nonce); err != nil {
		return 0, err
	}
	return nonce.Value, nil
}

func (k KVStore) NonceExists(ctx cosmos.Context, spenderPubKey common.PubKey, contractId uint64) bool {
	return k.has(ctx, k.GetKey(ctx, prefixNonce, GetNonceKey(spenderPubKey, contractId)))
}

func GetNonceKey(spenderPubKey common.PubKey, contractId uint64) string {
	return fmt.Sprintf("%s/%d", spenderPubKey, contractId)
}

func (k KVStore) SetNonce(ctx cosmos.Context, spenderPubKey common.PubKey, contractId uint64, nonce int64) error {
	if spenderPubKey.IsEmpty() {
		return errors.New("cannot save a nonce with an empty pubkey")
	}
	if !k.ContractExists(ctx, contractId) {
		return errors.New("cannot save a nonce for a non existing contract")
	}
	k.setNonce(ctx, k.GetKey(ctx, prefixNonce, GetNonceKey(spenderPubKey, contractId)), nonce)
	return nil
}

func (k KVStore) GetNonce(ctx cosmos.Context, spenderPubKey common.PubKey, contractId uint64) (int64, error) {
	key := GetNonceKey(spenderPubKey, contractId)
	return k.getNonce(ctx, k.GetKey(ctx, prefixNonce, key))
}
