package keeper

import (
	"context"

	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ClaimRecord(goCtx context.Context, req *types.QueryClaimRecordRequest) (*types.QueryClaimRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Process the query
	_ = ctx

	return &types.QueryClaimRecordResponse{}, nil
}
