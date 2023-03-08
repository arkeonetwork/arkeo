package keeper

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/stretchr/testify/require"
)

func TestContract(t *testing.T) {
	ctx, k := SetupKeeper(t)
	require.Error(t, k.SetContract(ctx, types.Contract{}))

	contract := types.NewContract(types.GetRandomPubKey(), common.BTCChain, types.GetRandomPubKey())
	contract.Id = 1
	err := k.SetContract(ctx, contract)
	require.NoError(t, err)

	contract, err = k.GetContract(ctx, contract.Id)
	require.NoError(t, err)
	require.Equal(t, contract.Chain, common.BTCChain)
	require.True(t, k.ContractExists(ctx, contract.Id))
	require.False(t, k.ContractExists(ctx, contract.Id+1))

	k.RemoveContract(ctx, contract.Id)
	require.False(t, k.ContractExists(ctx, contract.Id))
}

func TestContractExpirationSet(t *testing.T) {
	ctx, k := SetupKeeper(t)
	set := types.ContractExpirationSet{}
	require.Error(t, k.SetContractExpirationSet(ctx, set)) // empty asset should error

	set.Height = 100
	set.ContractSet = &types.ContractSet{}
	require.NoError(t, k.SetContractExpirationSet(ctx, set))

	contractId := uint64(100)
	set.ContractSet.ContractIds = append(set.ContractSet.ContractIds, contractId)

	require.NoError(t, k.SetContractExpirationSet(ctx, set))
	set, err := k.GetContractExpirationSet(ctx, set.Height)
	require.NoError(t, err)
	require.Equal(t, set.Height, int64(100))
	require.Len(t, set.ContractSet.ContractIds, 1)

	k.RemoveContractExpirationSet(ctx, 100)
	set, err = k.GetContractExpirationSet(ctx, set.Height)
	require.NoError(t, err)
	require.Nil(t, set.ContractSet)
}
