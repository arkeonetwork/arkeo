//go:build !regtest
// +build !regtest

package main

import (
	"os"

	"github.com/arkeonetwork/arkeo/app"
	"github.com/arkeonetwork/arkeo/arkeocli"
	"github.com/arkeonetwork/arkeo/x/claim/client/offchain"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/ignite/cli/ignite/pkg/cosmoscmd"
	"github.com/ignite/cli/ignite/pkg/xstrings"
)

func main() {
	rootCmd, _ := cosmoscmd.NewRootCmd(
		app.Name,
		app.AccountAddressPrefix,
		app.DefaultNodeHome,
		xstrings.NoDash(app.Name),
		app.ModuleBasics,
		app.New,
		// this line is used by starport scaffolding # root/arguments
	)
	// add in arkeo specific utilities
	rootCmd.AddCommand(arkeocli.GetArkeoCmd())
	rootCmd.AddCommand(offchain.OffChain())
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
