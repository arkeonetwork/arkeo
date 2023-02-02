package keeper_test

import (
	"testing"

	testkeeper "github.com/arkeonetwork/arkeo/testutil/keeper"

	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	ctx, k := testkeeper.ArkeoKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
