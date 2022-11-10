package main

import (
	"arkeo/app"
	"arkeo/common/cosmos"
	"arkeo/sentinel"
	"arkeo/sentinel/conf"
)

func main() {
	c := cosmos.GetConfig()
	c.SetBech32PrefixForAccount(app.AccountAddressPrefix, app.AccountAddressPrefix+"pub")

	config := conf.NewConfiguration()
	proxy := sentinel.NewProxy(config)
	proxy.Run()
}
