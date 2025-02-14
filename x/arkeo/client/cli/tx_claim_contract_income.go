package cli

import (
	"context"
	"encoding/hex"

	"cosmossdk.io/errors"

	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdClaimContractIncome() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-contract-income [contract-id] [nonce] [signature] [chain-id]",
		Short: "Broadcast message claimContractIncome",
		Args:  cobra.ExactArgs(5),
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

			chainId := args[3]

			node, err := clientCtx.GetNode()
			if err != nil {
				return errors.Wrapf(err, "failed to get node")
			}

			status, err := node.Status(context.Background())
			if err != nil {
				return errors.Wrapf(err, "failed to get node status")
			}

			signatureExpiry := status.SyncInfo.LatestBlockHeight + types.ExpirationDelta

			msg := types.NewMsgClaimContractIncome(
				clientCtx.GetFromAddress(),
				argContractId,
				argNonce,
				signature,
				chainId,
				signatureExpiry,
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
