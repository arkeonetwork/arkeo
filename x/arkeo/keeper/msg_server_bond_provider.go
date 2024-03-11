package keeper

import (
	"context"
	"fmt"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) BondProvider(goCtx context.Context, msg *types.MsgBondProvider) (*types.MsgBondProviderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.Logger().Info(
		"receive MsgBondProvider",
		"provider", msg.Provider,
		"service", msg.Service,
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
		return errors.Wrapf(types.ErrDisabledHandler, "bond provider")
	}

	// We allow providers to unbond WHILE active contracts are underway. This
	// is because A) users can cancel their owned contracts at any time, and B)
	// this is the way the provider signals to the service that they don't want
	// to open any new contracts (as there is a min bond requirement for new
	// contracts to be opened)

	return nil
}

func (k msgServer) BondProviderHandle(ctx cosmos.Context, msg *types.MsgBondProvider) error {
	service, err := common.NewService(msg.Service)
	if err != nil {
		return err
	}
	pk, _ := common.NewPubKey(msg.Provider)
	provider, err := k.GetProvider(ctx, pk, service)
	if err != nil {
		return err
	}
	addr, err := pk.GetMyAddress()
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
			return errors.Wrapf(types.ErrInsufficientFunds, "not enough bond to satisfy bond request: %d/%d", coins[0].Amount.Int64(), provider.Bond.Int64())
		}
		if err := k.SendFromModuleToAccount(ctx, types.ProviderName, addr, coins); err != nil {
			return err
		}
	default:
		return fmt.Errorf("dev error: bond is neither positive or negative")
	}
	provider.Bond = provider.Bond.Add(msg.Bond)
	if provider.Bond.IsZero() {
		k.RemoveProvider(ctx, provider.PubKey, provider.Service)
		return k.EmitBondProviderEvent(ctx, provider.Bond, msg)
	}

	provider.LastUpdate = ctx.BlockHeight()

	err = k.SetProvider(ctx, provider)
	if err == nil {
		return k.EmitBondProviderEvent(ctx, provider.Bond, msg)
	}
	return err
}
