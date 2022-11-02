package cli

import (
	"context"

	"mercury/x/mercury/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdListProviders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-providers",
		Short: "list all providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllProviderRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ProviderAll(context.Background(), params)
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

func CmdShowProvider() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-provider [pubkey] [chain]",
		Short: "shows a provider",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argPubKey := args[0]
			argChain := args[1]

			params := &types.QueryFetchProviderRequest{
				Pubkey: argPubKey,
				Chain:  argChain,
			}

			res, err := queryClient.FetchProvider(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
