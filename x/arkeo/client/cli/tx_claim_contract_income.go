package cli

import (
	"context"
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
		Use:   "claim-contract-income [contract-id] [nonce] [signature] [chain-id]",
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

			chainID := args[3]

			node, err := clientCtx.GetNode()
			if err != nil {
				return err
			}

			status, err := node.Status(context.Background())
			if err != nil {
				return err
			}
			currentBlock := status.SyncInfo.LatestBlockHeight

			// Get expiration delta from flags (default to 50 blocks if not specified)
			expirationDelta, err := cmd.Flags().GetInt64("expiration-delta")
			if err != nil {
				expirationDelta = 50
			}

			// Calculate expiration block
			expiresAtBlock := currentBlock + expirationDelta

			msg := types.NewMsgClaimContractIncome(
				clientCtx.GetFromAddress(),
				argContractId,
				argNonce,
				signature,
				chainID,
				expiresAtBlock,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().Int64("expiration-delta", 50, "number of blocks until expiration")

	return cmd
}
