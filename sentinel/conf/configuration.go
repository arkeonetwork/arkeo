package conf

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/arkeonetwork/arkeo/common"
)

type TLSConfiguration struct {
	Cert string `json:"tls_certificate"`
	Key  string `json:"tls_key"`
}

type Configuration struct {
	Moniker                     string           `json:"moniker"`
	Website                     string           `json:"website"`
	Description                 string           `json:"description"`
	Location                    string           `json:"location"`
	Port                        string           `json:"port"`
	SourceChain                 string           `json:"source_chain"` // base url for arkeo block chain
	EventStreamHost             string           `json:"event_stream_host"`
	ClaimStoreLocation          string           `json:"claim_store_location"`           // file location where claims are stored
	ContractConfigStoreLocation string           `json:"contract_config_store_location"` // file location where contract configurations are stored
	ProviderConfigStoreLocation string           `json:"provider_config_store_location"` // file location where provider configurations are stored
	ProviderPubKey              common.PubKey    `json:"provider_pubkey"`
	FreeTierRateLimit           int              `json:"free_tier_rate_limit"`
	TLS                         TLSConfiguration `json:"tls"`
	ArkeoAuthContractId         uint64           `json:"arkeo_auth_contract_id"`  // Contract ID for auth
	ArkeoAuthChainId            string           `json:"arkeo_auth_chain_id"`     // Chain ID for auth
	ArkeoAuthMnemonic           string           `json:"arkeo_auth_mnemonic"`     // Mnemonic phrase for signing
	ArkeoAuthNonceStore         string           `json:"arkeo_auth_nonce_store"`  // LevelDB path for nonce storage
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
	}
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
