package cli

import (
	"strconv"

	"mercury/common"
	"mercury/x/mercury/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdClaimContractIncome() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-contract-income [pubkey] [chain] [client] [nonce] [height]",
		Short: "Broadcast message claimContractIncome",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPubkey := args[0]
			argChain := args[1]
			argClient := args[2]
			argNonce, err := cast.ToInt64E(args[3])
			if err != nil {
				return err
			}
			argHeight, err := cast.ToInt64E(args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			chain, err := common.NewChain(argChain)
			if err != nil {
				return err
			}

			pubkey, err := common.NewPubKey(argPubkey)
			if err != nil {
				return err
			}

			msg := types.NewMsgClaimContractIncome(
				clientCtx.GetFromAddress().String(),
				pubkey,
				chain,
				argClient,
				argNonce,
				argHeight,
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
