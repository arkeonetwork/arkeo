package types

import (
	"cosmossdk.io/errors"
	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgCloseContract = "close_contract"

var _ sdk.Msg = &MsgCloseContract{}

func NewMsgCloseContract(creator cosmos.AccAddress, contractId uint64, client common.PubKey, delegate common.PubKey) *MsgCloseContract {
	return &MsgCloseContract{
		Creator:    creator.String(),
		ContractId: contractId,
		Client:     client,
		Delegate:   delegate,
	}
}

func (msg *MsgCloseContract) Route() string {
	return RouterKey
}

func (msg *MsgCloseContract) Type() string {
	return TypeMsgCloseContract
}

func (msg *MsgCloseContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(msg.Creator)}
}

func (msg *MsgCloseContract) MustGetSigner() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(msg.Creator)
}

func (msg *MsgCloseContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCloseContract) ValidateBasic() error {
	if msg == nil {
		return errors.Wrap(cosmos.ErrUnknownRequest("invalid close contract message"), "message cammot be empty")
	}
	if msg.Creator == "" {
		return errors.Wrapf(ErrCloseContractUnauthorized, "creator cannot be empty")
	}

	if msg.ContractId == 0 {
		return errors.Wrap(ErrContractNotFound, "invalid contract id")
	}

	if msg.Client == nil {
		return errors.Wrap(ErrCloseContractUnauthorized, "client id cannot be empty")
	}

	// verify client
	_, err := common.NewPubKey(msg.Client.String())
	if err != nil {
		return errors.Wrapf(ErrInvalidPubKey, "invalid pubkey (%s)", err)
	}

	signer := msg.MustGetSigner()
	client, err := msg.Client.GetMyAddress()
	if err != nil {
		return err
	}
	if !signer.Equals(client) {
		return errors.Wrapf(ErrInvalidPubKey, "Signer: %s, Client Address: %s", msg.GetSigners(), client)
	}

	return nil
}
