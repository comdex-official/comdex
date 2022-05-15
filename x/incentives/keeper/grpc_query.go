package keeper

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/comdex-official/comdex/x/incentives/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper.
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// Params queries the parameters of the incentives module.
func (k Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var params types.Params
	k.Keeper.paramSpace.GetParamSet(ctx, &params)
	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Querier) QueryEpochInfoByDuration(c context.Context, req *types.QueryEpochInfoByDurationRequest) (*types.QueryEpochInfoByDurationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	item, found := k.Keeper.GetEpochInfoByDuration(ctx, time.Second*time.Duration(req.DurationSeconds))
	if !found {
		return nil, status.Errorf(codes.NotFound, "epoch does not exist for given duration %d", req.DurationSeconds)
	}

	return &types.QueryEpochInfoByDurationResponse{
		Epoch: item,
	}, nil
}

func (k Querier) QueryAllEpochsInfo(c context.Context, req *types.QueryAllEpochsInfoRequest) (*types.QueryAllEpochsInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.EpochInfo
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(k.Store(ctx), types.EpochInfoByDurationKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.EpochInfo
			if err := k.cdc.Unmarshal(value, &item); err != nil {
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

	return &types.QueryAllEpochsInfoResponse{
		Epochs:     items,
		Pagination: pagination,
	}, nil
}

func (k Querier) QueryAllGauges(c context.Context, req *types.QueryAllGaugesRequest) (*types.QueryAllGaugesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.Gauge
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(k.Store(ctx), types.GaugeKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.Gauge
			if err := k.cdc.Unmarshal(value, &item); err != nil {
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

	return &types.QueryAllGaugesResponse{
		Gauges:     items,
		Pagination: pagination,
	}, nil
}
