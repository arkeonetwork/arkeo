package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClaimArkeo = "claim_arkeo"

var _ sdk.Msg = &MsgClaimArkeo{}

func NewMsgClaimArkeo(creator string) *MsgClaimArkeo {
	return &MsgClaimArkeo{
		Creator: creator,
	}
}

func (msg *MsgClaimArkeo) Route() string {
	return RouterKey
}

func (msg *MsgClaimArkeo) Type() string {
	return TypeMsgClaimArkeo
}

func (msg *MsgClaimArkeo) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClaimArkeo) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimArkeo) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
