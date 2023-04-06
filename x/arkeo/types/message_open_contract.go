package types

import (
	"cosmossdk.io/errors"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgOpenContract = "open_contract"

var _ sdk.Msg = &MsgOpenContract{}

func NewMsgOpenContract(creator string, provider common.PubKey, service string, client, delegate common.PubKey,
	userType UserType, meterType MeterType, duration, settlementDuration int64, rate types.Coin,
	deposit cosmos.Int, restrictions Restrictions,
) *MsgOpenContract {
	return &MsgOpenContract{
		Creator:            creator,
		Provider:           provider,
		Service:            service,
		UserType:           userType,
		MeterType:          meterType,
		Duration:           duration,
		Rate:               rate,
		Client:             client,
		Deposit:            deposit,
		Delegate:           delegate,
		SettlementDuration: settlementDuration,
		Restrictions:       &restrictions,
	}
}

func (msg *MsgOpenContract) Route() string {
	return RouterKey
}

func (msg *MsgOpenContract) Type() string {
	return TypeMsgOpenContract
}

func (msg *MsgOpenContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

func (msg *MsgOpenContract) MustGetSigner() sdk.AccAddress {
	return msg.Creator
}

func (msg *MsgOpenContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgOpenContract) GetSpender() common.PubKey {
	if !msg.Delegate.IsEmpty() {
		return msg.Delegate
	}
	return msg.Client
}

func (msg *MsgOpenContract) ValidateBasic() error {
	// verify pubkey
	_, err := common.NewPubKey(msg.Provider.String())
	if err != nil {
		return errors.Wrapf(ErrInvalidPubKey, "invalid pubkey (%s)", err)
	}

	// verify service
	_, err = common.NewService(msg.Service)
	if err != nil {
		return errors.Wrapf(ErrInvalidService, "invalid service (%s): %s", msg.Service, err)
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

	if err := msg.Rate.Validate(); err != nil {
		return errors.Wrapf(err, "invalid rate")
	}

	if !msg.Rate.Amount.IsPositive() {
		return errors.Wrapf(ErrOpenContractRate, "contract rate cannot be zero")
	}

	if msg.SettlementDuration < 0 {
		return errors.Wrapf(ErrInvalidModProviderSettlementDuration, "settlement duration cannot be negative")
	}

	return nil
}
