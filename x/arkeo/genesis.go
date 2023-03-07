package arkeo

import (
	"github.com/arkeonetwork/arkeo/x/arkeo/keeper"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)

	for _, provider := range genState.Providers {
		if err := k.SetProvider(ctx, provider); err != nil {
			ctx.Logger().Error("unable to set provider", "provider", provider.PubKey, "chain", provider.Chain, "error", err)
		}
	}

	for _, contract := range genState.Contracts {
		if err := k.SetContract(ctx, contract); err != nil {
			ctx.Logger().Error("unable to set contract", "provider", contract.Provider, "chain", contract.Chain, "client", contract.Client, "error", err)
		}
	}
	k.SetNextContractId(ctx, genState.NextContractId)

	for _, expirationSet := range genState.ContractExpirationSets {
		if err := k.SetContractExpirationSet(ctx, expirationSet); err != nil {
			ctx.Logger().Error("unable to set contract expiration set", "height", expirationSet.Height, "error", err)
		}
	}

	for _, userContractSet := range genState.UserContractSets {
		if err := k.SetUserContractSet(ctx, userContractSet); err != nil {
			ctx.Logger().Error("unable to set user contract set", "user", userContractSet.User, "error", err)
		}
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// providers
	iter := k.GetProviderIterator(ctx)
	for ; iter.Valid(); iter.Next() {
		var provider types.Provider
		if err := k.Cdc().Unmarshal(iter.Value(), &provider); err != nil {
			ctx.Logger().Error("unable to get provider", "provider", iter.Key(), "error", err)
			continue
		}
		genesis.Providers = append(genesis.Providers, provider)
	}
	iter.Close()

	// contracts
	iter = k.GetContractIterator(ctx)
	for ; iter.Valid(); iter.Next() {
		var contract types.Contract
		if err := k.Cdc().Unmarshal(iter.Value(), &contract); err != nil {
			ctx.Logger().Error("unable to get contract", "contract", iter.Key(), "error", err)
			continue
		}
		genesis.Contracts = append(genesis.Contracts, contract)
	}
	iter.Close()
	genesis.NextContractId = k.GetNextContractId(ctx)

	// contract expiration sets
	iter = k.GetContractExpirationSetIterator(ctx)
	for ; iter.Valid(); iter.Next() {
		var expirationSet types.ContractExpirationSet
		if err := k.Cdc().Unmarshal(iter.Value(), &expirationSet); err != nil {
			ctx.Logger().Error("unable to get contract expiration set", "contract", iter.Key(), "error", err)
			continue
		}
		genesis.ContractExpirationSets = append(genesis.ContractExpirationSets, expirationSet)
	}

	// user contract sets
	iter = k.GetUserContractSetIterator(ctx)
	for ; iter.Valid(); iter.Next() {
		var userContractSet types.UserContractSet
		if err := k.Cdc().Unmarshal(iter.Value(), &userContractSet); err != nil {
			ctx.Logger().Error("unable to get user contract set", "contract", iter.Key(), "error", err)
			continue
		}
		genesis.UserContractSets = append(genesis.UserContractSets, userContractSet)
	}

	return genesis
}
