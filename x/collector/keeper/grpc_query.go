package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (q *queryServer) QueryCollectorLookupByProduct(c context.Context, req *types.QueryCollectorLookupByProductRequest) (*types.QueryCollectorLookupByProductResponse, error) {
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
		return nil, status.Errorf(codes.NotFound, "Lookup table does not exist for product id %d", req.AppId)
	}

	return &types.QueryCollectorLookupByProductResponse{
		CollectorLookup: collectorLookupData.AssetRateInfo,
	}, nil
}

func (q *queryServer) QueryCollectorLookupByProductAndAsset(c context.Context, req *types.QueryCollectorLookupByProductAndAssetRequest) (*types.QueryCollectorLookupByProductAndAssetResponse, error) {
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
		return nil, status.Errorf(codes.NotFound, "Lookup table does not exist for product id %d", req.AppId)
	}

	return &types.QueryCollectorLookupByProductAndAssetResponse{
		CollectorLookup: collectorLookupData,
	}, nil
}

func (q *queryServer) QueryCollectorDataByProductAndAsset(c context.Context, req *types.QueryCollectorDataByProductAndAssetRequest) (*types.QueryCollectorDataByProductAndAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx = sdk.UnwrapSDKContext(c)
		collectorData types.CollectorData
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
	}
	collectormap, _ := q.GetAppidToAssetCollectorMapping(ctx,req.AppId)

	for _,data := range collectormap.AssetCollector{

		if data.AssetId == req.AssetId{
			collectorData.CollectedClosingFee = data.Collector.CollectedClosingFee
			collectorData.CollectedOpeningFee = data.Collector.CollectedOpeningFee
			collectorData.CollectedStabilityFee = data.Collector.CollectedStabilityFee
			collectorData.LiquidationRewardsCollected = data.Collector.LiquidationRewardsCollected
		}
	}

	return &types.QueryCollectorDataByProductAndAssetResponse{
		CollectorData: collectorData,
	}, nil
}

func (q *queryServer) QueryAuctionMappingForAppAndAsset(c context.Context, req *types.QueryAuctionMappingForAppAndAssetRequest) (*types.QueryAuctionMappingForAppAndAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx = sdk.UnwrapSDKContext(c)
		assetToAuctionLookup types.AssetIdToAuctionLookupTable
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
	}
	auctionData,_ := q.GetAuctionMappingForApp(ctx,req.AppId)
	for _, data := range auctionData.AssetIdToAuctionLookup{
		if data.AssetId == req.AssetId{
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


func (q *queryServer) QueryNetFeeCollectedForAppAndAsset(c context.Context, req *types.QueryNetFeeCollectedForAppAndAssetRequest) (*types.QueryNetFeeCollectedForAppAndAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx = sdk.UnwrapSDKContext(c)
		assetIdToFeeCollected types.AssetIdToFeeCollected
	)
	_, found := q.GetApp(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
	}
	fee, _ := q.GetNetFeeCollectedData(ctx,req.AppId)
	for _,data := range fee.AssetIdToFeeCollected{
		if data.AssetId == req.AssetId{
			assetIdToFeeCollected.AssetId = data.AssetId
			assetIdToFeeCollected.NetFeesCollected = data.NetFeesCollected
		}
	}


	return &types.QueryNetFeeCollectedForAppAndAssetResponse{
		AssetIdToFeeCollected: assetIdToFeeCollected,
	}, nil
}

func (k Keeper) QueryCollectorLookupByProduct(ctx context.Context, request *types.QueryCollectorLookupByProductRequest) (*types.QueryCollectorLookupByProductResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) QueryCollectorLookupByProductAndAsset(ctx context.Context, request *types.QueryCollectorLookupByProductAndAssetRequest) (*types.QueryCollectorLookupByProductAndAssetResponse, error) {
	//TODO implement me
	panic("implement me")
}
