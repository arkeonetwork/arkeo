package keeper_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/arkeonetwork/arkeo/testutil/sample"
	"github.com/arkeonetwork/arkeo/testutil/utils"
	"github.com/arkeonetwork/arkeo/x/claim/types"
)

func TestTransferClaim(t *testing.T) {
	msgServer, keepers, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	originalClaimAddr := sample.AccAddress()
	claimRecord := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        originalClaimAddr,
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 100),
	}
	require.NoError(t, keepers.ClaimKeeper.SetClaimRecord(sdkCtx, claimRecord))

	// mint coins to module account
	require.NoError(t, keepers.BankKeeper.MintCoins(sdkCtx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 10000))))

	creator, err := sdk.AccAddressFromBech32(originalClaimAddr)
	require.NoError(t, err)
	// get balance of arkeo address before claim
	balanceBefore := keepers.BankKeeper.GetBalance(sdkCtx, creator, types.DefaultClaimDenom)
	transferClaimMessage := &types.MsgTransferClaim{
		Creator:   originalClaimAddr,
		ToAddress: sample.AccAddress(),
	}
	_, err = msgServer.TransferClaim(ctx, transferClaimMessage)
	require.NoError(t, err)

	// check if claimrecord is updated
	claimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, originalClaimAddr, types.ARKEO)
	require.NoError(t, err)
	require.True(t, claimRecord.IsEmpty())

	// confirm we have a claimrecord for arkeo
	claimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, transferClaimMessage.ToAddress, types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, claimRecord.Address, transferClaimMessage.ToAddress)
	require.Equal(t, claimRecord.Chain, types.ARKEO)
	require.True(t, claimRecord.AmountClaim.IsZero()) // nothing to claim for claim action
	require.Equal(t, claimRecord.AmountVote, sdk.NewInt64Coin(types.DefaultClaimDenom, 100))
	require.Equal(t, claimRecord.AmountDelegate, sdk.NewInt64Coin(types.DefaultClaimDenom, 100))

	toAddr, err := sdk.AccAddressFromBech32(transferClaimMessage.ToAddress)
	require.NoError(t, err)
	// confirm balance increased by expected amount.
	balanceAfter := keepers.BankKeeper.GetBalance(sdkCtx, toAddr, types.DefaultClaimDenom)
	require.Equal(t, balanceAfter.Sub(balanceBefore), sdk.NewInt64Coin(types.DefaultClaimDenom, 100))

	// attempt to transfer it again
	_, err = msgServer.TransferClaim(ctx, transferClaimMessage)
	require.Error(t, err)
	require.True(t, errors.Is(err, sdkerrors.ErrUnknownRequest))

	// attempt to claim from arkeo should also fail!
	_, err = msgServer.ClaimArkeo(ctx, &types.MsgClaimArkeo{
		Creator: originalClaimAddr,
	})
	require.Error(t, err)
}

func TestOriginalClaimNotExistShouldFail(t *testing.T) {
	msgServer, keepers, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	originalClaimAddr := sample.AccAddress()

	// mint coins to module account
	require.NoError(t, keepers.BankKeeper.MintCoins(sdkCtx, types.ModuleName, sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 10000))))

	creator, err := sdk.AccAddressFromBech32(originalClaimAddr)
	require.NoError(t, err)
	// get balance of arkeo address before claim
	balanceBefore := keepers.BankKeeper.GetBalance(sdkCtx, creator, types.DefaultClaimDenom)
	transferClaimMessage := &types.MsgTransferClaim{
		Creator:   originalClaimAddr,
		ToAddress: sample.AccAddress(),
	}
	toAddr, err := sdk.AccAddressFromBech32(transferClaimMessage.ToAddress)
	require.NoError(t, err)
	toAddrBalanceBefore := keepers.BankKeeper.GetBalance(sdkCtx, toAddr, types.DefaultClaimDenom)

	_, err = msgServer.TransferClaim(ctx, transferClaimMessage)
	require.Error(t, err)

	// check if claimrecord is updated
	claimRecord, err := keepers.ClaimKeeper.GetClaimRecord(sdkCtx, originalClaimAddr, types.ARKEO)
	require.NoError(t, err)
	require.True(t, claimRecord.IsEmpty())

	// confirm we have a claimrecord for arkeo
	claimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, transferClaimMessage.ToAddress, types.ARKEO)
	require.NoError(t, err)
	require.True(t, claimRecord.IsEmpty())

	afterCreatorBalance := keepers.BankKeeper.GetBalance(sdkCtx, creator, types.DefaultClaimDenom)
	require.True(t, afterCreatorBalance.Equal(balanceBefore))
	toAddrAfterBalanceBefore := keepers.BankKeeper.GetBalance(sdkCtx, toAddr, types.DefaultClaimDenom)
	require.True(t, toAddrAfterBalanceBefore.Equal(toAddrBalanceBefore))
}

