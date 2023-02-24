package keeper

import (
	"context"

	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddClaim(goCtx context.Context, msg *types.MsgAddClaim) (*types.MsgAddClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgAddClaimResponse{}, nil
}
