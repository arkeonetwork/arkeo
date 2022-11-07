package keeper

import (
	"context"
	"mercury/common"
	"mercury/common/cosmos"
	"mercury/x/mercury/configs"
	"mercury/x/mercury/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) OpenContract(goCtx context.Context, msg *types.MsgOpenContract) (*types.MsgOpenContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.Logger().Info(
		"receive MsgOpenContract",
		"pubkey", msg.PubKey,
		"chain", msg.Chain,
		"client", msg.Client,
		"delegate", msg.Delegate,
		"contract type", msg.CType,
		"duration", msg.Duration,
		"rate", msg.Rate,
	)

	if err := k.OpenContractValidate(ctx, msg); err != nil {
		return nil, err
	}

	if err := k.OpenContractHandle(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgOpenContractResponse{}, nil
}

func (k msgServer) OpenContractValidate(ctx cosmos.Context, msg *types.MsgOpenContract) error {
	if k.FetchConfig(ctx, configs.HandlerOpenContract) > 0 {
		return sdkerrors.Wrapf(types.ErrDisabledHandler, "open contract")
	}

	chain, err := common.NewChain(msg.Chain)
	if err != nil {
		return err
	}
	provider, err := k.GetProvider(ctx, msg.PubKey, chain)
	if err != nil {
		return err
	}

	minBond := k.FetchConfig(ctx, configs.MinProviderBond)
	if provider.Bond.LT(cosmos.NewInt(minBond)) {
		return sdkerrors.Wrapf(types.ErrInvalidBond, "not enough provider bond to open a contract (%d/%d)", provider.Bond.Int64(), minBond)
	}

	if provider.Status != types.ProviderStatus_Online {
		return sdkerrors.Wrapf(types.ErrOpenContractBadProviderStatus, "has status %s", provider.Status.String())
	}

	if msg.Duration > provider.MaxContractDuration {
		return sdkerrors.Wrapf(types.ErrOpenContractDuration, "duration exceeds allowed maximum duration from provider")
	}

	if msg.Duration < provider.MinContractDuration {
		return sdkerrors.Wrapf(types.ErrOpenContractDuration, "duration below allowed minimum duration from provider")
	}

	switch msg.CType {
	case types.ContractType_Subscription:
		if msg.Rate != provider.SubscriptionRate {
			return sdkerrors.Wrapf(types.ErrOpenContractMismatchRate, "subscription %d (client) vs %d (provider)", msg.Rate, provider.SubscriptionRate)
		}
		if !cosmos.NewInt(msg.Rate * msg.Duration).Equal(msg.Deposit) {
			return sdkerrors.Wrapf(types.ErrOpenContractMismatchRate, "mismatch of rate*duration and deposit: %d * %d != %d", msg.Rate, msg.Duration, msg.Deposit.Int64())
		}
	case types.ContractType_PayAsYouGo:
		if msg.Rate != provider.SubscriptionRate {
			return sdkerrors.Wrapf(types.ErrOpenContractMismatchRate, "pay-as-you-go %d (client) vs %d (provider)", msg.Rate, provider.PayAsYouGoRate)
		}
	default:
		return sdkerrors.Wrapf(types.ErrInvalidContractType, "%s", msg.CType.String())
	}

	contract, err := k.GetContract(ctx, msg.PubKey, chain, msg.FetchSpender())
	if err != nil {
		return err
	}

	if contract.IsOpen(ctx.BlockHeight()) {
		return sdkerrors.Wrapf(types.ErrOpenContractAlreadyOpen, "expires in %d blocks", ctx.BlockHeight()-contract.Expiration())
	}

	return nil
}

func (k msgServer) OpenContractHandle(ctx cosmos.Context, msg *types.MsgOpenContract) error {
	openCost := k.FetchConfig(ctx, configs.OpenContractCost)
	if !openCost.IsZero() {
		if err := k.SendFromAccountToModule(ctx, msg.MustGetSigner(), types.ReserveName, cosmos.NewCoins(getCoin(openCost))); err != nil {
			return nil
		}
	}

	if err := k.SendFromAccountToModule(ctx, msg.MustGetSigner(), types.ContractName, getCoins(msg.Deposit)); err != nil {
		return nil
	}

	chain, err := common.NewChain(msg.Chain)
	if err != nil {
		return err
	}
	contract := types.NewContract(msg.PubKey, chain, msg.FetchSpender())
	contract.Client = msg.Client
	contract.Type = msg.CType
	contract.Height = ctx.BlockHeight()
	contract.Duration = msg.Duration
	contract.Rate = msg.Rate
	contract.Deposit = msg.Deposit

	exp := types.NewContractExpiration(msg.PubKey, chain, msg.FetchSpender())
	set, err := k.GetContractExpirationSet(ctx, contract.Expiration())
	if err != nil {
		return err
	}
	set.Contracts = append(set.Contracts, exp)
	err = k.SetContractExpirationSet(ctx, set)
	if err != nil {
		return err
	}

	err = k.SetContract(ctx, contract)
	if err != nil {
		return err
	}

	k.OpenContractEvent(ctx, openCost, confcontract)
	return nil
}
