package arkeocli

import (
	"github.com/spf13/cobra"
)

func GetArkeoCmd() *cobra.Command {
	arkeoCmd := &cobra.Command{
		Use:   "arkeo",
		Short: "arkeo subcommands",
	}
	arkeoCmd.AddCommand(newBondProviderCmd())
	arkeoCmd.AddCommand(newModProviderCmd())
	arkeoCmd.AddCommand(newOpenContractCmd())
	arkeoCmd.AddCommand(newShowPubkeyCmd())
	arkeoCmd.AddCommand(newClaimCmd())
	arkeoCmd.AddCommand(newCloseContractCmd())
	return arkeoCmd
}
