package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/liquidation/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = (*QueryServer)(nil)

type QueryServer struct {
	Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &QueryServer{
		Keeper: k,
	}
}

func (q *QueryServer) QueryParams(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	var (
		ctx    = sdk.UnwrapSDKContext(c)
		params = q.GetParams(ctx)
	)

	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

func (q *QueryServer) QueryLockedVault(c context.Context, req *types.QueryLockedVaultRequest) (*types.QueryLockedVaultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	item, found := q.GetLockedVault(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "locked-vault does not exist for id %d", req.Id)
	}

	return &types.QueryLockedVaultResponse{
		LockedVault: item,
	}, nil
}

func (q *QueryServer) QueryLockedVaults(c context.Context, req *types.QueryLockedVaultsRequest) (*types.QueryLockedVaultsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.LockedVault
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.LockedVaultKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.LockedVault
			if err := q.cdc.Unmarshal(value, &item); err != nil {
				return false, err
			}

			if accumulate {
				items = append(items, item)
			}

			return true, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryLockedVaultsResponse{
		LockedVaults: items,
		Pagination:   pagination,
	}, nil
}

func (q *QueryServer) QueryLockedVaultsHistory(c context.Context, req *types.QueryLockedVaultsHistoryRequest) (*types.QueryLockedVaultsHistoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.LockedVault
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.LockedVaultKeyHistory),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.LockedVault
			if err := q.cdc.Unmarshal(value, &item); err != nil {
				return false, err
			}

			if accumulate {
				items = append(items, item)
			}

			return true, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryLockedVaultsHistoryResponse{
		LockedVaultsHistory: items,
		Pagination:          pagination,
	}, nil
}

func (q *QueryServer) QueryUserLockedVaults(c context.Context, req *types.QueryUserLockedVaultsRequest) (*types.QueryUserLockedVaultsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.LockedVault
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.LockedVaultKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.LockedVault
			if err := q.cdc.Unmarshal(value, &item); err != nil {
				return false, err
			}

			if accumulate && item.Owner == req.UserAddress {
				items = append(items, item)
			}

			return true, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryUserLockedVaultsResponse{
		UserLockedVaults: items,
		Pagination:       pagination,
	}, nil
}

func (q *QueryServer) QueryUserLockedVaultsHistory(c context.Context, req *types.QueryUserLockedVaultsHistoryRequest) (*types.QueryUserLockedVaultsHistoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.LockedVault
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.LockedVaultKeyHistory),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.LockedVault
			if err := q.cdc.Unmarshal(value, &item); err != nil {
				return false, err
			}

			if accumulate && item.Owner == req.UserAddress {
				items = append(items, item)
			}

			return true, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryUserLockedVaultsHistoryResponse{
		UserLockedVaultsHistory: items,
		Pagination:              pagination,
	}, nil
}

func (q *QueryServer) QueryLockedVaultsPair(c context.Context, req *types.QueryLockedVaultsPairRequest) (*types.QueryLockedVaultsPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.LockedVault
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.LockedVaultKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.LockedVault
			if err := q.cdc.Unmarshal(value, &item); err != nil {
				return false, err
			}

			if accumulate && item.ExtendedPairId == req.PairId {
				items = append(items, item)
			}

			return true, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryLockedVaultsPairResponse{
		LockedVaultsPair: items,
		Pagination:       pagination,
	}, nil
}

func (q *QueryServer) QueryAppIds(c context.Context, _ *types.QueryAppIdsRequest) (*types.QueryAppIdsResponse, error) {
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item := q.GetAppIds(ctx)
	return &types.QueryAppIdsResponse{
		WhitelistedAppIds: item,
	}, nil
}
