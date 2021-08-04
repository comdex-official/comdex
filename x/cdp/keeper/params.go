package keeper

import (
	"fmt"
	"github.com/comdex-official/comdex/x/cdp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) GetLiquidationRatio(ctx sdk.Context, collateralType string) sdk.Dec {
	collateralParam, found := k.GetCollateralParam(ctx, collateralType)
	if !found {
		panic(fmt.Sprintf("collateral not found: %s", collateralType))
	}
	return collateralParam.LiquidationRatio
}

func (k Keeper) GetCollateralParam(ctx sdk.Context, collateralType string) (types.CollateralParam, bool) {
	params := k.GetParams(ctx)
	for _, collateralParam := range params.CollateralParams {
		if collateralParam.Type == collateralType {
			return collateralParam, true
		}
	}
	return types.CollateralParam{}, false
}
