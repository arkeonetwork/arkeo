package conf

import (
	"os"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type ConfigurationSuite struct{}

var _ = Suite(&ConfigurationSuite{})

func (ConfigurationSuite) TestConfiguration(c *C) {
	os.Setenv("MONIKER", "monkey")
	os.Setenv("WEBSITE", "webby")
	os.Setenv("DESCRIPTION", "dezy")
	os.Setenv("LOCATION", "locy")
	os.Setenv("PORT", "4000")
	os.Setenv("PROXY_HOST", "proxxy")
	os.Setenv("SOURCE_CHAIN", "sourcey")
	os.Setenv("PROVIDER_PUBKEY", "cosmospub1addwnpepqg3523h7e7ggeh6na2lsde6s394tqxnvufsz0urld6zwl8687ue9c3dasgu")
	os.Setenv("FREE_RATE_LIMIT", "99")
	os.Setenv("FREE_RATE_LIMIT_DURATION", "1h")
	os.Setenv("SUB_RATE_LIMIT", "98")
	os.Setenv("SUB_RATE_LIMIT_DURATION", "2m")
	os.Setenv("AS_GO_RATE_LIMIT", "97")
	os.Setenv("AS_GO_RATE_LIMIT_DURATION", "3h")
	os.Setenv("CLAIM_STORE_LOCATION", "clammy")

	config := NewConfiguration()

	c.Check(config.Moniker, Equals, "monkey")
	c.Check(config.Website, Equals, "webby")
	c.Check(config.Location, Equals, "locy")
	c.Check(config.Port, Equals, "4000")
	c.Check(config.ProxyHost, Equals, "proxxy")
	c.Check(config.SourceChain, Equals, "sourcey")
	c.Check(config.ProviderPubKey.String(), Equals, "cosmospub1addwnpepqg3523h7e7ggeh6na2lsde6s394tqxnvufsz0urld6zwl8687ue9c3dasgu")
	c.Check(config.FreeTierRateLimit, Equals, 99)
	c.Check(config.FreeTierRateLimitDuration, Equals, time.Hour)
	c.Check(config.SubTierRateLimit, Equals, 98)
	c.Check(config.SubTierRateLimitDuration, Equals, 2*time.Minute)
	c.Check(config.AsGoTierRateLimit, Equals, 97)
	c.Check(config.AsGoTierRateLimitDuration, Equals, 3*time.Hour)
	c.Check(config.ClaimStoreLocation, Equals, "clammy")
}
