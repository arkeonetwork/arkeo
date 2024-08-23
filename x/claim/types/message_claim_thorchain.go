package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/arkeonetwork/arkeo/common/cosmos"
)

const TypeMsgClaimThorchain = "claim_thorchain"

var _ sdk.Msg = &MsgClaimEth{}

func NewMsgClaimThorchain(creator cosmos.AccAddress, fromAddress string, toAddress string) *MsgClaimThorchain {
	return &MsgClaimThorchain{
		Creator:     creator.String(),
		FromAddress: fromAddress,
		ToAddress:   toAddress,
	}
}

func (msg *MsgClaimThorchain) Route() string {
	return RouterKey
}

func (msg *MsgClaimThorchain) Type() string {
	return TypeMsgClaimEth
}

func (msg *MsgClaimThorchain) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Creator)}
}

func (msg *MsgClaimThorchain) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimThorchain) ValidateBasic() error {
	return nil
}