func TestTransferClaimWithExistingClaimShouldMerge(t *testing.T) {
	msgServer, keepers, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// create valid eth claimrecords
	addrArkeo := utils.GetRandomArkeoAddress()

	claimRecord := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        addrArkeo.String(),
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
	}
	require.NoError(t, keepers.ClaimKeeper.SetClaimRecord(sdkCtx, claimRecord))

	arkeoToAddress := utils.GetRandomArkeoAddress()

	// create an arkeo claim record for the same user. This should be merged once they call claim.
	claimRecordArkeo := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        arkeoToAddress.String(),
		AmountClaim:    sdk.Coin{},
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 200),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 150),
	}
	require.NoError(t, keepers.ClaimKeeper.SetClaimRecord(sdkCtx, claimRecordArkeo))

	// mint coins to module account
	require.NoError(t, keepers.BankKeeper.MintCoins(sdkCtx, types.ModuleName,
		sdk.NewCoins(sdk.NewInt64Coin(types.DefaultClaimDenom, 10000))))

	// get balance of arkeo address before claim
	toAddressBalanceBefore := keepers.BankKeeper.GetBalance(sdkCtx, arkeoToAddress, types.DefaultClaimDenom)
	transferClaimMessage := types.MsgTransferClaim{
		Creator:   addrArkeo.String(),
		ToAddress: arkeoToAddress.String(),
	}
	result, err := msgServer.TransferClaim(ctx, &transferClaimMessage)
	require.NoError(t, err)
	require.NotNil(t, result)

	// check if claimrecord is updated
	claimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, addrArkeo.String(), types.ARKEO)
	require.NoError(t, err)
	require.True(t, claimRecord.IsEmpty())

	// confirm we have a claimrecord for arkeo
	claimRecord, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, arkeoToAddress.String(), types.ARKEO)
	require.NoError(t, err)
	require.Equal(t, claimRecord.Address, arkeoToAddress.String())
	require.Equal(t, claimRecord.Chain, types.ARKEO)
	require.True(t, claimRecord.AmountClaim.IsZero()) // nothing to claim for claim action
	require.Equal(t, claimRecord.AmountVote, sdk.NewInt64Coin(types.DefaultClaimDenom, 400))
	require.Equal(t, claimRecord.AmountDelegate, sdk.NewInt64Coin(types.DefaultClaimDenom, 350))

	// confirm balance increased by expected amount.
	balanceAfter := keepers.BankKeeper.GetBalance(sdkCtx, arkeoToAddress, types.DefaultClaimDenom)
	require.Equal(t, balanceAfter.Sub(toAddressBalanceBefore), sdk.NewInt64Coin(types.DefaultClaimDenom, 200))

	// attempt to transfer claim again to ensure it fails.
	result, err = msgServer.TransferClaim(ctx, &transferClaimMessage)
	require.Error(t, err)
	require.Nil(t, result)

	// attempt to claim from arkeo should also fail!
	resp, err := msgServer.ClaimArkeo(ctx, &types.MsgClaimArkeo{Creator: addrArkeo.String()})
	require.Error(t, err)
	require.Nil(t, resp)
}
