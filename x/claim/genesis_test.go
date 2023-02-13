package claim_test

import (
	"testing"
	"time"

	keepertest "github.com/arkeonetwork/arkeo/testutil/keeper"
	"github.com/arkeonetwork/arkeo/testutil/nullify"
	"github.com/arkeonetwork/arkeo/x/claim"
	"github.com/arkeonetwork/arkeo/x/claim/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {

	airdropStartTime := time.Now().UTC()
	claimParams := types.Params{
		AirdropStartTime:   airdropStartTime,
		DurationUntilDecay: types.DefaultDurationUntilDecay,
		DurationOfDecay:    types.DefaultDurationOfDecay,
		ClaimDenom:         types.DefaultClaimDenom,
	}
	genesisState := types.GenesisState{
		Params: claimParams,
	}

	k, ctx := keepertest.ClaimKeeper(t)
	claim.InitGenesis(ctx, k, genesisState)
	got := claim.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
