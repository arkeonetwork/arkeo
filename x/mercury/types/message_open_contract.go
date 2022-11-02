package types

import (
	"mercury/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgOpenContract = "open_contract"

var _ sdk.Msg = &MsgOpenContract{}

func NewMsgOpenContract(creator string, pubkey common.PubKey, chain common.Chain, cType ContractType, duration, rate int64) *MsgOpenContract {
	return &MsgOpenContract{
		Creator:  creator,
		PubKey:   pubkey,
		Chain:    chain,
		CType:    cType,
		Duration: duration,
		Rate:     rate,
	}
}

func (msg *MsgOpenContract) Route() string {
	return RouterKey
}

func (msg *MsgOpenContract) Type() string {
	return TypeMsgOpenContract
}

func (msg *MsgOpenContract) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgOpenContract) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

func (msg *MsgOpenContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgOpenContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// verify pubkey
	_, err = common.NewPubKey(msg.PubKey.String())
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidPubKey, "invalid pubkey (%s)", err)
	}

	// verify chain
	_, err = common.NewChain(msg.Chain.String())
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidChain, "invalid chain (%s): %s", msg.Chain, err)
	}

	// verify client
	_, err = common.NewPubKey(msg.Client.String())
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidPubKey, "invalid pubkey (%s)", err)
	}

	signer := msg.MustGetSigner()
	client, err := msg.Client.GetMyAddress()
	if err != nil {
		return err
	}
	if !signer.Equals(client) {
		return sdkerrors.Wrapf(ErrProviderBadSigner, "Signer: %s, Client Address: %s", msg.GetSigners(), client)
	}

	if msg.Duration <= 0 {
		return sdkerrors.Wrapf(ErrOpenContractDuration, "contract duration cannot be zero")
	}

	if msg.Rate <= 0 {
		return sdkerrors.Wrapf(ErrOpenContractRate, "contract rate cannot be zero")
	}

	return nil
}
