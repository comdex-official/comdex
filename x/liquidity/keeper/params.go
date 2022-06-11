package keeper

import (
	"github.com/comdex-official/comdex/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetParams returns the parameters for the liquidity module.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

// SetParams sets the parameters for the liquidity module.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetParams returns the parameters for the liquidity module.
func (k Keeper) GetGenericParams(ctx sdk.Context, appID uint64) (types.GenericParams, error) {
	genericParams, found := k.GetGenericLiquidityParams(ctx, appID)
	if !found {
		if ctx.IsCheckTx() {
			return types.GenericParams{}, status.Errorf(codes.NotFound, "params for app-id %d doesn't exist", appID)
		}
		_, found := k.assetKeeper.GetApp(ctx, appID)
		if !found {
			return types.GenericParams{}, sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", appID)
		}
		genericParams = types.DefaultGenericParams(appID)
		k.SetGenericParams(ctx, genericParams)
	}
	return genericParams, nil
}

// SetParams sets the parameters for the liquidity module.
func (k Keeper) SetGenericParams(ctx sdk.Context, genericParams types.GenericParams) {
	k.SetGenericLiquidityParams(ctx, genericParams)
}
