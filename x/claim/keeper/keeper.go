package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/arkeonetwork/arkeo/x/claim/types"
)

type (
	Keeper struct {
		cdc           codec.BinaryCodec
		storeKey      storetypes.StoreKey
		memKey        storetypes.StoreKey
		paramstore    paramtypes.Subspace
		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
		logger        log.Logger
	}
)

var UpdateParamsKey = []byte("UpdateParamsKey")
var TransferTokensKey = []byte("TransferTokensKey")

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	logger log.Logger,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		memKey:        memKey,
		paramstore:    ps,
		logger:        logger.With("module", fmt.Sprintf("x/%s", types.ModuleName)),
	}
}

func (k Keeper) Logger(ctx context.Context) log.Logger {
	return k.logger
}

func (k Keeper) AfterProposalVote(ctx context.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	_, err := k.ClaimCoinsForAction(sdkCtx, voterAddr.String(), types.ACTION_VOTE)
	if err != nil {
		k.Logger(ctx).Error("failed to claim coins for vote", "error", err.Error())
	}
}

func (k Keeper) AfterDelegationModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	_, err := k.ClaimCoinsForAction(sdkCtx, delAddr.String(), types.ACTION_DELEGATE)
	if err != nil {
		k.Logger(ctx).Error("failed to claim coins for delegate", "error", err.Error())
	}
	return nil
}

func (k Keeper) UpdateParams(ctx context.Context) bool {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(k.storeKey)
	if store.Has(UpdateParamsKey) {
		return false
	}

	var params types.Params
	k.paramstore.GetParamSet(sdkCtx, &params)

	params.DurationOfDecay = 37 * 24 * time.Hour    // 30 days as time.Duration
	params.DurationUntilDecay = 30 * 24 * time.Hour // 30 days as time.Duration

	k.paramstore.SetParamSet(sdkCtx, &params)
	store.Set(UpdateParamsKey, []byte{1})

	return true
}

func (k Keeper) MoveTokensFromReserveToClaimModule(ctx context.Context) bool {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(k.storeKey)
	if store.Has(TransferTokensKey) {
		return false
	}
	// Check if we've passed the end of the airdrop
	params := k.GetParams(sdkCtx)
	if sdkCtx.BlockTime().After(params.AirdropStartTime.Add(params.DurationUntilDecay).Add(params.DurationOfDecay)) {
		return false
	}

	reserveAddr := k.accountKeeper.GetModuleAddress(types.ReserveModuleName)
	reserve := k.bankKeeper.GetBalance(ctx, reserveAddr, types.DefaultClaimDenom)

	amountToSend := sdk.NewInt64Coin(types.DefaultClaimDenom, 23_000_000*int64(1e8))
	if reserve.Amount.GTE(amountToSend.Amount) {
		err := k.bankKeeper.SendCoinsFromModuleToModule(
			sdkCtx,
			types.ReserveModuleName,
			types.ClaimModuleName,
			sdk.NewCoins(amountToSend),
		)
		if err != nil {
			k.Logger(ctx).Error("failed to transfer tokens from reserve to claim module", "error", err.Error())
			return false
		}
		k.Logger(ctx).Info("transferred 23m ARKEO from reserve back to claim module")
		store.Set(TransferTokensKey, []byte{1})
		store.Delete(MoveClaimTokensKey)
	} else {
		k.Logger(ctx).Error("not enough reserve balance to transfer 23m tokens")
		return false
	}

	return true
}
