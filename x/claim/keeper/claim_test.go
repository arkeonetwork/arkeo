package keeper_test

import (
	"testing"

	testkeeper "github.com/arkeonetwork/arkeo/testutil/keeper"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGetClaimRecordForArkeo(t *testing.T) {
	keeper, ctx := testkeeper.ClaimKeeper(t)

	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String()
	addr2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String()
	addr3 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String()

	claimRecords := []types.ClaimRecord{
		{
			Chain:                  types.ARKEO,
			Address:                addr1,
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 100)),
			ActionCompleted:        []bool{false, false},
		},
		{
			Chain:                  types.ARKEO,
			Address:                addr2,
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 200)),
			ActionCompleted:        []bool{false, false},
		},
	}
	err := keeper.SetClaimRecords(ctx, claimRecords)
	require.NoError(t, err)

	coins1, err := keeper.GetUserTotalClaimable(ctx, addr1, types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, coins1, claimRecords[0].InitialClaimableAmount, coins1.String())

	coins2, err := keeper.GetUserTotalClaimable(ctx, addr2, types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, coins2, claimRecords[1].InitialClaimableAmount)

	coins3, err := keeper.GetUserTotalClaimable(ctx, addr3, types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, coins3, sdk.Coins{})

	// get rewards amount per action
	coins4, err := keeper.GetClaimableAmountForAction(ctx, addr1, types.ACTION_VOTE, types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, coins4.String(), sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 50)).String())

	// get completed activities
	claimRecord, err := keeper.GetClaimRecord(ctx, addr1, types.ARKEO)
	require.NoError(t, err)
	for i := range types.Action_name {
		require.False(t, claimRecord.ActionCompleted[i])
	}
}

func TestGetClaimRecordForMutlipleChains(t *testing.T) {
	keeper, ctx := testkeeper.ClaimKeeper(t)

	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String()
	addr2 := "0xDAFEA492D9c6733ae3d56b7Ed1ADB60692c98Bc5" // random eth address
	// addr3 := "thor18u55kxfudpy9q7mvhxzrh4xntjyukx420lt5fg" // random thorchain address

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
		// {
		// 	Chain:                  types.THORCHAIN,
		// 	Address:                addr3,
		// 	InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 500)),
		// 	ActionCompleted:        []bool{false, false},
		// },
	}
	err := keeper.SetClaimRecords(ctx, claimRecords)
	require.NoError(t, err)

	coins1, err := keeper.GetUserTotalClaimable(ctx, addr1, types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, coins1, claimRecords[0].InitialClaimableAmount, coins1.String())

	// user 1 should have no eth claim with an arkeo addy nor thor claims
	coins1, err = keeper.GetUserTotalClaimable(ctx, addr1, types.ETHEREUM)
	require.NoError(t, err)
	require.Equal(t, coins1, sdk.Coins{})
	coins1, err = keeper.GetUserTotalClaimable(ctx, addr1, types.THORCHAIN)
	require.NoError(t, err)
	require.Equal(t, coins1, sdk.Coins{})

	// user 2 should have no arkeo claim nor thor claims, only eth
	coins2, err := keeper.GetUserTotalClaimable(ctx, addr2, types.ETHEREUM)
	require.NoError(t, err)
	require.Equal(t, coins2, claimRecords[1].InitialClaimableAmount)

	coins2, err = keeper.GetUserTotalClaimable(ctx, addr2, types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, coins2, sdk.Coins{})
	coins2, err = keeper.GetUserTotalClaimable(ctx, addr2, types.THORCHAIN)
	require.NoError(t, err)
	require.Equal(t, coins2, sdk.Coins{})

	// user 3 should have no arkeo claim nor eth claims, only thor
	// coins3, err := keeper.GetUserTotalClaimable(ctx, addr3, types.ARKEO)
	// require.NoError(t, err)
	// require.Equal(t, coins3, sdk.Coins{})

	// coins3, err = keeper.GetUserTotalClaimable(ctx, addr3, types.ETHEREUM)
	// require.NoError(t, err)
	// require.Equal(t, coins3, sdk.Coins{})
	// // coins3, err = keeper.GetUserTotalClaimable(ctx, addr3, types.THORCHAIN)
	// // require.NoError(t, err)
	// // require.Equal(t, coins3, claimRecords[2].InitialClaimableAmount)
}

func TestSetClaimRecord(t *testing.T) {
	keeper, ctx := testkeeper.ClaimKeeper(t)

	// confirm setting a claim record with a bad eth address fails
	addr1Invalid := "0xDAFEA492D9c6733ae3d56b7Ed1ADB60692c98"  // random invalid eth address
	addr1Valid := "0xDAFEA492D9c6733ae3d56b7Ed1ADB60692c98Bc5" // random valid eth address
	claimRecord := types.ClaimRecord{
		Chain:                  types.ETHEREUM,
		Address:                addr1Invalid,
		InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 100)),
		ActionCompleted:        []bool{false, false},
	}
	err := keeper.SetClaimRecord(ctx, claimRecord)
	require.Error(t, err)

	claimRecord.Address = addr1Valid
	err = keeper.SetClaimRecord(ctx, claimRecord)
	require.NoError(t, err)

	// confirm setting a claim record with a bad arkeo address fails
	addr2Invalid := "0xDAFEA492D9c6733ae3d56b7Ed1ADB60692c98Bc5" // random eth address (should fail)
	addr2Valid := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()).String()
	claimRecord = types.ClaimRecord{
		Chain:                  types.ARKEO,
		Address:                addr2Invalid,
		InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 100)),
		ActionCompleted:        []bool{false, false},
	}
	err = keeper.SetClaimRecord(ctx, claimRecord)
	require.Error(t, err)

	claimRecord.Address = addr2Valid
	err = keeper.SetClaimRecord(ctx, claimRecord)
	require.NoError(t, err)
}
