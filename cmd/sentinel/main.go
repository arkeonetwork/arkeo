package main

import (
	"flag"
	"fmt"
	"github.com/arkeonetwork/arkeo/app"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/sentinel"
	"github.com/arkeonetwork/arkeo/sentinel/conf"
	"os"
)

func main() {
	configPath := flag.String("config", "", "Path to sentinel config YAML")
	flag.Parse()

	if *configPath == "" {
		fmt.Println("Error: --config flag is required")
		os.Exit(1)
	}

	c := cosmos.GetConfig()
	c.SetBech32PrefixForAccount(app.AccountAddressPrefix, app.AccountAddressPrefix+"pub")

	config, err := conf.LoadConfigurationFromFile(*configPath)
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}
	proxy, err := sentinel.NewProxy(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	proxy.Run()
}
