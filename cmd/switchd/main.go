package main

import (
	switchd "mercury/switch"
)

func main() {
	proxy := switchd.NewProxy()
	proxy.Run()
}
