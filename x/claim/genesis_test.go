package claim_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/arkeonetwork/arkeo/testutil/keeper"
	"github.com/arkeonetwork/arkeo/testutil/utils"
	"github.com/arkeonetwork/arkeo/x/claim"
	"github.com/arkeonetwork/arkeo/x/claim/types"
)

func TestGenesis(t *testing.T) {
	airdropStartTime := time.Now().UTC()
	claimParams := types.Params{
		AirdropStartTime:   airdropStartTime,
		DurationUntilDecay: types.DefaultDurationUntilDecay,
		DurationOfDecay:    types.DefaultDurationOfDecay,
		ClaimDenom:         types.DefaultClaimDenom,
	}
	genesisState := types.GenesisState{
		Params: claimParams,
	}

	testKeeepers, ctx := keepertest.CreateTestClaimKeepers(t)
	claim.InitGenesis(ctx, testKeeepers.ClaimKeeper, genesisState)
	got := claim.ExportGenesis(ctx, testKeeepers.ClaimKeeper)
	require.NotNil(t, got)

	addr1 := utils.GetRandomArkeoAddress().String()
	addr2 := utils.GetRandomArkeoAddress().String()
	ethAddr1 := "0xDAFEA492D9c6733ae3d56b7Ed1ADB60692c98Bc5" // random eth address

	testGenesis := types.GenesisState{
		ModuleAccountBalance: sdk.NewInt64Coin(types.DefaultClaimDenom, 750000000),
		Params: types.Params{
			AirdropStartTime:   airdropStartTime,
			DurationUntilDecay: types.DefaultDurationUntilDecay,
			DurationOfDecay:    types.DefaultDurationOfDecay,
			ClaimDenom:         types.DefaultClaimDenom,
		},
		ClaimRecords: []types.ClaimRecord{
			{
				Address:        addr1,
				Chain:          types.ARKEO,
				AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 1000000000),
				AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 1000000000),
				AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 1000000000),
			},
			{
				Address:        addr2,
				Chain:          types.ARKEO,
				AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 1000000000),
				AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 1500000000),
				AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 1500000000),
			},
			{
				Address:        ethAddr1,
				Chain:          types.ETHEREUM,
				AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 2000000000),
				AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 2000000000),
				AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 2000000000),
			},
		},
	}
	claim.InitGenesis(ctx, testKeeepers.ClaimKeeper, testGenesis)

	claimRecord, err := testKeeepers.ClaimKeeper.GetClaimRecord(ctx, addr2, types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, claimRecord, types.ClaimRecord{
		Address:        addr2,
		Chain:          types.ARKEO,
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 1000000000),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 1500000000),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 1500000000),
	})

	claimableAmount, err := testKeeepers.ClaimKeeper.GetClaimableAmountForAction(ctx, addr2, types.ACTION_VOTE, types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, claimableAmount, sdk.NewInt64Coin(types.DefaultClaimDenom, 1500000000))

	genesisExported := claim.ExportGenesis(ctx, testKeeepers.ClaimKeeper)
	require.Equal(t, genesisExported.Params, testGenesis.Params)
	require.ElementsMatch(t, genesisExported.ClaimRecords, testGenesis.ClaimRecords)
}
