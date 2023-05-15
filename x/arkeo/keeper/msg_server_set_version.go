package keeper

import (
	"context"

    "github.com/arkeonetwork/arkeo/x/arkeo/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func (k msgServer) SetVersion(goCtx context.Context,  msg *types.MsgSetVersion) (*types.MsgSetVersionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

    // TODO: Handling the message
    _ = ctx

	return &types.MsgSetVersionResponse{}, nil
}
