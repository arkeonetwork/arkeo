package configs

import "mercury/common"

// NewConfigValue010 get new instance of ConfigValue010
func NewConfigValue010() *ConfigVals {
	return &ConfigVals{
		int64values: map[ConfigName]int64{
			HandlerBondProvider: 0,                // enable/disable bond provider handler
			HandlerModProvider:  0,                // enable/disable mod provider handler
			HandlerOpenContract: 0,                // enable/disable open contract handler
			MaxContractLength:   5256000,          // one year
			OpenContractCost:    common.Tokens(1), // cost to open a contract
			MinProviderBond:     common.Tokens(1), // min bond for a data provider to be able to open contracts with
		},
		boolValues:   map[ConfigName]bool{},
		stringValues: map[ConfigName]string{},
	}
}
