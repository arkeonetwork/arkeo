package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfiguration(t *testing.T) {
	os.Setenv("MONIKER", "monkey")
	os.Setenv("WEBSITE", "webby")
	os.Setenv("DESCRIPTION", "dezy")
	os.Setenv("LOCATION", "locy")
	os.Setenv("PORT", "4000")
	os.Setenv("SOURCE_CHAIN", "sourcey")
	os.Setenv("EVENT_STREAM_HOST", "hosty")
	os.Setenv("PROVIDER_PUBKEY", "cosmospub1addwnpepqg3523h7e7ggeh6na2lsde6s394tqxnvufsz0urld6zwl8687ue9c3dasgu")
	os.Setenv("FREE_RATE_LIMIT", "99")
	os.Setenv("CLAIM_STORE_LOCATION", "clammy")
	os.Setenv("GAIA_RPC_ARCHIVE_HOST", "gaia-host")

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
}
