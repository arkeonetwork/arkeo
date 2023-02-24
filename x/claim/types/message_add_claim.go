package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddClaim = "add_claim"

var _ sdk.Msg = &MsgAddClaim{}

func NewMsgAddClaim(creator string, chain string, address string, amount sdk.Coins) *MsgAddClaim {
	return &MsgAddClaim{
		Creator: creator,
		Chain:   chain,
		Address: address,
		Amount:  amount,
	}
}

func (msg *MsgAddClaim) Route() string {
	return RouterKey
}

func (msg *MsgAddClaim) Type() string {
	return TypeMsgAddClaim
}

func (msg *MsgAddClaim) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddClaim) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
