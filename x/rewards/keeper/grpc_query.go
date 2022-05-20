package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = &Keeper{}

func (q *Keeper) QueryRewards(c context.Context, req *types.QueryRewardsRequest) (*types.QueryRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.InternalRewards
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.RewardsKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.InternalRewards
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

	return &types.QueryRewardsResponse{
		Rewards:    items,
		Pagination: pagination,
	}, nil
}

func (q *Keeper) QueryReward(c context.Context, req *types.QueryRewardRequest) (*types.QueryRewardResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetReward(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "asset does not exist for id %d", req.Id)
	}

	return &types.QueryRewardResponse{
		Reward: item,
	}, nil
}
