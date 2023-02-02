package keeper

import (
	"context"

	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ClaimEth(goCtx context.Context, msg *types.MsgClaimEth) (*types.MsgClaimEthResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgClaimEthResponse{}, nil
}
