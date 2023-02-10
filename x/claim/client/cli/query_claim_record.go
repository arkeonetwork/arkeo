package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/arkeonetwork/arkeo/x/claim/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdClaimRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-record [address] [chain]",
		Short: "Query claim-record",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			reqAddress := args[0]
			reqChain := args[1]

			// validate chain
			chainId, ok := types.Chain_value[strings.ToUpper(reqChain)]
			if !ok {
				return fmt.Errorf("invalid chain %s", reqChain)
			}
			chain := types.Chain(chainId)

			// validate address if valid based on chain
			if !types.IsValidAddress(reqAddress, chain) {
				return fmt.Errorf("invalid address %s", reqAddress)
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryClaimRecordRequest{
				Chain:   chain,
				Address: reqAddress,
			}

			res, err := queryClient.ClaimRecord(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
