package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/comdex-official/comdex/x/locking/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper.
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// Params queries the parameters of the liquidity module.
func (k Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var params types.Params
	k.Keeper.paramSpace.GetParamSet(ctx, &params)
	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Querier) QueryLockByID(c context.Context, req *types.QueryLockByIdRequest) (*types.QueryLockByIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := k.GetLockByID(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "lock does not exist for id %d", req.Id)
	}

	return &types.QueryLockByIdResponse{
		Lock: item,
	}, nil
}

func (k Querier) QueryLocksByOwner(c context.Context, req *types.QueryLocksByOwnerRequest) (*types.QueryLocksByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	lockIdsByOwner, found := k.GetLockByOwner(ctx, req.Owner)
	if !found {
		return nil, status.Errorf(codes.NotFound, "lock does not exist for address %s", req.Owner)
	}

	locksByOwner := []types.Lock{}
	for _, lockID := range lockIdsByOwner.LockIds {
		lock, _ := k.GetLockByID(ctx, lockID)
		locksByOwner = append(locksByOwner, lock)
	}

	return &types.QueryLocksByOwnerResponse{
		Locks: locksByOwner,
	}, nil
}

func (k Querier) QueryAllLocks(c context.Context, req *types.QueryAllLocksRequest) (*types.QueryAllLocksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.Lock
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(k.Store(ctx), types.LockKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.Lock
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

	return &types.QueryAllLocksResponse{
		Locks:      items,
		Pagination: pagination,
	}, nil
}

func (k Querier) QueryUnlockingByID(c context.Context, req *types.QueryUnlockingByIdRequest) (*types.QueryUnlockingByIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := k.GetUnlockingByID(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "unlocking does not exist for id %d", req.Id)
	}

	return &types.QueryUnlockingByIdResponse{
		Unlocking: item,
	}, nil
}

func (k Querier) QueryUnlockingsByOwner(c context.Context, req *types.QueryUnlockingsByOwnerRequest) (*types.QueryUnlockingsByOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	unlockIdsByOwner, found := k.GetUnlockByOwner(ctx, req.Owner)
	if !found {
		return nil, status.Errorf(codes.NotFound, "unlocking does not exist for address %s", req.Owner)
	}

	unlockingsByOwner := []types.Unlocking{}
	for _, lockID := range unlockIdsByOwner.UnlockingIds {
		lock, _ := k.GetUnlockingByID(ctx, lockID)
		unlockingsByOwner = append(unlockingsByOwner, lock)
	}

	return &types.QueryUnlockingsByOwnerResponse{
		Unlockings: unlockingsByOwner,
	}, nil
}

func (k Querier) QueryAllUnlockings(c context.Context, req *types.QueryAllUnlockingsRequest) (*types.QueryAllUnlockingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.Unlocking
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(k.Store(ctx), types.UnlockKeyPrefix),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.Unlocking
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

	return &types.QueryAllUnlockingsResponse{
		Unlockings: items,
		Pagination: pagination,
	}, nil
}
