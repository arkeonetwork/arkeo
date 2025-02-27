package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/arkeonetwork/arkeo/x/claim/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group claim queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdClaimRecord())
	cmd.AddCommand(CmdClaimableForAction())

	// this line is used by starport scaffolding # 1

	return cmd
}

func CmdClaimableForAction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claimable [address] [action] [chain]",
		Short: "Query claimable amount for a specific action",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			address := args[0]

			action, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			chain, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryClaimableForActionRequest{
				Address: address,
				Action:  types.Action(action),
				Chain:   types.Chain(chain),
			}

			res, err := queryClient.ClaimableForAction(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
