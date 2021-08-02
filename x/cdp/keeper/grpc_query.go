package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/cdp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) CDP(context context.Context, request *types.QueryCDPRequest) (*types.QueryCDPResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	ownerAddrs, err := sdk.AccAddressFromBech32(request.Owner)
	if err != nil {
		return &types.QueryCDPResponse{}, err
	}
	cdp, found := k.GetCdpByOwnerAndCollateralType(ctx, ownerAddrs, request.CollateralType)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "cdp not found")
	}

	return &types.QueryCDPResponse{Cdp: cdp}, nil
}

func (k Keeper) CDPs(ctx context.Context, request *types.QueryCDPsRequest) (*types.QueryCDPsResponse, error) {
	//TODO
	return nil, nil
}

func (k Keeper) CDPDeposits(ctx context.Context, request *types.QueryCDPDepositsRequest) (*types.QueryCDPDepositsResponse, error) {
	//TODO
	return nil, nil
}

func (k Keeper) CDPsByCollateralType(ctx context.Context, request *types.QueryCDPsByCollateralTypeRequest) (*types.QueryCDPsByCollateralTypeResponse, error) {
	//TODO
	return nil, nil
}

func (k Keeper) CDPsByCollateralizationRatio(ctx context.Context, request *types.QueryCDPsByCollateralizationRatioRequest) (*types.QueryCDPsByCollateralizationRatioResponse, error) {
	//TODO
	return nil, nil
}

func (k Keeper) Params(context context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}
