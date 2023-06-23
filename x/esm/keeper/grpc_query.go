package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/comdex-official/comdex/x/esm/types"
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

func (q QueryServer) QueryESMTriggerParams(c context.Context, req *types.QueryESMTriggerParamsRequest) (*types.QueryESMTriggerParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	item, found := q.GetESMTriggerParams(ctx, req.Id)
	if !found {
		return &types.QueryESMTriggerParamsResponse{}, nil
	}

	return &types.QueryESMTriggerParamsResponse{
		EsmTriggerParams: item,
	}, nil
}

func (q QueryServer) QueryESMStatus(c context.Context, req *types.QueryESMStatusRequest) (*types.QueryESMStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	item, found := q.GetESMStatus(ctx, req.Id)
	if !found {
		return &types.QueryESMStatusResponse{}, nil
	}

	return &types.QueryESMStatusResponse{
		EsmStatus: item,
	}, nil
}

func (q QueryServer) QueryCurrentDepositStats(c context.Context, req *types.QueryCurrentDepositStatsRequest) (*types.QueryCurrentDepositStatsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	item, found := q.GetCurrentDepositStats(ctx, req.Id)
	if !found {
		return &types.QueryCurrentDepositStatsResponse{}, nil
	}

	return &types.QueryCurrentDepositStatsResponse{
		CurrentDepositStats: item,
	}, nil
}

func (q QueryServer) QueryUsersDepositMapping(c context.Context, req *types.QueryUsersDepositMappingRequest) (*types.QueryUsersDepositMappingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	item, found := q.GetUserDepositByApp(ctx, req.Depositor, req.Id)
	if !found {
		return &types.QueryUsersDepositMappingResponse{}, nil
	}

	return &types.QueryUsersDepositMappingResponse{
		UsersDepositMapping: item,
	}, nil
}

func (q QueryServer) QueryDataAfterCoolOff(c context.Context, req *types.QueryDataAfterCoolOffRequest) (*types.QueryDataAfterCoolOffResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	item, found := q.GetDataAfterCoolOff(ctx, req.Id)
	if !found {
		return &types.QueryDataAfterCoolOffResponse{}, nil
	}

	return &types.QueryDataAfterCoolOffResponse{
		DataAfterCoolOff: item,
	}, nil
}

func (q QueryServer) QuerySnapshotPrice(c context.Context, req *types.QuerySnapshotPriceRequest) (*types.QuerySnapshotPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	price, found := q.GetSnapshotOfPrices(ctx, req.AppId, req.AssetId)
	if !found {
		return nil, types.ErrPriceNotFound
	}

	return &types.QuerySnapshotPriceResponse{
		Price: price,
	}, nil
}

func (q QueryServer) QueryAssetDataAfterCoolOff(c context.Context, req *types.QueryAssetDataAfterCoolOffRequest) (*types.QueryAssetDataAfterCoolOffResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	item := q.GetAllAssetToAmount(ctx, req.AppId)
	if item == nil {
		return &types.QueryAssetDataAfterCoolOffResponse{}, nil
	}

	return &types.QueryAssetDataAfterCoolOffResponse{
		AssetToAmount: item,
	}, nil
}
