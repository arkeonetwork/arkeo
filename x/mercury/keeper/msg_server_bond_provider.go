package keeper

import (
	"context"
	"fmt"
	"mercury/common"
	"mercury/common/cosmos"
	"mercury/x/mercury/configs"
	"mercury/x/mercury/types"

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
	if err := k.BondProviderValidate(ctx, msg); err != nil {
		return nil, err
	}

	if err := k.BondProviderHandle(ctx, msg); err != nil {
		return nil, err
	}

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

	coins := cosmos.NewCoins(cosmos.NewCoin(configs.Denom, msg.Bond.Abs()))

	switch {
	case msg.Bond.IsPositive():
		// provider is adding to their bond
		if err := k.SendFromAccountToModule(ctx, addr, types.ProviderName, coins); err != nil {
			return err
		}
	case msg.Bond.IsNegative():
		// provider is withdrawing their bond
		// ensure we provider bond is never negative
		if provider.Bond.LT(coin.Amount) {
			return sdkerrors.Wrapf(types.ErrInsufficientFunds, "not enough bond to satisfy bond request: %d/%d", coin.Amount.Int64(), provider.Bond.Int64())
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
