package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strings"
)

type DBConfig struct {
	DBHost         string `mapstructure:"db_host"`
	DBPort         uint   `mapstructure:"db_port"`
	DBUser         string `mapstructure:"db_user"`
	DBPass         string `mapstructure:"db_pass"`
	DBName         string `mapstructure:"db_name"`
	DBSSLMode      string `mapstructure:"db_ssl_mode"`
	DBPoolMaxConns int    `mapstructure:"db_pool_max_conns"`
	DBPoolMinConns int    `mapstructure:"db_pool_min_conns"`
}

// LoadFromEnv read config from environment variables
func LoadFromEnv(config any, configFilePathName string) error {
	if configFilePathName == "" {
		viper.SetConfigType("json")
		viper.SetConfigName("config")
		viper.AddConfigPath(".") // looking for config file in working directory
	} else {
		viper.SetConfigFile(configFilePathName)
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("fail to read config file,%w", err)
		}
		// trying to populate all config from environment variables
		// when the config file doesn't exist , then we feed all fields default value as empty
		// allow environment variable to override it
		buf, err := json.Marshal(config)
		if err != nil {
			return fmt.Errorf("fail to marshal default to json,err: %w", err)
		}
		if err := viper.ReadConfig(bytes.NewBuffer(buf)); err != nil {
			return fmt.Errorf("fail to read config file,err:%w", err)
		}
	}

	if err := viper.Unmarshal(config); err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}
	return nil
}
