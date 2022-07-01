package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/lend/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ types.QueryServer = (*queryServer)(nil)
)

type queryServer struct {
	Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &queryServer{
		Keeper: k,
	}
}

func (q queryServer) QueryLends(c context.Context, req *types.QueryLendsRequest) (*types.QueryLendsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.LendAsset
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.LendUserPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.LendAsset
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

	return &types.QueryLendsResponse{
		Lends:      items,
		Pagination: pagination,
	}, nil
}

func (q queryServer) QueryLend(c context.Context, req *types.QueryLendRequest) (*types.QueryLendResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetLend(ctx, req.Id)
	if !found {
		return &types.QueryLendResponse{}, nil
	}

	return &types.QueryLendResponse{
		Lend: item,
	}, nil
}

func (q queryServer) QueryAllLendByOwner(c context.Context, req *types.QueryAllLendByOwnerRequest) (*types.QueryAllLendByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var ( //nolint:prealloc
		ctx   = sdk.UnwrapSDKContext(c)
		lends []types.LendAsset
	)

	_, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Address is not correct")
	}

	userVaultAssetData, found := q.UserLends(ctx, req.Owner)
	if !found {
		return &types.QueryAllLendByOwnerResponse{}, nil
	}
	lends = append(lends, userVaultAssetData...)

	return &types.QueryAllLendByOwnerResponse{
		Lends: lends,
	}, nil
}

func (q queryServer) QueryAllLendByOwnerAndPool(c context.Context, req *types.QueryAllLendByOwnerAndPoolRequest) (*types.QueryAllLendByOwnerAndPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var ( //nolint:prealloc
		ctx   = sdk.UnwrapSDKContext(c)
		lends []types.LendAsset
	)

	_, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Address is not correct")
	}

	userVaultAssetData, found := q.LendIDByOwnerAndPool(ctx, req.Owner, req.PoolId)
	if !found {
		return &types.QueryAllLendByOwnerAndPoolResponse{}, nil
	}
	lends = append(lends, userVaultAssetData...)

	return &types.QueryAllLendByOwnerAndPoolResponse{
		Lends: lends,
	}, nil
}

func (q queryServer) QueryAllBorrowByOwnerAndPool(c context.Context, req *types.QueryAllBorrowByOwnerAndPoolRequest) (*types.QueryAllBorrowByOwnerAndPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var ( //nolint:prealloc
		ctx     = sdk.UnwrapSDKContext(c)
		borrows []types.BorrowAsset
	)

	_, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Address is not correct")
	}

	userVaultAssetData, found := q.BorrowIDByOwnerAndPool(ctx, req.Owner, req.PoolId)
	if !found {
		return &types.QueryAllBorrowByOwnerAndPoolResponse{}, nil
	}
	borrows = append(borrows, userVaultAssetData...)

	return &types.QueryAllBorrowByOwnerAndPoolResponse{
		Borrows: borrows,
	}, nil
}

func (q queryServer) QueryPairs(c context.Context, req *types.QueryPairsRequest) (*types.QueryPairsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.Extended_Pair
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.LendPairKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.Extended_Pair
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

	return &types.QueryPairsResponse{
		ExtendedPairs: items,
		Pagination:    pagination,
	}, nil
}

func (q queryServer) QueryPair(c context.Context, req *types.QueryPairRequest) (*types.QueryPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetLendPair(ctx, req.Id)
	if !found {
		return &types.QueryPairResponse{}, nil
	}

	return &types.QueryPairResponse{
		ExtendedPair: item,
	}, nil
}

func (q queryServer) QueryPools(c context.Context, req *types.QueryPoolsRequest) (*types.QueryPoolsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.Pool
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.PoolKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.Pool
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

	return &types.QueryPoolsResponse{
		Pools:      items,
		Pagination: pagination,
	}, nil
}

