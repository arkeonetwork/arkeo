package keeper

import (
	"context"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CloseContract(goCtx context.Context, msg *types.MsgCloseContract) (*types.MsgCloseContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.Logger().Info(
		"receive MsgCloseContract",
		"contract_id", msg.ContractId,
	)

	cacheCtx, commit := ctx.CacheContext()
	if err := k.CloseContractValidate(cacheCtx, msg); err != nil {
		ctx.Logger().Error("failed close contract validation", "err", err)
		return nil, err
	}

	if err := k.CloseContractHandle(cacheCtx, msg); err != nil {
		ctx.Logger().Error("failed close contract handler", "err", err)
		return nil, err
	}

	commit()
	return &types.MsgCloseContractResponse{}, nil
}

func (k msgServer) CloseContractValidate(ctx cosmos.Context, msg *types.MsgCloseContract) error {
	if k.FetchConfig(ctx, configs.HandlerCloseContract) > 0 {
		return errors.Wrapf(types.ErrDisabledHandler, "close contract")
	}

	contract, err := k.GetContract(ctx, msg.ContractId)
	if err != nil {
		return err
	}

	if contract.IsEmpty() {
		return errors.Wrapf(types.ErrContractNotFound, "id: %d", msg.ContractId)
	}

	signerAccountAddress := msg.MustGetSigner()

	clientPublicKey, err := common.NewPubKey(contract.Client.String())
	if err != nil {
		return err
	}

	clientAccountAddress, err := clientPublicKey.GetMyAddress()
	if err != nil {
		return err
	}

	if !signerAccountAddress.Equals(clientAccountAddress) {
		return errors.Wrapf(types.ErrCloseContractUnauthorized, "only the client can close the contract")
	}

	if contract.IsExpired(ctx.BlockHeight()) {
		return errors.Wrapf(types.ErrCloseContractAlreadyClosed, "closed %d", contract.Expiration())
	}

	return nil
}

func (k msgServer) CloseContractHandle(ctx cosmos.Context, msg *types.MsgCloseContract) error {
	contract, err := k.GetContract(ctx, msg.ContractId)
	if err != nil {
		return err
	}

	if contract.IsPayAsYouGo() {
		// add a new expiration return deposit to user
		newHeight := ctx.BlockHeight() + contract.SettlementDuration
		expirationSet, err := k.GetContractExpirationSet(ctx, newHeight)
		if err != nil {
			return err
		}
		expirationSet.Append(contract.Id)
		err = k.SetContractExpirationSet(ctx, expirationSet)
		if err != nil {
			return err
		}
	}

	_, err = k.mgr.SettleContract(ctx, contract, 0, contract.IsSubscription())
	if err != nil {
		return err
	}

	return k.EmitCloseContractEvent(ctx, &contract)
}
