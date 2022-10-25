package configs

// NewConfigValue010 get new instance of ConfigValue010
func NewConfigValue010() *ConfigVals {
	return &ConfigVals{
		int64values: map[ConfigName]int64{
			GasFee: 1_00000000, // number of tokens for gas fee
		},
		boolValues:   map[ConfigName]bool{},
		stringValues: map[ConfigName]string{},
	}
}
