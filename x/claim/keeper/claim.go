package keeper

import (
	"strings"

	"github.com/arkeonetwork/arkeo/x/claim/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
)

// SetClaimRecord sets a claim record for an address in store
func (k Keeper) SetClaimRecord(ctx sdk.Context, claimRecord types.ClaimRecord) error {
	// validate address if valid based on chain
	switch claimRecord.Chain {
	case types.ETHEREUM:
		if !types.IsValidEthAddress(claimRecord.Address) {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid eth address")
		}
	case types.ARKEO:
		_, err := sdk.AccAddressFromBech32(claimRecord.Address)
		if err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid arkeo address")
		}
	case types.THORCHAIN:
		return errors.New("thorchain not supported yet")
	}

	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, chainToStorePrefix(claimRecord.Chain))

	bz, err := k.cdc.Marshal(&claimRecord)
	if err != nil {
		return err
	}

	addr := []byte(strings.ToLower(claimRecord.Address))

	prefixStore.Set(addr, bz)
	return nil
}

// SetClaimables set claimable amount from balances object
func (k Keeper) SetClaimRecords(ctx sdk.Context, claimRecords []types.ClaimRecord) error {
	for _, claimRecord := range claimRecords {
		err := k.SetClaimRecord(ctx, claimRecord)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetClaimables get claimables for genesis export
func (k Keeper) GetClaimRecords(ctx sdk.Context, chain types.Chain) ([]types.ClaimRecord, error) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, chainToStorePrefix(chain))

	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	claimRecords := []types.ClaimRecord{}
	for ; iterator.Valid(); iterator.Next() {

		claimRecord := types.ClaimRecord{}

		err := k.cdc.Unmarshal(iterator.Value(), &claimRecord)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal claim record")
		}
		claimRecords = append(claimRecords, claimRecord)
	}
	return claimRecords, nil
}

// GetClaimRecord returns the claim record for a specific address
func (k Keeper) GetClaimRecord(ctx sdk.Context, addr string, chain types.Chain) (types.ClaimRecord, error) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, chainToStorePrefix(chain))
	addrBytes := []byte(strings.ToLower(addr))
	if !prefixStore.Has(addrBytes) {
		return types.ClaimRecord{}, nil
	}
	bz := prefixStore.Get(addrBytes)

	claimRecord := types.ClaimRecord{}
	err := k.cdc.Unmarshal(bz, &claimRecord)
	if err != nil {
		return types.ClaimRecord{}, err
	}

	return claimRecord, nil
}

// GetUserTotalClaimable returns the total claimable amount for an address
func (k Keeper) GetUserTotalClaimable(ctx sdk.Context, addr string, chain types.Chain) (sdk.Coins, error) {
	claimRecord, err := k.GetClaimRecord(ctx, addr, chain)
	if err != nil {
		return sdk.Coins{}, err
	}
	if claimRecord.Address == "" {
		return sdk.Coins{}, nil
	}

	totalClaimable := sdk.Coins{}
	for action := range types.Action_name {
		claimableForAction, err := k.GetClaimableAmountForAction(ctx, addr, types.Action(action), chain)
		if err != nil {
			return sdk.Coins{}, err
		}
		totalClaimable = totalClaimable.Add(claimableForAction...)
	}
	return totalClaimable, nil
}

