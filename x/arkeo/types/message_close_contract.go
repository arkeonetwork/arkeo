package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCloseContract = "close_contract"

var _ sdk.Msg = &MsgCloseContract{}

func NewMsgCloseContract(creator string, contractId uint64) *MsgCloseContract {
	return &MsgCloseContract{
		Creator:    creator,
		ContractId: contractId,
	}
}

func (msg *MsgCloseContract) Route() string {
	return RouterKey
}

func (msg *MsgCloseContract) Type() string {
	return TypeMsgCloseContract
}

func (msg *MsgCloseContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.MustGetSigner()}
}

func (msg *MsgCloseContract) MustGetSigner() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Creator)
	return addr
}

func (msg *MsgCloseContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCloseContract) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator")
	}
	return nil
}
