package keeper

import (
	"arkeo/common"
	"arkeo/common/cosmos"
	"arkeo/x/arkeo/configs"
	"arkeo/x/arkeo/types"
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) ClaimContractIncome(goCtx context.Context, msg *types.MsgClaimContractIncome) (*types.MsgClaimContractIncomeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.Logger().Info(
		"receive MsgClaimContractIncome",
		"pubkey", msg.PubKey,
		"chain", msg.Chain,
		"spender", msg.Spender,
		"nonce", msg.Nonce,
		"height", msg.Height,
	)

	cacheCtx, commit := ctx.CacheContext()
	if err := k.ClaimContractIncomeValidate(cacheCtx, msg); err != nil {
		return nil, err
	}

	if err := k.ClaimContractIncomeHandle(cacheCtx, msg); err != nil {
		return nil, err
	}

	commit()
	return &types.MsgClaimContractIncomeResponse{}, nil
}

func (k msgServer) ClaimContractIncomeValidate(ctx cosmos.Context, msg *types.MsgClaimContractIncome) error {
	if k.FetchConfig(ctx, configs.HandlerCloseContract) > 0 {
		return sdkerrors.Wrapf(types.ErrDisabledHandler, "close contract")
	}

	chain, err := common.NewChain(msg.Chain)
	if err != nil {
		return err
	}
	contract, err := k.GetContract(ctx, msg.PubKey, chain, msg.Spender)
	if err != nil {
		return err
	}

	if contract.Height != msg.Height {
		return sdkerrors.Wrapf(types.ErrClaimContractIncomeBadHeight, "contract height (%d) doesn't match msg height (%d)", contract.Height, msg.Height)
	}

	if contract.Nonce >= msg.Nonce {
		return sdkerrors.Wrapf(types.ErrClaimContractIncomeBadNonce, "contract nonce (%d) is greater than msg nonce (%d)", contract.Nonce, msg.Nonce)
	}

	if contract.IsClose(ctx.BlockHeight()) {
		return sdkerrors.Wrapf(types.ErrClaimContractIncomeClosed, "closed %d", contract.Expiration())
	}

	return nil
}

func (k msgServer) ClaimContractIncomeHandle(ctx cosmos.Context, msg *types.MsgClaimContractIncome) error {
	chain, err := common.NewChain(msg.Chain)
	if err != nil {
		return err
	}
	contract, err := k.GetContract(ctx, msg.PubKey, chain, msg.Spender)
	if err != nil {
		return err
	}

	_, err = k.mgr.SettleContract(ctx, contract, msg.Nonce, false)
	return err
}
