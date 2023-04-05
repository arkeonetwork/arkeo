package main

import (
	"flag"
	"fmt"

	"github.com/arkeonetwork/arkeo/app"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/common/logging"
	"github.com/arkeonetwork/arkeo/common/utils"
	"github.com/arkeonetwork/arkeo/directory/db"
	"github.com/arkeonetwork/arkeo/directory/indexer"
)

type Config struct {
	ArkeoApi            string `mapstructure:"ARKEO_API"`
	TendermintApi       string `mapstructure:"TENDERMINT_API"`
	TendermintWs        string `mapstructure:"TENDERMINT_WS"`
	ChainID             string `mapstructure:"CHAIN_ID"`
	Bech32PrefixAccAddr string `mapstructure:"BECH32_PREF_ACC_ADDR"`
	Bech32PrefixAccPub  string `mapstructure:"BECH32_PREF_ACC_PUB"`
	DBHost              string `mapstructure:"DB_HOST"`
	DBPort              uint   `mapstructure:"DB_PORT"`
	DBUser              string `mapstructure:"DB_USER"`
	DBPass              string `mapstructure:"DB_PASS"`
	DBName              string `mapstructure:"DB_NAME"`
	DBSSLMode           string `mapstructure:"DB_SSL_MODE"`
	DBPoolMaxConns      int    `mapstructure:"DB_POOL_MAX_CONNS"`
	DBPoolMinConns      int    `mapstructure:"DB_POOL_MIN_CONNS"`
}

var (
	log         = logging.WithoutFields()
	envPath     = flag.String("env", "", "path to env file (default: use os env)")
	configNames = []string{
		"ARKEO_API",
		"TENDERMINT_API",
		"TENDERMINT_WS",
		"CHAIN_ID",
		"BECH32_PREF_ACC_ADDR",
		"BECH32_PREF_ACC_PUB",
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

func main() {
	cosmosConfig := cosmos.GetConfig()
	cosmosConfig.SetBech32PrefixForAccount(app.AccountAddressPrefix, app.AccountAddressPrefix+"pub")

	log.Info("starting indexer")
	flag.Parse()
	c := &Config{}
	if *envPath == "" {
		if err := utils.LoadFromEnv(c, configNames...); err != nil {
			log.Panicf("failed to load config from env: %+v", err)
		}
	} else {
		if err := utils.Load(*envPath, c); err != nil {
			log.Panicf("failed to load config: %+v", err)
		}
	}

	app := indexer.NewIndexer(indexer.IndexerAppParams{
		ChainID:             c.ChainID,
		Bech32PrefixAccAddr: c.Bech32PrefixAccAddr,
		Bech32PrefixAccPub:  c.Bech32PrefixAccPub,
		ArkeoApi:            c.ArkeoApi,
		TendermintApi:       c.TendermintApi,
		TendermintWs:        c.TendermintWs,
		DBConfig: db.DBConfig{
			Host:         c.DBHost,
			Port:         c.DBPort,
			User:         c.DBUser,
			Pass:         c.DBPass,
			DBName:       c.DBName,
			PoolMaxConns: c.DBPoolMaxConns,
			PoolMinConns: c.DBPoolMinConns,
			SSLMode:      c.DBSSLMode,
		},
	})
	done, err := app.Run()
	if err != nil {
		panic(fmt.Sprintf("error starting indexer: %+v", err))
	}
	<-done
	log.Info("indexer complete")
}
