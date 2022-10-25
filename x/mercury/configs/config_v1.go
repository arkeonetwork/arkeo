package configs

// NewConfigValue010 get new instance of ConfigValue010
func NewConfigValue010() *ConfigVals {
	return &ConfigVals{
		int64values: map[ConfigName]int64{
			GasFee:              1_00000000, // number of tokens for gas fee
			RegisterProviderFee: 1_00000000, // number of tokens in fee to register a new provider
		},
		boolValues:   map[ConfigName]bool{},
		stringValues: map[ConfigName]string{},
	}
}
