package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/arkeonetwork/arkeo/common"
)

const TypeMsgCloseContract = "close_contract"

var _ sdk.Msg = &MsgCloseContract{}

func NewMsgCloseContract(creator string, pubkey common.PubKey, chain string, client, delegate common.PubKey) *MsgCloseContract {
	return &MsgCloseContract{
		Creator:  creator,
		PubKey:   pubkey,
		Chain:    chain,
		Client:   client,
		Delegate: delegate,
	}
}

func (msg *MsgCloseContract) Route() string {
	return RouterKey
}

func (msg *MsgCloseContract) Type() string {
	return TypeMsgCloseContract
}

func (msg *MsgCloseContract) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCloseContract) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

func (msg *MsgCloseContract) GetClientAddress() (sdk.AccAddress, error) {
	acc, err := msg.Client.GetMyAddress()
	if err == nil {
		return acc, nil
	}

	return sdk.AccAddressFromBech32(msg.Creator)
}

func (msg *MsgCloseContract) FetchSpender() common.PubKey {
	if !msg.Delegate.IsEmpty() {
		return msg.Delegate
	}
	return msg.Client
}

func (msg *MsgCloseContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCloseContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
