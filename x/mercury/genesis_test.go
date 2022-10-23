package mercury_test

import (
	keepertest "mercury/testutil/keeper"
	"mercury/testutil/nullify"
	"mercury/x/mercury"
	"mercury/x/mercury/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.MercuryKeeper(t)
	mercury.InitGenesis(ctx, k, genesisState)
	got := mercury.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
