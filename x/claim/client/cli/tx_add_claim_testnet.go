//go:build testnet

package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/claim/types"
)

func CmdAddClaim() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-claim [chain] [address] [amount]",
		Short: "Broadcast message add-claim",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argChain := args[0]
			chain, err := types.ChainFromString(argChain)
			if err != nil {
				return fmt.Errorf("invalid chain(%s),err: %w", argChain, err)
			}
			argAddress := args[1]
			addr, err := cosmos.AccAddressFromBech32(argAddress)
			if err != nil {
				return err
			}

			argAmount, err := cast.ToInt64E(args[2])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddClaim(
				clientCtx.GetFromAddress(),
				chain,
				addr,
				argAmount,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
