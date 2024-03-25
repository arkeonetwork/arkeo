package keeper_test

import (
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/testutil/utils"
	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClaimThorchainArkeo(t *testing.T) {
	msgServer, keepers, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cosmos.GetConfig().SetBech32PrefixForAccount("tarkeo", "tarkeopub")

	addrArkeo := utils.GetRandomArkeoAddress()
	claimRecord := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        addrArkeo.String(),
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
	}
	err := keepers.ClaimKeeper.SetClaimRecord(sdkCtx, claimRecord)
	require.NoError(t, err)

	thorClaimAddress := "tarkeo1dllfyp57l4xj5umqfcqy6c2l3xfk0qk6zpc3t7"

	thorClaimRecord := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        thorClaimAddress, // arkeo address derived from sender of thorchain tx "FA2768AEB52AE0A378372B48B10C5B374B25E8B2005C702AAD441B813ED2F174"
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
	}
	err = keepers.ClaimKeeper.SetClaimRecord(sdkCtx, thorClaimRecord)
	require.NoError(t, err)

	// mint coins to module account
	err = keepers.BankKeeper.MintCoins(sdkCtx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 10000)))
	require.NoError(t, err)

	// get balance of arkeo address before claim
	balanceBefore := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)

	claimMessage := types.MsgClaimArkeo{
		Creator: addrArkeo,
		ThorTx:  "FA2768AEB52AE0A378372B48B10C5B374B25E8B2005C702AAD441B813ED2F174",
	}
	_, err = msgServer.ClaimArkeo(ctx, &claimMessage)
	require.NoError(t, err)

	// check if claimrecord is updated
	thorClaimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, thorClaimAddress, types.ARKEO)
	require.NoError(t, err)
	require.True(t, thorClaimRecord.IsEmpty())

	claimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, addrArkeo.String(), types.ARKEO)
	require.NoError(t, err)
	require.True(t, !claimRecord.IsEmpty())

	require.Equal(t, claimRecord.Address, addrArkeo.String())
	require.Equal(t, claimRecord.Chain, types.ARKEO)
	require.True(t, claimRecord.AmountClaim.IsZero()) // nothing to claim for claim action
	require.Equal(t, claimRecord.AmountVote, sdk.NewInt64Coin(types.DefaultClaimDenom, 200))
	require.Equal(t, claimRecord.AmountDelegate, sdk.NewInt64Coin(types.DefaultClaimDenom, 200))

	// confirm balance increased by expected amount.
	balanceAfter := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)
	require.Equal(t, balanceAfter.Sub(balanceBefore), sdk.NewInt64Coin(types.DefaultClaimDenom, 200))

	// attempt to claim again to ensure it fails.
	_, err = msgServer.ClaimArkeo(ctx, &claimMessage)
	require.ErrorIs(t, err, types.ErrNoClaimableAmount)

	// ensure claim Arkeo fails from address with no claim record
	addrArkeo2 := utils.GetRandomArkeoAddress()
	claimMessage2 := types.MsgClaimArkeo{
		Creator: addrArkeo2,
	}
	_, err = msgServer.ClaimArkeo(ctx, &claimMessage2)
	require.ErrorIs(t, err, types.ErrNoClaimableAmount)
}

func TestClaimThorchainEth(t *testing.T) {
	msgServer, keepers, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	cosmos.GetConfig().SetBech32PrefixForAccount("tarkeo", "tarkeopub")

	// create valid eth claimrecords
	addrArkeo := utils.GetRandomArkeoAddress()
	addrEth, sigString, err := generateSignedEthClaim(addrArkeo.String(), "300")
	require.NoError(t, err)

	claimRecord := types.ClaimRecord{
		Chain:          types.ETHEREUM,
		Address:        addrEth,
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
	}
	err = keepers.ClaimKeeper.SetClaimRecord(sdkCtx, claimRecord)
	require.NoError(t, err)

	thorClaimAddress := "tarkeo1dllfyp57l4xj5umqfcqy6c2l3xfk0qk6zpc3t7"
	thorClaimRecord := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        thorClaimAddress, // arkeo address derived from sender of thorchain tx "FA2768AEB52AE0A378372B48B10C5B374B25E8B2005C702AAD441B813ED2F174"
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 500),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 500),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 500),
	}
	err = keepers.ClaimKeeper.SetClaimRecord(sdkCtx, thorClaimRecord)
	require.NoError(t, err)

	// mint coins to module account
	err = keepers.BankKeeper.MintCoins(sdkCtx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 10000)))
	require.NoError(t, err)

	// get balance of arkeo address before claim
	balanceBefore := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)

	claimMessage := types.MsgClaimEth{
		Creator:    addrArkeo,
		EthAddress: addrEth,
		Signature:  sigString,
		ThorTx:     "FA2768AEB52AE0A378372B48B10C5B374B25E8B2005C702AAD441B813ED2F174",
	}
	_, err = msgServer.ClaimEth(ctx, &claimMessage)
	require.NoError(t, err)

	// check if claimrecord is updated
	claimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, addrEth, types.ETHEREUM)
	require.NoError(t, err)
	require.True(t, claimRecord.IsEmpty())

	// confirm we have a claimrecord for arkeo
	claimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, addrArkeo.String(), types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, claimRecord.Address, addrArkeo.String())
	require.Equal(t, claimRecord.Chain, types.ARKEO)
	require.True(t, claimRecord.AmountClaim.IsZero()) // nothing to claim for claim action
	require.Equal(t, claimRecord.AmountVote, sdk.NewInt64Coin(types.DefaultClaimDenom, 600))
	require.Equal(t, claimRecord.AmountDelegate, sdk.NewInt64Coin(types.DefaultClaimDenom, 600))

	// confirm balance increased by expected amount.
	balanceAfter := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)
	require.Equal(t, balanceAfter.Sub(balanceBefore), sdk.NewInt64Coin(types.DefaultClaimDenom, 600))

	// attempt to claim again to ensure it fails.
	_, err = msgServer.ClaimEth(ctx, &claimMessage)
	require.Error(t, err)

	// attempt to claim from arkeo should also fail!
	_, err = msgServer.ClaimArkeo(ctx, &types.MsgClaimArkeo{Creator: addrArkeo})
	require.Error(t, err)
}
