package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = QueryServer{}

type QueryServer struct {
	Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &QueryServer{
		Keeper: k,
	}
}

func (q QueryServer) QueryLockedVault(c context.Context, req *types.QueryLockedVaultRequest) (*types.QueryLockedVaultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	item, found := q.GetLockedVault(ctx, req.AppId, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "locked-vault does not exist for id %d", req.Id)
	}

	return &types.QueryLockedVaultResponse{
		LockedVault: item,
	}, nil
}

func (q QueryServer) QueryLockedVaults(c context.Context, req *types.QueryLockedVaultsRequest) (*types.QueryLockedVaultsResponse, error) {
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

func (q QueryServer) QueryLiquidationWhiteListing(c context.Context, req *types.QueryLiquidationWhiteListingRequest) (*types.QueryLiquidationWhiteListingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	item, found := q.GetLiquidationWhiteListing(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "Liquidation WhiteListing does not exist for app_id %d", req.AppId)
	}

	return &types.QueryLiquidationWhiteListingResponse{
		LiquidationWhiteListing: item,
	}, nil
}

func (q QueryServer) QueryLiquidationWhiteListings(c context.Context, req *types.QueryLiquidationWhiteListingsRequest) (*types.QueryLiquidationWhiteListingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.LiquidationWhiteListing
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.LiquidationWhiteListingKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.LiquidationWhiteListing
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

	return &types.QueryLiquidationWhiteListingsResponse{
		LiquidationWhiteListings: items,
		Pagination:               pagination,
	}, nil
}

func (q QueryServer) QueryLockedVaultsHistory(c context.Context, req *types.QueryLockedVaultsHistoryRequest) (*types.QueryLockedVaultsHistoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.LockedVault
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.LockedVaultDataKeyHistory),
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

func (q QueryServer) QueryAppReserveFundsTxData(c context.Context, req *types.QueryAppReserveFundsTxDataRequest) (*types.QueryAppReserveFundsTxDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	item, found := q.GetAppReserveFundsTxData(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "AppReserveFunds Tx Data does not exist for app_id %d", req.AppId)
	}

	return &types.QueryAppReserveFundsTxDataResponse{
		AppReserveFundsTxData: item,
	}, nil
}
