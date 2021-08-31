package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/cdp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
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
		return &types.QueryCDPResponse{}, status.Error(codes.NotFound, "cdp not found")
	}

	return &types.QueryCDPResponse{Cdp: cdp}, nil
}

func (k Keeper) QueryCDPs(context context.Context, request *types.QueryCDPsRequest) (*types.QueryCDPsResponse, error) {

	var (
		cdps       []types.CDP
		pagination *query.PageResponse
		ctx        = sdk.UnwrapSDKContext(context)
	)

	ownerAddrs, err := sdk.AccAddressFromBech32(request.Owner)
	if err != nil {
		return &types.QueryCDPsResponse{}, err
	}

	ownerCDPList, found := k.GetOwnerCDPList(ctx, ownerAddrs)
	if !found {
		return &types.QueryCDPsResponse{}, nil
	}

	for _, ownedCdp := range ownerCDPList.OwnedCDPs {
		cdp, found := k.GetCDP(ctx, ownedCdp.Id)
		if found {
			cdps = append(cdps, cdp)
		}
	}

	return &types.QueryCDPsResponse{
		Cdps:       cdps,
		Pagination: pagination,
	}, nil

}

func (k Keeper) QueryCDPById(context context.Context, request *types.QueryCDPByIdRequest) (*types.QueryCDPByIdResponse, error) {

	ctx := sdk.UnwrapSDKContext(context)
	cdp, found := k.GetCDP(ctx, request.Id)

	if !found {
		return &types.QueryCDPByIdResponse{}, status.Error(codes.NotFound, "cdp not found")
	}

	return &types.QueryCDPByIdResponse{
		Cdp: cdp,
	}, nil
}

func (k Keeper) QueryParams(context context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	if request == nil {
		return &types.QueryParamsResponse{}, status.Error(codes.InvalidArgument, "empty request")
	}

	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}
