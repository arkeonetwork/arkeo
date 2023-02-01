package sentinel

import (
	"fmt"
	"strconv"

	"github.com/ArkeoNetwork/arkeo-protocol/common"
	"github.com/ArkeoNetwork/arkeo-protocol/common/cosmos"
	"github.com/ArkeoNetwork/arkeo-protocol/x/arkeo/types"
)

type ProviderBondEvent struct {
	PubKey       common.PubKey
	Chain        common.Chain
	BondRelative cosmos.Int
	BondAbsolute cosmos.Int
}

// nolint
func parseProviderBondEvent(input map[string]string) (ProviderBondEvent, error) {
	var err error
	var ok bool
	evt := ProviderBondEvent{}

	for k, v := range input {
		switch k {
		case "pubkey":
			evt.PubKey, err = common.NewPubKey(v)
			if err != nil {
				return evt, err
			}
		case "chain":
			evt.Chain, err = common.NewChain(v)
			if err != nil {
				return evt, err
			}
		case "bond_rel":
			evt.BondRelative, ok = cosmos.NewIntFromString(v)
			if !ok {
				return evt, fmt.Errorf("cannot parse %s as int", v)
			}
		case "bond_abs":
			evt.BondAbsolute, ok = cosmos.NewIntFromString(v)
			if !ok {
				return evt, fmt.Errorf("cannot parse %s as int", v)
			}
		}
	}

	return evt, nil
}

type ProviderModEvent struct {
	Provider types.Provider
}

// nolint
func parseProviderModEvent(input map[string]string) (ProviderModEvent, error) {
	var err error
	evt := ProviderModEvent{}

	for k, v := range input {
		switch k {
		case "pubkey":
			evt.Provider.PubKey, err = common.NewPubKey(v)
			if err != nil {
				return evt, err
			}
		case "chain":
			evt.Provider.Chain, err = common.NewChain(v)
			if err != nil {
				return evt, err
			}
		case "metadata_uri":
			evt.Provider.MetadataURI = v
		case "metadata_nonce":
			evt.Provider.MetadataNonce, err = strconv.ParseUint(v, 10, 64)
			if err != nil {
				return evt, err
			}
		case "status":
			evt.Provider.Status = types.ProviderStatus(types.ProviderStatus_value[v])
			if err != nil {
				return evt, err
			}
		case "min_contract_duration":
			evt.Provider.MinContractDuration, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return evt, err
			}
		case "max_contract_duration":
			evt.Provider.MaxContractDuration, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return evt, err
			}

		case "subscription_rate":
			evt.Provider.SubscriptionRate, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return evt, err
			}
		case "pay-as-you-go_rate":
			evt.Provider.PayAsYouGoRate, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return evt, err
			}
		}
	}

	return evt, nil
}

type OpenContract struct {
	Contract types.Contract
	OpenCost int64
}

func parseOpenContract(input map[string]string) (OpenContract, error) {
	var err error
	evt := OpenContract{Contract: types.Contract{}}

	for k, v := range input {
		switch k {
		case "pubkey":
			evt.Contract.ProviderPubKey, err = common.NewPubKey(v)
			if err != nil {
				return evt, err
			}
		case "chain":
			evt.Contract.Chain, err = common.NewChain(v)
			if err != nil {
				return evt, err
			}
		case "client":
			evt.Contract.Client, err = common.NewPubKey(v)
			if err != nil {
				return evt, err
			}
		case "delegate":
			evt.Contract.Delegate, err = common.NewPubKey(v)
			if err != nil {
				return evt, err
			}
		case "type":
			evt.Contract.Type = types.ContractType(types.ContractType_value[v])
			if err != nil {
				return evt, err
			}
		case "height":
			evt.Contract.Height, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return evt, err
			}
		case "duration":
			evt.Contract.Duration, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return evt, err
			}
		case "rate":
			evt.Contract.Rate, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return evt, err
			}
		case "OpenCost":
			evt.OpenCost, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return evt, err
			}
		}
	}

	return evt, nil
}

type CloseContract struct {
	Contract types.Contract
}

func parseCloseContract(input map[string]string) (CloseContract, error) {
	var err error
	evt := CloseContract{}

	for k, v := range input {
		switch k {
		case "pubkey":
			evt.Contract.ProviderPubKey, err = common.NewPubKey(v)
			if err != nil {
				return evt, err
			}
		case "chain":
			evt.Contract.Chain, err = common.NewChain(v)
			if err != nil {
				return evt, err
			}
		case "client":
			evt.Contract.Client, err = common.NewPubKey(v)
			if err != nil {
				return evt, err
			}
		case "delegate":
			evt.Contract.Delegate, err = common.NewPubKey(v)
			if err != nil {
				return evt, err
			}
		}
	}

	return evt, nil
}

type ClaimContractIncome struct {
	Contract types.Contract
	Paid     cosmos.Int
	Reserve  cosmos.Int
}

func parseClaimContractIncome(input map[string]string) (ClaimContractIncome, error) {
	var err error
	var ok bool
	evt := ClaimContractIncome{}

	for k, v := range input {
		switch k {
		case "pubkey":
			evt.Contract.ProviderPubKey, err = common.NewPubKey(v)
			if err != nil {
				return evt, err
			}
		case "chain":
			evt.Contract.Chain, err = common.NewChain(v)
			if err != nil {
				return evt, err
			}
		case "client":
			evt.Contract.Client, err = common.NewPubKey(v)
			if err != nil {
				return evt, err
			}
		case "delegate":
			evt.Contract.Delegate, err = common.NewPubKey(v)
			if err != nil {
				return evt, err
			}
		case "type":
			evt.Contract.Type = types.ContractType(types.ContractType_value[v])
			if err != nil {
				return evt, err
			}
		case "height":
			evt.Contract.Height, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return evt, err
			}
		case "nonce":
			evt.Contract.Nonce, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return evt, err
			}
		case "paid":
			evt.Paid, ok = cosmos.NewIntFromString(v)
			if !ok {
				return evt, fmt.Errorf("cannot parse %s as int", v)
			}
		case "reserve":
			evt.Reserve, ok = cosmos.NewIntFromString(v)
			if !ok {
				return evt, fmt.Errorf("cannot parse %s as int", v)
			}
		}
	}

	return evt, nil
}
