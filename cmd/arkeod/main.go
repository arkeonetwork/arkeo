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

	rootCmd, _ := NewRootCmd()
	// rootCmd, _ := NewRootCmd(
	// 	app.Name,
	// 	app.AccountAddressPrefix,
	// 	app.DefaultNodeHome,
	// 	xstrings.NoDash(app.Name),
	// 	app.ModuleBasics,
	// 	app.New,
	// 	// this line is used by starport scaffolding # root/arguments
	// )
	// add in arkeo specific utilities
	rootCmd.AddCommand(arkeocli.GetArkeoCmd())
	if err := svrcmd.Execute(rootCmd, "ARKEO", app.DefaultNodeHome); err != nil {
		fmt.Println(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
