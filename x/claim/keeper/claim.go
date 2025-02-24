package keeper

import (
	"sort"
	"context"
	"strings"

	sdkerror "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/claim/types"
)

// Key for tracking if end of airdrop transfer has occurred
var MoveClaimTokensKey = []byte("MoveClaimTokens")

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx context.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(k.storeKey)

	// Check if we've already done the transfer
	if store.Has(MoveClaimTokensKey) {
		return
	}

	params := k.GetParams(sdkCtx)
	// Check if we've passed the end of the airdrop
	if sdkCtx.BlockTime().After(params.AirdropStartTime.Add(params.DurationUntilDecay).Add(params.DurationOfDecay)) {
		err := k.TransferRemainingToReserve(sdkCtx)
		if err != nil {
			k.Logger(ctx).Error("failed to transfer remaining tokens to reserve", "error", err)
			return
		}

		store.Set(MoveClaimTokensKey, []byte{1})
	}
}

// SetClaimRecord sets a claim record for an address in store
func (k Keeper) SetClaimRecord(ctx sdk.Context, claimRecord types.ClaimRecord) error {
	// validate address if valid based on chain
	if !types.IsValidAddress(claimRecord.Address, claimRecord.Chain) {
		return sdkerror.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address %s for chain %s", claimRecord.Address, claimRecord.Chain.String())
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

func (k Keeper) GetAllClaimRecords(ctx sdk.Context) ([]types.ClaimRecord, error) {
	claimRecords := []types.ClaimRecord{}

	chains := make([]types.Chain, 0, len(types.Chain_name))
	for chain := range types.Chain_name {
		chains = append(chains, types.Chain(chain))
	}

	sort.Slice(chains, func(i, j int) bool {
		return int32(chains[i]) < int32(chains[j])
	})

	for _, chain := range chains {
		records, err := k.GetClaimRecords(ctx, chain)
		if err != nil {
			return nil, err
		}
		claimRecords = append(claimRecords, records...)
	}

	return claimRecords, nil
}

// GetClaimRecord returns the claim record for a specific address
func (k Keeper) GetClaimRecord(ctx sdk.Context, addr string, chain types.Chain) (types.ClaimRecord, error) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, chainToStorePrefix(chain))
	addrBytes := []byte(strings.ToLower(addr))
	if !prefixStore.Has(addrBytes) {
		return types.ClaimRecord{Chain: chain}, nil
	}
	bz := prefixStore.Get(addrBytes)

	claimRecord := types.ClaimRecord{}
	err := k.cdc.Unmarshal(bz, &claimRecord)
	if err != nil {
		return types.ClaimRecord{Chain: chain}, err
	}

	return claimRecord, nil
}

// GetUserTotalClaimable returns the total claimable amount for an address across all actions
func (k Keeper) GetUserTotalClaimable(ctx sdk.Context, addr string, chain types.Chain) (sdk.Coin, error) {
	claimRecord, err := k.GetClaimRecord(ctx, addr, chain)
	if err != nil {
		return sdk.Coin{}, err
	}
	if claimRecord.IsEmpty() {
		return sdk.Coin{}, nil
	}

	actions := make([]types.Action, 0, len(types.Action_name))
	for action := range types.Action_name {
		actions = append(actions, types.Action(action))
	}

	sort.Slice(actions, func(i, j int) bool {
		return int32(actions[i]) < int32(actions[j])
	})

	totalClaimable := sdk.NewCoin(claimRecord.AmountClaim.Denom, cosmos.ZeroInt())

	for _, action := range actions {
		claimableForAction, err := k.GetClaimableAmountForAction(ctx, addr, action, chain)
		if err != nil {
			return sdk.Coin{}, err
		}
		if claimableForAction.IsNil() {
			continue
		}
		totalClaimable = totalClaimable.AddAmount(claimableForAction.Amount)
	}
	return totalClaimable, nil
}