func (q queryServer) QueryPool(c context.Context, req *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetPool(ctx, req.Id)
	if !found {
		return &types.QueryPoolResponse{}, nil
	}

	return &types.QueryPoolResponse{
		Pool: item,
	}, nil
}

func (q queryServer) QueryAssetToPairMappings(c context.Context, req *types.QueryAssetToPairMappingsRequest) (*types.QueryAssetToPairMappingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.AssetToPairMapping
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.AssetToPairMappingKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.AssetToPairMapping
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

	return &types.QueryAssetToPairMappingsResponse{
		AssetToPairMappings: items,
		Pagination:          pagination,
	}, nil
}

func (q queryServer) QueryAssetToPairMapping(c context.Context, req *types.QueryAssetToPairMappingRequest) (*types.QueryAssetToPairMappingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetAssetToPair(ctx, req.AssetId, req.PoolId)
	if !found {
		return &types.QueryAssetToPairMappingResponse{}, nil
	}

	return &types.QueryAssetToPairMappingResponse{
		AssetToPairMapping: item,
	}, nil
}

func (q queryServer) QueryBorrows(c context.Context, req *types.QueryBorrowsRequest) (*types.QueryBorrowsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.BorrowAsset
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.BorrowPairKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.BorrowAsset
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

	return &types.QueryBorrowsResponse{
		Borrows:    items,
		Pagination: pagination,
	}, nil
}

func (q queryServer) QueryBorrow(c context.Context, req *types.QueryBorrowRequest) (*types.QueryBorrowResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetBorrow(ctx, req.Id)
	if !found {
		return &types.QueryBorrowResponse{}, nil
	}

	return &types.QueryBorrowResponse{
		Borrow: item,
	}, nil
}

func (q queryServer) QueryAllBorrowByOwner(c context.Context, req *types.QueryAllBorrowByOwnerRequest) (*types.QueryAllBorrowByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var ( //nolint:prealloc
		ctx     = sdk.UnwrapSDKContext(c)
		borrows []types.BorrowAsset
	)

	_, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Address is not correct")
	}

	userVaultAssetData, found := q.UserBorrows(ctx, req.Owner)
	if !found {
		return &types.QueryAllBorrowByOwnerResponse{}, nil
	}
	borrows = append(borrows, userVaultAssetData...)

	return &types.QueryAllBorrowByOwnerResponse{
		Borrows: borrows,
	}, nil
}

func (q queryServer) QueryAssetRatesStats(c context.Context, req *types.QueryAssetRatesStatsRequest) (*types.QueryAssetRatesStatsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.AssetRatesStats
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.AssetRatesStatsKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.AssetRatesStats
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

	return &types.QueryAssetRatesStatsResponse{
		AssetRatesStats: items,
		Pagination:      pagination,
	}, nil
}

func (q queryServer) QueryAssetRatesStat(c context.Context, req *types.QueryAssetRatesStatRequest) (*types.QueryAssetRatesStatResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetAssetRatesStats(ctx, req.Id)
	if !found {
		return &types.QueryAssetRatesStatResponse{}, nil
	}

	return &types.QueryAssetRatesStatResponse{
		AssetRatesStat: item,
	}, nil
}

func (q queryServer) QueryAssetStats(c context.Context, req *types.QueryAssetStatsRequest) (*types.QueryAssetStatsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	assetStatsData, found := q.AssetStatsByPoolIDAndAssetID(ctx, req.AssetId, req.PoolId)
	if !found {
		return &types.QueryAssetStatsResponse{}, nil
	}

	return &types.QueryAssetStatsResponse{
		AssetStats: assetStatsData,
	}, nil
}

func (q queryServer) QueryModuleBalance(c context.Context, req *types.QueryModuleBalanceRequest) (*types.QueryModuleBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	modBal, found := q.GetModuleBalanceByPoolID(ctx, req.PoolId)
	if !found {
		return &types.QueryModuleBalanceResponse{}, nil
	}

	return &types.QueryModuleBalanceResponse{
		ModuleBalance: modBal,
	}, nil
}
