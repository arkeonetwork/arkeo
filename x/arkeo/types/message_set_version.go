package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/arkeonetwork/arkeo/common/cosmos"
)

const TypeMsgSetVersion = "set_version"

var _ sdk.Msg = &MsgSetVersion{}

func NewMsgSetVersion(creator cosmos.AccAddress, version int64) *MsgSetVersion {
	return &MsgSetVersion{
		Creator: creator.String(),
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
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Creator)}
}

func (msg *MsgSetVersion) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetVersion) ValidateBasic() error {
	if msg.Version <= 0 {
		return ErrInvalidVersion
	}
	return nil
}
