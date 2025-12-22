package cli

import (
	"context"
	"strings"

	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/cosmos/cosmos-sdk/client"
)

// serviceExists checks the registry for a service name; returns true if found.
func serviceExists(clientCtx client.Context, name string) (bool, error) {
	queryClient := types.NewQueryClient(clientCtx)
	resp, err := queryClient.Service(context.Background(), &types.QueryServiceRequest{Name: name})
	if err != nil {
		return false, err
	}
	return resp != nil && strings.EqualFold(resp.Service.Name, name), nil
}

// loadServiceMap fetches all services from the registry and returns a name->id map.
func loadServiceMap(clientCtx client.Context) (map[string]uint64, error) {
	queryClient := types.NewQueryClient(clientCtx)
	resp, err := queryClient.AllServices(context.Background(), &types.QueryAllServicesRequest{})
	if err != nil {
		return nil, err
	}
	out := make(map[string]uint64, len(resp.Services))
	for _, svc := range resp.Services {
		out[strings.ToLower(svc.Name)] = uint64(svc.ServiceId)
	}
	return out, nil
}
