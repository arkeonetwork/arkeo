package conf

import (
	"fmt"
	"os"
	"strings"
)

type Configuration struct {
	Port        string
	ProxyHost   string
	SourceChain string // base url for arceo block chain
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
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

func NewConfiguration() Configuration {
	return Configuration{
		Port:        getEnv("PORT", "3636"),
		ProxyHost:   loadVarString("PROXY_HOST"),
		SourceChain: loadVarString("SOURCE_CHAIN"),
	}
}
