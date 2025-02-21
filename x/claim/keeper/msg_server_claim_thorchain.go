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
	if msg.Creator != "tarkeo1z02ke8639m47g9dfrheegr2u9zecegt5qvtj00" && msg.Creator != "arkeo1z02ke8639m47g9dfrheegr2u9zecegt50fjg7v" {
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

	amountClaim := sdk.NewInt64Coin(types.DefaultClaimDenom, toAddressClaimRecord.AmountClaim.Amount.Int64()+fromAddressClaimRecord.AmountClaim.Amount.Int64())
	amountVote := sdk.NewInt64Coin(types.DefaultClaimDenom, toAddressClaimRecord.AmountVote.Amount.Int64()+fromAddressClaimRecord.AmountVote.Amount.Int64())
	amountDelegate := sdk.NewInt64Coin(types.DefaultClaimDenom, toAddressClaimRecord.AmountDelegate.Amount.Int64()+fromAddressClaimRecord.AmountDelegate.Amount.Int64())

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
