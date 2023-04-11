package keeper

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {
	ctx, k := SetupKeeper(t)

	require.Error(t, k.SetProvider(ctx, types.Provider{})) // empty asset should error

	provider := types.NewProvider(types.GetRandomPubKey(), common.BTCService)
	provider.Bond = cosmos.NewInt(100)

	err := k.SetProvider(ctx, provider)
	require.NoError(t, err)
	provider, err = k.GetProvider(ctx, provider.PubKey, provider.Service)
	require.NoError(t, err)
	require.True(t, provider.Service.Equals(common.BTCService))
	require.True(t, k.ProviderExists(ctx, provider.PubKey, provider.Service))
	require.False(t, k.ProviderExists(ctx, provider.PubKey, common.ETHService))

	k.RemoveProvider(ctx, provider.PubKey, provider.Service)
	require.False(t, k.ProviderExists(ctx, provider.PubKey, provider.Service))
}
