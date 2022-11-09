package arkeo_test

import (
	keepertest "arkeo/testutil/keeper"
	"arkeo/testutil/nullify"
	"arkeo/x/arkeo"
	"arkeo/x/arkeo/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	ctx, k := keepertest.ArkeoKeeper(t)
	arkeo.InitGenesis(ctx, k, genesisState)
	got := arkeo.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
