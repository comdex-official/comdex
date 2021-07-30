package keeper

import (
	"fmt"
	"github.com/comdex-official/comdex/x/cdp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetParams(ctx sdk.Context) types.Params  {
	k.paramSpace.Get()
}

func (k Keeper) getLiquidationRatio(ctx sdk.Context, collateralType string) sdk.Dec {
	collateralParam, found := k.GetCollateral(ctx, collateralType)
	if !found{
		panic(fmt.Sprintf("collateral not found: %s", collateralType))
	}
	return collateralParam.

}

func (k Keeper) GetCollateral(ctx sdk.Context, collateralType string) (types.CollateralParam, bool) {
	params:= k.GetParams(ctx)
	for _,collateralParam:= range params.C

}