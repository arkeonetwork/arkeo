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
	StarWarsChain
	BaseChain
	BTCChain
	ETHChain
)

var ChainLookup = map[string]int32{
	"unknown":   0,
	"swapi.dev": 1, // star wars API for development purposes
	"github.com/arkeonetwork/arkeo-mainnet-fullnode": 2,
	"btc-mainnet-fullnode":                           3,
	"eth-mainnet-fullnode":                           4,
}

func (c Chain) String() string {
	switch c {
	case BaseChain:
		return "arkeo-mainnet-fullnode"
	case BTCChain:
		return "btc-mainnet-fullnode"
	case ETHChain:
		return "eth-mainnet-fullnode"
	case StarWarsChain:
		return "swapi.dev"
	default:
		return "unknown"
	}
}

// ChainNetwork is to indicate which chain environment
type ChainNetwork uint8

// NewChain create a new Chain
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
