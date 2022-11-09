package keeper_test

import (
	testkeeper "arkeo/testutil/keeper"
	"arkeo/x/arkeo/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	ctx, k := testkeeper.ArkeoKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
