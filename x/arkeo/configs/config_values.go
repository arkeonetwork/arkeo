package configs

import (
	"fmt"
)

// ConfigName the name we used to get constant values
type ConfigName int

const (
	HandlerBondProvider ConfigName = iota
	HandlerModProvider
	HandlerOpenContract
	HandlerCloseContract
	HandlerClaimContractIncome
	HandlerSetVersion
	MaxSupply
	MaxContractLength
	OpenContractCost
	MinProviderBond
	ReserveTax
	BlocksPerYear
	EmissionCurve
	ValidatorPayoutCycle
	VersionConsensus
)

var nameToString = map[ConfigName]string{
	HandlerBondProvider:        "HandlerBondProvider",
	HandlerModProvider:         "HandlerModProvider",
	HandlerOpenContract:        "HandlerOpenContract",
	HandlerCloseContract:       "HandlerCloseContract",
	HandlerClaimContractIncome: "HandlerClaimContractIncome",
	HandlerSetVersion:          "HandlerSetVersion",
	MaxSupply:                  "MaxSupply",
	MaxContractLength:          "MaxContractLength",
	OpenContractCost:           "OpenContractCost",
	MinProviderBond:            "MinProviderBond",
	ReserveTax:                 "ReserveTax",
	BlocksPerYear:              "BlocksPerYear",
	EmissionCurve:              "EmissionCurve",
	ValidatorPayoutCycle:       "ValidatorPayoutCycle",
	VersionConsensus:           "VersionConsensus",
}

// String implement fmt.stringer
func (cn ConfigName) String() string {
	val, ok := nameToString[cn]
	if !ok {
		return "NA"
	}
	return val
}

// ConfigValues define methods used to get constant values
type ConfigValues interface {
	fmt.Stringer
	GetInt64Value(name ConfigName) int64
	GetBoolValue(name ConfigName) bool
	GetStringValue(name ConfigName) string
}

// GetConfigValues will return an  implementation of ConfigValues which provide ways to get constant values
func GetConfigValues(ver int64) ConfigValues {
	if ver > 0 {
		return NewConfigValue010()
	}
	return nil
}
