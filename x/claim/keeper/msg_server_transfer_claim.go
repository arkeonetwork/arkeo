package keeper

import (
	"context"
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/arkeonetwork/arkeo/x/claim/types"
)

func (k msgServer) TransferClaim(goCtx context.Context, msg *types.MsgTransferClaim) (*types.MsgTransferClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// get original claim record
	originalClaim, err := k.GetClaimRecord(ctx, msg.Creator.String(), types.ARKEO)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get claim record for %s", msg.Creator)
	}

	if originalClaim.IsEmpty() || originalClaim.AmountClaim.IsZero() {
		return nil, errors.Wrapf(types.ErrNoClaimableAmount, "no claimable amount for %s", msg.Creator)
	}

	if !originalClaim.IsTransferable {
		return nil, errors.Wrapf(types.ErrClaimRecordNotTransferrable, "claim record for %s is not transferable", msg.Creator)
	}

	// create new arkeo claim
	arkeoClaim := types.ClaimRecord{
		Address:        msg.ToAddress.String(),
		Chain:          types.ARKEO,
		AmountClaim:    originalClaim.AmountClaim,
		AmountVote:     originalClaim.AmountVote,
		AmountDelegate: originalClaim.AmountDelegate,
		IsTransferable: false,
	}

	// set claim to completed
	originalClaim = setClaimableAmountForAllActions(originalClaim, sdk.Coin{})
	if err := k.SetClaimRecord(ctx, originalClaim); err != nil {
		return nil, errors.Wrapf(err, "failed to set claim record for %s", msg.Creator)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(sdk.AttributeKeySender, strings.ToLower(msg.Creator.String())),
			sdk.NewAttribute(sdk.AttributeKeyAmount, arkeoClaim.AmountClaim.String()),
		),
	})

	// see if there is an existing arkeo claim, so we can merge it
	existingArkeoClaim, err := k.GetClaimRecord(ctx, msg.ToAddress.String(), types.ARKEO)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get arkeo claim record for %s", msg.ToAddress)
	}

	arkeoClaim, err = mergeClaimRecords(existingArkeoClaim, arkeoClaim)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to merge claim records for %s", msg.ToAddress)
	}

	err = k.SetClaimRecord(ctx, arkeoClaim)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to set claim record for %s", msg.ToAddress)
	}

	_, err = k.ClaimCoinsForAction(ctx, msg.ToAddress.String(), types.ACTION_CLAIM)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to claim coins for %s", msg.ToAddress)
	}

	return &types.MsgTransferClaimResponse{}, nil
}
