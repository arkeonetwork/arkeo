package sentinel

import (
	"fmt"
	"strconv"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

type ProviderBondEvent struct {
	PubKey       common.PubKey
	Service      common.Service
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
		case "service":
			evt.Service, err = common.NewService(v)
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
		case "service":
			evt.Provider.Service, err = common.NewService(v)
			if err != nil {
				return evt, err
			}
		case "metadata_uri":
			evt.Provider.MetadataUri = v
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
			evt.Provider.SubscriptionRate, err = cosmos.ParseCoins(v)
			if err != nil {
				return evt, err
			}
		case "pay-as-you-go_rate":
			evt.Provider.PayAsYouGoRate, err = cosmos.ParseCoins(v)
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

type ClaimContractIncome struct {
	Contract types.Contract
	Paid     cosmos.Int
	Reserve  cosmos.Int
}

func parseContractSettlementEvent(input map[string]string) (ClaimContractIncome, error) {
	var err error
	var ok bool
	evt := ClaimContractIncome{}

	for k, v := range input {
		switch k {
		case "provider":
			evt.Contract.Provider, err = common.NewPubKey(v)
			if err != nil {
				return evt, err
			}
		case "contract_id":
			evt.Contract.Id, err = strconv.ParseUint(v, 10, 64)
			if err != nil {
				return evt, err
			}
		case "service":
			evt.Contract.Service, err = common.NewService(v)
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
