package sentinel

// TODO copied from sentinel
import (
	"time"
)

type Configuration struct {
	Nonce                     int64         `json:"-"`
	Moniker                   string        `json:"moniker"`
	Website                   string        `json:"website"`
	Description               string        `json:"description"`
	Location                  string        `json:"location"`
	Port                      string        `json:"port"`
	SourceChain               string        `json:"source_chain"` // base url for arceo block chain
	EventStreamHost           string        `json:"event_stream_host"`
	ClaimStoreLocation        string        `json:"claim_store_location"` // file location where claims are stored
	ProviderPubKey            string        `json:"provider_pubkey"`
	FreeTierRateLimit         int           `json:"free_tier_rate_limit"`
	FreeTierRateLimitDuration time.Duration `json:"free_tier_rate_limit_duration"`
	SubTierRateLimit          int           `json:"subscription_tier_rate_limit"`
	SubTierRateLimitDuration  time.Duration `json:"subscription_tier_rate_limit_duration"`
	AsGoTierRateLimit         int           `json:"pay_as_you_go_tier_rate_limit"`
	AsGoTierRateLimitDuration time.Duration `json:"pay_as_you_go_tier_rate_limit_duration"`
}

type Metadata struct {
	Configuration Configuration `json:"config"`
	Version       string        `json:"version" db:"version"`
}
