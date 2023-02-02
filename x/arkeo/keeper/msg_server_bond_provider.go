package keeper

import (
	"context"
	"fmt"

	"github.com/ArkeoNetwork/arkeo/common"
	"github.com/ArkeoNetwork/arkeo/common/cosmos"
	"github.com/ArkeoNetwork/arkeo/x/arkeo/configs"
	"github.com/ArkeoNetwork/arkeo/x/arkeo/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) BondProvider(goCtx context.Context, msg *types.MsgBondProvider) (*types.MsgBondProviderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.Logger().Info(
		"receive MsgBondProvider",
		"pubkey", msg.PubKey,
		"chain", msg.Chain,
		"bond", msg.Bond,
	)

	cacheCtx, commit := ctx.CacheContext()
	if err := k.BondProviderValidate(cacheCtx, msg); err != nil {
		ctx.Logger().Error("failed bond provider validation", "err", err)
		return nil, err
	}

	if err := k.BondProviderHandle(cacheCtx, msg); err != nil {
		ctx.Logger().Error("failed bond provider handle", "err", err)
		return nil, err
	}
	commit()

	return &types.MsgBondProviderResponse{}, nil
}

func (k msgServer) BondProviderValidate(ctx cosmos.Context, msg *types.MsgBondProvider) error {
	if k.FetchConfig(ctx, configs.HandlerBondProvider) > 0 {
		return sdkerrors.Wrapf(types.ErrDisabledHandler, "bond provider")
	}

	// We allow providers to unbond WHILE active contracts are underway. This
	// is because A) users can cancel their owned contracts at any time, and B)
	// this is the way the provider signals to the chain that they don't want
	// to open any new contracts (as there is a min bond requirement for new
	// contracts to be opened)

	return nil
}

func (k msgServer) BondProviderHandle(ctx cosmos.Context, msg *types.MsgBondProvider) error {
	chain, err := common.NewChain(msg.Chain)
	if err != nil {
		return err
	}
	provider, err := k.GetProvider(ctx, msg.PubKey, chain)
	if err != nil {
		return err
	}
	addr, err := msg.PubKey.GetMyAddress()
	if err != nil {
		return err
	}

	coins := getCoins(msg.Bond.Abs().Int64())

	switch {
	case msg.Bond.IsPositive():
		// provider is adding to their bond
		if err := k.SendFromAccountToModule(ctx, addr, types.ProviderName, coins); err != nil {
			return err
		}
	case msg.Bond.IsNegative():
		// provider is withdrawing their bond
		// ensure we provider bond is never negative
		if provider.Bond.LT(coins[0].Amount) {
			return sdkerrors.Wrapf(types.ErrInsufficientFunds, "not enough bond to satisfy bond request: %d/%d", coins[0].Amount.Int64(), provider.Bond.Int64())
		}
		if err := k.SendFromModuleToAccount(ctx, types.ProviderName, addr, coins); err != nil {
			return err
		}
	default:
		return fmt.Errorf("dev error: bond is neither positive or negative")
	}
	provider.Bond = provider.Bond.Add(msg.Bond)
	if provider.Bond.IsZero() {
		k.RemoveProvider(ctx, provider.PubKey, provider.Chain)
		k.BondProviderEvent(ctx, provider.Bond, msg)
		return nil
	}
	err = k.SetProvider(ctx, provider)
	if err == nil {
		k.BondProviderEvent(ctx, provider.Bond, msg)
	}
	return err
}
