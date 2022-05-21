package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ types.QueryServiceServer = (*queryServer)(nil)
)

type queryServer struct {
	Keeper
}

func NewQueryServiceServer(k Keeper) types.QueryServiceServer {
	return &queryServer{
		Keeper: k,
	}
}

func (q *queryServer) QueryCollectorLookupByProduct(c context.Context, req *types.QueryCollectorLookupByProductRequest) (*types.QueryCollectorLookupByProductResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
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
		CollectorLookup: collectorLookupData.AssetrateInfo,
	}, nil
}

func (q *queryServer) QueryCollectorLookupByProductAndAsset(c context.Context, req *types.QueryCollectorLookupByProductAndAssetRequest) (*types.QueryCollectorLookupByProductAndAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var( 
		ctx   = sdk.UnwrapSDKContext(c)
		collectorData types.CollectorLookupTable
	)
		_, found := q.GetApp(ctx, req.AppId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "product does not exist for id %d", req.AppId)
		}

		collectorLookupData, found := q.GetCollectorLookupByAsset(ctx, req.AppId, req.AssetId)
		if !found {
			return nil, status.Errorf(codes.NotFound, "Lookup table does not exist for product id %d", req.AppId)
		}

		for _, data := range collectorLookupData.AssetrateInfo{
			if data.CollectorAssetId == req.AssetId {
				collectorData = *data
			}
		}
	

	return &types.QueryCollectorLookupByProductAndAssetResponse{
		CollectorLookup: &collectorData,
	}, nil
}
