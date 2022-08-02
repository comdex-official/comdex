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
	_ types.QueryServer = QueryServer{}
)

type QueryServer struct {
	Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &QueryServer{
		Keeper: k,
	}
}

func (q QueryServer) QueryLends(c context.Context, req *types.QueryLendsRequest) (*types.QueryLendsResponse, error) {
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

func (q QueryServer) QueryLend(c context.Context, req *types.QueryLendRequest) (*types.QueryLendResponse, error) {
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

func (q QueryServer) QueryAllLendByOwner(c context.Context, req *types.QueryAllLendByOwnerRequest) (*types.QueryAllLendByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
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

func (q QueryServer) QueryAllLendByOwnerAndPool(c context.Context, req *types.QueryAllLendByOwnerAndPoolRequest) (*types.QueryAllLendByOwnerAndPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
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

func (q QueryServer) QueryAllBorrowByOwnerAndPool(c context.Context, req *types.QueryAllBorrowByOwnerAndPoolRequest) (*types.QueryAllBorrowByOwnerAndPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
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

func (q QueryServer) QueryPairs(c context.Context, req *types.QueryPairsRequest) (*types.QueryPairsResponse, error) {
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

func (q QueryServer) QueryPair(c context.Context, req *types.QueryPairRequest) (*types.QueryPairResponse, error) {
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

func (q QueryServer) QueryPools(c context.Context, req *types.QueryPoolsRequest) (*types.QueryPoolsResponse, error) {
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

func (q QueryServer) QueryPool(c context.Context, req *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
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

func (q QueryServer) QueryAssetToPairMappings(c context.Context, req *types.QueryAssetToPairMappingsRequest) (*types.QueryAssetToPairMappingsResponse, error) {
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

func (q QueryServer) QueryAssetToPairMapping(c context.Context, req *types.QueryAssetToPairMappingRequest) (*types.QueryAssetToPairMappingResponse, error) {
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

func (q QueryServer) QueryBorrows(c context.Context, req *types.QueryBorrowsRequest) (*types.QueryBorrowsResponse, error) {
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

func (q QueryServer) QueryBorrow(c context.Context, req *types.QueryBorrowRequest) (*types.QueryBorrowResponse, error) {
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

func (q QueryServer) QueryAllBorrowByOwner(c context.Context, req *types.QueryAllBorrowByOwnerRequest) (*types.QueryAllBorrowByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
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

func (q QueryServer) QueryAssetRatesStats(c context.Context, req *types.QueryAssetRatesStatsRequest) (*types.QueryAssetRatesStatsResponse, error) {
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

func (q QueryServer) QueryAssetRatesStat(c context.Context, req *types.QueryAssetRatesStatRequest) (*types.QueryAssetRatesStatResponse, error) {
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

func (q QueryServer) QueryAssetStats(c context.Context, req *types.QueryAssetStatsRequest) (*types.QueryAssetStatsResponse, error) {
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

func (q QueryServer) QueryModuleBalance(c context.Context, req *types.QueryModuleBalanceRequest) (*types.QueryModuleBalanceResponse, error) {
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

func (q QueryServer) QueryDepositStats(c context.Context, req *types.QueryDepositStatsRequest) (*types.QueryDepositStatsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	depositStatsData, found := q.GetDepositStats(ctx)
	if !found {
		return &types.QueryDepositStatsResponse{}, nil
	}

	return &types.QueryDepositStatsResponse{
		DepositStats: depositStatsData,
	}, nil
}

func (q QueryServer) QueryUserDepositStats(c context.Context, req *types.QueryUserDepositStatsRequest) (*types.QueryUserDepositStatsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	userDepositStatsData, found := q.GetUserDepositStats(ctx)
	if !found {
		return &types.QueryUserDepositStatsResponse{}, nil
	}

	return &types.QueryUserDepositStatsResponse{
		UserDepositStats: userDepositStatsData,
	}, nil
}

func (q QueryServer) QueryReserveDepositStats(c context.Context, req *types.QueryReserveDepositStatsRequest) (*types.QueryReserveDepositStatsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	reserveDepositStatsData, found := q.GetReserveDepositStats(ctx)
	if !found {
		return &types.QueryReserveDepositStatsResponse{}, nil
	}

	return &types.QueryReserveDepositStatsResponse{
		ReserveDepositStats: reserveDepositStatsData,
	}, nil
}

func (q QueryServer) QueryBuyBackDepositStats(c context.Context, req *types.QueryBuyBackDepositStatsRequest) (*types.QueryBuyBackDepositStatsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	buyBackDepositStatsData, found := q.GetBuyBackDepositStats(ctx)
	if !found {
		return &types.QueryBuyBackDepositStatsResponse{}, nil
	}

	return &types.QueryBuyBackDepositStatsResponse{
		BuyBackDepositStats: buyBackDepositStatsData,
	}, nil
}

func (q QueryServer) QueryBorrowStats(c context.Context, req *types.QueryBorrowStatsRequest) (*types.QueryBorrowStatsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	borrowStatsData, found := q.GetBorrowStats(ctx)
	if !found {
		return &types.QueryBorrowStatsResponse{}, nil
	}

	return &types.QueryBorrowStatsResponse{
		BorrowStats: borrowStatsData,
	}, nil
}

func (q QueryServer) QueryAuctionParams(c context.Context, req *types.QueryAuctionParamRequest) (*types.QueryAuctionParamResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetAddAuctionParamsData(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "Auction Params not exist for id %d", req.AppId)
	}

	return &types.QueryAuctionParamResponse{
		AuctionParams: item,
	}, nil
}
