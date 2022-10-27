package keeper

import (
	"context"

	"mercury/common/cosmos"
	"mercury/x/mercury/configs"
	"mercury/x/mercury/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CloseContract(goCtx context.Context, msg *types.MsgCloseContract) (*types.MsgCloseContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.Logger().Info(
		"receive MsgCloseContract",
		"pubkey", msg.PubKey,
		"chain", msg.Chain,
		"client", msg.Client,
	)

	if err := k.CloseContractValidate(ctx, msg); err != nil {
		return nil, err
	}

	if err := k.CloseContractHandle(ctx, msg); err != nil {
		return nil, err
	}
	return &types.MsgCloseContractResponse{}, nil
}

func (k msgServer) CloseContractValidate(ctx cosmos.Context, msg *types.MsgCloseContract) error {
	if k.FetchConfig(ctx, configs.HandlerCloseContract) > 0 {
		return sdkerrors.Wrapf(types.ErrDisabledHandler, "close contract")
	}

	client, err := msg.GetClientAddress()
	if err != nil {
		return err
	}

	contract, err := k.GetContract(ctx, msg.PubKey, msg.Chain, client)
	if err != nil {
		return err
	}

	if contract.IsClose(ctx.BlockHeight()) {
		return sdkerrors.Wrapf(types.ErrCloseContractAlreadyClosed, "closed %d", contract.Expiration())
	}

	return nil
}

func (k msgServer) CloseContractHandle(ctx cosmos.Context, msg *types.MsgCloseContract) error {
	client, err := msg.GetClientAddress()
	if err != nil {
		return err
	}

	contract, err := k.GetContract(ctx, msg.PubKey, msg.Chain, client)
	if err != nil {
		return err
	}

	_, err = k.SettleContract(ctx, contract, true)
	if err != nil {
		return err
	}

	k.CloseContractEvent(ctx, msg)
	return nil
}

func (k msgServer) CloseContractEvent(ctx cosmos.Context, msg *types.MsgCloseContract) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeCloseContract,
				sdk.NewAttribute("pubkey", msg.PubKey.String()),
				sdk.NewAttribute("chain", msg.Chain.String()),
				sdk.NewAttribute("client", msg.Client),
			),
		},
	)
}
