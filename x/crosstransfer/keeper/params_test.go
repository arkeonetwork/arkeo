package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "arkeo/testutil/keeper"
	"arkeo/x/crosstransfer/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.CrosstransferKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
