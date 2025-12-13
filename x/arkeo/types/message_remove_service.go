package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRemoveService = "remove_service"

var _ sdk.Msg = &MsgRemoveService{}

func NewMsgRemoveService(creator, name string) *MsgRemoveService {
	return &MsgRemoveService{
		Creator: creator,
		Name:    name,
	}
}

func (msg *MsgRemoveService) Route() string { return RouterKey }

func (msg *MsgRemoveService) Type() string { return TypeMsgRemoveService }

func (msg *MsgRemoveService) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveService) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg *MsgRemoveService) ValidateBasic() error {
	if _, err := cosmos.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if strings.TrimSpace(msg.Name) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "name cannot be empty")
	}
	return nil
}
