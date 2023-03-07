package types

import (
	"fmt"

	"cosmossdk.io/errors"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClaimContractIncome = "claim_contract_income"

var _ sdk.Msg = &MsgClaimContractIncome{}

func NewMsgClaimContractIncome(creator string, contractId uint64, spender common.PubKey, nonce int64, sig []byte) *MsgClaimContractIncome {
	return &MsgClaimContractIncome{
		Creator:    creator,
		ContractId: contractId,
		Spender:    spender,
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

func (msg *MsgClaimContractIncome) GetSpenderAddress() (sdk.AccAddress, error) {
	acc, err := msg.Spender.GetMyAddress()
	if err == nil {
		return acc, nil
	}

	return sdk.AccAddressFromBech32(msg.Creator)
}

func (msg *MsgClaimContractIncome) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClaimContractIncome) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

func (msg *MsgClaimContractIncome) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimContractIncome) GetBytesToSign() []byte {
	return GetBytesToSign(msg.ContractId, msg.Spender, msg.Nonce)
}

func GetBytesToSign(contractId uint64, spender common.PubKey, nonce int64) []byte {
	return []byte(fmt.Sprintf("%d:%s:%d", contractId, spender, nonce))
}

func (msg *MsgClaimContractIncome) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	// anyone can make the claim on a contract, but of course the payout would only happen to the provider

	// verify spender pubkey
	_, err = common.NewPubKey(msg.Spender.String())
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidPubKey, "invalid spender pubkey (%s)", err)
	}

	if len(msg.Signature) > 100 {
		return errors.Wrap(ErrClaimContractIncomeInvalidSignature, "too long")
	}

	if msg.Nonce <= 0 {
		return errors.Wrap(ErrClaimContractIncomeBadNonce, "")
	}

	pk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, msg.Spender.String())
	if err != nil {
		return err
	}

	if !pk.VerifySignature(msg.GetBytesToSign(), msg.Signature) {
		return errors.Wrap(ErrClaimContractIncomeInvalidSignature, "")
	}

	return nil
}
