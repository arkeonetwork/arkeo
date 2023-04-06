package types

import (
	"github.com/arkeonetwork/arkeo/common/cosmos"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgTransferClaim = "transfer_claim"

var _ sdk.Msg = &MsgTransferClaim{}

func NewMsgTransferClaim(creator, toAddress cosmos.AccAddress) *MsgTransferClaim {
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
	return []sdk.AccAddress{msg.Creator}
}

func (msg *MsgTransferClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTransferClaim) ValidateBasic() error {
	return nil
}
