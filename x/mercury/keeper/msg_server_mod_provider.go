package keeper

import (
	"context"
	"mercury/common/cosmos"
	"mercury/x/mercury/configs"
	"mercury/x/mercury/types"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) ModProvider(goCtx context.Context, msg *types.MsgModProvider) (*types.MsgModProviderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.Logger().Info(
		"receive MsgModProvider",
		"pubkey", msg.PubKey,
		"chain", msg.Chain,
		"metatadata uri", msg.MetadataURI,
		"metadata nonce", msg.MetadataNonce,
		"status", msg.Status,
		"min contract duration", msg.MinContractDuration,
		"max contract duration", msg.MaxContractDuration,
		"subscription rate", msg.SubscriptionRate,
		"pay-as-you-go rate", msg.PayAsYouGoRate,
	)
	if err := k.ModProviderValidate(ctx, msg); err != nil {
		return nil, err
	}

	if err := k.ModProviderHandle(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgModProviderResponse{}, nil
}

func (k msgServer) ModProviderValidate(ctx cosmos.Context, msg *types.MsgModProvider) error {
	if k.FetchConfig(ctx, configs.HandlerModProvider) > 0 {
		return sdkerrors.Wrapf(types.ErrDisabledHandler, "mod provider")
	}
	maxContractDuration := k.FetchConfig(ctx, configs.MaxContractLength)
	if maxContractDuration > 0 {
		if msg.MaxContractDuration > maxContractDuration {
			return sdkerrors.Wrapf(types.ErrInvalidModProviderMaxContractDuration, "max contract duration is too long (%d/%d)", msg.MaxContractDuration, maxContractDuration)
		}
		if msg.MinContractDuration > maxContractDuration {
			return sdkerrors.Wrapf(types.ErrInvalidModProviderMinContractDuration, "min contract duration is too long (%d/%d)", msg.MaxContractDuration, maxContractDuration)
		}
	}

	provider, err := k.GetProvider(ctx, msg.PubKey, msg.Chain)
	if err != nil {
		return err
	}
	if provider.Bond.IsZero() {
		return sdkerrors.Wrapf(types.ErrInvalidModProviderNoBond, "bond cannot be zero")
	}

	return nil
}

func (k msgServer) ModProviderHandle(ctx cosmos.Context, msg *types.MsgModProvider) error {
	provider, err := k.GetProvider(ctx, msg.PubKey, msg.Chain)
	if err != nil {
		return err
	}

	// update metadata URI
	if len(msg.MetadataURI) > 0 {
		provider.MetadataURI = msg.MetadataURI
	}

	// update metadata nonce
	if provider.MetadataNonce < msg.MetadataNonce {
		provider.MetadataNonce = msg.MetadataNonce
	}

	// update status
	provider.Status = msg.Status

	// update contract durations
	provider.MinContractDuration = msg.MinContractDuration
	provider.MaxContractDuration = msg.MaxContractDuration

	// update contract rates
	provider.SubscriptionRate = msg.SubscriptionRate
	provider.PayAsYouGoRate = msg.PayAsYouGoRate

	provider.LastUpdate = ctx.BlockHeight()

	if err := k.SetProvider(ctx, provider); err != nil {
		return err
	}
	k.ModProviderEvent(ctx, provider)
	return nil
}

func (k msgServer) ModProviderEvent(ctx cosmos.Context, provider types.Provider) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeProviderMod,
				sdk.NewAttribute("pubkey", provider.PubKey.String()),
				sdk.NewAttribute("chain", provider.Chain.String()),
				sdk.NewAttribute("metadata_uri", provider.MetadataURI),
				sdk.NewAttribute("metadata_nonce", strconv.FormatUint(provider.MetadataNonce, 10)),
				sdk.NewAttribute("status", provider.Status.String()),
				sdk.NewAttribute("min_contract_duration", strconv.FormatInt(provider.MinContractDuration, 10)),
				sdk.NewAttribute("max_contract_duration", strconv.FormatInt(provider.MaxContractDuration, 10)),
				sdk.NewAttribute("subscription_rate", strconv.FormatInt(provider.SubscriptionRate, 10)),
				sdk.NewAttribute("pay-as-you-go_rate", strconv.FormatInt(provider.PayAsYouGoRate, 10)),
			),
		},
	)
}
