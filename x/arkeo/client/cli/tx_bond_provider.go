package cli

import (
	"fmt"
	"strings"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdBondProvider() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bond-provider [pubkey] [service] [bond]",
		Short: "Broadcast message bondProvider",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPubkey := args[0]
			argService := strings.ToLower(args[1])
			argBond := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pubkey, err := common.NewPubKey(argPubkey)
			if err != nil {
				return err
			}

			bond, ok := cosmos.NewIntFromString(argBond)
			if !ok {
				return fmt.Errorf("bad bond amount: %s", argBond)
			}

			if ok, _ := serviceExists(clientCtx, argService); !ok {
				cmd.Printf("warning: service %s not found in registry; proceeding anyway\n", argService)
			}

			msg := types.NewMsgBondProvider(
				clientCtx.GetFromAddress(),
				pubkey,
				argService,
				bond,
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
