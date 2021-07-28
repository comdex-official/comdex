package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/cdp/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) CDP(ctx context.Context, request *types.QueryCDPRequest) (*types.QueryCDPResponse, error) {
	//TODO
	return nil, nil
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

func (k Keeper) Params(ctx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	//TODO
	return nil, nil
}
