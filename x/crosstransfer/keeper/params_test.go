package keeper_test

import (
	"arkeo/x/crosstransfer/types"
	"testing"

	testkeeper "arkeo/testutil/keeper"

	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.CrosstransferKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
