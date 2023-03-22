package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type DBConfig struct {
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         uint   `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPass         string `mapstructure:"DB_PASS"`
	DBName         string `mapstructure:"DB_NAME"`
	DBSSLMode      string `mapstructure:"DB_SSL_MODE"`
	DBPoolMaxConns int    `mapstructure:"DB_POOL_MAX_CONNS"`
	DBPoolMinConns int    `mapstructure:"DB_POOL_MIN_CONNS"`
}

var (
	dbConfigNames = []string{
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASS",
		"DB_NAME",
		"DB_SSL_MODE",
		"DB_POOL_MAX_CONNS",
		"DB_POOL_MIN_CONNS",
	}
)

func ReadDBConfig(envPath string) *DBConfig {
	c := &DBConfig{}
	if envPath == "" {
		if err := LoadFromEnv(c, dbConfigNames...); err != nil {
			log.Panicf("failed to load db config from env: %+v", err)
		}
	} else {
		if err := Load(envPath, c); err != nil {
			log.Panicf("failed to load db config: %+v", err)
		}
	}
	return c
}

// Load reads in a file at the specified path and unmarshals the values into your config struct.
func Load(path string, config interface{}) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrapf(err, "failed to read config from file path: %s", path)
	}

	if err := viper.Unmarshal(config); err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}

	return nil
}

// Load from env names given in keys
func LoadFromEnv(config interface{}, keys ...string) error {
	envVars := make(map[string]interface{})
	for _, key := range keys {
		val, ok := os.LookupEnv(key)
		if !ok {
			return errors.Errorf("%s environment variable not set", key)
		}

		envVars[key] = val
	}

	jsonStr, err := json.Marshal(envVars)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal envVars: %v", envVars)
	}

	viper.SetConfigType("json")

	if err := viper.ReadConfig(bytes.NewBuffer(jsonStr)); err != nil {
		return errors.Wrapf(err, "failed to read json: %s", jsonStr)
	}

	if err := viper.Unmarshal(config); err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}

	return nil
}
