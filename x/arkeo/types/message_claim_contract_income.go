package types

import (
	"fmt"

	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/arkeonetwork/arkeo/common/cosmos"
)

const TypeMsgClaimContractIncome = "claim_contract_income"

var _ sdk.Msg = &MsgClaimContractIncome{}

func NewMsgClaimContractIncome(creator cosmos.AccAddress, contractId uint64, nonce int64, sig []byte, chainId string, signatureExpiry int64) *MsgClaimContractIncome {
	return &MsgClaimContractIncome{
		Creator:            creator.String(),
		ContractId:         contractId,
		Nonce:              nonce,
		Signature:          sig,
		ChainId:            chainId,
		SignatureExpiresAt: signatureExpiry,
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
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimContractIncome) GetBytesToSign() []byte {
	return GetBytesToSign(msg.ContractId, msg.Nonce, msg.ChainId)
}

func GetBytesToSign(contractId uint64, nonce int64, chainId string) []byte {
	return []byte(fmt.Sprintf("%d:%d:%s", contractId, nonce, chainId))
}

func (msg *MsgClaimContractIncome) ValidateBasic() error {
	// anyone can make the claim on a contract, but of course the payout would only happen to the provider

	if len(msg.Signature) > 100 {
		return errors.Wrap(ErrClaimContractIncomeInvalidSignature, "too long")
	}

	if msg.Nonce <= 0 {
		return errors.Wrap(ErrClaimContractIncomeBadNonce, "")
	}

	if len(msg.ChainId) == 0 {
		return errors.Wrap(ErrInvalidChainId, "chain id is not specified")
	}

	if msg.SignatureExpiresAt <= 0 {
		return errors.Wrap(ErrSignatureExpired, "expiry on signature is not set")
	}
	return nil
}
