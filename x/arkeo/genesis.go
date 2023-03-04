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
			ctx.Logger().Error("unable to set contract", "provider", contract.ProviderPubKey, "chain", contract.Chain, "client", contract.Client, "error", err)
		}
	}
	k.SetNextContractId(ctx, genState.NextContractId)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

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
	return genesis
}
