package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/arkeonetwork/arkeo/app"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/common/logging"
	"github.com/arkeonetwork/arkeo/common/utils"
	"github.com/arkeonetwork/arkeo/directory/api"
	"github.com/arkeonetwork/arkeo/directory/db"
)

type Config struct {
	ApiListenAddr  string `mapstructure:"API_LISTEN"`
	ApiStaticDir   string `mapstructure:"API_STATIC_DIR"`
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
	log         = logging.WithoutFields()
	envPath     = flag.String("env", "", "path to env file (default: use os env)")
	configNames = []string{
		"API_LISTEN",
		"API_STATIC_DIR",
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	log.Info("starting api")

	cosmos.GetConfig().SetBech32PrefixForAccount(app.AccountAddressPrefix, app.AccountAddressPrefix+"pub")

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
	// TODO determine config mechanism
	api := api.NewApiService(api.ApiServiceParams{
		ListenAddr: c.ApiListenAddr,
		StaticDir:  c.ApiStaticDir,
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
	done, err := api.Start()
	if err != nil {
		panic(fmt.Sprintf("error starting api service: %+v", err))
	}
	<-done
	log.Info("api complete")
}
