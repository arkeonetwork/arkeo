package keeper_test

import (
	testkeeper "mercury/testutil/keeper"
	"mercury/x/mercury/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	ctx, k := testkeeper.MercuryKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
