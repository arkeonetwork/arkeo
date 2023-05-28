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
	log     = logging.WithoutFields()
	envPath = flag.String("env", "", "path to env file (default: use os env)")
)

func main() {
	log.Info("starting indexer")
	cosmos.GetConfig().SetBech32PrefixForAccount(app.AccountAddressPrefix, app.AccountAddressPrefix+"pub")
	flag.Parse()
	var c indexer.ServiceParams
	if err := utils.LoadFromEnv(&c, *envPath); err != nil {
		log.Panicf("failed to load config from env: %+v", err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	indexApp := indexer.NewIndexer(c)
	done, err := indexApp.Run()
	if err != nil {
		panic(fmt.Sprintf("error starting indexer: %+v", err))
	}
	<-done
	log.Info("indexer complete")
}
