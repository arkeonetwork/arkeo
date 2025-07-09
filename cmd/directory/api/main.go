package main

import (
	"flag"
	"fmt"

	"github.com/arkeonetwork/arkeo/app"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/common/logging"
	"github.com/arkeonetwork/arkeo/common/utils"
	"github.com/arkeonetwork/arkeo/directory/api"
)

var (
	log        = logging.WithoutFields()
	configPath = flag.String("config", "", "path to config file")
)

func main() {
	log.Info("starting api")
	cosmos.GetConfig().SetBech32PrefixForAccount(app.AccountAddressPrefix, app.AccountAddressPrefix+"pub")
	flag.Parse()
	var c api.ServiceParams
	if err := utils.LoadFromEnv(&c, *configPath); err != nil {
		log.Panicf("failed to load config from yaml: %+v", err)
	}
	apiService := api.NewApiService(c)
	done, err := apiService.Start()
	if err != nil {
		panic(fmt.Sprintf("error starting api service: %+v", err))
	}
	<-done
	log.Info("api complete")
}
