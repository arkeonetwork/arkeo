package keeper

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/pkg/errors"

	"github.com/arkeonetwork/arkeo/x/claim/types"
)

func (k msgServer) ClaimEth(goCtx context.Context, msg *types.MsgClaimEth) (*types.MsgClaimEthResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// get eth claim
	ethClaim, err := k.GetClaimRecord(ctx, msg.EthAddress, types.ETHEREUM)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get claim record for %s", msg.EthAddress)
	}

	if ethClaim.IsEmpty() || ethClaim.AmountClaim.IsZero() {
		return nil, errors.Wrapf(types.ErrNoClaimableAmount, "no claimable amount for %s", msg.EthAddress)
	}
	totalAmountClaimable := getInitialClaimableAmountTotal(ethClaim)

	// Store the amounts before we modify the claims
	ethClaimAmount := ethClaim.AmountClaim.Amount.Int64()

	// validate signature
	isValid, err := IsValidClaimSignature(msg.EthAddress, msg.Creator,
		totalAmountClaimable.Amount.String(), msg.Signature)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidSignature, "failed to validate signature for %s", msg.EthAddress)
	}

	if !isValid {
		// this shouldn't happen without an error, but just in case
		return nil, errors.Wrapf(types.ErrInvalidSignature, "failed to validate signature for %s", msg.EthAddress)
	}

	// create new arkeo claim
	arkeoClaim := types.ClaimRecord{
		Address:        msg.Creator,
		Chain:          types.ARKEO,
		AmountClaim:    ethClaim.AmountClaim,
		AmountVote:     ethClaim.AmountVote,
		AmountDelegate: ethClaim.AmountDelegate,
	}

	// set eth claim to completed
	ethClaim = setClaimableAmountForAllActions(ethClaim, sdk.Coin{})
	err = k.SetClaimRecord(ctx, ethClaim)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to set claim record for %s", msg.EthAddress)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaimFromEth,
			sdk.NewAttribute(sdk.AttributeKeySender, strings.ToLower(msg.EthAddress)),
			sdk.NewAttribute(sdk.AttributeKeyAmount, arkeoClaim.AmountClaim.String()),
		),
	})

	// see if there is an existing arkeo claim so we can merge it
	existingArkeoClaim, err := k.GetClaimRecord(ctx, msg.Creator, types.ARKEO)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get arkeo claim record for %s", msg.Creator)
	}

	// Get existing claim amount, defaulting to 0 if no existing claim
	existingArkeoClaimAmount := int64(0)
	if !existingArkeoClaim.IsEmpty() && !existingArkeoClaim.AmountClaim.IsZero() {
		existingArkeoClaimAmount = existingArkeoClaim.AmountClaim.Amount.Int64()
	}

	arkeoClaim, err = mergeClaimRecords(existingArkeoClaim, arkeoClaim)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to merge claim records for %s", msg.Creator)
	}

	err = k.SetClaimRecord(ctx, arkeoClaim)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to set claim record for %s", msg.Creator)
	}

	// call claim on arkeo to claim arkeo (note: this could CLAIM for all tokens that are now merged)
	_, err = k.ClaimCoinsForAction(ctx, msg.Creator, types.ACTION_CLAIM)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to claim coins for %s", msg.Creator)
	}

	return &types.MsgClaimEthResponse{
		EthAddress:       msg.EthAddress,
		ArkeoAddress:     msg.Creator,
		EthClaimAmount:   ethClaimAmount,
		ArkeoClaimAmount: existingArkeoClaimAmount,
	}, nil
}

func GenerateClaimTypedDataBytes(ethAddress, arkeoAddress, amount string) ([]byte, error) {
	claimEthAddress := common.HexToAddress(ethAddress)
	signerTypedData := apitypes.TypedData{
		Types:       types.EIP712Types,
		PrimaryType: "Claim",
		Domain:      types.EIP712Domain,
		Message: apitypes.TypedDataMessage{
			"address":      claimEthAddress.String(),
			"arkeoAddress": arkeoAddress,
			"amount":       amount,
		},
	}
	typedDataHash, err := signerTypedData.HashStruct(signerTypedData.PrimaryType, signerTypedData.Message)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to hash struct")
	}
	domainSeparator, err := signerTypedData.HashStruct("EIP712Domain", signerTypedData.Domain.Map())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to hash domain")
	}

	return []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash))), nil
}

func IsValidClaimSignature(ethAddress, arkeoAdddress, amount, signature string) (bool, error) {
	rawData, err := GenerateClaimTypedDataBytes(ethAddress, arkeoAdddress, amount)
	if err != nil {
		return false, errors.Wrapf(err, "failed to generate claim typed data bytes")
	}
	rawDataHash := crypto.Keccak256Hash(rawData)
	sigHex, err := hexDecode(signature)
	if err != nil {
		return false, errors.Wrapf(err, "failed to hex decode signature")
	}

	if len(sigHex) != crypto.SignatureLength {
		return false, fmt.Errorf("signature must be %d bytes long", crypto.SignatureLength)
	}

	if len(sigHex) != 65 {
		return false, fmt.Errorf("invalid signature length: %d", len(sigHex))
	}

	if sigHex[crypto.RecoveryIDOffset] == 27 || sigHex[crypto.RecoveryIDOffset] == 28 {
		sigHex[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
	}

	pubKeyRaw, err := crypto.Ecrecover(rawDataHash.Bytes(), sigHex)
	if err != nil {
		return false, errors.Wrapf(err, "failed to recover public key from signature")
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyRaw)
	if err != nil {
		return false, errors.Wrapf(err, "failed to unmarshal public key from signature")
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	if !bytes.Equal(common.HexToAddress(ethAddress).Bytes(), recoveredAddr.Bytes()) {
		return false, errors.New("signature does not match address")
	}
	return true, nil
}

// HexDecode returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func hexDecode(s string) ([]byte, error) {
	if has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return hex.DecodeString(s)
}

// Has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(str string) bool {
	return strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X")
}

func mergeClaimRecords(claimA, claimB types.ClaimRecord) (types.ClaimRecord, error) {
	if claimA.IsEmpty() {
		return claimB, nil
	}

	if claimB.IsEmpty() {
		return claimA, nil
	}

	if claimA.Address != claimB.Address {
		return types.ClaimRecord{}, errors.New("cannot merge claims for different addresses")
	}

	if claimA.Chain != claimB.Chain {
		return types.ClaimRecord{}, errors.New("cannot merge claims for different chains")
	}

	if claimA.AmountClaim.IsNil() || !claimA.AmountClaim.IsValid() {
		claimA.AmountClaim = claimB.AmountClaim
	} else {
		claimA.AmountClaim = claimA.AmountClaim.Add(claimB.AmountClaim)
	}

	if claimA.AmountDelegate.IsNil() || !claimA.AmountDelegate.IsValid() {
		claimA.AmountDelegate = claimB.AmountDelegate
	} else {
		claimA.AmountDelegate = claimA.AmountDelegate.Add(claimB.AmountDelegate)
	}

	if claimA.AmountVote.IsNil() || !claimA.AmountVote.IsValid() {
		claimA.AmountVote = claimB.AmountVote
	} else {
		claimA.AmountVote = claimA.AmountVote.Add(claimB.AmountVote)
	}

	return claimA, nil
}
