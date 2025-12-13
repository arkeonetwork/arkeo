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

	"encoding/base64"
	"math/big"
)

var secpN, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
var secpHalfN = new(big.Int).Rsh(secpN, 1)

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
		// Verify with the contract's spender (client) pubkey using preimage "<cid>:<nonce>:<chain_id>".
		pk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, contract.GetSpender().String())
		if err != nil {
			return err
		}

		pre := fmt.Sprintf("%d:%d:%s", msg.ContractId, msg.Nonce, ctx.ChainID())
		digest := sha256.Sum256([]byte(pre))

		if len(msg.Signature) != 64 {
			return errors.Wrap(types.ErrClaimContractIncomeInvalidSignature, "signature must be 64 bytes (r||s)")
		}

		// Unpack r||s (64 bytes, big-endian)
		r := new(big.Int).SetBytes(msg.Signature[:32])
		s := new(big.Int).SetBytes(msg.Signature[32:])
		highS := s.Cmp(secpHalfN) == 1

		ctx.Logger().Info("claim signature verification debug",
			"preimage", pre,
			"digest_hex", fmt.Sprintf("%x", digest[:]),
			"signature_len", len(msg.Signature),
			"r_hex", fmt.Sprintf("%064x", r),
			"s_hex", fmt.Sprintf("%064x", s),
			"s_high", highS,
		)

		sigHex := fmt.Sprintf("%064x%064x", r, s)
		ctx.Logger().Info("claim sig hex (r||s)",
			"contract_id", msg.ContractId,
			"nonce", msg.Nonce,
			"sig_hex", sigHex,
		)

		pkB64 := base64.StdEncoding.EncodeToString(pk.Bytes())
		sigHexFull := fmt.Sprintf("%064x%064x", r, s)
		ctx.Logger().Info("claim sig verify inputs (keeper)",
			"contract_id", msg.ContractId,
			"nonce", msg.Nonce,
			"preimage", pre,
			"digest_hex", fmt.Sprintf("%x", digest[:]),
			"pk_b64", pkB64,
			"sig_hex", sigHexFull,
		)

		preNoChain := fmt.Sprintf("%d:%d:", msg.ContractId, msg.Nonce)

		// Try multiple verification paths for compatibility:
		// 1) raw preimage with chain-id
		// 2) sha256(preimage with chain-id)
		// 3) raw preimage without chain-id
		// 4) sha256(preimage without chain-id)
		ok := pk.VerifySignature([]byte(pre), msg.Signature) ||
			pk.VerifySignature(digest[:], msg.Signature) ||
			pk.VerifySignature([]byte(preNoChain), msg.Signature) ||
			pk.VerifySignature(sha256.Sum256([]byte(preNoChain))[:], msg.Signature)

		if !ok && highS {
			// normalize to low-S for dev/local testing only
			s.Sub(secpN, s)
			rb := r.FillBytes(make([]byte, 32))
			sb := s.FillBytes(make([]byte, 32))
			norm := append(rb, sb...)
			ctx.Logger().Info("claim sig normalized to low-S", "nonce", msg.Nonce)
			ok = pk.VerifySignature([]byte(pre), norm) ||
				pk.VerifySignature(digest[:], norm) ||
				pk.VerifySignature([]byte(preNoChain), norm) ||
				pk.VerifySignature(sha256.Sum256([]byte(preNoChain))[:], norm)
			if ok {
				ctx.Logger().Info("claim sig normalized verification succeeded",
					"contract_id", msg.ContractId,
					"nonce", msg.Nonce,
					"preimage", pre,
					"digest_hex", fmt.Sprintf("%x", digest[:]),
					"r_hex", fmt.Sprintf("%064x", r),
					"s_hex", fmt.Sprintf("%064x", s),
					"s_high", true,
					"normalized", true,
				)
			}
		}

		if !ok {
			ctx.Logger().Error("claim signature verify failed",
				"contract_id", msg.ContractId,
				"nonce", msg.Nonce,
				"spender", contract.GetSpender().String(),
				"preimage", pre,
				"digest_hex", fmt.Sprintf("%x", digest[:]),
				"pk_b64", pkB64,
				"sig_hex", sigHexFull,
				"s_high", highS,
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
