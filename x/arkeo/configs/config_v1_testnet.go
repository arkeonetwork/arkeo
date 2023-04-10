//go:build testnet
// +build testnet

// For internal testing and mockneting
package configs

import "github.com/arkeonetwork/arkeo/common"

func init() {
	int64Overrides = map[ConfigName]int64{
		MaxSupply: common.Tokens(1_000_000_000),
	}
}
