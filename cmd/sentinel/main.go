package main

import (
	"fmt"

	"github.com/arkeonetwork/arkeo/app"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/sentinel"
	"github.com/arkeonetwork/arkeo/sentinel/conf"
)

func main() {
	c := cosmos.GetConfig()
	c.SetBech32PrefixForAccount(app.AccountAddressPrefix, app.AccountAddressPrefix+"pub")

	config := conf.NewConfiguration()
	proxy, err := sentinel.NewProxy(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	proxy.Run()
}
