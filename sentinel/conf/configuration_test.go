package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupRequiredEnvVars(t *testing.T) {
	t.Helper()
	os.Setenv("MONIKER", "monkey")
	os.Setenv("WEBSITE", "webby")
	os.Setenv("DESCRIPTION", "dezy")
	os.Setenv("LOCATION", "locy")
	os.Setenv("PORT", "4000")
	os.Setenv("SOURCE_CHAIN", "sourcey")
	os.Setenv("PROVIDER_HUB_URI", "hubby")
	os.Setenv("EVENT_STREAM_HOST", "hosty")
	os.Setenv("PROVIDER_PUBKEY", "cosmospub1addwnpepqg3523h7e7ggeh6na2lsde6s394tqxnvufsz0urld6zwl8687ue9c3dasgu")
	os.Setenv("FREE_RATE_LIMIT", "99")
	os.Setenv("CLAIM_STORE_LOCATION", "clammy")
	os.Setenv("CONTRACT_CONFIG_STORE_LOCATION", "configy")
	os.Setenv("PROVIDER_CONFIG_STORE_LOCATION", "providy")
}

func TestConfiguration(t *testing.T) {
	setupRequiredEnvVars(t)
	os.Setenv("TRUSTED_PROXY_IPS", "10.0.0.1,192.168.0.0/16,172.16.0.1")

	config := NewConfiguration()

	require.Equal(t, config.Moniker, "monkey")
	require.Equal(t, config.Website, "webby")
	require.Equal(t, config.Location, "locy")
	require.Equal(t, config.Port, "4000")
	require.Equal(t, config.SourceChain, "sourcey")
	require.Equal(t, config.EventStreamHost, "hosty")
	require.Equal(t, config.ProviderPubKey.String(), "cosmospub1addwnpepqg3523h7e7ggeh6na2lsde6s394tqxnvufsz0urld6zwl8687ue9c3dasgu")
	require.Equal(t, config.FreeTierRateLimit, 99)
	require.Equal(t, config.ClaimStoreLocation, "clammy")
	require.Equal(t, config.ContractConfigStoreLocation, "configy")
	require.Equal(t, config.ProviderConfigStoreLocation, "providy")
	require.Equal(t, config.TrustedProxyIPs, []string{"10.0.0.1", "192.168.0.0/16", "172.16.0.1"})
}
