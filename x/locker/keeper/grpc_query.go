package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/locker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServiceServer = (*queryServer)(nil)

type queryServer struct {
	Keeper
}

func NewQueryServiceServer(k Keeper) types.QueryServiceServer {
	return &queryServer{
		Keeper: k,
	}
}

func (q *queryServer) QueryParams(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	var (
		ctx    = sdk.UnwrapSDKContext(c)
		params = q.GetParams(ctx)
	)

	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

func (q *queryServer) QueryLockerInfo(c context.Context, req *types.QueryLockerInfoRequest) (*types.QueryLockerInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)
	item, found := q.GetLocker(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "locker-info does not exist for id %d", req.Id)
	}

	return &types.QueryLockerInfoResponse{
		LockerInfo: item,
	}, nil
}

func (q *queryServer) QueryLockersByProductToAssetID(c context.Context, request *types.QueryLockersByProductToAssetIDRequest) (*types.QueryLockersByProductToAssetIDResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
	//ctx = sdk.UnwrapSDKContext(c)
	)

	return &types.QueryLockersByProductToAssetIDResponse{
		LockerInfo: nil,
	}, nil
}

func (q *queryServer) QueryLockerInfoByProductID(c context.Context, request *types.QueryLockerInfoByProductIDRequest) (*types.QueryLockerInfoByProductIDResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
	//ctx = sdk.UnwrapSDKContext(c)
	)

	return &types.QueryLockerInfoByProductIDResponse{
		LockerInfo: nil,
	}, nil

}

func (q *queryServer) QueryTotalDepositByAssetID(c context.Context, request *types.QueryTotalDepositByAssetIDRequest) (*types.QueryTotalDepositByAssetIDResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
	//ctx = sdk.UnwrapSDKContext(c)
	)

	return &types.QueryTotalDepositByAssetIDResponse{
		TotalDeposit: 0,
	}, nil
}

func (q *queryServer) QueryTotalDepositByProductAssetID(c context.Context, request *types.QueryTotalDepositByProductAssetIDRequest) (*types.QueryTotalDepositByProductAssetIDResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
	//ctx = sdk.UnwrapSDKContext(c)
	)

	return &types.QueryTotalDepositByProductAssetIDResponse{
		TotalDeposit: 0,
	}, nil
}

func (q *queryServer) QueryOwnerLockerByProductID(ctx context.Context, request *types.QueryOwnerLockerByProductIDRequest) (*types.QueryOwnerLockerByProductIDResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
	//ctx = sdk.UnwrapSDKContext(c)
	)

	return &types.QueryOwnerLockerByProductIDResponse{
		LockerInfo: nil,
	}, nil
}

func (q *queryServer) QueryOwnerLockerByProductToAssetID(ctx context.Context, request *types.QueryOwnerLockerByProductToAssetIDRequest) (*types.QueryOwnerLockerByProductToAssetIDResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
	//ctx = sdk.UnwrapSDKContext(c)
	)

	return &types.QueryOwnerLockerByProductToAssetIDResponse{
		LockerInfo: nil,
	}, nil
}

func (q *queryServer) QueryLockerCountByProductID(ctx context.Context, request *types.QueryLockerCountByProductIDRequest) (*types.QueryLockerCountByProductIDResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
	//ctx = sdk.UnwrapSDKContext(c)
	)

	return &types.QueryLockerCountByProductIDResponse{
		TotalCount: 0,
	}, nil

}

func (q *queryServer) QueryLockerCountByProductToAssetID(ctx context.Context, request *types.QueryLockerCountByProductToAssetIDRequest) (*types.QueryLockerCountByProductToAssetIDResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
	//ctx = sdk.UnwrapSDKContext(c)
	)

	return &types.QueryLockerCountByProductToAssetIDResponse{
		TotalCount: 0,
	}, nil
}

func (q *queryServer) QueryWhiteListedAssetIDsByProductID(ctx context.Context, request *types.QueryWhiteListedAssetIDsByProductIDRequest) (*types.QueryWhiteListedAssetIDsByProductIDResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
	//ctx = sdk.UnwrapSDKContext(c)
	)

	return &types.QueryWhiteListedAssetIDsByProductIDResponse{
		AssetIds: nil,
	}, nil
}

func (q *queryServer) QueryWhiteListedAssetByAllProduct(ctx context.Context, request *types.QueryWhiteListedAssetByAllProductRequest) (*types.QueryWhiteListedAssetByAllProductResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
	//ctx = sdk.UnwrapSDKContext(c)
	)

	return &types.QueryWhiteListedAssetByAllProductResponse{
		Asset: nil,
	}, nil
}
