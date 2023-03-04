package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
)

const TypeMsgBondProvider = "bond_provider"

var _ sdk.Msg = &MsgBondProvider{}

func NewMsgBondProvider(creator string, provider common.PubKey, chain string, bond cosmos.Int) *MsgBondProvider {
	return &MsgBondProvider{
		Creator:  creator,
		Provider: provider,
		Chain:    chain,
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
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBondProvider) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

func (msg *MsgBondProvider) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBondProvider) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// verify pubkey
	_, err = common.NewPubKey(msg.Provider.String())
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

	// verify chain
	_, err = common.NewChain(msg.Chain)
	if err != nil {
		return errors.Wrapf(ErrInvalidChain, "invalid chain (%s): %s", msg.Chain, err)
	}

	if msg.Bond.IsNil() || msg.Bond.IsZero() {
		return errors.Wrapf(ErrInvalidBond, "bond cannot be set to zero")
	}

	return nil
}
