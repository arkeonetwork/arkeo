package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
)

const TypeMsgBondProvider = "bond_provider"

var _ sdk.Msg = &MsgBondProvider{}

func NewMsgBondProvider(creator cosmos.AccAddress, provider common.PubKey, service string, bond cosmos.Int) *MsgBondProvider {
	return &MsgBondProvider{
		Creator:  creator,
		Provider: provider,
		Service:  service,
		Bond:     bond,
	}
}

func (msg *MsgBondProvider) Route() string {
	return RouterKey
}

func (msg *MsgBondProvider) Type() string {
	return TypeMsgBondProvider
}

func (msg *MsgBondProvider) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

func (msg *MsgBondProvider) MustGetSigner() sdk.AccAddress {
	return msg.Creator
}

func (msg *MsgBondProvider) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBondProvider) ValidateBasic() error {
	// verify pubkey
	_, err := common.NewPubKey(msg.Provider.String())
	if err != nil {
		return errors.Wrapf(ErrInvalidPubKey, "invalid pubkey (%s): %s", msg.Provider, err)
	}

	signer := msg.MustGetSigner()
	provider, err := msg.Provider.GetMyAddress()
	if err != nil {
		return err
	}
	if !signer.Equals(provider) {
		return errors.Wrapf(ErrProviderBadSigner, "Signer: %s, Provider Address: %s", msg.GetSigners(), provider)
	}

	// verify service
	_, err = common.NewService(msg.Service)
	if err != nil {
		return errors.Wrapf(ErrInvalidService, "invalid service (%s): %s", msg.Service, err)
	}

	if msg.Bond.IsNil() || msg.Bond.IsZero() {
		return errors.Wrapf(ErrInvalidBond, "bond cannot be set to zero")
	}

	return nil
}
