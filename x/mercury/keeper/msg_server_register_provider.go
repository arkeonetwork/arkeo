package keeper

import (
	"context"
	"mercury/common/cosmos"
	"mercury/x/mercury/configs"
	"mercury/x/mercury/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) RegisterProvider(goCtx context.Context, msg *types.MsgRegisterProvider) (*types.MsgRegisterProviderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.RegisterProviderValidate(ctx, msg); err != nil {
		return nil, err
	}

	if err := k.RegisterProviderHandle(ctx, msg); err != nil {
		return nil, err
	}

	return nil, nil
}

func (k msgServer) RegisterProviderValidate(ctx cosmos.Context, msg *types.MsgRegisterProvider) error {
	// Verify signer is the same as the bidder (this validates both the bidder and signer addresses)
	signer := msg.MustGetSigner()
	provider, err := msg.Pubkey.GetMyAddress()
	if err != nil {
		return err
	}

	if !signer.Equals(provider) {
		return sdkerrors.Wrapf(types.ErrProviderBadSigner, "Signer: %s, Provider Address: %s", msg.GetSigners(), provider)
	}

	if k.ProviderExists(ctx, msg.Pubkey, msg.Chain) {
		return sdkerrors.Wrapf(types.ErrProviderAlreadyExists, "Provider already exists: %s", msg.Pubkey)
	}

	if err := k.hasCoins(ctx, provider, configs.GasFee, configs.RegisterProviderFee); err != nil {
		return err
	}

	return nil
}

func (k msgServer) RegisterProviderHandle(ctx cosmos.Context, msg *types.MsgRegisterProvider) error {
	// pay the fee
	addr, err := msg.Pubkey.GetMyAddress()
	if err != nil {
		return err
	}

	fee := getCoins(k.getFee(ctx, configs.GasFee, configs.RegisterProviderFee))
	if err := k.SendFromAccountToModule(ctx, addr, types.ReserveName, fee); err != nil {
		return err
	}

	provider := types.NewProvider(msg.Pubkey, msg.Chain)
	return k.SetProvider(ctx, provider)
}
