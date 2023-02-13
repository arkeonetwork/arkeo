package keeper

import (
	"context"

	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ClaimArkeo(goCtx context.Context, msg *types.MsgClaimArkeo) (*types.MsgClaimArkeoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgClaimArkeoResponse{}, nil
}
