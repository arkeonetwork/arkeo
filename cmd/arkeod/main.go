//go:build !regtest
// +build !regtest

package main

import (
	"fmt"
	"os"

	"github.com/arkeonetwork/arkeo/app"
	"github.com/arkeonetwork/arkeo/arkeocli"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {

	// create new root command
	rootCmd, _ := NewRootCmd()
	rootCmd.AddCommand(arkeocli.GetArkeoCmd())
	if err := svrcmd.Execute(rootCmd, "ARKEO", app.DefaultNodeHome); err != nil {
		fmt.Println(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
