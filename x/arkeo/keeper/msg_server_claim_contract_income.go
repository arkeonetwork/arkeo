package keeper

import (
	"context"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ClaimContractIncome(goCtx context.Context, msg *types.MsgClaimContractIncome) (*types.MsgClaimContractIncomeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.Logger().Info(
		"receive MsgClaimContractIncome",
		"contract_id", msg.ContractId,
		"nonce", msg.Nonce,
		"height", msg.Height,
	)

	cacheCtx, commit := ctx.CacheContext()
	if err := k.ClaimContractIncomeValidate(cacheCtx, msg); err != nil {
		ctx.Logger().Error("failed claim contract validation", "err", err)
		return nil, err
	}

	if err := k.ClaimContractIncomeHandle(ctx, msg); err != nil {
		ctx.Logger().Error("failed claim contract handler", "err", err)
		return nil, err
	}
	commit()

	return &types.MsgClaimContractIncomeResponse{}, nil
}

func (k msgServer) ClaimContractIncomeValidate(ctx cosmos.Context, msg *types.MsgClaimContractIncome) error {
	if k.FetchConfig(ctx, configs.HandlerCloseContract) > 0 {
		return errors.Wrapf(types.ErrDisabledHandler, "close contract")
	}

	contract, err := k.GetContract(ctx, msg.ContractId)
	if err != nil {
		return err
	}

	if contract.Height != msg.Height {
		return errors.Wrapf(types.ErrClaimContractIncomeBadHeight, "contract height (%d) doesn't match msg height (%d)", contract.Height, msg.Height)
	}

	if contract.Nonce >= msg.Nonce {
		return errors.Wrapf(types.ErrClaimContractIncomeBadNonce, "contract nonce (%d) is greater than msg nonce (%d)", contract.Nonce, msg.Nonce)
	}

	if contract.IsClose(ctx.BlockHeight()) {
		return errors.Wrapf(types.ErrClaimContractIncomeClosed, "closed %d", contract.Expiration())
	}

	return nil
}

func (k msgServer) ClaimContractIncomeHandle(ctx cosmos.Context, msg *types.MsgClaimContractIncome) error {
	contract, err := k.GetContract(ctx, msg.ContractId)
	if err != nil {
		return err
	}

	_, err = k.mgr.SettleContract(ctx, contract, msg.Nonce, false)
	return err
}
