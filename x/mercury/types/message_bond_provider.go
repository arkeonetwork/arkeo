package types

import (
	fmt "fmt"
	"mercury/common"
	"mercury/common/cosmos"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgBondProvider = "bond_provider"

var _ sdk.Msg = &MsgBondProvider{}

func NewMsgBondProvider(creator string, pubkey common.PubKey, chain common.Chain, bond cosmos.Int) *MsgBondProvider {
	return &MsgBondProvider{
		Creator: creator,
		PubKey:  pubkey,
		Chain:   chain,
		Bond:    bond,
	}
}

func (msg *MsgBondProvider) Route() string {
	return RouterKey
}

func (msg *MsgBondProvider) Type() string {
	return TypeMsgBondProvider
}

func (msg *MsgBondProvider) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBondProvider) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

func (msg *MsgBondProvider) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBondProvider) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		fmt.Println("got here")
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	signer := msg.MustGetSigner()
	provider, err := msg.PubKey.GetMyAddress()
	if err != nil {
		return err
	}
	if !signer.Equals(provider) {
		return sdkerrors.Wrapf(ErrProviderBadSigner, "Signer: %s, Provider Address: %s", msg.GetSigners(), provider)
	}

	if msg.Bond.IsNil() || msg.Bond.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidBond, "bond cannot be set to zero")
	}

	return nil
}
