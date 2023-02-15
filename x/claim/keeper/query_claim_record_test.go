package keeper_test

import (
	"testing"

	testkeeper "github.com/arkeonetwork/arkeo/testutil/keeper"
	"github.com/arkeonetwork/arkeo/testutil/utils"
	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestClaimRecord(t *testing.T) {
	keepers, ctx := testkeeper.CreateTestClaimKeepers(t)

	addr1 := utils.GetRandomArkeoAddress().String()
	addr2 := "0xDAFEA492D9c6733ae3d56b7Ed1ADB60692c98Bc5" // random eth address

	claimRecords := []types.ClaimRecord{
		{
			Chain:          types.ARKEO,
			Address:        addr1,
			AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
			AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
			AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		},
		{
			Chain:          types.ETHEREUM,
			Address:        addr2,
			AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
			AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
			AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
		},
	}
	err := keepers.ClaimKeeper.SetClaimRecords(ctx, claimRecords)
	require.NoError(t, err)

	req := types.QueryClaimRecordRequest{
		Address: addr1,
		Chain:   types.ARKEO,
	}
	resp, err := keepers.ClaimKeeper.ClaimRecord(ctx, &req)
	require.NoError(t, err)
	require.Equal(t, *resp.ClaimRecord, claimRecords[0])

	req = types.QueryClaimRecordRequest{
		Address: addr2,
		Chain:   types.ETHEREUM,
	}
	resp, err = keepers.ClaimKeeper.ClaimRecord(ctx, &req)
	require.NoError(t, err)
	require.Equal(t, *resp.ClaimRecord, claimRecords[1])

	req = types.QueryClaimRecordRequest{
		Address: "invalid address",
		Chain:   types.ETHEREUM,
	}
	resp, _ = keepers.ClaimKeeper.ClaimRecord(ctx, &req)
	require.Equal(t, *resp.ClaimRecord, types.ClaimRecord{})
}
