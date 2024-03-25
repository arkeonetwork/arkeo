package keeper

import (
	"context"
	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

func (k msgServer) ClaimArkeo(goCtx context.Context, msg *types.MsgClaimArkeo) (*types.MsgClaimArkeoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	arkeoClaimRecord, err := k.GetClaimRecord(ctx, msg.Creator.String(), types.ARKEO)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get claim record for %s", msg.Creator)
	}

	if msg.ThorTx != "" {
		arkeoClaimRecord, err = k.updateThorClaimRecord(ctx, msg.Creator.String(), msg.ThorTx, arkeoClaimRecord)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get claim record for %s", msg.ThorTx)
		}
	}

	if arkeoClaimRecord.IsEmpty() || arkeoClaimRecord.AmountClaim.IsZero() {
		return nil, errors.Wrapf(types.ErrNoClaimableAmount, "no claimable amount for %s", msg.Creator)
	}

	_, err = k.ClaimCoinsForAction(ctx, msg.Creator.String(), types.ACTION_CLAIM)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to claim coins for %s", msg.Creator)
	}

	return &types.MsgClaimArkeoResponse{}, nil
}
