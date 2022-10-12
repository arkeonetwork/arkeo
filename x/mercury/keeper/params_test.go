package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "mercury/testutil/keeper"
	"mercury/x/mercury/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.MercuryKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
