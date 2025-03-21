package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

var (
	EIP712Types = apitypes.Types{
		"Claim": []apitypes.Type{
			{Name: "address", Type: "address"},
			{Name: "arkeoAddress", Type: "string"},
			{Name: "amount", Type: "string"},
		},
		"EIP712Domain": []apitypes.Type{
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
		},
	}

	EIP712Domain = apitypes.TypedDataDomain{
		Name:    "ArkdropClaim",
		Version: "1",
		ChainId: math.NewHexOrDecimal256(1),
	}
)

// isValidEthAddress checks if the provided string is a valid address or not.
func IsValidEthAddress(address string) bool {
	return common.IsHexAddress(address)
}
