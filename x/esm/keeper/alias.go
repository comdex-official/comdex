package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) GetApp(ctx sdk.Context, id uint64) (app assettypes.AppData, found bool) {
	return k.asset.GetApp(ctx, id)
}

func (k *Keeper) GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool) {
	return k.asset.GetAsset(ctx, id)
}

func (k *Keeper) GetApps(ctx sdk.Context) (apps []assettypes.AppData, found bool) {
	return k.asset.GetApps(ctx)
}

func (k *Keeper) GetAssetsForOracle(ctx sdk.Context) (assets []assettypes.Asset) {
	return k.asset.GetAssetsForOracle(ctx)
}

func (k *Keeper) GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool) {
	return k.market.GetPriceForAsset(ctx, id)
}

func (k *Keeper) BurnTokensForApp(ctx sdk.Context, appMappingID uint64, assetID uint64, amount sdk.Int) error {
	return k.tokenmint.BurnTokensForApp(ctx, appMappingID, assetID, amount)
}