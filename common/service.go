package common

import (
	"fmt"
	"strings"
)

type (
	Service  int32
	Services []Service
)

const (
	EmptyService    Service = iota
	StarWarsService Service = 1
	BTCService      Service = 10
	ETHService      Service = 20
)

var ServiceLookup = map[string]int32{
	"unknown":                      0,
	"swapi.dev":                    1, // star wars API for development purposes
	"arkeo-mainnet-fullnode":       2,
	"avax-mainnet-fullnode":        3,
	"avax-mainnet-archivenode":     4,
	"bch-mainnet-fullnode":         5,
	"bch-mainnet-lightnode":        6,
	"bnb-mainnet-fullnode":         7,
	"bsc-mainnet-fullnode":         8,
	"bsc-mainnet-archivenode":      9,
	"btc-mainnet-fullnode":         10,
	"btc-mainnet-lightnode":        11,
	"cardano-mainnet-relaynode":    12,
	"gaia-mainnet-rpc":             13,
	"doge-mainnet-fullnode":        14,
	"doge-mainnet-lightnode":       15,
	"etc-mainnet-archivenode":      16,
	"etc-mainnet-fullnode":         17,
	"etc-mainnet-lightnode":        18,
	"eth-mainnet-archivenode":      19,
	"eth-mainnet-fullnode":         20,
	"eth-mainnet-lightnode":        21,
	"ltc-mainnet-fullnode":         22,
	"ltc-mainnet-lightnode":        23,
	"optimism-mainnet-fullnode":    24,
	"osmosis-mainnet-fullnode":     25,
	"polkadot-mainnet-fullnode":    26,
	"polkadot-mainnet-lightnode":   27,
	"polkadot-mainnet-archivenode": 28,
	"polygon-mainnet-fullnode":     29,
	"polygon-mainnet-archivenode":  30,
	"sol-mainnet-fullnode":         31,
	"thorchain-mainnet-fullnode":   32,
	"bch-mainnet-unchained":        33,
	"btc-mainnet-unchained":        34,
	"bsc-mainnet-unchained":        36,
	"gaia-mainnet-unchained":       38,
	"doge-mainnet-unchained":       39,
	"eth-mainnet-unchained":        40,
	"avax-mainnet-unchained":       41,
	"ltc-mainnet-unchained":        42,
	"osmosis-mainnet-unchained":    43,
	"thorchain-mainnet-unchained":  44,
	"optimism-mainnet-unchained":   45,
	"gaia-mainnet-grpc":            46,
}

var ServiceReverseLookup = map[Service]string{
	0:  "unknown",
	1:  "swapi.dev",
	2:  "arkeo-mainnet-fullnode",
	3:  "avax-mainnet-fullnode",
	4:  "avax-mainnet-archivenode",
	5:  "bch-mainnet-fullnode",
	6:  "bch-mainnet-lightnode",
	7:  "bnb-mainnet-fullnode",
	8:  "bsc-mainnet-fullnode",
	9:  "bsc-mainnet-archivenode",
	10: "btc-mainnet-fullnode",
	11: "btc-mainnet-lightnode",
	12: "cardano-mainnet-relaynode",
	13: "gaia-mainnet-rpc",
	14: "doge-mainnet-fullnode",
	15: "doge-mainnet-lightnode",
	16: "etc-mainnet-archivenode",
	17: "etc-mainnet-fullnode",
	18: "etc-mainnet-lightnode",
	19: "eth-mainnet-archivenode",
	20: "eth-mainnet-fullnode",
	21: "eth-mainnet-lightnode",
	22: "ltc-mainnet-fullnode",
	23: "ltc-mainnet-lightnode",
	24: "optimism-mainnet-fullnode",
	25: "osmosis-mainnet-fullnode",
	26: "polkadot-mainnet-fullnode",
	27: "polkadot-mainnet-lightnode",
	28: "polkadot-mainnet-archivenode",
	29: "polygon-mainnet-fullnode",
	30: "polygon-mainnet-archivenode",
	31: "sol-mainnet-fullnode",
	32: "thorchain-mainnet-fullnode",
	33: "bch-mainnet-unchained",
	34: "btc-mainnet-unchained",
	35: "bnb-mainnet-unchained",
	36: "bsc-mainnet-unchained",
	37: "bsc-mainnet-unchained",
	38: "gaia-mainnet-unchained",
	39: "doge-mainnet-unchained",
	40: "eth-mainnet-unchained",
	41: "avax-mainnet-unchained",
	42: "ltc-mainnet-unchained",
	43: "osmosis-mainnet-unchained",
	44: "thorchain-mainnet-unchained",
	45: "optimism-mainnet-unchained",
	46: "gaia-mainnet-grpc",
}

func (service Service) String() string {
	if r, ok := ServiceReverseLookup[service]; ok {
		return r
	}
	return "unknown"
}

// NewService create a new service
func NewService(serviceId string) (Service, error) {
	service := ServiceLookup[strings.ToLower(serviceId)]
	if service == 0 {
		return Service(service), fmt.Errorf("service not found (%s)", serviceId)
	}
	return Service(service), nil
}

// Equals compare two services to see whether they represent the same service
func (c Service) Equals(c2 Service) bool {
	return strings.EqualFold(c.String(), c2.String())
}

// IsEmpty is to determinate whether the service is empty
func (service Service) IsEmpty() bool {
	return strings.TrimSpace(service.String()) == ""
}
