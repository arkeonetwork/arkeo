package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"

	"github.com/arkeonetwork/arkeo/x/claim/types"
)

func (k msgServer) ClaimThorchain(goCtx context.Context, msg *types.MsgClaimThorchain) (*types.MsgClaimThorchainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	k.Logger(ctx).Info(msg.Creator)

	// Add check for matching addresses
	if msg.FromAddress == msg.ToAddress {
		return nil, errors.Wrapf(types.ErrInvalidAddress, "from address and to address cannot be the same: %s", msg.FromAddress)
	}

	// only allow thorchain claim server address to call this function
	if msg.Creator != "arkeo1zsafqx0qk6rp2vvs97n9udylquj7mfktg7s7sr" {
		return nil, errors.Wrapf(types.ErrInvalidCreator, "Invalid Creator %s", msg.Creator)
	}

	fromAddressClaimRecord, err := k.GetClaimRecord(ctx, msg.FromAddress, types.ARKEO)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get claim record for %s", msg.FromAddress)
	}
	if fromAddressClaimRecord.IsEmpty() || (fromAddressClaimRecord.AmountClaim.IsZero() && fromAddressClaimRecord.AmountVote.IsZero() && fromAddressClaimRecord.AmountDelegate.IsZero()) {
		return nil, errors.Wrapf(types.ErrNoClaimableAmount, "no claimable amount for %s", msg.FromAddress)
	}

	toAddressClaimRecord, err := k.GetClaimRecord(ctx, msg.ToAddress, types.ARKEO)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get claim record for %s", msg.ToAddress)
	}

	// Get amounts from toAddress, defaulting to 0 if nil
	toClaimAmount := int64(0)
	if !toAddressClaimRecord.AmountClaim.IsNil() {
		toClaimAmount = toAddressClaimRecord.AmountClaim.Amount.Int64()
	}
	toVoteAmount := int64(0)
	if !toAddressClaimRecord.AmountVote.IsNil() {
		toVoteAmount = toAddressClaimRecord.AmountVote.Amount.Int64()
	}
	toDelegateAmount := int64(0)
	if !toAddressClaimRecord.AmountDelegate.IsNil() {
		toDelegateAmount = toAddressClaimRecord.AmountDelegate.Amount.Int64()
	}

	// Get amounts from fromAddress, defaulting to 0 if nil
	fromClaimAmount := int64(0)
	if !fromAddressClaimRecord.AmountClaim.IsNil() {
		fromClaimAmount = fromAddressClaimRecord.AmountClaim.Amount.Int64()
	}
	fromVoteAmount := int64(0)
	if !fromAddressClaimRecord.AmountVote.IsNil() {
		fromVoteAmount = fromAddressClaimRecord.AmountVote.Amount.Int64()
	}
	fromDelegateAmount := int64(0)
	if !fromAddressClaimRecord.AmountDelegate.IsNil() {
		fromDelegateAmount = fromAddressClaimRecord.AmountDelegate.Amount.Int64()
	}

	amountClaim := sdk.NewInt64Coin(types.DefaultClaimDenom, toClaimAmount+fromClaimAmount)
	amountVote := sdk.NewInt64Coin(types.DefaultClaimDenom, toVoteAmount+fromVoteAmount)
	amountDelegate := sdk.NewInt64Coin(types.DefaultClaimDenom, toDelegateAmount+fromDelegateAmount)

	combinedClaimRecord := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        msg.ToAddress,
		AmountClaim:    amountClaim,
		AmountVote:     amountVote,
		AmountDelegate: amountDelegate,
		IsTransferable: toAddressClaimRecord.IsTransferable,
	}
	emptyClaimRecord := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        msg.FromAddress,
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 0),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 0),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 0),
	}

	err = k.SetClaimRecord(ctx, emptyClaimRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to set empty claim record for %s: %w", msg.FromAddress, err)
	}
	err = k.SetClaimRecord(ctx, combinedClaimRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to set combined claim record for %s: %w", msg.ToAddress, err)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeThorchainDelegate,
			sdk.NewAttribute(types.AttributeKeyFromAddress, msg.FromAddress),
			sdk.NewAttribute(types.AttributeKeyToAddress, msg.ToAddress),
		),
	})

	return &types.MsgClaimThorchainResponse{
		FromAddress: msg.FromAddress,
		ToAddress:   msg.ToAddress,
	}, nil
}
