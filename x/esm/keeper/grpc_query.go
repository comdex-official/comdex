package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/esm/types"
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

func (q queryServer) QueryESMTriggerParams(c context.Context, req *types.QueryESMTriggerParamsRequest) (*types.QueryESMTriggerParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetESMTriggerParams(ctx, req.Id)
	if !found {
		return &types.QueryESMTriggerParamsResponse{}, nil
	}

	return &types.QueryESMTriggerParamsResponse{
		EsmTriggerParams: item,
	}, nil
}

func (q queryServer) QueryESMStatus(c context.Context, req *types.QueryESMStatusRequest) (*types.QueryESMStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetESMStatus(ctx, req.Id)
	if !found {
		return &types.QueryESMStatusResponse{}, nil
	}

	return &types.QueryESMStatusResponse{
		EsmStatus: item,
	}, nil
}

func (q queryServer) QueryCurrentDepositStats(c context.Context, req *types.QueryCurrentDepositStatsRequest) (*types.QueryCurrentDepositStatsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetCurrentDepositStats(ctx, req.Id)
	if !found {
		return &types.QueryCurrentDepositStatsResponse{}, nil
	}

	return &types.QueryCurrentDepositStatsResponse{
		CurrentDepositStats: item,
	}, nil
}

func (q queryServer) QueryUsersDepositMapping(c context.Context, req *types.QueryUsersDepositMappingRequest) (*types.QueryUsersDepositMappingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetUserDepositByApp(ctx, req.Depositor, req.Id)
	if !found {
		return &types.QueryUsersDepositMappingResponse{}, nil
	}

	return &types.QueryUsersDepositMappingResponse{
		UsersDepositMapping: item,
	}, nil
}
