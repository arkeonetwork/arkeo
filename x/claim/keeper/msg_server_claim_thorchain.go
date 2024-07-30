package keeper

import (
	"fmt"
	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

// Verify and update the claim record based on the memo of the thorchain tx
func (k msgServer) ClaimThorchain(ctx sdk.Context, msg *types.MsgClaimThorchain) (types.ClaimRecord, error) {

	thorClaimRecord, err := k.GetClaimRecord(ctx, originalArkeoAddress, types.ARKEO)
	if err != nil {
		return types.ClaimRecord{}, errors.Wrapf(err, "failed to get claim record for %s", originalArkeoAddress)
	}
	if thorClaimRecord.IsEmpty() || thorClaimRecord.AmountClaim.IsZero() {
		return types.ClaimRecord{}, errors.Wrapf(types.ErrNoClaimableAmount, "no claimable amount for %s", originalArkeoAddress)
	}

	combinedClaimRecord := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        creator,
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, arkeoClaimRecord.AmountClaim.Amount.Int64()+thorClaimRecord.AmountClaim.Amount.Int64()),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, arkeoClaimRecord.AmountVote.Amount.Int64()+thorClaimRecord.AmountVote.Amount.Int64()),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, arkeoClaimRecord.AmountDelegate.Amount.Int64()+thorClaimRecord.AmountDelegate.Amount.Int64()),
	}
	emptyClaimRecord := types.ClaimRecord{
		Chain:          types.ARKEO,
		Address:        thorDerivedArkeoAddress,
		AmountClaim:    sdk.NewInt64Coin(types.DefaultClaimDenom, 0),
		AmountVote:     sdk.NewInt64Coin(types.DefaultClaimDenom, 0),
		AmountDelegate: sdk.NewInt64Coin(types.DefaultClaimDenom, 0),
	}
	err = k.SetClaimRecord(ctx, emptyClaimRecord)
	if err != nil {
		return types.ClaimRecord{}, fmt.Errorf("failed to set empty claim record for %s: %w", thorDerivedArkeoAddress, err)
	}
	err = k.SetClaimRecord(ctx, combinedClaimRecord)
	if err != nil {
		return types.ClaimRecord{}, fmt.Errorf("failed to set combined claim record for %s: %w", creator, err)
	}

	newClaimRecord, err := k.GetClaimRecord(ctx, creator, types.ARKEO)
	if err != nil {
		return types.ClaimRecord{}, errors.Wrapf(err, "failed to get claim record for %s", creator)
	}
	return newClaimRecord, nil
}
