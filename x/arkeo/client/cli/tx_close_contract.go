package cli

import (
	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdCloseContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close-contract [contract-id] [client-pubkey] [delegate-optional]",
		Short: "Broadcast message closeContract",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argContractId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argClient := args[1]

			cl, err := common.NewPubKey(argClient)
			if err != nil {
				return err
			}

			delegate := common.EmptyPubKey
			if len(args) > 2 {
				delegate, err = common.NewPubKey(args[2])
				if err != nil {
					return err
				}
			}

			msg := types.NewMsgCloseContract(
				clientCtx.GetFromAddress(),
				argContractId,
				cl,
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
