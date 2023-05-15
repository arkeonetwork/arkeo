package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetVersion = "set_version"

var _ sdk.Msg = &MsgSetVersion{}

func NewMsgSetVersion(creator string, version int32) *MsgSetVersion {
	return &MsgSetVersion{
		Creator: creator,
		Version: version,
	}
}

func (msg *MsgSetVersion) Route() string {
	return RouterKey
}

func (msg *MsgSetVersion) Type() string {
	return TypeMsgSetVersion
}

func (msg *MsgSetVersion) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetVersion) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetVersion) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
