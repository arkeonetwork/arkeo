package cli

import (
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdActiveContract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "active-contract [spender] [provider] [chain]",
		Short: "Query active_contract",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqSpender := args[0]
			reqProvider := args[1]
			reqChain := args[2]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryActiveContractRequest{
				Spender:  reqSpender,
				Provider: reqProvider,
				Chain:    reqChain,
			}

			res, err := queryClient.ActiveContract(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
