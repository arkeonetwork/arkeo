package cli

import (
	"encoding/hex"

	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdClaimContractIncome() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-contract-income [contract-id] [nonce] [signature]",
		Short: "Broadcast message claimContractIncome",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argContractId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}

			argNonce, err := cast.ToInt64E(args[1])
			if err != nil {
				return err
			}
			signature, err := hex.DecodeString(args[2])
			if err != nil {
				return err
			}
			msg := types.NewMsgClaimContractIncome(
				clientCtx.GetFromAddress().String(),
				argContractId,
				argNonce,
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
