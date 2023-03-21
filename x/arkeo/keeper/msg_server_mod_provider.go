package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func (k msgServer) ModProvider(goCtx context.Context, msg *types.MsgModProvider) (*types.MsgModProviderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.Logger().Info(
		"receive MsgModProvider",
		"proivder", msg.Provider,
		"service", msg.Service,
		"metatadata uri", msg.MetadataUri,
		"metadata nonce", msg.MetadataNonce,
		"status", msg.Status,
		"min contract duration", msg.MinContractDuration,
		"max contract duration", msg.MaxContractDuration,
		"subscription rate", msg.SubscriptionRate,
		"pay-as-you-go rate", msg.PayAsYouGoRate,
		"settlement duration", msg.SettlementDuration,
		"support-pay-as-you-go", msg.SupportPayAsYouGo,
	)

	cacheCtx, commit := ctx.CacheContext()
	if err := k.ModProviderValidate(cacheCtx, msg); err != nil {
		ctx.Logger().Error("failed mod provider validation", "err", err)
		return nil, err
	}

	if err := k.ModProviderHandle(cacheCtx, msg); err != nil {
		ctx.Logger().Error("failed mod provider handle", "err", err)
		return nil, err
	}

	commit()

	return &types.MsgModProviderResponse{}, nil
}

func (k msgServer) ModProviderValidate(ctx cosmos.Context, msg *types.MsgModProvider) error {
	if k.FetchConfig(ctx, configs.HandlerModProvider) > 0 {
		return errors.Wrapf(types.ErrDisabledHandler, "mod provider")
	}
	maxContractDuration := k.FetchConfig(ctx, configs.MaxContractLength)
	if maxContractDuration > 0 {
		if msg.MaxContractDuration > maxContractDuration {
			return errors.Wrapf(types.ErrInvalidModProviderMaxContractDuration, "max contract duration is too long (%d/%d)", msg.MaxContractDuration, maxContractDuration)
		}
		if msg.MinContractDuration > maxContractDuration {
			return errors.Wrapf(types.ErrInvalidModProviderMinContractDuration, "min contract duration is too long (%d/%d)", msg.MaxContractDuration, maxContractDuration)
		}
	}

	service, err := common.NewService(msg.Service)
	if err != nil {
		return err
	}
	provider, err := k.GetProvider(ctx, msg.Provider, service)
	if err != nil {
		return err
	}
	if provider.Bond.IsZero() {
		return errors.Wrapf(types.ErrInvalidModProviderNoBond, "bond cannot be zero")
	}

	return nil
}

func (k msgServer) ModProviderHandle(ctx cosmos.Context, msg *types.MsgModProvider) error {
	service, err := common.NewService(msg.Service)
	if err != nil {
		return err
	}
	provider, err := k.GetProvider(ctx, msg.Provider, service)
	if err != nil {
		return err
	}

	// update metadata URI
	if len(msg.MetadataUri) > 0 {
		provider.MetadataUri = msg.MetadataUri
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
	provider.SettlementDuration = msg.SettlementDuration
	provider.SupportPayAsYouGo = msg.SupportPayAsYouGo
	provider.LastUpdate = ctx.BlockHeight()

	if err := k.SetProvider(ctx, provider); err != nil {
		return err
	}
	k.ModProviderEvent(ctx, &provider)
	return nil
}
