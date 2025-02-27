package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/claim/types"
)

// ClaimableForAction returns the claimable amount for a specific action by an address
func (k Keeper) ClaimableForAction(goCtx context.Context, req *types.QueryClaimableForActionRequest) (*types.QueryClaimableForActionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	k.logger.Info("CLAIMABLE FOR ACTION", "address", req.Address, "action", req.Action, "chain", req.Chain)

	ctx := sdk.UnwrapSDKContext(goCtx)

	claimableAmount, err := k.GetClaimableAmountForAction(ctx, req.Address, req.Action, req.Chain)
	if err != nil {
		// If the error is airdrop ended, return zero amount instead of an error
		if err == types.ErrAirdropEnded {
			return &types.QueryClaimableForActionResponse{
				Amount: sdk.NewCoin(k.GetParams(ctx).ClaimDenom, cosmos.ZeroInt()),
			}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	k.logger.Info("CLAIMABLE AMOUNT", "amount", claimableAmount)
	// If the amount is nil or not initialized, return zero amount
	if claimableAmount.IsNil() {
		k.logger.Info("CLAIMABLE AMOUNT IS NIL")
		claimableAmount = sdk.NewCoin(k.GetParams(ctx).ClaimDenom, cosmos.ZeroInt())
	}

	return &types.QueryClaimableForActionResponse{
		Amount: claimableAmount,
	}, nil
}
