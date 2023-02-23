package keeper

import (
	"context"

	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) TransferClaim(goCtx context.Context, msg *types.MsgTransferClaim) (*types.MsgTransferClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgTransferClaimResponse{}, nil
}
