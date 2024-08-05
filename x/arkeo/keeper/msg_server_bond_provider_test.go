package keeper

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)

	s := newMsgServer(k, sk)

	// setup
	providerPubKey := types.GetRandomPubKey()
	acct, err := providerPubKey.GetMyAddress()
	require.NoError(t, err)
	require.NoError(t, k.MintAndSendToAccount(ctx, acct, getCoin(common.Tokens(10))))

	// Add to bond
	msg := types.MsgBondProvider{
		Creator:  acct.String(),
		Provider: providerPubKey,
		Service:  common.BTCService.String(),
		Bond:     cosmos.NewInt(common.Tokens(8)),
	}
	require.NoError(t, s.BondProviderHandle(ctx, &msg))
	// check balance as drawn down by two
	bal := k.GetBalance(ctx, acct)
	require.Equal(t, bal.AmountOf(configs.Denom).Int64(), common.Tokens(2))
	// check that provider now exists
	require.True(t, k.ProviderExists(ctx, msg.Provider, common.BTCService))
	provider, err := k.GetProvider(ctx, msg.Provider, common.BTCService)
	require.NoError(t, err)
	require.Equal(t, provider.Bond.Int64(), common.Tokens(8))

	// remove too much bond
	msg.Bond = cosmos.NewInt(common.Tokens(-20))
	err = s.BondProviderHandle(ctx, &msg)
	require.ErrorIs(t, err, types.ErrInsufficientFunds)

	// check balance hasn't changed
	bal = k.GetBalance(ctx, acct)
	require.Equal(t, bal.AmountOf(configs.Denom).Int64(), common.Tokens(2))

	// check provider has same bond
	provider, err = k.GetProvider(ctx, msg.Provider, common.BTCService)
	require.NoError(t, err)
	require.Equal(t, provider.Bond.Int64(), common.Tokens(8))

	// remove all bond
	msg.Bond = cosmos.NewInt(common.Tokens(-8))
	err = s.BondProviderHandle(ctx, &msg)
	require.NoError(t, err)

	bal = k.GetBalance(ctx, acct) // check balance
	require.Equal(t, bal.AmountOf(configs.Denom).Int64(), common.Tokens(10))
	require.False(t, k.ProviderExists(ctx, msg.Provider, common.BTCService)) // should be removed
}
