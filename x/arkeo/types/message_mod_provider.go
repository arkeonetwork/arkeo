package types

import (
	"cosmossdk.io/errors"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgModProvider = "mod_provider"

var _ sdk.Msg = &MsgModProvider{}

func NewMsgModProvider(creator string, provider common.PubKey, service, metadataUri string,
	metadataNonce uint64, status ProviderStatus, minContractDuration,
	maxContractDuration int64, subscriptionRate, payAsYouGoRate types.Coins, settlementDuration int64,
) *MsgModProvider {
	return &MsgModProvider{
		Creator:             creator,
		Provider:            provider,
		Service:             service,
		MetadataUri:         metadataUri,
		MetadataNonce:       metadataNonce,
		Status:              status,
		MinContractDuration: minContractDuration,
		MaxContractDuration: maxContractDuration,
		SubscriptionRate:    subscriptionRate,
		PayAsYouGoRate:      payAsYouGoRate,
		SettlementDuration:  settlementDuration,
	}
}

func (msg *MsgModProvider) Route() string {
	return RouterKey
}

func (msg *MsgModProvider) Type() string {
	return TypeMsgModProvider
}

func (msg *MsgModProvider) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgModProvider) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}

func (msg *MsgModProvider) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgModProvider) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// verify pubkey
	_, err = common.NewPubKey(msg.Provider.String())
	if err != nil {
		return errors.Wrapf(ErrInvalidPubKey, "invalid provider pubkey (%s): %s", msg.Provider, err)
	}

	// verify service
	_, err = common.NewService(msg.Service)
	if err != nil {
		return errors.Wrapf(ErrInvalidService, "invalid service (%s): %s", msg.Service, err)
	}

	signer := msg.MustGetSigner()
	provider, err := msg.Provider.GetMyAddress()
	if err != nil {
		return err
	}
	if !signer.Equals(provider) {
		return errors.Wrapf(ErrProviderBadSigner, "Signer: %s, Provider Address: %s", msg.GetSigners(), provider)
	}

	// test metadataURI
	/*
		Disabling URI parsing check due to a potential that the underlying golang code may change its behavior between golang versions. We can assume data providers are giving valid URIs, because if they aren't, they won't be able to make income
		if _, err := url.ParseRequestURI(msg.MetadataURI); err != nil {
			return errors.Wrapf(ErrInvalidModProviderMetdataURI, "(%s)", err)
		}
	*/
	// Ensure URIs don't get too long and cause chain bloat
	if len(msg.MetadataUri) > 100 {
		return errors.Wrapf(ErrInvalidModProviderMetdataURI, "length is too long (%d/100)", len(msg.MetadataUri))
	}

	// check durations
	if msg.MinContractDuration <= 0 {
		return errors.Wrapf(ErrInvalidModProviderMinContractDuration, "min contraction duration cannot be zero")
	}

	if msg.MinContractDuration > msg.MaxContractDuration {
		return errors.Wrapf(ErrInvalidModProviderMinContractDuration, "min contract duration is too long (%d/%d)", msg.MaxContractDuration, msg.MaxContractDuration)
	}

	if msg.SettlementDuration < 0 {
		return errors.Wrapf(ErrInvalidModProviderSettlementDuration, "settlement duration cannot be negative")
	}

	subRate := cosmos.NewCoins(msg.SubscriptionRate...)
	if err := subRate.Validate(); err != nil {
		return errors.Wrapf(err, "invalid subscription rate")
	}

	if !subRate.IsAllPositive() {
		return errors.Wrapf(ErrInvalidModProviderRate, "all subscription rates must be positive")
	}

	payRate := cosmos.NewCoins(msg.PayAsYouGoRate...)
	if err := payRate.Validate(); err != nil {
		return errors.Wrapf(err, "invalid subscription rate")
	}

	if !payRate.IsAllPositive() {
		return errors.Wrapf(ErrInvalidModProviderRate, "all pay-as-you-go rates must be positive")
	}

	return nil
}
