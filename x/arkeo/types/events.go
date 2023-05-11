package types

import (
	"github.com/arkeonetwork/arkeo/common/cosmos"
)

const (
	EventTypeBondProvider    = "arkeo.arkeo.EventBondProvider"
	EventTypeModProvider     = "arkeo.arkeo.EventModProvider"
	EventTypeOpenContract    = "arkeo.arkeo.EventOpenContract"
	EventTypeSettleContract  = "arkeo.arkeo.EventSettleContract"
	EventTypeCloseContract   = "arkeo.arkeo.EventCloseContract"
	EventTypeValidatorPayout = "arkeo.arkeo.EventValidatorPayout"
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

func NewValidatorPayoutEvent(acc cosmos.AccAddress, reward cosmos.Int) EventValidatorPayout {
	return EventValidatorPayout{
		Validator: acc,
		Reward:    reward,
	}
}
