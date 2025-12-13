package keeper

import (
	"context"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k KVStore) ProviderAll(c context.Context, req *types.QueryAllProviderRequest) (*types.QueryAllProviderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var providers []types.Provider
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	providerStore := prefix.NewStore(store, types.KeyPrefix(prefixProvider.String()))

	pageRes, err := query.Paginate(providerStore, req.Pagination, func(key, value []byte) error {
		var provider types.Provider
		if err := k.cdc.Unmarshal(value, &provider); err != nil {
			return err
		}

		providers = append(providers, provider)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllProviderResponse{Provider: providers, Pagination: pageRes}, nil
}

func (k KVStore) FetchProvider(c context.Context, req *types.QueryFetchProviderRequest) (*types.QueryFetchProviderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	pk, err := common.NewPubKey(req.Pubkey)
	if err != nil {
		return nil, status.Error(codes.NotFound, "pubkey not found")
	}

	service, _, err := k.ResolveServiceEnum(ctx, req.Service)
	if err != nil {
		return nil, status.Error(codes.NotFound, "service not found")
	}

	val, err := k.GetProvider(ctx, pk, service)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	if val.MetadataNonce <= 0 {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryFetchProviderResponse{Provider: val}, nil
}
