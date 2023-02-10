package cli

import (
	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdCloseContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close-contract [pubkey] [chain] [client] [delegate-optional]",
		Short: "Broadcast message closeContract",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPubkey := args[0]
			argChain := args[1]
			argClient := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pubkey, err := common.NewPubKey(argPubkey)
			if err != nil {
				return err
			}

			client, err := common.NewPubKey(argClient)
			if err != nil {
				return err
			}

			delegate := common.EmptyPubKey
			if len(args) > 3 {
				delegate, err = common.NewPubKey(args[3])
				if err != nil {
					return err
				}
			}

			msg := types.NewMsgCloseContract(
				clientCtx.GetFromAddress().String(),
				pubkey,
				argChain,
				client,
				delegate,
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