// GetClaimable returns claimable amount for a specific action done by an address
func (k Keeper) GetClaimableAmountForAction(ctx sdk.Context, addr string, action types.Action, chain types.Chain) (sdk.Coin, error) {
	claimRecord, err := k.GetClaimRecord(ctx, addr, chain)
	if err != nil {
		return sdk.Coin{}, err
	}

	if claimRecord.IsEmpty() {
		return sdk.Coin{}, nil
	}

	params := k.GetParams(ctx)

	// If we are before the start time, do nothing.
	// This case _shouldn't_ occur on chain, since the
	// start time ought to be chain start time.
	if ctx.BlockTime().Before(params.AirdropStartTime) {
		return sdk.Coin{}, nil
	}

	initalClaimableAmount := getInitialClaimableAmount(claimRecord, action)
	if initalClaimableAmount.IsNil() || initalClaimableAmount.Amount.IsZero() {
		return sdk.Coin{}, nil
	}

	elapsedAirdropTime := ctx.BlockTime().Sub(params.AirdropStartTime)
	// Are we early enough in the airdrop s.t. theres no decay?
	if elapsedAirdropTime <= params.DurationUntilDecay {
		return initalClaimableAmount, nil
	}

	// The entire airdrop has completed
	if elapsedAirdropTime > params.DurationUntilDecay+params.DurationOfDecay {
		return sdk.Coin{}, types.ErrAirdropEnded
	}

	// Positive, since goneTime > params.DurationUntilDecay
	decayTime := elapsedAirdropTime - params.DurationUntilDecay
	decayPercent := cosmos.NewDec(decayTime.Nanoseconds()).QuoInt64(params.DurationOfDecay.Nanoseconds())
	claimablePercent := sdkmath.LegacyOneDec().Sub(decayPercent)

	claimableAmount := initalClaimableAmount.Amount.Mul(claimablePercent.Mul(cosmos.NewDec(10000)).RoundInt()).QuoRaw(10000)
	claimableCoin := sdk.NewCoin(initalClaimableAmount.Denom, claimableAmount)

	return claimableCoin, nil
}

// TransferRemainingToReserve transfers remaining funds to the reserve module when airdrop period ends
func (k Keeper) TransferRemainingToReserve(ctx sdk.Context) error {
	remainingAmount := k.GetModuleAccountBalance(ctx)

	if remainingAmount.IsZero() {
		return nil
	}

	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, "arkeo-reserve", sdk.NewCoins(remainingAmount))
}

// GetModuleAccountBalance gets the airdrop coin balance of module account
func (k Keeper) GetModuleAccountAddress(ctx sdk.Context) sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}

// GetModuleAccountBalance gets the airdrop coin balance of module account
func (k Keeper) GetModuleAccountBalance(ctx sdk.Context) sdk.Coin {
	moduleAccAddr := k.GetModuleAccountAddress(ctx)
	params := k.GetParams(ctx)
	return k.bankKeeper.GetBalance(ctx, moduleAccAddr, params.ClaimDenom)
}

// ClaimCoins remove claimable amount entry and transfer it to user's account
func (k Keeper) ClaimCoinsForAction(ctx sdk.Context, addr string, action types.Action) (sdk.Coin, error) {
	claimableAmount, err := k.GetClaimableAmountForAction(ctx, addr, action, types.ARKEO)
	if err != nil {
		return claimableAmount, err
	}

	if claimableAmount.IsNil() || claimableAmount.IsZero() {
		return claimableAmount, nil
	}

	claimRecord, err := k.GetClaimRecord(ctx, addr, types.ARKEO)
	if err != nil {
		return sdk.Coin{}, err
	}

	accountAddress, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return sdk.Coin{}, err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, accountAddress, sdk.NewCoins(claimableAmount))
	if err != nil {
		return sdk.Coin{}, err
	}

	claimRecord = setClaimableAmountForAction(claimRecord, action, sdk.Coin{}) // set to nil/zero to mark as completed.
	err = k.SetClaimRecord(ctx, claimRecord)
	if err != nil {
		return sdk.Coin{}, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(sdk.AttributeKeySender, addr),
			sdk.NewAttribute(sdk.AttributeKeyAmount, claimableAmount.String()),
		),
	})

	return claimableAmount, nil
}

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
	default:
		return []byte{}
	}
}

func getInitialClaimableAmountTotal(claim types.ClaimRecord) sdk.Coin {
	totalAmount := sdk.NewCoin(claim.AmountClaim.Denom, cosmos.ZeroInt())
	totalAmount = totalAmount.Add(claim.AmountClaim)
	totalAmount = totalAmount.Add(claim.AmountDelegate)
	totalAmount = totalAmount.Add(claim.AmountVote)
	return totalAmount
}

func getInitialClaimableAmount(claim types.ClaimRecord, action types.Action) sdk.Coin {
	switch action {
	case types.ACTION_CLAIM:
		return claim.AmountClaim
	case types.ACTION_DELEGATE:
		return claim.AmountDelegate
	case types.ACTION_VOTE:
		return claim.AmountVote
	default:
		return sdk.Coin{}
	}
}

func setClaimableAmountForAction(claim types.ClaimRecord, action types.Action, amount sdk.Coin) types.ClaimRecord {
	switch action {
	case types.ACTION_CLAIM:
		claim.AmountClaim = amount
	case types.ACTION_DELEGATE:
		claim.AmountDelegate = amount
	case types.ACTION_VOTE:
		claim.AmountVote = amount
	}
	return claim
}

func setClaimableAmountForAllActions(claim types.ClaimRecord, amount sdk.Coin) types.ClaimRecord {
	claim.AmountClaim = amount
	claim.AmountDelegate = amount
	claim.AmountVote = amount
	return claim
}
