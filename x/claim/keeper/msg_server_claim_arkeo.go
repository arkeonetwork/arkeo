package keeper

import (
	"context"

	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

func (k msgServer) ClaimArkeo(goCtx context.Context, msg *types.MsgClaimArkeo) (*types.MsgClaimArkeoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	arkeoClaim, err := k.GetClaimRecord(ctx, msg.Creator, types.ARKEO)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get claim record for %s", msg.Creator)
	}

	if arkeoClaim == (types.ClaimRecord{}) || arkeoClaim.AmountClaim.IsNil() || arkeoClaim.AmountClaim.IsZero() {
		return nil, errors.Wrapf(err, "no claimable amount for %s", msg.Creator)
	}

	_, err = k.ClaimCoinsForAction(ctx, msg.Creator, types.ACTION_CLAIM)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to claim coins for %s", msg.Creator)
	}

	return &types.MsgClaimArkeoResponse{}, nil
}
