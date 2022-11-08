package conf

import (
	"fmt"
	"mercury/common"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

type Configuration struct {
	Moniker                   string        `json:"moniker"`
	Website                   string        `json:"website"`
	Description               string        `json:"description"`
	Location                  string        `json:"location"`
	Port                      string        `json:"port"`
	ProxyHost                 string        `json:"proxy_host"`
	SourceChain               string        `json:"source_chain"`         // base url for arceo block chain
	ClaimStoreLocation        string        `json:"claim_store_location"` // file location where claims are stored
	ProviderPubKey            common.PubKey `json:"provider_pubkey"`
	FreeTierRateLimit         int           `json:"free_tier_rate_limit"`
	FreeTierRateLimitDuration time.Duration `json:"free_tier_rate_limit_duration"`
	SubTierRateLimit          int           `json:"subscription_tier_rate_limit"`
	SubTierRateLimitDuration  time.Duration `json:"subscription_tier_rate_limit_duration"`
	AsGoTierRateLimit         int           `json:"pay_as_you_go_tier_rate_limit"`
	AsGoTierRateLimitDuration time.Duration `json:"pay_as_you_go_tier_rate_limit_duration"`
}

// Simple helper function to read an environment or return a default value
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func loadVarString(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("%s env var is not set", key))
	}
	return strings.TrimSpace(val)
}

func loadVarPubKey(key string) common.PubKey {
	val, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("%s env var is not set", key))
	}
	pk, err := common.NewPubKey(val)
	if err != nil {
		panic(fmt.Errorf("env var %s is not a pubkey: %s", key, err))
	}
	return pk
}

func loadVarInt(key string) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("%s env var is not set", key))
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Errorf("env var %s is not an integer: %s", key, err))
	}
	return i
}

func loadVarDuration(key string) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("%s env var is not set", key))
	}
	dur, err := time.ParseDuration(val)
	if err != nil {
		panic(fmt.Errorf("env var %s is not a duration: %s", key, err))
	}
	return dur
}

func NewConfiguration() Configuration {
	return Configuration{
		Moniker:                   loadVarString("MONIKER"),
		Website:                   loadVarString("WEBSITE"),
		Description:               loadVarString("DESCRIPTION"),
		Location:                  loadVarString("LOCATION"),
		Port:                      loadVarString("PORT"),
		ProxyHost:                 loadVarString("PROXY_HOST"),
		SourceChain:               loadVarString("SOURCE_CHAIN"),
		ProviderPubKey:            loadVarPubKey("PROVIDER_PUBKEY"),
		FreeTierRateLimit:         loadVarInt("FREE_RATE_LIMIT"),
		FreeTierRateLimitDuration: loadVarDuration("FREE_RATE_LIMIT_DURATION"),
		SubTierRateLimit:          loadVarInt("SUB_RATE_LIMIT"),
		SubTierRateLimitDuration:  loadVarDuration("SUB_RATE_LIMIT_DURATION"),
		AsGoTierRateLimit:         loadVarInt("AS_GO_RATE_LIMIT"),
		AsGoTierRateLimitDuration: loadVarDuration("AS_GO_RATE_LIMIT_DURATION"),
		ClaimStoreLocation:        loadVarString("CLAIM_STORE_LOCATION"),
	}
}

func (c Configuration) Print() {
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintln(writer, fmt.Sprintf("%s\t%+v", "Moniker", c.Moniker))
	fmt.Fprintln(writer, fmt.Sprintf("%s\t%+v", "Website", c.Website))
	fmt.Fprintln(writer, fmt.Sprintf("%s\t%+v", "Description", c.Description))
	fmt.Fprintln(writer, fmt.Sprintf("%s\t%+v", "Location", c.Location))
	fmt.Fprintln(writer, fmt.Sprintf("%s\t%+v", "Port", c.Port))
	fmt.Fprintln(writer, fmt.Sprintf("%s\t%+v", "ProxyHost", c.ProxyHost))
	fmt.Fprintln(writer, fmt.Sprintf("%s\t%+v", "SourceChain", c.SourceChain))
	fmt.Fprintln(writer, fmt.Sprintf("%s\t%+v", "Provider PubKey", c.ProviderPubKey))
	fmt.Fprintln(writer, fmt.Sprintf("%s\t%+v", "Claim Store Location", c.ClaimStoreLocation))
	fmt.Fprintln(writer, fmt.Sprintf("%s\t%d requests per %+v", "Free Tier Rate Limit", c.FreeTierRateLimit, c.FreeTierRateLimitDuration))
	fmt.Fprintln(writer, fmt.Sprintf("%s\t%d requests per %+v", "Subscription Rate Limit", c.SubTierRateLimit, c.SubTierRateLimitDuration))
	fmt.Fprintln(writer, fmt.Sprintf("%s\t%d requests per %+v", "Pay-As-You-Go Rate Limit", c.AsGoTierRateLimit, c.AsGoTierRateLimitDuration))
	writer.Flush()
}
