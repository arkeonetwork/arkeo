package crosstransfer_test

import (
	"arkeo/testutil/nullify"
	"arkeo/x/crosstransfer"
	"arkeo/x/crosstransfer/types"
	"testing"

	keepertest "arkeo/testutil/keeper"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.CrosstransferKeeper(t)
	crosstransfer.InitGenesis(ctx, *k, genesisState)
	got := crosstransfer.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	// this line is used by starport scaffolding # genesis/test/assert
}
