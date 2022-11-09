package main

import (
	"arkeo/sentinel"
	"arkeo/sentinel/conf"
)

func main() {
	config := conf.NewConfiguration()
	proxy := sentinel.NewProxy(config)
	proxy.Run()
}
