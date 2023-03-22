package types

import (
	"cosmossdk.io/errors"

	"github.com/arkeonetwork/arkeo/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgModProvider = "mod_provider"

var _ sdk.Msg = &MsgModProvider{}

func NewMsgModProvider(creator string, provider common.PubKey, service, metadataUri string,
	metadataNonce uint64, status ProviderStatus, minContractDuration,
	maxContractDuration, settlementDuration int64, rates []*ContractRate,
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
		Rates:               rates,
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

	// confirm no duplicated rates and valid rate
	for i, rate := range msg.Rates {
		if rate.Rate <= 0 {
			return errors.Wrapf(ErrInvalidModProviderRate, "rate cannot be equal to or less than zero")
		}
		for ii, rateMatch := range msg.Rates {
			if i == ii {
				continue
			}

			if rateMatch.MeterType == rate.MeterType && rateMatch.UserType == rate.UserType {
				return errors.Wrapf(ErrInvalidModProviderDuplicateContractRates, "duplicated rate for meter type %s with user type %s", rate.MeterType, rate.UserType)
			}
		}
	}

	return nil
}
