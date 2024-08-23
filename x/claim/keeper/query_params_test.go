package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	testkeeper "github.com/arkeonetwork/arkeo/testutil/keeper"
	"github.com/arkeonetwork/arkeo/x/claim/types"
)

func TestParamsQuery(t *testing.T) {
	keepers, ctx := testkeeper.CreateTestClaimKeepers(t)
	wctx := sdk.WrapSDKContext(ctx) //nolint:staticcheck
	params := types.DefaultParams()
	keepers.ClaimKeeper.SetParams(ctx, params)

	response, err := keepers.ClaimKeeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
