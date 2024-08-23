package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/arkeonetwork/arkeo/common/cosmos"
)

const TypeMsgClaimArkeo = "claim_arkeo"

var _ sdk.Msg = &MsgClaimArkeo{}

func NewMsgClaimArkeo(creator cosmos.AccAddress) *MsgClaimArkeo {
	return &MsgClaimArkeo{
		Creator: creator.String(),
	}
}

func (msg *MsgClaimArkeo) Route() string {
	return RouterKey
}

func (msg *MsgClaimArkeo) Type() string {
	return TypeMsgClaimArkeo
}

func (msg *MsgClaimArkeo) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Creator)}
}

func (msg *MsgClaimArkeo) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimArkeo) ValidateBasic() error {
	return nil
}
