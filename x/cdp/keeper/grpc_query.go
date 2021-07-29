package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/cdp/types"
)

var _ types.QueryServiceServer = Keeper{}

func (k Keeper) QueryCDP(ctx context.Context, request *types.QueryCDPRequest) (*types.QueryCDPResponse, error) {
	//TODO
	return nil, nil
}

func (k Keeper) QueryCDPs(ctx context.Context, request *types.QueryCDPsRequest) (*types.QueryCDPsResponse, error) {
	//TODO
	return nil, nil
}

func (k Keeper) QueryCDPDeposits(ctx context.Context, request *types.QueryCDPDepositsRequest) (*types.QueryCDPDepositsResponse, error) {
	//TODO
	return nil, nil
}

func (k Keeper) QueryCDPsByCollateralType(ctx context.Context, request *types.QueryCDPsByCollateralTypeRequest) (*types.QueryCDPsByCollateralTypeResponse, error) {
	//TODO
	return nil, nil
}

func (k Keeper) QueryCDPsByCollateralizationRatio(ctx context.Context, request *types.QueryCDPsByCollateralizationRatioRequest) (*types.QueryCDPsByCollateralizationRatioResponse, error) {
	//TODO
	return nil, nil
}

func (k Keeper) QueryParams(ctx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	//TODO
	return nil, nil
}
