package keeper

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/stretchr/testify/require"
)

func TestNonces(t *testing.T) {
	ctx, k := SetupKeeper(t)

	// invalid contract id should error
	nonceSet := int64(20)
	require.Error(t, k.SetNonce(ctx, types.GetRandomPubKey(), 20, nonceSet))

	contract := types.NewContract(types.GetRandomPubKey(), common.BTCService, types.GetRandomPubKey())
	contract.Id = 10
	err := k.SetContract(ctx, contract)
	require.NoError(t, err)
	require.Error(t, k.SetNonce(ctx, common.EmptyPubKey, contract.Id, nonceSet)) // empty pub key should error

	spenderPubKey := types.GetRandomPubKey()
	err = k.SetNonce(ctx, spenderPubKey, contract.Id, nonceSet)
	require.NoError(t, err)
	nonceReturned, err := k.GetNonce(ctx, spenderPubKey, contract.Id)
	require.NoError(t, err)
	require.Equal(t, nonceReturned, nonceSet)
	require.True(t, k.NonceExists(ctx, spenderPubKey, contract.Id))
	require.False(t, k.NonceExists(ctx, spenderPubKey, contract.Id+1))

}
