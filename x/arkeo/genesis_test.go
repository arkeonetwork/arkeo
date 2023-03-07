package arkeo_test

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	keepertest "github.com/arkeonetwork/arkeo/testutil/keeper"

	"github.com/arkeonetwork/arkeo/testutil/nullify"
	"github.com/arkeonetwork/arkeo/x/arkeo"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	ctx, k := keepertest.ArkeoKeeper(t)
	arkeo.InitGenesis(ctx, k, genesisState)
	got := arkeo.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}

func TestGenesisWithContracts(t *testing.T) {
	ctx, k := keepertest.ArkeoKeeper(t)

	providerPubkey := types.GetRandomPubKey()
	user1PubKey := types.GetRandomPubKey()
	user2PubKey := types.GetRandomPubKey()

	// create provider
	provider := types.NewProvider(providerPubkey, common.BTCChain)
	provider.Status = types.ProviderStatus_ONLINE
	provider.LastUpdate = 100
	k.SetProvider(ctx, provider)

	// create contracts
	contracts := []types.Contract{
		{
			Provider: providerPubkey,
			Chain:    common.BTCChain,
			Client:   user1PubKey,
			Duration: 100,
			Rate:     100,
			Id:       0,
			Deposit:  cosmos.NewInt(500),
			Paid:     cosmos.ZeroInt(),
		},
		{
			Provider: providerPubkey,
			Chain:    common.ETHChain,
			Client:   user1PubKey,
			Duration: 100,
			Rate:     100,
			Id:       1,
			Deposit:  cosmos.NewInt(500),
			Paid:     cosmos.ZeroInt(),
		},
		{
			Provider: providerPubkey,
			Chain:    common.BTCChain,
			Client:   user2PubKey,
			Duration: 150,
			Rate:     100,
			Id:       2,
			Deposit:  cosmos.NewInt(200),
			Paid:     cosmos.ZeroInt(),
		},
	}

	for _, contract := range contracts {
		err := k.SetContract(ctx, contract)
		require.NoError(t, err)
	}

	// create user contract sets
	user1ContractSet := types.UserContractSet{
		User: user1PubKey,
		ContractSet: &types.ContractSet{
			ContractIds: []uint64{0, 1},
		},
	}
	user2ContractSet := types.UserContractSet{
		User: user2PubKey,
		ContractSet: &types.ContractSet{
			ContractIds: []uint64{2},
		},
	}

	err := k.SetUserContractSet(ctx, user1ContractSet)
	require.NoError(t, err)
	err = k.SetUserContractSet(ctx, user2ContractSet)
	require.NoError(t, err)

	// create contract expiration sets.
	contractExpirationSet1 := types.ContractExpirationSet{
		Height: 100,
		ContractSet: &types.ContractSet{
			ContractIds: []uint64{0, 1},
		},
	}

	contractExpirationSet2 := types.ContractExpirationSet{
		Height: 150,
		ContractSet: &types.ContractSet{
			ContractIds: []uint64{2},
		},
	}
	err = k.SetContractExpirationSet(ctx, contractExpirationSet1)
	require.NoError(t, err)

	err = k.SetContractExpirationSet(ctx, contractExpirationSet2)
	require.NoError(t, err)

	exportedGenesis := arkeo.ExportGenesis(ctx, k)
	require.NotNil(t, exportedGenesis)

	// check if the exported genesis is the same as the state we just set.
	require.Equal(t, exportedGenesis.Contracts, contracts)
	require.Equal(t, exportedGenesis.UserContractSets, []types.UserContractSet{user1ContractSet, user2ContractSet})
	require.Equal(t, exportedGenesis.ContractExpirationSets, []types.ContractExpirationSet{contractExpirationSet1, contractExpirationSet2})
}
