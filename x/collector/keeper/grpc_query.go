package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ types.QueryServer = (*QueryServer)(nil)
)

type QueryServer struct {
	Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &QueryServer{
		Keeper: k,
	}
}

func (q *QueryServer) QueryCollectorLookupByApp(c context.Context, req *types.QueryCollectorLookupByAppRequest) (*types.QueryCollectorLookupByAppResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
	}

	collectorLookupData, found := q.GetCollectorLookupTable(ctx, req.AppId)
	if !found {
		return &types.QueryCollectorLookupByAppResponse{}, nil
	}

	return &types.QueryCollectorLookupByAppResponse{
		CollectorLookup: collectorLookupData.AssetRateInfo,
	}, nil
}

func (q *QueryServer) QueryCollectorLookupByAppAndAsset(c context.Context, req *types.QueryCollectorLookupByAppAndAssetRequest) (*types.QueryCollectorLookupByAppAndAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
	}

	collectorLookupData, found := q.GetCollectorLookupByAsset(ctx, req.AppId, req.AssetId)
	if !found {
		return &types.QueryCollectorLookupByAppAndAssetResponse{}, nil
	}

	return &types.QueryCollectorLookupByAppAndAssetResponse{
		CollectorLookup: collectorLookupData,
	}, nil
}

func (q *QueryServer) QueryCollectorDataByAppAndAsset(c context.Context, req *types.QueryCollectorDataByAppAndAssetRequest) (*types.QueryCollectorDataByAppAndAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx           = sdk.UnwrapSDKContext(c)
		collectorData types.CollectorData
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
	}
	collectorMap, _ := q.GetAppidToAssetCollectorMapping(ctx, req.AppId)

	for _, data := range collectorMap.AssetCollector {
		if data.AssetId == req.AssetId {
			collectorData.CollectedClosingFee = data.Collector.CollectedClosingFee
			collectorData.CollectedOpeningFee = data.Collector.CollectedOpeningFee
			collectorData.CollectedStabilityFee = data.Collector.CollectedStabilityFee
			collectorData.LiquidationRewardsCollected = data.Collector.LiquidationRewardsCollected
		}
	}

	return &types.QueryCollectorDataByAppAndAssetResponse{
		CollectorData: collectorData,
	}, nil
}

func (q *QueryServer) QueryAuctionMappingForAppAndAsset(c context.Context, req *types.QueryAuctionMappingForAppAndAssetRequest) (*types.QueryAuctionMappingForAppAndAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx                  = sdk.UnwrapSDKContext(c)
		assetToAuctionLookup types.AssetIdToAuctionLookupTable
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
	}
	auctionData, _ := q.GetAuctionMappingForApp(ctx, req.AppId)
	for _, data := range auctionData.AssetIdToAuctionLookup {
		if data.AssetId == req.AssetId {
			assetToAuctionLookup.AssetId = data.AssetId
			assetToAuctionLookup.AssetOutOraclePrice = data.AssetOutOraclePrice
			assetToAuctionLookup.AssetOutPrice = data.AssetOutPrice
			assetToAuctionLookup.IsAuctionActive = data.IsAuctionActive
			assetToAuctionLookup.IsDebtAuction = data.IsDebtAuction
			assetToAuctionLookup.IsSurplusAuction = data.IsSurplusAuction
		}
	}

	return &types.QueryAuctionMappingForAppAndAssetResponse{
		AssetIdToAuctionLookupTable: assetToAuctionLookup,
	}, nil
}

func (q *QueryServer) QueryNetFeeCollectedForAppAndAsset(c context.Context, req *types.QueryNetFeeCollectedForAppAndAssetRequest) (*types.QueryNetFeeCollectedForAppAndAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx                   = sdk.UnwrapSDKContext(c)
		assetIDToFeeCollected types.AssetIdToFeeCollected
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
	}
	fee, _ := q.GetNetFeeCollectedData(ctx, req.AppId)
	for _, data := range fee.AssetIdToFeeCollected {
		if data.AssetId == req.AssetId {
			assetIDToFeeCollected.AssetId = data.AssetId
			assetIDToFeeCollected.NetFeesCollected = data.NetFeesCollected
		}
	}

	return &types.QueryNetFeeCollectedForAppAndAssetResponse{
		AssetIdToFeeCollected: assetIDToFeeCollected,
	}, nil
}
