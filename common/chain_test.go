package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChain(t *testing.T) {
	chn, err := NewChain("btc-mainnet-fullnode")
	require.NoError(t, err)
	require.True(t, chn.Equals(BTCChain))
	require.False(t, chn.IsEmpty())
	require.Equal(t, chn.String(), "btc-mainnet-fullnode")

	chn, err = NewChain("swapi.dev")
	require.NoError(t, err)
	require.True(t, chn.Equals(StarWarsChain))
	require.False(t, chn.IsEmpty())
	require.Equal(t, chn.String(), "swapi.dev")

	_, err = NewChain("B") // invalid
	require.Error(t, err)
}
