package main

import (
	"mercury/sentinel"
	"mercury/sentinel/conf"
)

func main() {
	config := conf.NewConfiguration()
	proxy := sentinel.NewProxy(config)
	proxy.Run()
}
