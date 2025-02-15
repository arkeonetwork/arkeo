package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:                 DefaultParams(),
		Providers:              make([]Provider, 0),
		Contracts:              make([]Contract, 0),
		NextContractId:         0,
		ContractExpirationSets: make([]ContractExpirationSet, 0),
		UserContractSets:       make([]UserContractSet, 0),
		Version:                1,
		ValidatorVersions:      make([]ValidatorVersion, 0),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	if err := gs.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	seenProviders := make(map[string]bool)
	for _, provider := range gs.Providers {
		key := fmt.Sprintf("%s-%s", provider.PubKey, provider.Service)
		if seenProviders[key] {
			return fmt.Errorf("duplicate provider found: %s", key)
		}
		seenProviders[key] = true
	}

	seenContracts := make(map[uint64]bool)
	for _, contract := range gs.Contracts {
		if seenContracts[contract.Id] {
			return fmt.Errorf("duplicate contract ID found: %d", contract.Id)
		}
		seenContracts[contract.Id] = true
	}

	seenValidators := make(map[string]bool)
	for _, vv := range gs.ValidatorVersions {
		if seenValidators[vv.ValidatorAddress] {
			return fmt.Errorf("duplicate validator address in versions: %s", vv.ValidatorAddress)
		}
		seenValidators[vv.ValidatorAddress] = true

		// Validate validator address format
		if _, err := sdk.ValAddressFromBech32(vv.ValidatorAddress); err != nil {
			return fmt.Errorf("invalid validator address %s: %w", vv.ValidatorAddress, err)
		}

		// Validate version number
		if vv.Version <= 0 {
			return fmt.Errorf("invalid version for validator %s: %d", vv.ValidatorAddress, vv.Version)
		}
	}

	maxContractId := uint64(0)
	for _, contract := range gs.Contracts {
		if contract.Id > maxContractId {
			maxContractId = contract.Id
		}
	}
	if gs.NextContractId < maxContractId {
		return fmt.Errorf("NextContractId (%d) must be greater than the highest contract ID (%d)",
			gs.NextContractId, maxContractId)
	}

	// Validate version
	if gs.Version < 0 {
		return fmt.Errorf("invalid version: %d", gs.Version)
	}

	return nil
}
