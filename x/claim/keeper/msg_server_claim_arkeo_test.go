package keeper_test

import (
	"testing"

	"github.com/arkeonetwork/arkeo/testutil/utils"
	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestClaimArkeo(t *testing.T) {
	msgServer, keepers, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

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

	// mint coins to module account
	err = keepers.BankKeeper.MintCoins(sdkCtx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 10000)))
	require.NoError(t, err)

	// get balance of arkeo address before claim
	balanceBefore := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)

	claimMessage := types.MsgClaimArkeo{
		Creator: addrArkeo.String(),
	}
	_, err = msgServer.ClaimArkeo(ctx, &claimMessage)
	require.NoError(t, err)

	// check if claimrecord is updated
	claimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, addrArkeo.String(), types.ARKEO)
	require.NoError(t, err)
	require.True(t, !claimRecord.IsEmpty())

	require.Equal(t, claimRecord.Address, addrArkeo.String())
	require.Equal(t, claimRecord.Chain, types.ARKEO)
	require.True(t, claimRecord.AmountClaim.IsZero()) // nothing to claim for claim action
	require.Equal(t, claimRecord.AmountVote, sdk.NewInt64Coin(types.DefaultClaimDenom, 100))
	require.Equal(t, claimRecord.AmountDelegate, sdk.NewInt64Coin(types.DefaultClaimDenom, 100))

	// confirm balance increased by expected amount.
	balanceAfter := keepers.BankKeeper.GetBalance(sdkCtx, addrArkeo, types.DefaultClaimDenom)
	require.Equal(t, balanceAfter.Sub(balanceBefore), sdk.NewInt64Coin(types.DefaultClaimDenom, 100))

	// attempt to claim again to ensure it fails.
	_, err = msgServer.ClaimArkeo(ctx, &claimMessage)
	require.Error(t, err)

	// ensure claim Arkeo fails from address with no claim record
	addrArkeo2 := utils.GetRandomArkeoAddress()
	claimMessage2 := types.MsgClaimArkeo{
		Creator: addrArkeo2.String(),
	}
	_, err = msgServer.ClaimArkeo(ctx, &claimMessage2)
	require.Error(t, err)
}
