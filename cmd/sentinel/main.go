package main

import (
	"github.com/ArkeoNetwork/arkeo/app"
	"github.com/ArkeoNetwork/arkeo/common/cosmos"
	"github.com/ArkeoNetwork/arkeo/sentinel"
	"github.com/ArkeoNetwork/arkeo/sentinel/conf"
)

func main() {
	c := cosmos.GetConfig()
	c.SetBech32PrefixForAccount(app.AccountAddressPrefix, app.AccountAddressPrefix+"pub")

	config := conf.NewConfiguration()
	proxy := sentinel.NewProxy(config)
	proxy.Run()
}
