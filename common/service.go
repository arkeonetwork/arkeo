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

var ServiceReverseLookup = map[Service]string{
	0: "unknown",
	1: "swapi.dev",
	2: "github.com/arkeonetwork/arkeo-mainnet-fullnode",
	3: "btc-mainnet-fullnode",
	4: "eth-mainnet-fullnode",
	5: "gaia-mainnet-rpc-archive",
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
