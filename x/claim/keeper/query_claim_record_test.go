package keeper_test

import (
	"testing"

	testkeeper "github.com/arkeonetwork/arkeo/testutil/keeper"
	"github.com/arkeonetwork/arkeo/x/claim/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestClaimRecord(t *testing.T) {
	keeper, ctx := testkeeper.ClaimKeeper(t)

	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String()
	addr2 := "0xDAFEA492D9c6733ae3d56b7Ed1ADB60692c98Bc5" // random eth address

	claimRecords := []types.ClaimRecord{
		{
			Chain:                  types.ARKEO,
			Address:                addr1,
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 100)),
			ActionCompleted:        []bool{false, false},
		},
		{
			Chain:                  types.ETHEREUM,
			Address:                addr2,
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 200)),
			ActionCompleted:        []bool{false, false},
		},
	}
	err := keeper.SetClaimRecords(ctx, claimRecords)
	require.NoError(t, err)

	req := types.QueryClaimRecordRequest{
		Address: addr1,
		Chain:   types.ARKEO,
	}
	resp, err := keeper.ClaimRecord(ctx, &req)
	require.NoError(t, err)
	require.Equal(t, *resp.ClaimRecord, claimRecords[0])

	req = types.QueryClaimRecordRequest{
		Address: addr2,
		Chain:   types.ETHEREUM,
	}
	resp, err = keeper.ClaimRecord(ctx, &req)
	require.NoError(t, err)
	require.Equal(t, *resp.ClaimRecord, claimRecords[1])

	req = types.QueryClaimRecordRequest{
		Address: "invalid address",
		Chain:   types.ETHEREUM,
	}
	resp, _ = keeper.ClaimRecord(ctx, &req)
	require.Equal(t, *resp.ClaimRecord, types.ClaimRecord{})
}
