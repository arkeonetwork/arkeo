package keeper

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"

	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/pkg/errors"
)

func (k msgServer) ClaimEth(goCtx context.Context, msg *types.MsgClaimEth) (*types.MsgClaimEthResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// 1. get eth claim
	ethClaim, err := k.GetClaimRecord(ctx, msg.EthAddress, types.ETHEREUM)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get claim record for %s", msg.EthAddress)
	}

	if ethClaim.InitialClaimableAmount.IsZero() {
		return nil, errors.Wrapf(err, "no claimable amount for %s", msg.EthAddress)
	}

	// 2. check if already claimed
	if ethClaim.ActionCompleted[types.FOREIGN_CHAIN_ACTION_CLAIM] {
		return nil, errors.Wrapf(err, "already claimed for %s", msg.EthAddress)
	}

	// 3. validate signature
	isValid, err := IsValidClaimSignature(msg.EthAddress, msg.Creator,
		ethClaim.InitialClaimableAmount.AmountOf(types.DefaultClaimDenom).String(), msg.Signature)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to validate signature for %s", msg.EthAddress)
	}

	if !isValid {
		// this shouldn't happen without an error, but just in case
		return nil, errors.New("invalid signature")
	}

	// set eth claim to completed
	ethClaim.ActionCompleted[types.FOREIGN_CHAIN_ACTION_CLAIM] = true
	err = k.SetClaimRecord(ctx, ethClaim)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to set claim record for %s", msg.EthAddress)
	}

	// create new arkeo claim
	arkeoClaim := types.ClaimRecord{
		Address:                msg.Creator,
		Chain:                  types.ARKEO,
		InitialClaimableAmount: ethClaim.InitialClaimableAmount,
		ActionCompleted:        []bool{false, false},
	}
	err = k.SetClaimRecord(ctx, arkeoClaim)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to set claim record for %s", msg.Creator)
	}

	return &types.MsgClaimEthResponse{}, nil
}

func GenerateClaimTypedDataBytes(ethAddress string, arkeoAddress string, amount string) ([]byte, error) {
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

func IsValidClaimSignature(ethAddress string, arkeoAdddress string, amount string, signature string) (bool, error) {
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

	// if sigHex[crypto.RecoveryIDOffset] != 27 && sigHex[crypto.RecoveryIDOffset] != 28 {
	// 	return false, fmt.Errorf("invalid recovery id: %d", sigHex[64])
	// }
	// sigHex[64] -= 27

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
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}
