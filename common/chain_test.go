package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	service, err := NewService("btc-mainnet-fullnode")
	require.NoError(t, err)
	require.True(t, service.Equals(BTCService))
	require.False(t, service.IsEmpty())
	require.Equal(t, service.String(), "btc-mainnet-fullnode")

	service, err = NewService("swapi.dev")
	require.NoError(t, err)
	require.True(t, service.Equals(StarWarsService))
	require.False(t, service.IsEmpty())
	require.Equal(t, service.String(), "swapi.dev")

	_, err = NewService("B") // invalid
	require.Error(t, err)
}
