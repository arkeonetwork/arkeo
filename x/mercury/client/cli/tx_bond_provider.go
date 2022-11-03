package cli

import (
	"fmt"
	"mercury/common"
	"mercury/common/cosmos"
	"mercury/x/mercury/types"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdBondProvider() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bond-provider [pubkey] [chain] [add] [remove]",
		Short: "Broadcast message bondProvider",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPubkey := args[0]
			argChain := args[1]
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

			msg := types.NewMsgBondProvider(
				clientCtx.GetFromAddress().String(),
				pubkey,
				argChain,
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
