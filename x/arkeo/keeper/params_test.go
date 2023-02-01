package keeper_test

import (
	testkeeper "arkeo/testutil/keeper"
	"testing"

	"github.com/ArkeoNetwork/arkeo-protocol/x/arkeo/types"

	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	ctx, k := testkeeper.ArkeoKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
