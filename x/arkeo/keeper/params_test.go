package keeper_test

import (
	"testing"

	testkeeper "github.com/ArkeoNetwork/arkeo-protocol/testutil/keeper"

	"github.com/ArkeoNetwork/arkeo-protocol/x/arkeo/types"

	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	ctx, k := testkeeper.ArkeoKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
