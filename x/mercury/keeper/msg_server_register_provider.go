package keeper

import (
	"context"

    "mercury/x/mercury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func (k msgServer) RegisterProvider(goCtx context.Context,  msg *types.MsgRegisterProvider) (*types.MsgRegisterProviderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

    // TODO: Handling the message
    _ = ctx

	return &types.MsgRegisterProviderResponse{}, nil
}
