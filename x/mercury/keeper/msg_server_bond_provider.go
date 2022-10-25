package keeper

import (
	"context"
	"fmt"
	"mercury/common/cosmos"
	"mercury/x/mercury/configs"
	"mercury/x/mercury/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) BondProvider(goCtx context.Context, msg *types.MsgBondProvider) (*types.MsgBondProviderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.BondProviderValidate(ctx, msg); err != nil {
		return nil, err
	}

	if err := k.BondProviderHandle(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgBondProviderResponse{}, nil
}

func (k msgServer) BondProviderValidate(ctx cosmos.Context, msg *types.MsgBondProvider) error {
	return nil
}

func (k msgServer) BondProviderHandle(ctx cosmos.Context, msg *types.MsgBondProvider) error {
	provider, err := k.GetProvider(ctx, msg.PubKey, msg.Chain)
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
		if err := k.SendFromAccountToModule(ctx, addr, types.ProviderName, coins); err != nil {
			return err
		}
	case msg.Bond.IsNegative():
		// ensure we provider bond is never negative
		if provider.Bond.LT(msg.Bond.Abs()) {
			return sdkerrors.Wrapf(types.ErrInsufficientFunds, "not enough bond to satisfy bond request: %d/%d", msg.Bond.Int64(), provider.Bond.Int64())
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

func (k msgServer) BondProviderEvent(ctx cosmos.Context, bond cosmos.Int, msg *types.MsgBondProvider) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeProviderBond,
				sdk.NewAttribute("pubkey", msg.PubKey.String()),
				sdk.NewAttribute("chain", msg.Chain.String()),
				sdk.NewAttribute("bond_rel", msg.Bond.String()),
				sdk.NewAttribute("bond_abs", bond.String()),
			),
		},
	)
}
