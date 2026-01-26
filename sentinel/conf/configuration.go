package conf

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/arkeonetwork/arkeo/common"
	"gopkg.in/yaml.v2"
)

type TLSConfiguration struct {
	Cert string `json:"tls_certificate,omitempty"`
	Key  string `json:"tls_key,omitempty"`
}

type ServiceConfig struct {
	Name    string `json:"name" yaml:"name"`
	Id      int    `json:"id" yaml:"id"`
	Type    string `json:"type" yaml:"type"`
	RpcUrl  string `json:"rpc_url" yaml:"rpc_url,omitempty"`
	RpcUser string `json:"rpc_user,omitempty" yaml:"rpc_user,omitempty"`
	RpcPass string `json:"rpc_pass,omitempty" yaml:"rpc_pass,omitempty"`
}

type Configuration struct {
	Moniker                     string           `json:"moniker,omitempty"`
	Website                     string           `json:"website,omitempty"`
	Description                 string           `json:"description,omitempty"`
	Location                    string           `json:"location,omitempty"`
	Port                        string           `json:"port,omitempty"`
	SourceChain                 string           `json:"source_chain,omitempty"` // base url for arkeo block chain
	HubProviderURI              string           `json:"hub_provider_uri,omitempty"`
	EventStreamHost             string           `json:"event_stream_host,omitempty"`
	ClaimStoreLocation          string           `json:"claim_store_location,omitempty"`           // file location where claims are stored
	ContractConfigStoreLocation string           `json:"contract_config_store_location,omitempty"` // file location where contract configurations are stored
	ProviderConfigStoreLocation string           `json:"provider_config_store_location,omitempty"` // file location where provider configurations are stored
	ProviderPubKey              common.PubKey    `json:"provider_pubkey,omitempty"`
	FreeTierRateLimit           int              `json:"free_tier_rate_limit,omitempty"`
	TLS                         TLSConfiguration `json:"tls"`
	Services                    []ServiceConfig  `json:"services" yaml:"services"`
	ArkeoAuthContractId         uint64           `json:"arkeo_auth_contract_id,omitempty"` // Contract ID for auth
	ArkeoAuthChainId            string           `json:"arkeo_auth_chain_id,omitempty"`    // Chain ID for auth
	ArkeoAuthMnemonic           string           `json:"arkeo_auth_mnemonic,omitempty"`    // Mnemonic phrase for signing
	ArkeoAuthNonceStore         string           `json:"arkeo_auth_nonce_store,omitempty"` // LevelDB path for nonce storage
	TrustedProxyIPs             []string         `json:"trusted_proxy_ips,omitempty"`       // List of trusted proxy IPs/CIDRs for X-Forwarded-For header validation
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

func loadVarIntOptional(key string, defaultValue int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Errorf("env var %s is not an integer: %s", key, err))
	}
	return i
}

func NewTLSConfiguration() TLSConfiguration {
	return TLSConfiguration{
		Cert: getEnv("TLS_CERT", ""),
		Key:  getEnv("TLS_KEY", ""),
	}
}

func (c TLSConfiguration) HasTLS() bool {
	return len(c.Cert) > 0 && len(c.Key) > 0
}

func NewConfiguration() Configuration {
	return Configuration{
		Moniker:                     loadVarString("MONIKER"),
		Website:                     loadVarString("WEBSITE"),
		Description:                 loadVarString("DESCRIPTION"),
		Location:                    loadVarString("LOCATION"),
		Port:                        getEnv("PORT", "3636"),
		SourceChain:                 loadVarString("SOURCE_CHAIN"),
		HubProviderURI:              loadVarString("PROVIDER_HUB_URI"),
		EventStreamHost:             loadVarString("EVENT_STREAM_HOST"),
		ProviderPubKey:              loadVarPubKey("PROVIDER_PUBKEY"),
		FreeTierRateLimit:           loadVarInt("FREE_RATE_LIMIT"),
		ClaimStoreLocation:          loadVarString("CLAIM_STORE_LOCATION"),
		ContractConfigStoreLocation: loadVarString("CONTRACT_CONFIG_STORE_LOCATION"),
		TLS:                         NewTLSConfiguration(),
		ProviderConfigStoreLocation: loadVarString("PROVIDER_CONFIG_STORE_LOCATION"),
		ArkeoAuthContractId:         uint64(loadVarIntOptional("ARKEO_AUTH_CONTRACT_ID", 0)),
		ArkeoAuthChainId:            getEnv("ARKEO_AUTH_CHAIN_ID", ""),
		ArkeoAuthMnemonic:           getEnv("ARKEO_AUTH_MNEMONIC", ""),
		ArkeoAuthNonceStore:         getEnv("ARKEO_AUTH_NONCE_STORE", ""),
		TrustedProxyIPs:             loadTrustedProxyIPs(),
	}
}

// loadTrustedProxyIPs loads trusted proxy IPs from environment variable.
// Format: comma-separated list of IPs or CIDRs (e.g., "10.0.0.1,192.168.0.0/16")
// If not set, returns empty slice (defaults to localhost-only).
func loadTrustedProxyIPs() []string {
	envVal := getEnv("TRUSTED_PROXY_IPS", "")
	if envVal == "" {
		return []string{} // Empty = default to localhost-only
	}
	// Split by comma and trim whitespace
	ips := strings.Split(envVal, ",")
	result := make([]string, 0, len(ips))
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		if ip != "" {
			result = append(result, ip)
		}
	}
	return result
}

