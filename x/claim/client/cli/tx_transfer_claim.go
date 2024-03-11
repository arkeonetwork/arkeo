package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/claim/types"
)

func CmdTransferClaim() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-claim [to-address]",
		Short: "Broadcast message transfer-claim",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argToAddress := args[0]
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			toAddress, err := cosmos.AccAddressFromBech32(argToAddress)
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferClaim(
				clientCtx.GetFromAddress().String(),
				toAddress.String(),
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
