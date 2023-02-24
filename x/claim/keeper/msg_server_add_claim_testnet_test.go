//go:build testnet

package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/arkeonetwork/arkeo/testutil/utils"
	"github.com/arkeonetwork/arkeo/x/claim/types"
)

func TestAddClaim(t *testing.T) {
	msgServer, keepers, ctx := setupMsgServer(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	creator := utils.GetRandomArkeoAddress()
	target := utils.GetRandomArkeoAddress()
	addClaimMsg := types.NewMsgAddClaim(creator.String(), types.ARKEO, target.String(), 1000_00000000)
	result, err := msgServer.AddClaim(ctx, addClaimMsg)
	require.NotNil(t, result)
	require.Nil(t, err)

	cr, err := keepers.ClaimKeeper.GetClaimRecord(sdkCtx, target.String(), types.ARKEO)
	require.Nil(t, err)
	require.False(t, cr.IsEmpty())
	require.Equal(t, cr.Address, target.String())
	require.Equal(t, cr.Chain, types.ARKEO)
	require.Equal(t, cr.IsTransferable, true)
	require.True(t, cr.AmountClaim.Equal(sdk.NewCoin(types.DefaultClaimDenom, sdk.NewInt(addClaimMsg.Amount))))
	require.True(t, cr.AmountVote.Equal(sdk.NewCoin(types.DefaultClaimDenom, sdk.NewInt(addClaimMsg.Amount))))
	require.True(t, cr.AmountDelegate.Equal(sdk.NewCoin(types.DefaultClaimDenom, sdk.NewInt(addClaimMsg.Amount))))

	addClaimMsg = types.NewMsgAddClaim(creator.String(), types.ETHEREUM, utils.GetRandomETHAddress(), 1000_00000000)
	result, err = msgServer.AddClaim(ctx, addClaimMsg)
	require.NotNil(t, result)
	require.Nil(t, err)
	cr, err = keepers.ClaimKeeper.GetClaimRecord(sdkCtx, target.String(), types.ETHEREUM)
	require.Nil(t, err)
	require.False(t, cr.IsEmpty())
	require.Equal(t, cr.Address, target.String())
	require.Equal(t, cr.Chain, types.ETHEREUM)
	require.Equal(t, cr.IsTransferable, false)
	require.True(t, cr.AmountClaim.Equal(sdk.NewCoin(types.DefaultClaimDenom, sdk.NewInt(addClaimMsg.Amount))))
	require.True(t, cr.AmountVote.Equal(sdk.NewCoin(types.DefaultClaimDenom, sdk.NewInt(addClaimMsg.Amount))))
	require.True(t, cr.AmountDelegate.Equal(sdk.NewCoin(types.DefaultClaimDenom, sdk.NewInt(addClaimMsg.Amount))))
}