func (c Configuration) Print() {
	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintln(writer, "Moniker\t", c.Moniker)
	fmt.Fprintln(writer, "Website\t", c.Website)
	fmt.Fprintln(writer, "Description\t", c.Description)
	fmt.Fprintln(writer, "Location\t", c.Location)
	fmt.Fprintln(writer, "Port\t", c.Port)
	fmt.Fprintln(writer, "TLS Certificate\t", c.TLS.Cert)
	fmt.Fprintln(writer, "TLS Key\t", c.TLS.Key)
	fmt.Fprintln(writer, "Source Chain\t", c.SourceChain)
	fmt.Fprintln(writer, "Event Stream Host\t", c.EventStreamHost)
	fmt.Fprintln(writer, "Provider PubKey\t", c.ProviderPubKey)
	fmt.Fprintln(writer, "Claim Store Location\t", c.ClaimStoreLocation)
	fmt.Fprintln(writer, "Contract Config Store Location\t", c.ContractConfigStoreLocation)
	fmt.Fprintln(writer, "Free Tier Rate Limit\t", fmt.Sprintf("%d requests per 1m", c.FreeTierRateLimit))
	fmt.Fprintln(writer, "Provider Config Store Location\t", c.ProviderConfigStoreLocation)

	if c.ArkeoAuthContractId > 0 {
		fmt.Fprintln(writer, "Arkeo Auth Contract ID\t", c.ArkeoAuthContractId)
		fmt.Fprintln(writer, "Arkeo Auth Chain ID\t", c.ArkeoAuthChainId)
		fmt.Fprintln(writer, "Arkeo Auth Configured\t", c.ArkeoAuthMnemonic != "")
		fmt.Fprintln(writer, "Arkeo Auth Nonce Store\t", c.ArkeoAuthNonceStore)
	}

	writer.Flush()
}

// LoadConfigurationFromFile loads the configuration from a YAML file and applies environment variable overrides.
func LoadConfigurationFromFile(filename string) (Configuration, error) {
	var cfg Configuration
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	// Optional: Environment variable overrides for legacy/env-style fields
	// (the below block is optional and safe; comment it out if you don't want override behavior)
	overrideString := func(env, val string) string {
		if v := os.Getenv(env); v != "" {
			return v
		}
		return val
	}
	overrideInt := func(env string, val int) int {
		if v := os.Getenv(env); v != "" {
			i, err := strconv.Atoi(v)
			if err == nil {
				return i
			}
		}
		return val
	}
	overrideUint64 := func(env string, val uint64) uint64 {
		if v := os.Getenv(env); v != "" {
			i, err := strconv.ParseUint(v, 10, 64)
			if err == nil {
				return i
			}
		}
		return val
	}

	cfg.Moniker = overrideString("MONIKER", cfg.Moniker)
	cfg.Website = overrideString("WEBSITE", cfg.Website)
	cfg.Description = overrideString("DESCRIPTION", cfg.Description)
	cfg.Location = overrideString("LOCATION", cfg.Location)
	cfg.Port = overrideString("PORT", cfg.Port)
	cfg.SourceChain = overrideString("SOURCE_CHAIN", cfg.SourceChain)
	cfg.HubProviderURI = overrideString("PROVIDER_HUB_URI", cfg.HubProviderURI)
	cfg.EventStreamHost = overrideString("EVENT_STREAM_HOST", cfg.EventStreamHost)
	cfg.ClaimStoreLocation = overrideString("CLAIM_STORE_LOCATION", cfg.ClaimStoreLocation)
	cfg.ContractConfigStoreLocation = overrideString("CONTRACT_CONFIG_STORE_LOCATION", cfg.ContractConfigStoreLocation)
	cfg.ProviderConfigStoreLocation = overrideString("PROVIDER_CONFIG_STORE_LOCATION", cfg.ProviderConfigStoreLocation)
	cfg.FreeTierRateLimit = overrideInt("FREE_RATE_LIMIT", cfg.FreeTierRateLimit)
	// ProviderPubKey override (optional, if you want):
	if v := os.Getenv("PROVIDER_PUBKEY"); v != "" {
		pk, err := common.NewPubKey(v)
		if err == nil {
			cfg.ProviderPubKey = pk
		}
	}
	// TLS overrides
	if v := os.Getenv("TLS_CERT"); v != "" {
		cfg.TLS.Cert = v
	}
	if v := os.Getenv("TLS_KEY"); v != "" {
		cfg.TLS.Key = v
	}
	cfg.ArkeoAuthContractId = overrideUint64("ArkeoAuthContractId", cfg.ArkeoAuthContractId)
	cfg.ArkeoAuthChainId = overrideString("ArkeoAuthChainId", cfg.ArkeoAuthChainId)
	cfg.ArkeoAuthMnemonic = overrideString("ArkeoAuthMnemonic", cfg.ArkeoAuthMnemonic)
	cfg.ArkeoAuthNonceStore = overrideString("ArkeoAuthNonceStore", cfg.ArkeoAuthNonceStore)
	
	// TrustedProxyIPs: if set in env, override; otherwise use YAML value
	if envVal := os.Getenv("TRUSTED_PROXY_IPS"); envVal != "" {
		cfg.TrustedProxyIPs = loadTrustedProxyIPs()
	} else if len(cfg.TrustedProxyIPs) == 0 {
		// If not in YAML either, default to empty (localhost-only)
		cfg.TrustedProxyIPs = []string{}
	}

	return cfg, nil
}
