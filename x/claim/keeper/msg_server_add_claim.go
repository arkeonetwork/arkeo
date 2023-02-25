//go:build !testnet

package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/arkeonetwork/arkeo/x/claim/types"
)

func (k msgServer) AddClaim(goCtx context.Context, msg *types.MsgAddClaim) (*types.MsgAddClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx.Logger().Info("received add-claim message,not support")
	return &types.MsgAddClaimResponse{}, nil
}
