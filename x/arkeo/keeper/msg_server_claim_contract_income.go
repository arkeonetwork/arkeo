package keeper

import (
	"context"
	"crypto/sha256"
	"fmt"

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

	if contract.Nonce >= msg.Nonce {
		return errors.Wrapf(types.ErrClaimContractIncomeBadNonce, "contract nonce (%d) is greater than msg nonce (%d)", contract.Nonce, msg.Nonce)
	}

	if contract.IsSettled(ctx.BlockHeight()) {
		return errors.Wrapf(types.ErrClaimContractIncomeClosed, "settled on block: %d", contract.SettlementPeriodEnd())
	}

	// open subscription contracts do NOT need to verify the signature
	if !(contract.IsSubscription() && contract.IsOpenAuthorization()) {
		// Verify with the contract's spender (client) pubkey using SHA-256("<cid>:<nonce>:")
		pk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, contract.GetSpender().String())
		if err != nil {
			return err
		}
		pre := fmt.Sprintf("%d:%d:", msg.ContractId, msg.Nonce)
		digest := sha256.Sum256([]byte(pre))

		ctx.Logger().Info("claim signature verification debug",
			"preimage", pre,
			"digest_hex", fmt.Sprintf("%x", digest[:]),
			"signature_len", len(msg.Signature),
		)

		if !pk.VerifySignature(digest[:], msg.Signature) {
			ctx.Logger().Error("claim signature verify failed",
				"contract_id", msg.ContractId,
				"nonce", msg.Nonce,
				"spender", contract.GetSpender().String(),
			)
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
