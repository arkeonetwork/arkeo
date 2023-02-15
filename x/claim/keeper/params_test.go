package keeper_test

import (
	"testing"

	testkeeper "github.com/arkeonetwork/arkeo/testutil/keeper"
	"github.com/arkeonetwork/arkeo/x/claim/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.ClaimKeeper(t)
	params := types.DefaultParams()
	params.ClaimDenom = "Test!"
	k.SetParams(ctx, params)
	got := k.GetParams(ctx)
	require.EqualValues(t, params, got)
}
