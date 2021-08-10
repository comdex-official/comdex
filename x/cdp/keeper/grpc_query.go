package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/cdp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServiceServer = Keeper{}

func (k Keeper) QueryCDP(context context.Context, request *types.QueryCDPRequest) (*types.QueryCDPResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	ownerAddrs, err := sdk.AccAddressFromBech32(request.Owner)
	if err != nil {
		return &types.QueryCDPResponse{}, err
	}
	cdp, found := k.GetCDPByOwnerAndCollateralType(ctx, ownerAddrs, request.CollateralType)
	if !found {
		return nil, status.Error(codes.NotFound, "cdp not found")
	}

	return &types.QueryCDPResponse{Cdp: cdp}, nil
}

func (k Keeper) QueryCDPs(ctx context.Context, request *types.QueryCDPsRequest) (*types.QueryCDPsResponse, error) {
	//TODO
	return nil, nil
}

func (k Keeper) QueryCDPById(ctx context.Context, request *types.QueryCDPByIdRequest) (*types.QueryCDPByIdResponse, error) {
	//TODO
	return nil, nil
}

func (k Keeper) QueryParams(context context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}
