package arkeocli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func newCloseContractCmd() *cobra.Command {
	closeContractCmd := &cobra.Command{
		Use:   "close-contract",
		Short: "close a contract",
		Args:  cobra.ExactArgs(0),
		RunE:  runCloseContractCmd,
	}

	flags.AddTxFlagsToCmd(closeContractCmd)
	closeContractCmd.Flags().Uint64("contract-id", 0, "id of contract")
	closeContractCmd.Flags().String("provider-pubkey", "", "provider pubkey")
	closeContractCmd.Flags().String("client-pubkey", "", "client pubkey")
	closeContractCmd.Flags().String("delegate-pubkey", "", "delegate pubkey")
	closeContractCmd.Flags().String("service", "", "service name")
	return closeContractCmd
}

func runCloseContractCmd(cmd *cobra.Command, args []string) (err error) {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return err
	}

	key, err := ensureKeys(cmd)
	if err != nil {
		return err
	}
	addr, err := key.GetAddress()
	if err != nil {
		return
	}
	clientCtx = clientCtx.WithFromName(key.Name).WithFromAddress(addr)
	if err = client.SetCmdClientContext(cmd, clientCtx); err != nil {
		return
	}

	contract, err := getContract(cmd)
	if err != nil {
		return err
	}

	argClientPubkey, _ := cmd.Flags().GetString("client-pubkey")
	if argClientPubkey == "" {
		argClientPubkey, err = toPubkey(cmd, addr)
		if err != nil {
			return err
		}
	}

	argDelegatePubkey, _ := cmd.Flags().GetString("delegate-pubkey")
	if argDelegatePubkey == "" {
		argDelegatePubkey, err = promptForArg(cmd, "Specify delegate pubkey (leave blank to use client key): ")
		if err != nil {
			return err
		}
	}

	msg := types.NewMsgCloseContract(
		clientCtx.GetFromAddress(),
		contract.Id,
		common.PubKey(argClientPubkey),
		common.PubKey(argDelegatePubkey),
	)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
}
