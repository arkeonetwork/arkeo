package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgTransferClaim = "transfer_claim"

var _ sdk.Msg = &MsgTransferClaim{}

func NewMsgTransferClaim(creator string, toAddress string) *MsgTransferClaim {
	return &MsgTransferClaim{
		Creator:   creator,
		ToAddress: toAddress,
	}
}

func (msg *MsgTransferClaim) Route() string {
	return RouterKey
}

func (msg *MsgTransferClaim) Type() string {
	return TypeMsgTransferClaim
}

func (msg *MsgTransferClaim) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgTransferClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTransferClaim) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
