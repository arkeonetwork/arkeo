package types

import (
	fmt "fmt"

	"github.com/ArkeoNetwork/arkeo/common"
	"github.com/ArkeoNetwork/arkeo/common/cosmos"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClaimContractIncome = "claim_contract_income"

var _ sdk.Msg = &MsgClaimContractIncome{}

func NewMsgClaimContractIncome(creator string, pubkey common.PubKey, chain string, spender common.PubKey, nonce, height int64, sig []byte) *MsgClaimContractIncome {
	return &MsgClaimContractIncome{
		Creator:   creator,
		PubKey:    pubkey,
		Chain:     chain,
		Spender:   spender,
		Nonce:     nonce,
		Height:    height,
		Signature: sig,
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

func (msg *MsgClaimContractIncome) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// verify pubkey
	_, err = common.NewPubKey(msg.PubKey.String())
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidPubKey, "invalid pubkey (%s): %s", msg.PubKey, err)
	}

	// anyone can make the claim on a contract, but of course the payout would only happen to the provider

	// verify chain
	_, err = common.NewChain(msg.Chain)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidChain, "invalid chain (%s): %s", msg.Chain, err)
	}

	// verify spender pubkey
	_, err = common.NewPubKey(msg.Spender.String())
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "invalid spender pubkey (%s)", err)
	}

	if len(msg.Signature) > 100 {
		return sdkerrors.Wrap(ErrClaimContractIncomeInvalidSignature, "too long")
	}

	if msg.Height <= 0 {
		return sdkerrors.Wrap(ErrClaimContractIncomeBadHeight, "")
	}

	if msg.Nonce <= 0 {
		return sdkerrors.Wrap(ErrClaimContractIncomeBadNonce, "")
	}

	pk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, msg.Spender.String())
	if err != nil {
		return err
	}
	bites := []byte(fmt.Sprintf("%s:%s:%s:%d:%d", msg.PubKey, msg.Chain, msg.Spender, msg.Height, msg.Nonce))
	if !pk.VerifySignature(bites, msg.Signature) {
		return sdkerrors.Wrap(ErrClaimContractIncomeInvalidSignature, "")
	}

	return nil
}
