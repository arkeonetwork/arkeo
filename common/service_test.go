package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	var service Service = 0
	require.Equal(t, "unknown", service.String())
	service = MockService
	require.Equal(t, "mock", service.String())
	service = BTCService
	require.Equal(t, "btc-mainnet-fullnode", service.String())
	service = ETHService
	require.Equal(t, "eth-mainnet-fullnode", service.String())
	service = 600
	require.Equal(t, "unknown", service.String())
}
