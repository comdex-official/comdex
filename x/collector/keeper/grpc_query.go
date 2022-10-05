package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/comdex-official/comdex/x/collector/types"
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

func (q QueryServer) QueryCollectorLookupByApp(c context.Context, req *types.QueryCollectorLookupByAppRequest) (*types.QueryCollectorLookupByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)

	collectorLookupData, found := q.GetCollectorLookupTableByApp(ctx, req.AppId)
	if !found {
		return &types.QueryCollectorLookupByAppResponse{}, nil
	}

	return &types.QueryCollectorLookupByAppResponse{
		CollectorLookup: collectorLookupData,
	}, nil
}

func (q QueryServer) QueryCollectorLookupByAppAndAsset(c context.Context, req *types.QueryCollectorLookupByAppAndAssetRequest) (*types.QueryCollectorLookupByAppAndAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
	}

	collectorLookupData, found := q.GetCollectorLookupTable(ctx, req.AppId, req.AssetId)
	if !found {
		return &types.QueryCollectorLookupByAppAndAssetResponse{}, nil
	}

	return &types.QueryCollectorLookupByAppAndAssetResponse{
		CollectorLookup: collectorLookupData,
	}, nil
}

func (q QueryServer) QueryCollectorDataByAppAndAsset(c context.Context, req *types.QueryCollectorDataByAppAndAssetRequest) (*types.QueryCollectorDataByAppAndAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
	}
	collectorMap, _ := q.GetAppidToAssetCollectorMapping(ctx, req.AppId, req.AssetId)

	return &types.QueryCollectorDataByAppAndAssetResponse{
		CollectorData: collectorMap.Collector,
	}, nil
}

func (q QueryServer) QueryAuctionMappingForAppAndAsset(c context.Context, req *types.QueryAuctionMappingForAppAndAssetRequest) (*types.QueryAuctionMappingForAppAndAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
	}
	auctionData, _ := q.GetAuctionMappingForApp(ctx, req.AppId, req.AssetId)

	return &types.QueryAuctionMappingForAppAndAssetResponse{
		AssetIdToAuctionLookupTable: auctionData,
	}, nil
}

func (q QueryServer) QueryNetFeeCollectedForAppAndAsset(c context.Context, req *types.QueryNetFeeCollectedForAppAndAssetRequest) (*types.QueryNetFeeCollectedForAppAndAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	ctx := sdk.UnwrapSDKContext(c)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
	}
	feeData, _ := q.GetNetFeeCollectedData(ctx, req.AppId, req.AssetId)

	return &types.QueryNetFeeCollectedForAppAndAssetResponse{
		AssetIdToFeeCollected: feeData,
	}, nil
}
