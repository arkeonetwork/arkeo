//go:build !testnet

package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdAddClaim() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-claim [chain] [address] [amount]",
		Short: "Broadcast message add-claim",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return fmt.Errorf("add-claim command only available on testnet build")
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
