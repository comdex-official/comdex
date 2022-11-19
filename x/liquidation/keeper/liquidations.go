package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/petrichormoney/petri/x/liquidation/types"
)

// SetLiquidationOffsetHolder stores the LiquidationOffsetHolder.
func (k Keeper) SetLiquidationOffsetHolder(ctx sdk.Context, liquidatonPrefix string, liquidationOffsetHolder types.LiquidationOffsetHolder) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalLiquidationOffsetHolder(k.cdc, liquidationOffsetHolder)
	store.Set(
		types.GetLiquidationOffsetHolderKey(liquidationOffsetHolder.AppId, liquidatonPrefix),
		bz,
	)
}

// GetLiquidationOffsetHolder returns liquidationOffsetHolder object for the given app id, pool id and farmer.
func (k Keeper) GetLiquidationOffsetHolder(ctx sdk.Context, appID uint64, liquidatonPrefix string) (liquidationOffsetHolder types.LiquidationOffsetHolder, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLiquidationOffsetHolderKey(appID, liquidatonPrefix))
	if bz == nil {
		return
	}
	liquidationOffsetHolder = types.MustUnmarshalLiquidationOffsetHolder(k.cdc, bz)
	return liquidationOffsetHolder, true
}
