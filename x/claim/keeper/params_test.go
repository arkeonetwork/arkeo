package keeper_test

import (
	"testing"

	testkeeper "github.com/arkeonetwork/arkeo/testutil/keeper"
	"github.com/arkeonetwork/arkeo/x/claim/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	keepers, ctx := testkeeper.CreateTestClaimKeepers(t)
	params := types.DefaultParams()
	params.ClaimDenom = "Test!"
	keepers.ClaimKeeper.SetParams(ctx, params)
	got := keepers.ClaimKeeper.GetParams(ctx)
	require.EqualValues(t, params, got)
}
