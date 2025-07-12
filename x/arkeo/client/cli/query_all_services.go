package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func CmdAllServices() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-services",
		Short: "Query all service enums and descriptions",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllServicesRequest{}

			res, err := queryClient.AllServices(context.Background(), params)
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

	return cmd
}
