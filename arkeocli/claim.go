package arkeocli

import (
	"strconv"

	"cosmossdk.io/errors"
	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
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

	spenderPubkeyRaw, err := key.GetPubKey()
	if err != nil {
		return err
	}
	spenderPubkey, err := common.NewPubKeyFromCrypto(spenderPubkeyRaw)
	if err != nil {
		return err
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

	clientPubkey := contract.GetDelegate()
	if clientPubkey == "" {
		clientPubkey = contract.GetClient()
	}

	creatorAddr, err := clientPubkey.GetMyAddress()
	if err != nil {
		return err
	}
	creator := creatorAddr.String()

	signBytes := types.GetBytesToSign(contract.Id, spenderPubkey, nonce)
	signature, _, err := clientCtx.Keyring.Sign(key.Name, signBytes)
	if err != nil {
		return errors.Wrapf(err, "error signing")
	}

	msg := types.NewMsgClaimContractIncome(
		creator,
		contract.Id,
		spenderPubkey,
		nonce,
		signature,
	)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
}
