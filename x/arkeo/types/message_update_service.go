package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateService = "update_service"

var _ sdk.Msg = &MsgUpdateService{}

func NewMsgUpdateService(creator string, id uint64, name, description, serviceType string) *MsgUpdateService {
	return &MsgUpdateService{
		Creator:     creator,
		Id:          id,
		Name:        name,
		Description: description,
		ServiceType: serviceType,
	}
}

func (msg *MsgUpdateService) Route() string { return RouterKey }

func (msg *MsgUpdateService) Type() string { return TypeMsgUpdateService }

func (msg *MsgUpdateService) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateService) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg *MsgUpdateService) ValidateBasic() error {
	if _, err := cosmos.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Id == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "id must be non-zero")
	}
	if strings.TrimSpace(msg.Name) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "name cannot be empty")
	}
	return nil
}
