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
	)

	cacheCtx, commit := ctx.CacheContext()
	if err := k.HandlerClaimContractIncome(cacheCtx, msg); err != nil {
		ctx.Logger().Error("failed to handle claim contract income", "err", err)
		return nil, err
	}
	commit()

	return &types.MsgClaimContractIncomeResponse{}, nil
}

func (k msgServer) HandlerClaimContractIncome(ctx cosmos.Context, msg *types.MsgClaimContractIncome) error {
	// validate contract
	if k.FetchConfig(ctx, configs.HandlerClaimContractIncome) > 0 {
		return errors.Wrapf(types.ErrDisabledHandler, "Claim Contract Income")
	}

	contract, err := k.GetContract(ctx, msg.ContractId)

	if err != nil {
		return err
	}

	if msg.ChainId != ctx.ChainID() {
		return errors.Wrap(types.ErrInvalidChainId, "chain id mismatch")
	}

	if contract.Nonce >= msg.Nonce {
		return errors.Wrapf(types.ErrClaimContractIncomeBadNonce, "contract nonce (%d) is greater than msg nonce (%d)", contract.Nonce, msg.Nonce)
	}

	if contract.IsSettled(ctx.BlockHeight()) {
		return errors.Wrapf(types.ErrClaimContractIncomeClosed, "settled on block: %d", contract.SettlementPeriodEnd())
	}

	if msg.SignatureExpiresAt <= ctx.BlockHeight() {
		return errors.Wrapf(types.ErrSignatureExpired, "signature expired at block: (%d)", ctx.BlockHeight())
	}

	// open subscription contracts do NOT need to verify the signature
	if !(contract.IsSubscription() && contract.IsOpenAuthorization()) {
		pk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, contract.GetSpender().String())
		if err != nil {
			return err
		}
		if !pk.VerifySignature(msg.GetBytesToSign(), msg.Signature) {
			return errors.Wrap(types.ErrClaimContractIncomeInvalidSignature, "signature mismatch")
		}
	}

	// excute settlement
	_, err = k.mgr.SettleContract(ctx, contract, msg.Nonce, false)
	if err != nil {
		return err
	}
	return nil
}
