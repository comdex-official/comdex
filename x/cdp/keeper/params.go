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

func (k Keeper) getLiquidationRatio(ctx sdk.Context, collateralType string) sdk.Dec {
	collateralParam, found := k.GetCollateral(ctx, collateralType)
	if !found {
		panic(fmt.Sprintf("collateral not found: %s", collateralType))
	}
	return collateralParam.LiquidationRatio

}

func (k Keeper) GetDebtParam(ctx sdk.Context, denom string) (types.DebtParam, bool) {
	dp := k.GetParams(ctx).DebtParam
	if dp.Denom == denom {
		return dp, true
	}
	return types.DebtParam{}, false
}

func (k Keeper) GetCollateralTypePrefix(ctx sdk.Context, collateralType string) (byte, bool)  {
	params:= k.GetParams(ctx)

	for _, cp := range params.CollateralParams{
		if cp.Type == collateralType{
			return cp.Prefix, true
		}
	}
	return 0x00, false
}


func (k Keeper) GetCollateral(ctx sdk.Context, collateralType string) (types.CollateralParam, bool) {
	params := k.GetParams(ctx)
	for _, collateralParam := range params.CollateralParams {
		if collateralParam.Type == collateralType {
			return collateralParam, true
		}
	}
	return types.CollateralParam{}, false
}
