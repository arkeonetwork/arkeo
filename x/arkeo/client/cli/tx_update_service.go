package cli

import (
	"fmt"

	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdUpdateService() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-service [id] [name] [description] [type]",
		Short: "Update an existing service in the registry (authority only)",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := cast.ToUint64E(args[0])
			if err != nil {
				return fmt.Errorf("invalid id: %w", err)
			}
			name := args[1]
			desc := args[2]
			svcType := args[3]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateService(clientCtx.GetFromAddress().String(), id, name, desc, svcType)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
