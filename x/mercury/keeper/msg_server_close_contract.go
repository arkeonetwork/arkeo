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

func (k msgServer) CloseContract(goCtx context.Context, msg *types.MsgCloseContract) (*types.MsgCloseContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.Logger().Info(
		"receive MsgCloseContract",
		"pubkey", msg.PubKey,
		"chain", msg.Chain,
		"client", msg.Client,
	)

	cacheCtx, commit := ctx.CacheContext()
	if err := k.CloseContractValidate(cacheCtx, msg); err != nil {
		return nil, err
	}

	if err := k.CloseContractHandle(cacheCtx, msg); err != nil {
		return nil, err
	}

	commit()
	return &types.MsgCloseContractResponse{}, nil
}

func (k msgServer) CloseContractValidate(ctx cosmos.Context, msg *types.MsgCloseContract) error {
	if k.FetchConfig(ctx, configs.HandlerCloseContract) > 0 {
		return sdkerrors.Wrapf(types.ErrDisabledHandler, "close contract")
	}

	chain, err := common.NewChain(msg.Chain)
	if err != nil {
		return err
	}
	contract, err := k.GetContract(ctx, msg.PubKey, chain, msg.Client)
	if err != nil {
		return err
	}

	if contract.IsClose(ctx.BlockHeight()) {
		return sdkerrors.Wrapf(types.ErrCloseContractAlreadyClosed, "closed %d", contract.Expiration())
	}

	provider, err := contract.ProviderPubKey.GetMyAddress()
	if err != nil {
		return err
	}
	if contract.Type == types.ContractType_PayAsYouGo && !provider.Equals(msg.MustGetSigner()) {
		// clients are not allowed to cancel a pay-as-you-go contract as it
		// could be a way to game providers. IE, the client make 1,000 requests
		// and before the provider can claim the rewards, the client cancels
		// the contract.
		return sdkerrors.Wrapf(types.ErrCloseContractUnauthorized, "client cannot cancel a pay-as-you-go contract")
	}

	return nil
}

func (k msgServer) CloseContractHandle(ctx cosmos.Context, msg *types.MsgCloseContract) error {
	chain, err := common.NewChain(msg.Chain)
	if err != nil {
		return err
	}
	contract, err := k.GetContract(ctx, msg.PubKey, chain, msg.Client)
	if err != nil {
		return err
	}

	_, err = k.mgr.SettleContract(ctx, contract, 0, true)
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
				sdk.NewAttribute("chain", msg.Chain),
				sdk.NewAttribute("client", msg.Client.String()),
			),
		},
	)
}
