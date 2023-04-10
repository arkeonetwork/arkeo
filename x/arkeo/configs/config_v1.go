package configs

import "github.com/arkeonetwork/arkeo/common"

// NewConfigValue010 get new instance of ConfigValue010
func NewConfigValue010() *ConfigVals {
	return &ConfigVals{
		int64values: map[ConfigName]int64{
			HandlerBondProvider:        0,                          // enable/disable bond provider handler
			HandlerModProvider:         0,                          // enable/disable mod provider handler
			HandlerOpenContract:        0,                          // enable/disable open contract handler
			HandlerCloseContract:       0,                          // enable/disable close contract handler
			HandlerClaimContractIncome: 0,                          // enable/disable claim contract income handler
			MaxContractLength:          5256000,                    // one year
			MaxSupply:                  common.Tokens(121_000_000), // max supply of tokens
			OpenContractCost:           common.Tokens(1),           // cost to open a contract
			MinProviderBond:            common.Tokens(1),           // min bond for a data provider to be able to open contracts with
			ReserveTax:                 1000,                       // reserve income off provider income, in basis points
			BlocksPerYear:              5256666,                    // blocks per year
			EmissionCurve:              4,                          // rate in which the reserve is depleted to pay validators
			ValidatorPayoutCycle:       1,                          // how often validators are paid out rewards
		},
		boolValues:   map[ConfigName]bool{},
		stringValues: map[ConfigName]string{},
	}
}
