package keeper

import (
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/x/params/types"
)

// LiquidationBatchSize defines the batch size for each ABCI interaction
func (k Keeper) LiquidationBatchSize(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyLiquidationBatchSize, &res)
	return
}

// GetParams returns the parameters for the liquidation module.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	return types.NewParams(
		k.LiquidationBatchSize(ctx),
	)
}

// SetParams set the params.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
