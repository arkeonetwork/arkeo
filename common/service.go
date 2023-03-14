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
	EmptyService Service = iota
	StarWarsService
	BaseService
	BTCService
	ETHService
	GAIAChainRPCArchiveService
)

var ServiceLookup = map[string]int32{
	"unknown":   0,
	"swapi.dev": 1, // star wars API for development purposes
	"github.com/arkeonetwork/arkeo-mainnet-fullnode": 2,
	"btc-mainnet-fullnode":                           3,
	"eth-mainnet-fullnode":                           4,
	"gaia-mainnet-rpc-archive":                       5,
}

func (c Service) String() string {
	switch c {
	case BaseService:
		return "arkeo-mainnet-fullnode"
	case BTCService:
		return "btc-mainnet-fullnode"
	case ETHService:
		return "eth-mainnet-fullnode"
	case StarWarsService:
		return "swapi.dev"
	case GAIAChainRPCArchiveService:
		return "gaia-mainnet-rpc-archive"
	default:
		return "unknown"
	}
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
