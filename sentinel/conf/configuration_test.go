package conf

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConfiguration(t *testing.T) {
	os.Setenv("MONIKER", "monkey")
	os.Setenv("WEBSITE", "webby")
	os.Setenv("DESCRIPTION", "dezy")
	os.Setenv("LOCATION", "locy")
	os.Setenv("PORT", "4000")
	os.Setenv("PROXY_HOST", "proxxy")
	os.Setenv("SOURCE_CHAIN", "sourcey")
	os.Setenv("EVENT_STREAM_HOST", "hosty")
	os.Setenv("PROVIDER_PUBKEY", "cosmospub1addwnpepqg3523h7e7ggeh6na2lsde6s394tqxnvufsz0urld6zwl8687ue9c3dasgu")
	os.Setenv("FREE_RATE_LIMIT", "99")
	os.Setenv("FREE_RATE_LIMIT_DURATION", "1h")
	os.Setenv("SUB_RATE_LIMIT", "98")
	os.Setenv("SUB_RATE_LIMIT_DURATION", "2m")
	os.Setenv("AS_GO_RATE_LIMIT", "97")
	os.Setenv("AS_GO_RATE_LIMIT_DURATION", "3h")
	os.Setenv("CLAIM_STORE_LOCATION", "clammy")
	os.Setenv("GAIA_RPC_ARCHIVE_HOST", "gaia-host")

	config := NewConfiguration()

	require.Equal(t, config.Moniker, "monkey")
	require.Equal(t, config.Website, "webby")
	require.Equal(t, config.Location, "locy")
	require.Equal(t, config.Port, "4000")
	require.Equal(t, config.ProxyHost, "proxxy")
	require.Equal(t, config.SourceChain, "sourcey")
	require.Equal(t, config.EventStreamHost, "hosty")
	require.Equal(t, config.ProviderPubKey.String(), "cosmospub1addwnpepqg3523h7e7ggeh6na2lsde6s394tqxnvufsz0urld6zwl8687ue9c3dasgu")
	require.Equal(t, config.FreeTierRateLimit, 99)
	require.Equal(t, config.FreeTierRateLimitDuration, time.Hour)
	require.Equal(t, config.SubTierRateLimit, 98)
	require.Equal(t, config.SubTierRateLimitDuration, 2*time.Minute)
	require.Equal(t, config.AsGoTierRateLimit, 97)
	require.Equal(t, config.AsGoTierRateLimitDuration, 3*time.Hour)
	require.Equal(t, config.ClaimStoreLocation, "clammy")
}
