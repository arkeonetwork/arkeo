package keeper

import (
	"context"
	"fmt"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

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

	pageRes, err := query.Paginate(contractStore, req.Pagination, func(key, value []byte) error {
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

	val, err := k.GetContract(ctx, req.ContractId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if val.Height <= 0 {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryFetchContractResponse{Contract: val}, nil
}

func (k KVStore) ActiveContract(goCtx context.Context, req *types.QueryActiveContractRequest) (*types.QueryActiveContractResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	fmt.Printf(">>>>>>>>>> REQ: %+v\n", req)
	providerPubKey, err := common.NewPubKey(req.Provider)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid provider pubkey")
	}
	service, err := common.NewService(req.Service)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid service")
	}
	spenderPubKey, err := common.NewPubKey(req.Spender)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid spender pubkey")
	}

	activeContract, err := k.GetActiveContractForUser(ctx, spenderPubKey, providerPubKey, service)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}

	if activeContract.IsEmpty() {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryActiveContractResponse{Contract: activeContract}, nil
}
