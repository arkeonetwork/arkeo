package types

import (
	"cosmossdk.io/errors"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgOpenContract = "open_contract"

var _ sdk.Msg = &MsgOpenContract{}

func NewMsgOpenContract(creator string, pubkey common.PubKey, chain string, client, delegate common.PubKey, contractType ContractType, duration, rate int64, deposit cosmos.Int) *MsgOpenContract {
	return &MsgOpenContract{
		Creator:      creator,
		PubKey:       pubkey,
		Chain:        chain,
		ContractType: contractType,
		Duration:     duration,
		Rate:         rate,
		Client:       client,
		Deposit:      deposit,
		Delegate:     delegate,
	}
}

func (msg *MsgOpenContract) Route() string {
	return RouterKey
}

func (msg *MsgOpenContract) Type() string {
	return TypeMsgOpenContract
}

func (msg *MsgOpenContract) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgOpenContract) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

func (msg *MsgOpenContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgOpenContract) FetchSpender() common.PubKey {
	if !msg.Delegate.IsEmpty() {
		return msg.Delegate
	}
	return msg.Client
}

func (msg *MsgOpenContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// verify pubkey
	_, err = common.NewPubKey(msg.PubKey.String())
	if err != nil {
		return errors.Wrapf(ErrInvalidPubKey, "invalid pubkey (%s)", err)
	}

	// verify chain
	_, err = common.NewChain(msg.Chain)
	if err != nil {
		return errors.Wrapf(ErrInvalidChain, "invalid chain (%s): %s", msg.Chain, err)
	}

	// verify client
	_, err = common.NewPubKey(msg.Client.String())
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

	if msg.Duration <= 0 {
		return errors.Wrapf(ErrOpenContractDuration, "contract duration cannot be zero")
	}

	if msg.Rate <= 0 {
		return errors.Wrapf(ErrOpenContractRate, "contract rate cannot be zero")
	}

	return nil
}
