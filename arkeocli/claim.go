package arkeocli

import (
	"context"
	"strconv"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/spf13/cobra"

	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func newClaimCmd() *cobra.Command {
	claimCmd := &cobra.Command{
		Use:   "claim",
		Short: "claim accrued contract income",
		Args:  cobra.ExactArgs(0),
		RunE:  runClaimCmd,
	}

	flags.AddTxFlagsToCmd(claimCmd)
	claimCmd.Flags().Uint64("contract-id", 0, "id of contract")
	claimCmd.Flags().String("provider-pubkey", "", "provider pubkey")
	claimCmd.Flags().String("client-pubkey", "", "client pubkey")
	claimCmd.Flags().String("service", "", "service name")
	claimCmd.Flags().Int64("nonce", 0, "requests claimed (must increment each call)")
	return claimCmd
}

func runClaimCmd(cmd *cobra.Command, args []string) (err error) {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	key, err := ensureKeys(cmd)
	if err != nil {
		return err
	}
	spenderAddress, err := key.GetAddress()
	if err != nil {
		return
	}

	clientCtx = clientCtx.WithFromName(key.Name).WithFromAddress(spenderAddress)
	if err = client.SetCmdClientContext(cmd, clientCtx); err != nil {
		return
	}

	contract, err := getContract(cmd)
	if err != nil {
		return err
	}

	nonce, err := cmd.Flags().GetInt64("nonce")
	if err != nil {
		return err
	}
	if nonce == 0 {
		nonceString, err := promptForArg(cmd, "Specify nonce: ")
		if err != nil {
			return err
		}
		nonce, err = strconv.ParseInt(nonceString, 10, 64)
		if err != nil {
			return err
		}
	}

	chainId, err := cmd.Flags().GetString("chain-id")
	if err != nil {
		return err
	}

	if len(chainId) == 0 {
		chainId, err = promptForArg(cmd, "specify the chain id:")
		if err != nil {
			return err
		}
	}

	clientPubkey := contract.GetDelegate()
	if clientPubkey.IsEmpty() {
		clientPubkey = contract.GetClient()
	}

	creatorAddr, err := clientPubkey.GetMyAddress()
	if err != nil {
		return err
	}

	node, err := clientCtx.GetNode()
	if err != nil {
		return errors.Wrapf(err, "failed to get node")
	}

	status, err := node.Status(context.Background())
	if err != nil {
		return errors.Wrapf(err, "failed to get node status")
	}

	signatureExpiry := status.SyncInfo.LatestBlockHeight + types.ExpirationDelta

	signBytes := types.GetBytesToSign(contract.Id, nonce, chainId)
	signature, _, err := clientCtx.Keyring.Sign(key.Name, signBytes, signing.SignMode_SIGN_MODE_DIRECT)
	if err != nil {
		return errors.Wrapf(err, "error signing")
	}

	msg := types.NewMsgClaimContractIncome(
		creatorAddr,
		contract.Id,
		nonce,
		signature,
		chainId,
		signatureExpiry,
	)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
}
