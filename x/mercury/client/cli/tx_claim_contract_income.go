package cli

import (
	"encoding/hex"
	"mercury/common"
	"mercury/x/mercury/types"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdClaimContractIncome() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-contract-income [pubkey] [chain] [client] [nonce] [height] [signature]",
		Short: "Broadcast message claimContractIncome",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPubkey := args[0]
			argChain := args[1]
			argClient := args[2]
			argSignature := args[5]
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

			pubkey, err := common.NewPubKey(argPubkey)
			if err != nil {
				return err
			}

			client, err := common.NewPubKey(argClient)
			if err != nil {
				return err
			}

			signature, err := hex.DecodeString(argSignature)
			if err != nil {
				return err
			}

			msg := types.NewMsgClaimContractIncome(
				clientCtx.GetFromAddress().String(),
				pubkey,
				argChain,
				client,
				argNonce,
				argHeight,
				signature,
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
