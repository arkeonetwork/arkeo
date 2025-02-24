package types

import (
	"fmt"

	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/arkeonetwork/arkeo/common/cosmos"
)

const TypeMsgClaimContractIncome = "claim_contract_income"

var _ sdk.Msg = &MsgClaimContractIncome{}

func NewMsgClaimContractIncome(creator cosmos.AccAddress, contractId uint64, nonce int64, sig []byte) *MsgClaimContractIncome {
	return &MsgClaimContractIncome{
		Creator:    creator.String(),
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
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Creator)}
}

func (msg *MsgClaimContractIncome) MustGetSigner() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(msg.Creator)
}

func (msg *MsgClaimContractIncome) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimContractIncome) GetBytesToSign(chainId string) []byte {
	return GetBytesToSign(msg.ContractId, msg.Nonce, chainId)
}

func GetBytesToSign(contractId uint64, nonce int64, chainID string) []byte {
	return []byte(fmt.Sprintf("%d:%d:%s", contractId, nonce, chainID))
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
