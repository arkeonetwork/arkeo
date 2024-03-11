package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgTransferClaim = "transfer_claim"

var _ sdk.Msg = &MsgTransferClaim{}

func NewMsgTransferClaim(creator, toAddress string) *MsgTransferClaim {
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
	addr, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{addr}
}

func (msg *MsgTransferClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTransferClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator")
	}
	return nil
}
