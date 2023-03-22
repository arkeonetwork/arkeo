package types

import (
	"encoding/json"
	"strconv"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeProviderBond       = "provider_bond"
	EventTypeProviderMod        = "provider_mod"
	EventTypeOpenContract       = "open_contract"
	EventTypeCloseContract      = "close_contract"
	EventTypeContractSettlement = "contract_settlement"
	EventTypeValidatorPayout    = "validator_payout"
)

func NewOpenContractEvent(openCost int64, contract *Contract) sdk.Event {
	return sdk.NewEvent(
		EventTypeOpenContract,
		sdk.NewAttribute("provider", contract.Provider.String()),
		sdk.NewAttribute("contract_id", strconv.FormatUint(contract.Id, 10)),
		sdk.NewAttribute("service", contract.Service.String()),
		sdk.NewAttribute("client", contract.Client.String()),
		sdk.NewAttribute("delegate", contract.Delegate.String()),
		sdk.NewAttribute("meter_type", contract.MeterType.String()),
		sdk.NewAttribute("user_type", contract.UserType.String()),
		sdk.NewAttribute("height", strconv.FormatInt(contract.Height, 10)),
		sdk.NewAttribute("duration", strconv.FormatInt(contract.Duration, 10)),
		sdk.NewAttribute("rate", strconv.FormatInt(contract.Rate, 10)),
		sdk.NewAttribute("open_cost", strconv.FormatInt(openCost, 10)),
		sdk.NewAttribute("deposit", contract.Deposit.String()),
		sdk.NewAttribute("settlement_duration", strconv.FormatInt(contract.SettlementDuration, 10)),
	)
}

func NewContractSettlementEvent(debt cosmos.Int, valIncome cosmos.Int, contract *Contract) sdk.Event {
	nonce := int64(0)
	if contract.Nonces != nil {
		// for now we only support one nonce, but we can extend once we handle
		// downstream dependencies
		nonce = contract.Nonces[contract.GetSpender()]
	}
	return sdk.NewEvent(
		EventTypeContractSettlement,
		sdk.NewAttribute("provider", contract.Provider.String()),
		sdk.NewAttribute("contract_id", strconv.FormatUint(contract.Id, 10)),
		sdk.NewAttribute("service", contract.Service.String()),
		sdk.NewAttribute("client", contract.Client.String()),
		sdk.NewAttribute("delegate", contract.Delegate.String()),
		sdk.NewAttribute("meter_type", contract.MeterType.String()),
		sdk.NewAttribute("user_type", contract.UserType.String()),
		sdk.NewAttribute("height", strconv.FormatInt(contract.Height, 10)),
		sdk.NewAttribute("paid", debt.String()),
		sdk.NewAttribute("reserve", valIncome.String()),
		sdk.NewAttribute("nonce", strconv.FormatInt(nonce, 10)),
	)

}

func NewCloseContractEvent(contract *Contract) sdk.Event {
	return sdk.NewEvent(
		EventTypeCloseContract,
		sdk.NewAttribute("contract_id", strconv.FormatUint(contract.Id, 10)),
		sdk.NewAttribute("provider", contract.Provider.String()),
		sdk.NewAttribute("service", contract.Service.String()),
		sdk.NewAttribute("client", contract.Client.String()),
		sdk.NewAttribute("delegate", contract.Delegate.String()),
	)
}

func NewBondProviderEvent(bond cosmos.Int, msg *MsgBondProvider) sdk.Event {
	return sdk.NewEvent(
		EventTypeProviderBond,
		sdk.NewAttribute("provider", msg.Provider.String()),
		sdk.NewAttribute("service", msg.Service),
		sdk.NewAttribute("bond_rel", msg.Bond.String()),
		sdk.NewAttribute("bond_abs", bond.String()),
	)
}

func NewModProviderEvent(provider *Provider) sdk.Event {
	return sdk.NewEvent(
		EventTypeProviderMod,
		sdk.NewAttribute("provider", provider.PubKey.String()),
		sdk.NewAttribute("service", provider.Service.String()),
		sdk.NewAttribute("metadata_uri", provider.MetadataUri),
		sdk.NewAttribute("metadata_nonce", strconv.FormatUint(provider.MetadataNonce, 10)),
		sdk.NewAttribute("status", provider.Status.String()),
		sdk.NewAttribute("min_contract_duration", strconv.FormatInt(provider.MinContractDuration, 10)),
		sdk.NewAttribute("max_contract_duration", strconv.FormatInt(provider.MaxContractDuration, 10)),
		sdk.NewAttribute("bond", provider.Bond.String()),
		sdk.NewAttribute("settlement_duration", strconv.FormatInt(provider.SettlementDuration, 10)),
		getRatesAttribute(provider.Rates),
	)
}

func NewValidatorPayoutEvent(acc cosmos.AccAddress, reward cosmos.Int) sdk.Event {
	return sdk.NewEvent(
		EventTypeValidatorPayout,
		sdk.NewAttribute("validator", acc.String()),
		sdk.NewAttribute("reward", reward.String()),
	)
}

func getRatesAttribute(rates []*ContractRate) sdk.Attribute {
	rateString, err := json.Marshal(rates)
	if err != nil {
		return sdk.NewAttribute("rates", "error")
	}
	return sdk.NewAttribute("rates", string(rateString))
}
