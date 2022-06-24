package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	collecortypes "github.com/comdex-official/comdex/x/collector/types"
	"github.com/comdex-official/comdex/x/locker/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetLockerProductAssetMapping(ctx sdk.Context, appMappingID uint64) (lockerProductMapping types.LockerProductAssetMapping, found bool) {
	return k.locker.GetLockerProductAssetMapping(ctx, appMappingID)
}

func (k Keeper) GetAppidToAssetCollectorMapping(ctx sdk.Context, appID uint64) (appAssetCollectorData collecortypes.AppIdToAssetCollectorMapping, found bool) {
	return k.collector.GetAppidToAssetCollectorMapping(ctx, appID)
}

func (k Keeper) GetCollectorLookupTable(ctx sdk.Context, appID uint64) (collectorLookup collecortypes.CollectorLookup, found bool) {
	return k.collector.GetCollectorLookupTable(ctx, appID)
}

func (k Keeper) GetAppToDenomsMapping(ctx sdk.Context, appID uint64) (appToDenom collecortypes.AppToDenomsMapping, found bool) {
	return k.collector.GetAppToDenomsMapping(ctx, appID)
}

func (k Keeper) GetLocker(ctx sdk.Context, lockerID string) (locker types.Locker, found bool) {
	return k.locker.GetLocker(ctx, lockerID)
}

func (k Keeper) GetSupply(ctx sdk.Context, denom string) sdk.Coin {
	return k.bank.GetSupply(ctx, denom)
}

func (k Keeper) GetLockerLookupTable(ctx sdk.Context, appMappingID uint64) (lockerLookupData types.LockerLookupTable, found bool) {
	return k.locker.GetLockerLookupTable(ctx, appMappingID)
}

func (k Keeper) GetCollectorLookupByAsset(ctx sdk.Context, appID, assetID uint64) (collectorLookup collecortypes.CollectorLookupTable, found bool) {
	return k.collector.GetCollectorLookupByAsset(ctx, appID, assetID)
}

func (k Keeper) UpdateLocker(ctx sdk.Context, locker types.Locker) {
	k.locker.UpdateLocker(ctx, locker)
}

func (k *Keeper) GetAppExtendedPairVaultMapping(ctx sdk.Context, appMappingID uint64) (appExtendedPairVaultData vaulttypes.AppExtendedPairVaultMapping, found bool) {
	return k.vault.GetAppExtendedPairVaultMapping(ctx, appMappingID)
}

func (k *Keeper) CalculateCollaterlizationRatio(ctx sdk.Context, extendedPairVaultID uint64, amountIn sdk.Int, amountOut sdk.Int) (sdk.Dec, error) {
	return k.vault.CalculateCollaterlizationRatio(ctx, extendedPairVaultID, amountIn, amountOut)
}

func (k *Keeper) GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool) {
	return k.asset.GetAsset(ctx, id)
}

func (k *Keeper) GetVault(ctx sdk.Context, id string) (vault vaulttypes.Vault, found bool) {
	return k.vault.GetVault(ctx, id)
}

func (k *Keeper) DeleteVault(ctx sdk.Context, id string) {
	k.vault.DeleteVault(ctx, id)
}

func (k *Keeper) UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx sdk.Context, counter uint64, vaultData vaulttypes.Vault) {
	k.vault.UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx, counter, vaultData)
}

func (k *Keeper) UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, valutLookupData vaulttypes.AppExtendedPairVaultMapping, extendedPairID uint64, amount sdk.Int, changeType bool) {
	k.vault.UpdateCollateralLockedAmountLockerMapping(ctx, valutLookupData, extendedPairID, amount, changeType)
}

func (k *Keeper) UpdateUserVaultExtendedPairMapping(ctx sdk.Context, extendedPairID uint64, userAddress string, appMappingID uint64) {
	k.vault.UpdateUserVaultExtendedPairMapping(ctx, extendedPairID, userAddress, appMappingID)
}

func (k *Keeper) DeleteAddressFromAppExtendedPairVaultMapping(ctx sdk.Context, extendedPairID uint64, userVaultID string, appMappingID uint64) {
	k.vault.DeleteAddressFromAppExtendedPairVaultMapping(ctx, extendedPairID, userVaultID, appMappingID)
}
func (k *Keeper) GetPairsVault(ctx sdk.Context, id uint64) (pairs assettypes.ExtendedPairVault, found bool) {
	return k.asset.GetPairsVault(ctx, id)
}

func (k *Keeper) SetVault(ctx sdk.Context, vault vaulttypes.Vault) {
	k.vault.SetVault(ctx, vault)
}

func (k *Keeper) BurnCoin(ctx sdk.Context, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return nil
	}

	return k.bank.BurnCoins(ctx, name, sdk.NewCoins(coin))
}

func (k *Keeper) MintCoin(ctx sdk.Context, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return nil
	}

	return k.bank.MintCoins(ctx, name, sdk.NewCoins(coin))
}

func (k *Keeper) SendCoinFromAccountToModule(ctx sdk.Context, address sdk.AccAddress, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return nil
	}

	return k.bank.SendCoinsFromAccountToModule(ctx, address, name, sdk.NewCoins(coin))
}

func (k *Keeper) SendCoinFromModuleToAccount(ctx sdk.Context, name string, address sdk.AccAddress, coin sdk.Coin) error {
	if coin.IsZero() {
		return nil
	}

	return k.bank.SendCoinsFromModuleToAccount(ctx, name, address, sdk.NewCoins(coin))
}
func (k *Keeper) SendCoinFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, coin sdk.Coins) error {
	if coin.IsZero() {
		return nil
	}
	return k.bank.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, coin)
}

func (k *Keeper) SpendableCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins {
	return k.bank.SpendableCoins(ctx, address)
}

func (k *Keeper) GetNetFeeCollectedData(ctx sdk.Context, appID uint64) (netFeeData collecortypes.NetFeeCollectedData, found bool) {
	return k.collector.GetNetFeeCollectedData(ctx, appID)
}

func (k *Keeper) SetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, fee sdk.Int) error {
	return k.collector.SetNetFeeCollectedData(ctx, appID, assetID, fee)
}

func (k *Keeper) DecreaseNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, fee sdk.Int) error {
	return k.collector.DecreaseNetFeeCollectedData(ctx, appID, assetID, fee)
}

func (k *Keeper) SetLockerTotalRewardsByAssetAppWise(ctx sdk.Context, lockerRewardsMapping types.LockerTotalRewardsByAssetAppWise) error {
	return k.locker.SetLockerTotalRewardsByAssetAppWise(ctx, lockerRewardsMapping)
}
func (k *Keeper) GetLockerTotalRewardsByAssetAppWise(ctx sdk.Context, appID, assetID uint64) (lockerRewardsMapping types.LockerTotalRewardsByAssetAppWise, found bool) {
	return k.locker.GetLockerTotalRewardsByAssetAppWise(ctx, appID, assetID)
}
