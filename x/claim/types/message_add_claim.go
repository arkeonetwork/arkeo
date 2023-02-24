package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddClaim = "add_claim"

var _ sdk.Msg = &MsgAddClaim{}

func NewMsgAddClaim(creator string, chain Chain, address string, amount int64) *MsgAddClaim {
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
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid  address (%s)", err)
	}
	_, ok := Chain_value[msg.Chain.String()]
	if !ok {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid chain(%s),err: %s", msg.Chain, err)
	}
	if msg.Amount <= 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "amount should larger than 0")
	}
	return nil
}