// GetClaimable returns claimable amount for a specific action done by an address
func (k Keeper) GetClaimableAmountForAction(ctx sdk.Context, addr string, action types.Action, chain types.Chain) (sdk.Coins, error) {
	claimRecord, err := k.GetClaimRecord(ctx, addr, chain)
	if err != nil {
		return nil, err
	}

	if claimRecord.Address == "" {
		return sdk.Coins{}, nil
	}

	// if action already completed, nothing is claimable
	if claimRecord.ActionCompleted[action] {
		return sdk.Coins{}, nil
	}

	params := k.GetParams(ctx)

	// If we are before the start time, do nothing.
	// This case _shouldn't_ occur on chain, since the
	// start time ought to be chain start time.
	if ctx.BlockTime().Before(params.AirdropStartTime) {
		return sdk.Coins{}, nil
	}

	InitialClaimablePerAction := sdk.Coins{}
	for _, coin := range claimRecord.InitialClaimableAmount {
		InitialClaimablePerAction = InitialClaimablePerAction.Add(
			sdk.NewCoin(coin.Denom,
				coin.Amount.QuoRaw(int64(len(types.Action_name))),
			),
		)
	}

	elapsedAirdropTime := ctx.BlockTime().Sub(params.AirdropStartTime)
	// Are we early enough in the airdrop s.t. theres no decay?
	if elapsedAirdropTime <= params.DurationUntilDecay {
		return InitialClaimablePerAction, nil
	}

	// The entire airdrop has completed
	if elapsedAirdropTime > params.DurationUntilDecay+params.DurationOfDecay {
		return sdk.Coins{}, nil
	}

	// Positive, since goneTime > params.DurationUntilDecay
	decayTime := elapsedAirdropTime - params.DurationUntilDecay
	decayPercent := sdk.NewDec(decayTime.Nanoseconds()).QuoInt64(params.DurationOfDecay.Nanoseconds())
	claimablePercent := sdk.OneDec().Sub(decayPercent)

	claimableCoins := sdk.Coins{}
	for _, coin := range InitialClaimablePerAction {
		claimableAmount := coin.Amount.Mul(claimablePercent.Mul(sdk.NewDec(10000)).RoundInt()).QuoRaw(10000)
		claimableCoins = claimableCoins.Add(sdk.NewCoin(coin.Denom, claimableAmount))
	}

	return claimableCoins, nil
}

// // ClaimCoins remove claimable amount entry and transfer it to user's account
// func (k Keeper) ClaimCoinsForAction(ctx sdk.Context, addr sdk.AccAddress, action types.Action) (sdk.Coins, error) {
// 	claimableAmount, err := k.GetClaimableAmountForAction(ctx, addr, action)
// 	if err != nil {
// 		return claimableAmount, err
// 	}

// 	if claimableAmount.Empty() {
// 		return claimableAmount, nil
// 	}

// 	claimRecord, err := k.GetClaimRecord(ctx, addr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, claimableAmount)
// 	if err != nil {
// 		return nil, err
// 	}

// 	claimRecord.ActionCompleted[action] = true

// 	err = k.SetClaimRecord(ctx, claimRecord)
// 	if err != nil {
// 		return claimableAmount, err
// 	}

// 	ctx.EventManager().EmitEvents(sdk.Events{
// 		sdk.NewEvent(
// 			types.EventTypeClaim,
// 			sdk.NewAttribute(sdk.AttributeKeySender, addr.String()),
// 			sdk.NewAttribute(sdk.AttributeKeyAmount, claimableAmount.String()),
// 		),
// 	})

// 	return claimableAmount, nil
// }

// // FundRemainingsToCommunity fund remainings to the community when airdrop period end
// func (k Keeper) fundRemainingsToCommunity(ctx sdk.Context) error {
// 	moduleAccAddr := k.GetModuleAccountAddress(ctx)
// 	amt := k.GetModuleAccountBalance(ctx)
// 	return k.distrKeeper.FundCommunityPool(ctx, sdk.NewCoins(amt), moduleAccAddr)
// }

func chainToStorePrefix(chain types.Chain) []byte {
	switch chain {
	case types.ARKEO:
		return []byte(types.ClaimRecordsArkeoStorePrefix)
	case types.ETHEREUM:
		return []byte(types.ClaimRecordsEthStorePrefix)
	case types.THORCHAIN:
		return []byte(types.ClaimRecordsThorStorePrefix)
	default:
		return []byte{}
	}
}
