package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/comdex-official/comdex/x/asset/types"
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

func (q *queryServer) QueryPairs(c context.Context, req *types.QueryPairsRequest) (*types.QueryPairsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items      []types.Pair
		pagination *query.PageResponse
		ctx        = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.PairKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.Pair
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
		Pairs:      items,
		Pagination: pagination,
	}, nil
}

func (q *queryServer) QueryPair(c context.Context, req *types.QueryPairRequest) (*types.QueryPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetPair(ctx, req.ID)
	if !found {
		return nil, status.Errorf(codes.NotFound, "pair does not exist for id %d", req.ID)
	}

	return &types.QueryPairResponse{
		Pair: item,
	}, nil
}
