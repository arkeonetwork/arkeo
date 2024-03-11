package types

import (
	"cosmossdk.io/errors"
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
	addr, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{addr}
}

func (msg *MsgClaimArkeo) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimArkeo) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator")
	}
	return nil
}
