package types

import (
	"mercury/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClaimContractIncome = "claim_contract_income"

var _ sdk.Msg = &MsgClaimContractIncome{}

func NewMsgClaimContractIncome(creator string, pubkey common.PubKey, chain common.Chain, client string, nonce, height int64) *MsgClaimContractIncome {
	return &MsgClaimContractIncome{
		Creator: creator,
		PubKey:  pubkey,
		Chain:   chain,
		Client:  client,
		Nonce:   nonce,
		Height:  height,
	}
}

func (msg *MsgClaimContractIncome) Route() string {
	return RouterKey
}

func (msg *MsgClaimContractIncome) Type() string {
	return TypeMsgClaimContractIncome
}

func (msg *MsgClaimContractIncome) GetClientAddress() (sdk.AccAddress, error) {
	acc, err := sdk.AccAddressFromBech32(msg.Client)
	if err == nil {
		return acc, nil
	}

	return sdk.AccAddressFromBech32(msg.Creator)
}

func (msg *MsgClaimContractIncome) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClaimContractIncome) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

func (msg *MsgClaimContractIncome) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimContractIncome) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// verify pubkey
	_, err = common.NewPubKey(msg.PubKey.String())
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidPubKey, "invalid pubkey (%s): %s", msg.PubKey, err)
	}

	signer := msg.MustGetSigner()
	provider, err := msg.PubKey.GetMyAddress()
	if err != nil {
		return err
	}
	if !signer.Equals(provider) {
		return sdkerrors.Wrapf(ErrProviderBadSigner, "Signer: %s, Provider Address: %s", msg.GetSigners(), provider)
	}

	// verify chain
	_, err = common.NewChain(msg.Chain.String())
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidChain, "invalid chain (%s): %s", msg.Chain, err)
	}

	// verify client address
	_, err = sdk.AccAddressFromBech32(msg.Client)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid client address (%s)", err)
	}

	if msg.Height <= 0 {
		return sdkerrors.Wrap(ErrClaimContractIncomeBadHeight, "")
	}

	if msg.Nonce < 0 {
		return sdkerrors.Wrap(ErrClaimContractIncomeBadNonce, "")
	}

	// TODO: verify cryptographic signature of claim

	return nil

}
