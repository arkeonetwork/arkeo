package keeper

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/stretchr/testify/require"
)

func TestModProviderValidate(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)

	s := newMsgServer(k, sk)

	// setup
	pubkey := types.GetRandomPubKey()

	provider := types.NewProvider(pubkey, common.BTCService)
	provider.Bond = cosmos.NewInt(500)
	require.NoError(t, k.SetProvider(ctx, provider))

	// happy path
	msg := types.MsgModProvider{
		Provider:            provider.PubKey,
		Service:             provider.Service.String(),
		MinContractDuration: 10,
		MaxContractDuration: 500,
		Status:              types.ProviderStatus_ONLINE,
	}
	require.NoError(t, s.ModProviderValidate(ctx, &msg))

	// bad min duration
	msg.MinContractDuration = 5256000 * 2
	err := s.ModProviderValidate(ctx, &msg)
	require.ErrorIs(t, err, types.ErrInvalidModProviderMinContractDuration)

	// bad max duration
	msg.MinContractDuration = 10
	msg.MaxContractDuration = 5256000 * 2
	err = s.ModProviderValidate(ctx, &msg)
	require.ErrorIs(t, err, types.ErrInvalidModProviderMaxContractDuration)
}

func TestModProviderHandle(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)

	s := newMsgServer(k, sk)

	// setup
	pubkey := types.GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	require.NoError(t, err)

	sRates, err := cosmos.ParseCoins("11uarkeo")
	require.NoError(t, err)
	pRates, err := cosmos.ParseCoins("12uarkeo")
	require.NoError(t, err)

	require.NoError(t, err)
	rates := []*types.ContractRate{
		{
			MeterType: types.MeterType_PAY_PER_BLOCK,
			UserType:  types.UserType_SINGLE_USER,
			Rates:     sRates,
		},
		{
			MeterType: types.MeterType_PAY_PER_CALL,
			UserType:  types.UserType_SINGLE_USER,
			Rates:     pRates,
		},
	}

	// happy path
	msg := types.MsgModProvider{
		Creator:             acct,
		Provider:            pubkey,
		Service:             common.BTCService.String(),
		MetadataUri:         "foobar",
		MetadataNonce:       3,
		MinContractDuration: 10,
		MaxContractDuration: 500,
		Status:              types.ProviderStatus_ONLINE,
		Rates:               rates,
	}

	require.NoError(t, s.ModProviderHandle(ctx, &msg))

	provider, err := k.GetProvider(ctx, msg.Provider, common.BTCService)
	require.NoError(t, err)
	require.Equal(t, provider.MetadataUri, "foobar")
	require.Equal(t, provider.MetadataNonce, uint64(3))
	require.Equal(t, provider.MinContractDuration, int64(10))
	require.Equal(t, provider.MaxContractDuration, int64(500))
	require.Equal(t, provider.Status, types.ProviderStatus_ONLINE)
	require.ElementsMatch(t, provider.Rates, rates)
}
