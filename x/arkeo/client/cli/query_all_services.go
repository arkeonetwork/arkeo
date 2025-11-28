package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func CmdAllServices() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-services",
		Short: "Query all service enums and descriptions",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllServicesRequest{}

			res, err := queryClient.AllServices(cmd.Context(), params)
			if err != nil {
				return err
			}

			cmd.Println("Arkeo Supported Provider Service List:")
			for _, svc := range res.Services {
				cmd.Printf("- %s : %d (%s)\n", svc.Name, svc.ServiceId, svc.Description)
			}
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
