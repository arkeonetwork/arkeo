package types

import (
	fmt "fmt"

	"cosmossdk.io/errors"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgOpenContract = "open_contract"

var _ sdk.Msg = &MsgOpenContract{}

func NewMsgOpenContract(creator, provider, service, client, delegate string, contractType ContractType, duration, settlementDuration int64, rate types.Coin, deposit cosmos.Int, authorization ContractAuthorization, qpm int64) *MsgOpenContract {
	return &MsgOpenContract{
		Creator:            creator,
		Provider:           provider,
		Service:            service,
		ContractType:       contractType,
		Duration:           duration,
		Rate:               rate,
		Client:             client,
		Deposit:            deposit,
		Delegate:           delegate,
		SettlementDuration: settlementDuration,
		Authorization:      authorization,
		QueriesPerMinute:   qpm,
	}
}

func (msg *MsgOpenContract) Route() string {
	return RouterKey
}

func (msg *MsgOpenContract) Type() string {
	return TypeMsgOpenContract
}

func (msg *MsgOpenContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.MustGetSigner()}
}

func (msg *MsgOpenContract) MustGetSigner() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Creator)
	return addr
}

func (msg *MsgOpenContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgOpenContract) GetSpender() common.PubKey {
	if len(msg.Delegate) > 0 {
		pk, _ := common.NewPubKey(msg.Delegate)
		return pk
	}
	pk, _ := common.NewPubKey(msg.Client)
	return pk
}

func (msg *MsgOpenContract) ValidateBasic() error {
	// verify pubkey
	_, err := common.NewPubKey(msg.Provider)
	if err != nil {
		return errors.Wrapf(ErrInvalidPubKey, "invalid pubkey (%s)", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator")
	}
	// verify service
	_, err = common.NewService(msg.Service)
	if err != nil {
		return errors.Wrapf(ErrInvalidService, "invalid service (%s): %s", msg.Service, err)
	}

	// verify client
	pk, err := common.NewPubKey(msg.Client)
	if err != nil {
		return errors.Wrapf(ErrInvalidPubKey, "invalid pubkey (%s)", err)
	}

	signer := msg.MustGetSigner()
	client, err := pk.GetMyAddress()
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

	if msg.QueriesPerMinute <= 0 {
		return fmt.Errorf("queries per minute must be greater than zero")
	}

	if !msg.Rate.Amount.IsPositive() {
		return errors.Wrapf(ErrOpenContractRate, "contract rate cannot be zero")
	}

	if msg.SettlementDuration < 0 {
		return errors.Wrapf(ErrInvalidModProviderSettlementDuration, "settlement duration cannot be negative")
	}

	// cannot open pay-as-you-go contract and be "open" authorization. The
	// reason for this is a pay-as-you-go contract that is open allows anyone
	// (including the data provider) to completely empty the contract tokens to
	// the data provider without providing any data to the contract owner. It
	// would be too easy to "rug" contract owners.
	if msg.ContractType == ContractType_PAY_AS_YOU_GO && msg.Authorization == ContractAuthorization_OPEN {
		return errors.Wrapf(ErrInvalidAuthorization, "pay-as-you-go contract cannot use open authorization")
	}

	return nil
}
