package keeper

import (
	"context"

	"mercury/common"
	"mercury/x/mercury/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k KVStore) ContractAll(c context.Context, req *types.QueryAllContractRequest) (*types.QueryAllContractResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var contracts []types.Contract
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	contractStore := prefix.NewStore(store, types.KeyPrefix(prefixContract.String()))

	pageRes, err := query.Paginate(contractStore, req.Pagination, func(key []byte, value []byte) error {
		var contract types.Contract
		if err := k.cdc.Unmarshal(value, &contract); err != nil {
			return err
		}

		contracts = append(contracts, contract)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllContractResponse{Contract: contracts, Pagination: pageRes}, nil
}

func (k KVStore) FetchContract(c context.Context, req *types.QueryFetchContractRequest) (*types.QueryFetchContractResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	pk, err := common.NewPubKey(req.Pubkey)
	if err != nil {
		return nil, status.Error(codes.NotFound, "pubkey not found")
	}

	chain, err := common.NewChain(req.Chain)
	if err != nil {
		return nil, status.Error(codes.NotFound, "chain not found")
	}

	client, err := common.NewPubKey(req.Client)
	if err != nil {
		return nil, status.Error(codes.NotFound, "client not found")
	}

	val, err := k.GetContract(ctx, pk, chain, client)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	if val.Height <= 0 {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryFetchContractResponse{Contract: val}, nil
}
