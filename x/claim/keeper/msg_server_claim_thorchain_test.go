package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/arkeonetwork/arkeo/testutil/utils"
	"github.com/arkeonetwork/arkeo/x/claim/types"
)

func TestClaimThorchainMainnetAddress(t *testing.T) {
	msgServer, keepers, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("arkeo", "arkeopub")

	arkeoServerAddress, err := sdk.AccAddressFromBech32("arkeo1darsrv7x7kr26geefyn4cljn3m9y5yg5c84xgv")
	require.NoError(t, err)

	fromAddr := utils.GetRandomArkeoAddress()
	toAddr := utils.GetRandomArkeoAddress()

	claimRecordFrom := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        fromAddr.String(),
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
	}
	claimRecordTo := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        toAddr.String(),
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
	}
	err = keepers.ClaimKeeper.SetClaimRecord(sdkCtx, claimRecordFrom)
	require.NoError(t, err)
	err = keepers.ClaimKeeper.SetClaimRecord(sdkCtx, claimRecordTo)
	require.NoError(t, err)

	// mint coins to module account
	err = keepers.BankKeeper.MintCoins(sdkCtx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 10000)))
	require.NoError(t, err)

	invalidClaimMessage := types.MsgClaimThorchain{
		Creator:     fromAddr.String(),
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
	}
	_, err = msgServer.ClaimThorchain(ctx, &invalidClaimMessage)
	require.ErrorIs(t, types.ErrInvalidCreator, err)

	claimMessage := types.MsgClaimThorchain{
		Creator:     arkeoServerAddress.String(),
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
	}
	response, err := msgServer.ClaimThorchain(ctx, &claimMessage)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, fromAddr.String(), response.FromAddress)
	require.Equal(t, toAddr.String(), response.ToAddress)

	// check if claimrecord is updated
	claimRecordFrom, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, fromAddr.String(), types.ARKEO)
	require.NoError(t, err)
	require.True(t, claimRecordFrom.IsEmpty())

	claimRecordTo, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, toAddr.String(), types.ARKEO)
	require.NoError(t, err)
	require.True(t, !claimRecordTo.IsEmpty())

	require.Equal(t, claimRecordTo.Address, toAddr.String())
	require.Equal(t, claimRecordTo.Chain, types.ARKEO)
	require.Equal(t, claimRecordTo.AmountClaim, sdk.NewInt64Coin(types.DefaultClaimDenom, 200))
	require.Equal(t, claimRecordTo.AmountVote, sdk.NewInt64Coin(types.DefaultClaimDenom, 200))
	require.Equal(t, claimRecordTo.AmountDelegate, sdk.NewInt64Coin(types.DefaultClaimDenom, 200))

	// attempt to claim again to ensure it fails.
	_, err = msgServer.ClaimThorchain(ctx, &claimMessage)
	require.ErrorIs(t, err, types.ErrNoClaimableAmount)
}

func TestClaimThorchainFailureCases(t *testing.T) {
	msgServer, keepers, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("arkeo", "arkeopub")

	arkeoServerAddress, err := sdk.AccAddressFromBech32("arkeo1zsafqx0qk6rp2vvs97n9udylquj7mfktg7s7sr")
	require.NoError(t, err)

	fromAddr := utils.GetRandomArkeoAddress()
	toAddr := utils.GetRandomArkeoAddress()

	// Test case 1: Same from and to address
	sameAddressMsg := types.MsgClaimThorchain{
		Creator:     arkeoServerAddress.String(),
		FromAddress: fromAddr.String(),
		ToAddress:   fromAddr.String(),
	}
	_, err = msgServer.ClaimThorchain(ctx, &sameAddressMsg)
	require.ErrorIs(t, types.ErrInvalidAddress, err)

	// Test case 2: Empty claim record for from address
	emptyFromMsg := types.MsgClaimThorchain{
		Creator:     arkeoServerAddress.String(),
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
	}
	_, err = msgServer.ClaimThorchain(ctx, &emptyFromMsg)
	require.ErrorIs(t, types.ErrNoClaimableAmount, err)

	// Test case 3: Zero amount claim record
	zeroClaimRecord := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        fromAddr.String(),
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 0),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 0),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 0),
	}
	err = keepers.ClaimKeeper.SetClaimRecord(sdkCtx, zeroClaimRecord)
	require.NoError(t, err)

	zeroAmountMsg := types.MsgClaimThorchain{
		Creator:     arkeoServerAddress.String(),
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
	}
	_, err = msgServer.ClaimThorchain(ctx, &zeroAmountMsg)
	require.ErrorIs(t, types.ErrNoClaimableAmount, err)
}
