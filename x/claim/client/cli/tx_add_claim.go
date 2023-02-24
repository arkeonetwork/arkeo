//go:build !testnet

package cli

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdAddClaim() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-claim [chain] [address] [amount]",
		Short: "Broadcast message add-claim",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			panic("add-claim is only available for testnet build!")
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
