package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRegisterProvider = "register_provider"

var _ sdk.Msg = &MsgRegisterProvider{}

func NewMsgRegisterProvider(creator string, chain string, pubkey string) *MsgRegisterProvider {
  return &MsgRegisterProvider{
		Creator: creator,
    Chain: chain,
    Pubkey: pubkey,
	}
}

func (msg *MsgRegisterProvider) Route() string {
  return RouterKey
}

func (msg *MsgRegisterProvider) Type() string {
  return TypeMsgRegisterProvider
}

func (msg *MsgRegisterProvider) GetSigners() []sdk.AccAddress {
  creator, err := sdk.AccAddressFromBech32(msg.Creator)
  if err != nil {
    panic(err)
  }
  return []sdk.AccAddress{creator}
}

func (msg *MsgRegisterProvider) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterProvider) ValidateBasic() error {
  _, err := sdk.AccAddressFromBech32(msg.Creator)
  	if err != nil {
  		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
  	}
  return nil
}

