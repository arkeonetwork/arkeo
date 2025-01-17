package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"

	"github.com/arkeonetwork/arkeo/x/claim/types"
)

func (k msgServer) ClaimArkeo(goCtx context.Context, msg *types.MsgClaimArkeo) (*types.MsgClaimArkeoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	arkeoClaimRecord, err := k.GetClaimRecord(ctx, msg.Creator, types.ARKEO)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get claim record for %s", msg.Creator)
	}

	if arkeoClaimRecord.IsEmpty() || arkeoClaimRecord.AmountClaim.IsZero() {
		return nil, errors.Wrapf(types.ErrNoClaimableAmount, "no claimable amount for %s", msg.Creator)
	}

	_, err = k.ClaimCoinsForAction(ctx, msg.Creator, types.ACTION_CLAIM)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to claim coins for %s", msg.Creator)
	}

	claimedAmount := arkeoClaimRecord.AmountClaim

	return &types.MsgClaimArkeoResponse{
		Address: msg.Creator,
		Amount:  claimedAmount.Amount.Int64(),
	}, nil
}
