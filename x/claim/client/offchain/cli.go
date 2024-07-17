package offchain

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func OffChain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "off-chain",
		Short: "Off-chain utilities.",
		Long:  `Utilities for off-chain data.`,
	}

	cmd.AddCommand(
		CmdThorChainTxFetachOffline(),
		CmdSignMessage(),
	)

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdThorChainTxFetachOffline() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "fetch-thor-tx-data",
		Short: "Fetch Thorchain Tx Data",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			txHash := args[0]

			data, err := fetchThorChainTxData(txHash)
			if err != nil {
				return err
			}

			fmt.Println(string(data))

			return nil

		},
	}

	return cmd

}

func CmdSignMessage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign [data]",
		Short: "Sign Off chain Data ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			data := args[0]
			signMessage(data)

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
