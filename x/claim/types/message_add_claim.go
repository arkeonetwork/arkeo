package types

import (
	"cosmossdk.io/errors"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddClaim = "add_claim"

var _ sdk.Msg = &MsgAddClaim{}

func NewMsgAddClaim(creator cosmos.AccAddress, chain Chain, address string, amount int64) *MsgAddClaim {
	return &MsgAddClaim{
		Creator: creator.String(),
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
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Creator)}
}

func (msg *MsgAddClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddClaim) ValidateBasic() error {
	_, ok := Chain_value[msg.Chain.String()]
	if !ok {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid chain(%s)", msg.Chain)
	}
	if !IsValidAddress(msg.Address, msg.Chain) {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid address")
	}
	if msg.Amount <= 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "amount should larger than 0")
	}
	return nil
}
