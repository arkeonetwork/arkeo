package types

import (
	"github.com/arkeonetwork/arkeo/common/cosmos"
)

const (
	EventTypeProviderBond       = "arkeo.arkeo.EventProviderBond"
	EventTypeProviderMod        = "arkeo.arkeo.EventProviderMod"
	EventTypeOpenContract       = "arkeo.arkeo.EventOpenContract"
	EventTypeCloseContract      = "arkeo.arkeo.Event.CloseContract"
	EventTypeContractSettlement = "arkeo.arkeo.EventContractSettlement"
	EventTypeValidatorPayout    = "arkeo.arkeo.EventValidatorPayout"
)

func NewOpenContractEvent(openCost int64, contract *Contract) EventOpenContract {
	return EventOpenContract{
		Provider:           contract.Provider,
		ContractId:         contract.Id,
		Service:            contract.Service.String(),
		Client:             contract.Client,
		Delegate:           contract.Delegate,
		Type:               contract.Type,
		Height:             contract.Height,
		Duration:           contract.Duration,
		Rate:               contract.Rate,
		OpenCost:           openCost,
		Deposit:            contract.Deposit,
		SettlementDuration: contract.SettlementDuration,
		Authorization:      contract.Authorization,
		QueriesPerMinute:   contract.QueriesPerMinute,
	}
}

func NewContractSettlementEvent(debt, valIncome cosmos.Int, contract *Contract) EventSettleContract {
	return EventSettleContract{
		Provider:   contract.Provider,
		ContractId: contract.Id,
		Service:    contract.Service.String(),
		Client:     contract.Client,
		Delegate:   contract.Delegate,
		Type:       contract.Type,
		Nonce:      contract.Nonce,
		Height:     contract.Height,
		Paid:       debt,
		Reserve:    valIncome,
	}
}

func NewCloseContractEvent(contract *Contract) EventCloseContract {
	return EventCloseContract{
		ContractId: contract.Id,
		Provider:   contract.Provider,
		Service:    contract.Service.String(),
		Client:     contract.Client,
		Delegate:   contract.Delegate,
	}
}

func NewBondProviderEvent(bond cosmos.Int, msg *MsgBondProvider) EventBondProvider {
	return EventBondProvider{
		Provider: msg.Provider,
		Service:  msg.Service,
		BondRel:  msg.Bond,
		BondAbs:  bond,
	}
}

func NewValidatorPayoutEvent(acc cosmos.AccAddress, reward cosmos.Int) ValidatorPayoutEvent {
	return ValidatorPayoutEvent{
		Validator: acc,
		Reward:    reward,
	}
}
