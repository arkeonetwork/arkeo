package common

import (
	"fmt"
	"strings"
)

type (
	Chain  int32
	Chains []Chain
)

const (
	EmptyChain Chain = iota
	BaseChain
	BTCChain
	ETHChain
)

var ChainLookup = map[string]int32{
	"unknown":         0,
	"mercury-mainnet": 1,
	"btc-mainnet":     2,
	"eth-mainnet":     3,
}

func (c Chain) String() string {
	switch c {
	case BaseChain:
		return "mercury-mainnet"
	case BTCChain:
		return "btc-mainnet"
	case ETHChain:
		return "eth-mainnet"
	default:
		return "unknown"
	}
}

// ChainNetwork is to indicate which chain environment
type ChainNetwork uint8

// NewChain create a new Chain and default the siging_algo to Secp256k1
func NewChain(chainID string) (Chain, error) {
	chain := ChainLookup[strings.ToLower(chainID)]
	if chain == 0 {
		return Chain(chain), fmt.Errorf("Chain not found (%s)", chainID)
	}
	return Chain(chain), nil
}

// Equals compare two chain to see whether they represent the same chain
func (c Chain) Equals(c2 Chain) bool {
	return strings.EqualFold(c.String(), c2.String())
}

// IsEmpty is to determinate whether the chain is empty
func (c Chain) IsEmpty() bool {
	return strings.TrimSpace(c.String()) == ""
}
