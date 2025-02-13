package types

import (
	"fmt"

	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/arkeonetwork/arkeo/common/cosmos"
)

const TypeMsgClaimContractIncome = "claim_contract_income"

var _ sdk.Msg = &MsgClaimContractIncome{}

func NewMsgClaimContractIncome(creator cosmos.AccAddress, contractId uint64, nonce int64, sig []byte, chainId string, expiresAtBlock int64) *MsgClaimContractIncome {
	return &MsgClaimContractIncome{
		Creator:                 creator.String(),
		ContractId:              contractId,
		Nonce:                   nonce,
		Signature:               sig,
		ChainId:                 chainId,
		SignatureExpiresAtBlock: expiresAtBlock,
	}
}

func (msg *MsgClaimContractIncome) Route() string {
	return RouterKey
}

func (msg *MsgClaimContractIncome) Type() string {
	return TypeMsgClaimContractIncome
}

func (msg *MsgClaimContractIncome) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Creator)}
}

func (msg *MsgClaimContractIncome) MustGetSigner() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(msg.Creator)
}

func (msg *MsgClaimContractIncome) GetSignBytes() []byte {
	return ModuleCdc.MustMarshalJSON(msg)
}

func (msg *MsgClaimContractIncome) GetBytesToSign() []byte {
	return GetBytesToSign(msg.ContractId, msg.Nonce, msg.ChainId, msg.SignatureExpiresAtBlock)
}

func GetBytesToSign(contractId uint64, nonce int64, chainId string, expiresAtBlock int64) []byte {
	return []byte(fmt.Sprintf("%d:%d:%s:%d", contractId, nonce, chainId, expiresAtBlock))
}

func (msg *MsgClaimContractIncome) ValidateBasic() error {
	// anyone can make the claim on a contract, but of course the payout would only happen to the provider

	if len(msg.Signature) > 100 {
		return errors.Wrap(ErrClaimContractIncomeInvalidSignature, "too long")
	}

	if len(msg.ChainId) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "chain ID cannot be empty")
	}

	if msg.Nonce <= 0 {
		return errors.Wrap(ErrClaimContractIncomeBadNonce, "nonce must be greater than zero")
	}

	if msg.SignatureExpiresAtBlock <= 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "expiration block height must be positive")
	}

	return nil
}
