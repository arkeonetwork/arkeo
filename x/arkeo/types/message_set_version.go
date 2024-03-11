package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetVersion = "set_version"

var _ sdk.Msg = &MsgSetVersion{}

func NewMsgSetVersion(creator string, version int64) *MsgSetVersion {
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
	addr, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{addr}
}

func (msg *MsgSetVersion) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetVersion) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator")
	}
	if msg.Version <= 0 {
		return ErrInvalidVersion
	}
	return nil
}
