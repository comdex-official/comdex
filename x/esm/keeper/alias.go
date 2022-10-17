package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
)

func (k Keeper) GetApp(ctx sdk.Context, id uint64) (app assettypes.AppData, found bool) {
	return k.asset.GetApp(ctx, id)
}

func (k Keeper) GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool) {
	return k.asset.GetAsset(ctx, id)
}

func (k Keeper) GetApps(ctx sdk.Context) (apps []assettypes.AppData, found bool) {
	return k.asset.GetApps(ctx)
}

func (k Keeper) GetAssets(ctx sdk.Context) (assets []assettypes.Asset) {
	return k.asset.GetAssets(ctx)
}

func (k Keeper) GetTwa(ctx sdk.Context, id uint64) (twa markettypes.TimeWeightedAverage, found bool) {
	return k.market.GetTwa(ctx, id)
}

func (k Keeper) GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool) {
	return k.asset.GetPair(ctx, id)
}

func (k Keeper) GetPairsVault(ctx sdk.Context, id uint64) (pairs assettypes.ExtendedPairVault, found bool) {
	return k.asset.GetPairsVault(ctx, id)
}

func (k Keeper) DeleteVault(ctx sdk.Context, id uint64) {
	k.vault.DeleteVault(ctx, id)
}

func (k Keeper) DeleteAddressFromAppExtendedPairVaultMapping(ctx sdk.Context, extendedPairID uint64, userVaultID uint64, appMappingID uint64) {
	k.vault.DeleteAddressFromAppExtendedPairVaultMapping(ctx, extendedPairID, userVaultID, appMappingID)
}

func (k Keeper) BurnTokensForApp(ctx sdk.Context, appMappingID uint64, assetID uint64, amount sdk.Int) error {
	return k.tokenmint.BurnTokensForApp(ctx, appMappingID, assetID, amount)
}

func (k Keeper) GetVaults(ctx sdk.Context) (vaults []vaulttypes.Vault) {
	return k.vault.GetVaults(ctx)
}

func (k Keeper) GetStableMintVaults(ctx sdk.Context) (stableVaults []vaulttypes.StableMintVault) {
	return k.vault.GetStableMintVaults(ctx)
}

func (k Keeper) GetAssetForDenom(ctx sdk.Context, denom string) (asset assettypes.Asset, found bool) {
	return k.asset.GetAssetForDenom(ctx, denom)
}

func (k Keeper) GetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64) (netFeeData collectortypes.AppAssetIdToFeeCollectedData, found bool) {
	return k.collector.GetNetFeeCollectedData(ctx, appID, assetID)
}

func (k Keeper) GetAppNetFeeCollectedData(ctx sdk.Context, appID uint64) (netFeeData []collectortypes.AppAssetIdToFeeCollectedData, found bool) {
	return k.collector.GetAppNetFeeCollectedData(ctx, appID)
}
