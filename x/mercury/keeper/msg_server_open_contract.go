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

func (k msgServer) OpenContract(goCtx context.Context, msg *types.MsgOpenContract) (*types.MsgOpenContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.Logger().Info(
		"receive MsgOpenContract",
		"pubkey", msg.PubKey,
		"chain", msg.Chain,
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

	provider, err := k.GetProvider(ctx, msg.PubKey, msg.Chain)
	if err != nil {
		return err
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

	minBond := k.FetchConfig(ctx, configs.MinProviderBond)
	if provider.Bond.LT(cosmos.NewInt(minBond)) {
		return sdkerrors.Wrapf(types.ErrInvalidBond, "not enough provider bond to open a contract (%d/%d)", provider.Bond.Int64(), minBond)
	}

	contract, err := k.GetContract(ctx, msg.PubKey, msg.Chain, msg.MustGetSigner())
	if err != nil {
		return err
	}

	if contract.IsOpen(ctx.BlockHeight()) {
		return sdkerrors.Wrapf(types.ErrOpenContractAlreadyOpen, "expires in %d blocks", ctx.BlockHeight()-contract.Expiration())
	}

	return nil
}

func (k msgServer) OpenContractHandle(ctx cosmos.Context, msg *types.MsgOpenContract) error {
	cost := getCoin(k.FetchConfig(ctx, configs.OpenContractCost)).AddAmount(msg.Deposit)
	if !cost.IsZero() {
		if err := k.SendFromAccountToModule(ctx, msg.MustGetSigner(), types.ReserveName, cosmos.NewCoins(cost)); err != nil {
			return nil
		}
	}

	contract := types.NewContract(msg.PubKey, msg.Chain, msg.MustGetSigner())
	contract.Type = msg.CType
	contract.Height = ctx.BlockHeight()
	contract.Duration = msg.Duration
	contract.Rate = msg.Rate
	contract.Deposit = msg.Deposit

	err := k.SetContract(ctx, contract)
	if err != nil {
		return err
	}

	k.OpenContractEvent(ctx, contract)
	return nil
}

func (k msgServer) OpenContractEvent(ctx cosmos.Context, contract types.Contract) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeOpenContract,
				sdk.NewAttribute("pubkey", contract.ProviderPubKey.String()),
				sdk.NewAttribute("chain", contract.Chain.String()),
				sdk.NewAttribute("client", contract.ClientAddress.String()),
				sdk.NewAttribute("type", contract.Type.String()),
				sdk.NewAttribute("height", strconv.FormatInt(contract.Height, 10)),
				sdk.NewAttribute("duration", strconv.FormatInt(contract.Duration, 10)),
				sdk.NewAttribute("rate", strconv.FormatInt(contract.Rate, 10)),
			),
		},
	)
}
