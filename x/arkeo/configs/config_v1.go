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
			HandlerSetVersion:          0,                          // enable/disable set version handler
			MaxContractLength:          5256000,                    // one year
			MaxSupply:                  common.Tokens(121_000_000), // max supply of tokens
			OpenContractCost:           20_000_000,                 // cost to open a contract (was common.Tokens(1))
			MinProviderBond:            common.Tokens(1),           // min bond for a data provider to be able to open contracts with
			ReserveTax:                 1000,                       // reserve income off provider income, in basis points
			BlocksPerYear:              6311520,                    // blocks per year
			EmissionCurve:              10,                         // rate in which the reserve is depleted to pay validators
			ValidatorPayoutCycle:       1,                          // how often validators are paid out rewards
			VersionConsensus:           90,                         // out of 100, percentage of nodes on a specific version before it is accepted
		},
		boolValues:   map[ConfigName]bool{},
		stringValues: map[ConfigName]string{},
	}
}
