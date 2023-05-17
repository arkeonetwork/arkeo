package keeper

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/stretchr/testify/require"
)

func TestValidateSetVersion(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	s := newMsgServer(k, sk)

	cosmos.GetConfig().SetBech32PrefixForAccount("arkeo", "arkeopub")
	cosmos.GetConfig().SetBech32PrefixForValidator("varkeo", "varkeopub")

	// setup
	providerPubKey := types.GetRandomPubKey()
	acct, err := providerPubKey.GetMyAddress()
	require.NoError(t, err)

	msg := types.NewMsgSetVersion(acct, 15)
	require.NoError(t, msg.ValidateBasic())
	require.NoError(t, s.SetVersionValidate(ctx, msg))

	valAddr := cosmos.ValAddress(acct)
	k.SetVersionForAddress(ctx, valAddr, 100)
	require.Error(t, s.SetVersionValidate(ctx, msg))
}

func TestHandleSetVersion(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	s := newMsgServer(k, sk)

	cosmos.GetConfig().SetBech32PrefixForAccount("arkeo", "arkeopub")
	cosmos.GetConfig().SetBech32PrefixForValidator("varkeo", "varkeopub")

	// setup
	providerPubKey := types.GetRandomPubKey()
	acct, err := providerPubKey.GetMyAddress()
	require.NoError(t, err)

	msg := types.NewMsgSetVersion(acct, 15)
	require.NoError(t, msg.ValidateBasic())
	require.NoError(t, s.SetVersionHandle(ctx, msg))

	valAddr := cosmos.ValAddress(acct)
	currentVersion := k.GetVersionForAddress(ctx, valAddr)
	require.Equal(t, int64(15), currentVersion)
}
