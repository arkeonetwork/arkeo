package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/arkeonetwork/arkeo/common/cosmos"
)

const TypeMsgTransferClaim = "transfer_claim"

var _ sdk.Msg = &MsgTransferClaim{}

func NewMsgTransferClaim(creator, toAddress cosmos.AccAddress) *MsgTransferClaim {
	return &MsgTransferClaim{
		Creator:   creator.String(),
		ToAddress: toAddress.String(),
	}
}

func (msg *MsgTransferClaim) Route() string {
	return RouterKey
}

func (msg *MsgTransferClaim) Type() string {
	return TypeMsgTransferClaim
}

func (msg *MsgTransferClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Creator)}
}

func (msg *MsgTransferClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTransferClaim) ValidateBasic() error {
	return nil
}
