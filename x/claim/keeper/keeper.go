package keeper

import (
	"context"
	"fmt"

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
