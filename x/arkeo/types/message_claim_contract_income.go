package types

import (
	"fmt"

	"cosmossdk.io/errors"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgClaimContractIncome = "claim_contract_income"

var _ sdk.Msg = &MsgClaimContractIncome{}

func NewMsgClaimContractIncome(creator cosmos.AccAddress, contractId uint64, nonce int64, sig []byte) *MsgClaimContractIncome {
	return &MsgClaimContractIncome{
		Creator:    creator,
		ContractId: contractId,
		Nonce:      nonce,
		Signature:  sig,
	}
}

func (msg *MsgClaimContractIncome) Route() string {
	return RouterKey
}

func (msg *MsgClaimContractIncome) Type() string {
	return TypeMsgClaimContractIncome
}

func (msg *MsgClaimContractIncome) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

func (msg *MsgClaimContractIncome) MustGetSigner() sdk.AccAddress {
	return msg.Creator
}

func (msg *MsgClaimContractIncome) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimContractIncome) GetBytesToSign() []byte {
	return GetBytesToSign(msg.ContractId, msg.Nonce)
}

func GetBytesToSign(contractId uint64, nonce int64) []byte {
	return []byte(fmt.Sprintf("%d:%d", contractId, nonce))
}

func (msg *MsgClaimContractIncome) ValidateBasic() error {
	// anyone can make the claim on a contract, but of course the payout would only happen to the provider

	if len(msg.Signature) > 100 {
		return errors.Wrap(ErrClaimContractIncomeInvalidSignature, "too long")
	}

	if msg.Nonce <= 0 {
		return errors.Wrap(ErrClaimContractIncomeBadNonce, "")
	}

	return nil
}
