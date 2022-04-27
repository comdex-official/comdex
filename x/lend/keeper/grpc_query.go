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

func (q *queryServer) QueryAssets(c context.Context, req *types.QueryAssetsRequest) (*types.QueryAssetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		items []types.Asset
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.WhitelistedAssetKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.Asset
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

	return &types.QueryAssetsResponse{
		Assets:     items,
		Pagination: pagination,
	}, nil
}

func (q *queryServer) QueryAsset(c context.Context, req *types.QueryAssetRequest) (*types.QueryAssetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetAsset(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "asset does not exist for id %d", req.Id)
	}

	return &types.QueryAssetResponse{
		Asset: item,
	}, nil
}

func (q *queryServer) QueryAssetPerDenom(c context.Context, req *types.QueryAssetPerDenomRequest) (*types.QueryAssetPerDenomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetAssetForDenom(ctx, req.Denom)
	if !found {
		return nil, status.Errorf(codes.NotFound, "asset does not exist for id %d", req.Denom)
	}

	return &types.QueryAssetPerDenomResponse{
		Asset: item,
	}, nil
}
