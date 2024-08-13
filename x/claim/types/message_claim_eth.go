package types

import (
	"github.com/arkeonetwork/arkeo/common/cosmos"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgClaimEth = "claim_eth"

var _ sdk.Msg = &MsgClaimEth{}

func NewMsgClaimEth(creator cosmos.AccAddress, ethAdress, signature string) *MsgClaimEth {
	return &MsgClaimEth{
		Creator:    creator.String(),
		EthAddress: ethAdress,
		Signature:  signature,
	}
}

func (msg *MsgClaimEth) Route() string {
	return RouterKey
}

func (msg *MsgClaimEth) Type() string {
	return TypeMsgClaimEth
}

func (msg *MsgClaimEth) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Creator)}
}

func (msg *MsgClaimEth) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimEth) ValidateBasic() error {
	return nil
}
