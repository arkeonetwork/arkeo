package main

import (
	"github.com/ArkeoNetwork/arkeo-protocol/app"
	"github.com/ArkeoNetwork/arkeo-protocol/common/cosmos"
	"github.com/ArkeoNetwork/arkeo-protocol/sentinel"
	"github.com/ArkeoNetwork/arkeo-protocol/sentinel/conf"
)

func main() {
	c := cosmos.GetConfig()
	c.SetBech32PrefixForAccount(app.AccountAddressPrefix, app.AccountAddressPrefix+"pub")

	config := conf.NewConfiguration()
	proxy := sentinel.NewProxy(config)
	proxy.Run()
}
