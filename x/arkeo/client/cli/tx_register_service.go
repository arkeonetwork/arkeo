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

func CmdRegisterService() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-service [id] [name] [description] [type]",
		Short: "Register a new service in the registry (authority only)",
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

			msg := types.NewMsgRegisterService(clientCtx.GetFromAddress().String(), id, name, desc, svcType)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
