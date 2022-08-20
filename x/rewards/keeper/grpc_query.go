package keeper

import (
	"context"
	"time"

	"github.com/comdex-official/comdex/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = &Keeper{}

// Params queries the parameters of the incentives module.
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var params types.Params
	k.paramstore.GetParamSet(ctx, &params)
	return &types.QueryParamsResponse{Params: params}, nil
}

// QueryEpochInfoByDuration queries the epoch info for the given duration of seconds.
func (k Keeper) QueryEpochInfoByDuration(c context.Context, req *types.QueryEpochInfoByDurationRequest) (*types.QueryEpochInfoByDurationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	item, found := k.GetEpochInfoByDuration(ctx, time.Second*time.Duration(req.DurationSeconds))
	if !found {
		return nil, status.Errorf(codes.NotFound, "epoch does not exist for given duration %ds", req.DurationSeconds)
	}

	return &types.QueryEpochInfoByDurationResponse{
		Epoch: item,
	}, nil
}

// QueryAllEpochsInfo queries all the epochs available.
func (k Keeper) QueryAllEpochsInfo(c context.Context, req *types.QueryAllEpochsInfoRequest) (*types.QueryAllEpochsInfoResponse, error) {
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

// QueryAllGauges queries all the gauges available.
func (k Keeper) QueryAllGauges(c context.Context, req *types.QueryAllGaugesRequest) (*types.QueryAllGaugesResponse, error) {
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

// QueryGaugeByID queries a gauge by specific ID.
func (k Keeper) QueryGaugeByID(c context.Context, req *types.QueryGaugeByIdRequest) (*types.QueryGaugeByIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	item, found := k.GetGaugeByID(ctx, req.GaugeId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "gauge does not exist for given id %d", req.GaugeId)
	}

	return &types.QueryGaugeByIdResponse{
		Gauge: item,
	}, nil
}

// QueryGaugeByDuration queries gauges for the given duration.
func (k Keeper) QueryGaugeByDuration(c context.Context, req *types.QueryGaugesByDurationRequest) (*types.QueryGaugeByDurationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	gaugesIDsByTriggerDuration, found := k.GetGaugeIdsByTriggerDuration(ctx, time.Second*time.Duration(req.DurationSeconds))
	if !found {
		return nil, status.Errorf(codes.NotFound, "no gauges for given duration %ds", req.DurationSeconds)
	}

	var gauges []types.Gauge

	for _, gaugeID := range gaugesIDsByTriggerDuration.GaugeIds {
		gauge, found := k.GetGaugeByID(ctx, gaugeID)
		if !found {
			continue
		}
		gauges = append(gauges, gauge)
	}

	return &types.QueryGaugeByDurationResponse{
		Gauges: gauges,
	}, nil
}

func (k Keeper) QueryRewards(c context.Context, req *types.QueryRewardsRequest) (*types.QueryRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.InternalRewards
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(k.Store(ctx), types.RewardsKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.InternalRewards
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

	return &types.QueryRewardsResponse{
		Rewards:    items,
		Pagination: pagination,
	}, nil
}

func (k Keeper) QueryReward(c context.Context, req *types.QueryRewardRequest) (*types.QueryRewardResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := k.GetRewardByApp(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "asset does not exist for id %d", req.Id)
	}

	return &types.QueryRewardResponse{
		Reward: item,
	}, nil
}

func (k Keeper) QueryExternalRewardsLockers(c context.Context, req *types.QueryExternalRewardsLockersRequest) (*types.QueryExternalRewardsLockersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	items := k.GetExternalRewardsLockers(ctx)

	return &types.QueryExternalRewardsLockersResponse{
		LockerExternalRewards: items,
	}, nil
}

func (k Keeper) QueryExternalRewardVaults(c context.Context, req *types.QueryExternalRewardVaultsRequest) (*types.QueryExternalRewardVaultsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	items := k.GetExternalRewardVaults(ctx)

	return &types.QueryExternalRewardVaultsResponse{
		VaultExternalRewards: items,
	}, nil
}

func (k Keeper) QueryWhitelistedAppIdsVault(c context.Context, req *types.QueryWhitelistedAppIdsVaultRequest) (*types.QueryWhitelistedAppIdsVaultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	items := k.GetAppIDs(ctx)

	return &types.QueryWhitelistedAppIdsVaultResponse{
		WhitelistedAppIdsVault: items,
	}, nil
}
