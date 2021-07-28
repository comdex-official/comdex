package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/cdp/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) GetCDP(ctx context.Context, request *types.GetCDPRequest) (*types.GetCDPResponse, error) {
	//TODO
	return nil, nil
}

func (k Keeper) GetCDPs(ctx context.Context, request *types.GetCDPsRequest) (*types.GetCDPsResponse, error) {
	//TODO
	return nil, nil
}

func (k Keeper) GetCDPDeposits(ctx context.Context, request *types.GetCDPDepositsRequest) (*types.GetCDPDepositsResponse, error) {
	//TODO
	return nil, nil
}

func (k Keeper) GetCDPsByCollateralType(ctx context.Context, request *types.GetCDPsByCollateralTypeRequest) (*types.GetCDPsByCollateralTypeResponse, error) {
	//TODO
	return nil, nil
}

func (k Keeper) GetCDPsByCollateralizationRatio(ctx context.Context, request *types.GetCDPsByCollateralizationRatioRequest) (*types.GetCDPsByCollateralizationRatioResponse, error) {
	//TODO
	return nil, nil
}

func (k Keeper) GetParams(ctx context.Context, request *types.GetParamsRequest) (*types.GetParamsResponse, error) {
	//TODO
	return nil, nil
}
