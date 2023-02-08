package keeper

import (
	"strings"

	"github.com/arkeonetwork/arkeo/x/claim/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/pkg/errors"
)

// SetClaimRecord sets a claim record for an address in store
func (k Keeper) SetClaimRecord(ctx sdk.Context, claimRecord types.ClaimRecord) error {
	// validate address if valid based on chain
	if !types.IsValidAddress(claimRecord.Address, claimRecord.Chain) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address for chain %s", claimRecord.Chain.String())
	}

	if len(claimRecord.ActionCompleted) != len(types.Action_name) {
		return errors.New("claim record not properly initialized, length of ActionCompleted must match the length of Actions")
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
	for chain := range types.Chain_name {
		records, err := k.GetClaimRecords(ctx, types.Chain(chain))
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
func (k Keeper) GetUserTotalClaimable(ctx sdk.Context, addr string, chain types.Chain) (sdk.Coin, error) {
	claimRecord, err := k.GetClaimRecord(ctx, addr, chain)
	if err != nil {
		return sdk.Coin{}, err
	}
	if claimRecord.Address == "" {
		return sdk.Coin{}, nil
	}

	totalClaimable := sdk.NewCoin(claimRecord.InitialClaimableAmount.Denom, sdk.ZeroInt())
	for action := range types.Action_name {
		claimableForAction, err := k.GetClaimableAmountForAction(ctx, addr, types.Action(action), chain)
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

	if claimRecord.Address == "" {
		return sdk.Coin{}, nil
	}

	if len(claimRecord.ActionCompleted) != len(types.Action_name) {
		return sdk.Coin{}, errors.New("claim record is not initialized")
	}

	// if action already completed, nothing is claimable
	if claimRecord.ActionCompleted[action] {
		return sdk.Coin{}, nil
	}

	params := k.GetParams(ctx)

	// If we are before the start time, do nothing.
	// This case _shouldn't_ occur on chain, since the
	// start time ought to be chain start time.
	if ctx.BlockTime().Before(params.AirdropStartTime) {
		return sdk.Coin{}, nil
	}

	initalClaimableAmount := sdk.NewCoin(claimRecord.InitialClaimableAmount.Denom,
		claimRecord.InitialClaimableAmount.Amount.QuoRaw(int64(len(types.Action_name))))

	elapsedAirdropTime := ctx.BlockTime().Sub(params.AirdropStartTime)
	// Are we early enough in the airdrop s.t. theres no decay?
	if elapsedAirdropTime <= params.DurationUntilDecay {
		return initalClaimableAmount, nil
	}

	// The entire airdrop has completed
	if elapsedAirdropTime > params.DurationUntilDecay+params.DurationOfDecay {
		return sdk.Coin{}, nil
	}

	// Positive, since goneTime > params.DurationUntilDecay
	decayTime := elapsedAirdropTime - params.DurationUntilDecay
	decayPercent := sdk.NewDec(decayTime.Nanoseconds()).QuoInt64(params.DurationOfDecay.Nanoseconds())
	claimablePercent := sdk.OneDec().Sub(decayPercent)

	claimableAmount := initalClaimableAmount.Amount.Mul(claimablePercent.Mul(sdk.NewDec(10000)).RoundInt()).QuoRaw(10000)
	claimableCoin := sdk.NewCoin(initalClaimableAmount.Denom, claimableAmount)

	return claimableCoin, nil
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

	if claimableAmount.IsNil() {
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

	claimRecord.ActionCompleted[action] = true

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

// CreateModuleAccount creates module account and mint coins to it
func (k Keeper) CreateModuleAccount(ctx sdk.Context, amount sdk.Coin) {
	moduleAcc := authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Minter)
	k.accountKeeper.SetModuleAccount(ctx, moduleAcc)
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(amount))
	if err != nil {
		panic(err) // module can not be set up correctly, should panic?
	}
}

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
