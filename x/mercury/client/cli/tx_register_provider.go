package cli

import (
    "strconv"
	
	"github.com/spf13/cobra"
    "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"mercury/x/mercury/types"
)

var _ = strconv.Itoa(0)

func CmdRegisterProvider() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-provider [chain] [pubkey]",
		Short: "Broadcast message registerProvider",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
      		 argChain := args[0]
             argPubkey := args[1]
            
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRegisterProvider(
				clientCtx.GetFromAddress().String(),
				argChain,
				argPubkey,
				
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