package arkeocli

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var (
	arkeoCmd = &cobra.Command{
		Use:   "arkeo",
		Short: "arkeo subcommands",
	}
)

func GetArkeoCmd() *cobra.Command {
	flags.AddTxFlagsToCmd(bondProviderCmd)
	bondProviderCmd.Flags().StringP("pubkey", "p", "", "provider pubkey")
	bondProviderCmd.Flags().StringP("chain", "c", "", "provider chain")
	bondProviderCmd.Flags().String("bond", "", "provider bond amount")
	arkeoCmd.AddCommand(bondProviderCmd)

	return arkeoCmd
}
