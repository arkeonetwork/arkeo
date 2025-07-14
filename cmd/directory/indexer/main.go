package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/arkeonetwork/arkeo/app"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/common/logging"
	"github.com/arkeonetwork/arkeo/common/utils"
	"github.com/arkeonetwork/arkeo/directory/indexer"
)

var (
	log        = logging.WithoutFields()
	configPath = flag.String("config", "", "Path to config file")
)

func main() {
	log.Info("starting indexer")
	cosmos.GetConfig().SetBech32PrefixForAccount(app.AccountAddressPrefix, app.AccountAddressPrefix+"pub")
	flag.Parse()
	if *configPath == "" {
		log.Panic("No config file provided. Use --config <path>")
	}
	var c indexer.ServiceParams
	if err := utils.LoadFromEnv(&c, *configPath); err != nil {
		log.Panicf("failed to load config from file: %+v", err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	indexApp, err := indexer.NewIndexer(c)
	if err != nil {
		panic(err)
	}
	if err := indexApp.Run(); err != nil {
		panic(fmt.Sprintf("error starting indexer: %+v", err))
	}
	<-quit
	log.Info("receive signal to shutdown indexer")
	if err := indexApp.Close(); err != nil {
		panic(err)
	}
	log.Info("indexer shutdown successfully")
}
