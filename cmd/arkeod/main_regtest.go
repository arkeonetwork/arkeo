//go:build regtest
// +build regtest

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/arkeonetwork/arkeo/app"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd, _ := NewRootCmd()

	// for coverage data we need to exit main without allowing the server to call os.Exit

	syn := make(chan error)
	go func() {
		syn <- svrcmd.Execute(rootCmd, "", app.DefaultNodeHome)
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGUSR1)
	select {
	case <-sig:
	case err := <-syn:
		if err != nil {
			os.Exit(1)
		}
	}
}
