package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
)


func (k *Keeper) GetApps(ctx sdk.Context) (assettypes.AppMapping, bool) {
	return k.asset.GetApps(ctx)
}

func (k *Keeper) HasAssetForDenom(ctx sdk.Context, id string) (bool) {
	return k.asset.HasAssetForDenom(ctx, id)
}

func (k *Keeper) HasAsset(ctx sdk.Context, id uint64) (bool) {
	return k.asset.HasAsset(ctx, id)
}