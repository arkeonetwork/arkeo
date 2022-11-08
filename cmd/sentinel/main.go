package main

import (
	"mercury/sentinel"
)

func main() {
	proxy := sentinel.NewProxy()
	proxy.Run()
}
