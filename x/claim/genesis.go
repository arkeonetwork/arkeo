package claim

import (
	"github.com/arkeonetwork/arkeo/x/claim/keeper"
	"github.com/arkeonetwork/arkeo/x/claim/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
	k.SetClaimRecords(ctx, genState.ClaimRecords)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	//genesis.ModuleAccountBalance = k.GetModuleAccountBalance(ctx)
	claimRecords, err := k.GetAllClaimRecords(ctx)
	if err != nil {
		panic(err)
	}
	genesis.ClaimRecords = claimRecords
	return genesis
}
