package cli

import (
	"context"

	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdListContracts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-contracts",
		Short: "list all contracts",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllContractRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ContractAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-contract [provider] [chain] [client]",
		Short: "shows a contract",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argPubKey := args[0]
			argChain := args[1]
			argClient := args[2]

			params := &types.QueryFetchContractRequest{
				Pubkey: argPubKey,
				Chain:  argChain,
				Client: argClient,
			}

			res, err := queryClient.FetchContract(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
