package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClaimEth = "claim_eth"

var _ sdk.Msg = &MsgClaimEth{}

func NewMsgClaimEth(creator, ethAdress, signature string) *MsgClaimEth {
	return &MsgClaimEth{
		Creator:    creator,
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
	addr, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{addr}
}

func (msg *MsgClaimEth) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimEth) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator")
	}
	return nil
}
