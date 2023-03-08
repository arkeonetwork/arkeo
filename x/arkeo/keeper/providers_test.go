package keeper

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {
	ctx, k := SetupKeeper(t)

	require.Error(t, k.SetProvider(ctx, types.Provider{})) // empty asset should error

	provider := types.NewProvider(types.GetRandomPubKey(), common.BTCChain)

	err := k.SetProvider(ctx, provider)
	require.NoError(t, err)
	provider, err = k.GetProvider(ctx, provider.PubKey, provider.Chain)
	require.NoError(t, err)
	require.True(t, provider.Chain.Equals(common.BTCChain))
	require.True(t, k.ProviderExists(ctx, provider.PubKey, provider.Chain))
	require.False(t, k.ProviderExists(ctx, provider.PubKey, common.ETHChain))

	k.RemoveProvider(ctx, provider.PubKey, provider.Chain)
	require.False(t, k.ProviderExists(ctx, provider.PubKey, provider.Chain))
}
