package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	"github.com/comdex-official/comdex/x/locker/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetLockerProductAssetMapping(ctx sdk.Context, appMappingID uint64) (lockerProductMapping types.LockerProductAssetMapping, found bool) {
	return k.locker.GetLockerProductAssetMapping(ctx, appMappingID)
}

func (k Keeper) GetAppidToAssetCollectorMapping(ctx sdk.Context, appID uint64) (appAssetCollectorData collectortypes.AppIdToAssetCollectorMapping, found bool) {
	return k.collector.GetAppidToAssetCollectorMapping(ctx, appID)
}

func (k Keeper) GetCollectorLookupTable(ctx sdk.Context, appID uint64) (collectorLookup collectortypes.CollectorLookup, found bool) {
	return k.collector.GetCollectorLookupTable(ctx, appID)
}

func (k Keeper) GetAppToDenomsMapping(ctx sdk.Context, appID uint64) (appToDenom collectortypes.AppToDenomsMapping, found bool) {
	return k.collector.GetAppToDenomsMapping(ctx, appID)
}

func (k Keeper) GetLocker(ctx sdk.Context, lockerID uint64) (locker types.Locker, found bool) {
	return k.locker.GetLocker(ctx, lockerID)
}

func (k Keeper) GetSupply(ctx sdk.Context, denom string) sdk.Coin {
	return k.bank.GetSupply(ctx, denom)
}

func (k Keeper) GetLockerLookupTable(ctx sdk.Context, appMappingID uint64) (lockerLookupData types.LockerLookupTable, found bool) {
	return k.locker.GetLockerLookupTable(ctx, appMappingID)
}

func (k Keeper) GetCollectorLookupByAsset(ctx sdk.Context, appID, assetID uint64) (collectorLookup collectortypes.CollectorLookupTable, found bool) {
	return k.collector.GetCollectorLookupByAsset(ctx, appID, assetID)
}

func (k Keeper) UpdateLocker(ctx sdk.Context, locker types.Locker) {
	k.locker.UpdateLocker(ctx, locker)
}

func (k Keeper) GetAppExtendedPairVaultMapping(ctx sdk.Context, appMappingID uint64) (appExtendedPairVaultData vaulttypes.AppExtendedPairVaultMapping, found bool) {
	return k.vault.GetAppExtendedPairVaultMapping(ctx, appMappingID)
}

func (k Keeper) CalculateCollaterlizationRatio(ctx sdk.Context, extendedPairVaultID uint64, amountIn sdk.Int, amountOut sdk.Int) (sdk.Dec, error) {
	return k.vault.CalculateCollaterlizationRatio(ctx, extendedPairVaultID, amountIn, amountOut)
}

func (k Keeper) GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool) {
	return k.asset.GetAsset(ctx, id)
}

func (k Keeper) GetVault(ctx sdk.Context, id uint64) (vault vaulttypes.Vault, found bool) {
	return k.vault.GetVault(ctx, id)
}

func (k Keeper) DeleteVault(ctx sdk.Context, id uint64) {
	k.vault.DeleteVault(ctx, id)
}

func (k Keeper) UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx sdk.Context, counter uint64, vaultData vaulttypes.Vault) {
	k.vault.UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx, counter, vaultData)
}

func (k Keeper) UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, vaultLookupData vaulttypes.AppExtendedPairVaultMapping, extendedPairID uint64, amount sdk.Int, changeType bool) {
	k.vault.UpdateCollateralLockedAmountLockerMapping(ctx, vaultLookupData, extendedPairID, amount, changeType)
}

func (k Keeper) UpdateUserVaultExtendedPairMapping(ctx sdk.Context, extendedPairID uint64, userAddress string, appMappingID uint64) {
	k.vault.UpdateUserVaultExtendedPairMapping(ctx, extendedPairID, userAddress, appMappingID)
}

func (k Keeper) DeleteAddressFromAppExtendedPairVaultMapping(ctx sdk.Context, extendedPairID uint64, userVaultID uint64, appMappingID uint64) {
	k.vault.DeleteAddressFromAppExtendedPairVaultMapping(ctx, extendedPairID, userVaultID, appMappingID)
}
func (k Keeper) GetPairsVault(ctx sdk.Context, id uint64) (pairs assettypes.ExtendedPairVault, found bool) {
	return k.asset.GetPairsVault(ctx, id)
}

func (k Keeper) SetVault(ctx sdk.Context, vault vaulttypes.Vault) {
	k.vault.SetVault(ctx, vault)
}

func (k Keeper) BurnCoin(ctx sdk.Context, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return rewardstypes.BurnCoinValueInRewardsIsZero
	}

	return k.bank.BurnCoins(ctx, name, sdk.NewCoins(coin))
}

func (k Keeper) MintCoin(ctx sdk.Context, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return rewardstypes.MintCoinValueInRewardsIsZero
	}

	return k.bank.MintCoins(ctx, name, sdk.NewCoins(coin))
}

func (k Keeper) SendCoinFromAccountToModule(ctx sdk.Context, address sdk.AccAddress, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return rewardstypes.SendCoinsFromAccountToModuleInRewardsIsZero
	}

	return k.bank.SendCoinsFromAccountToModule(ctx, address, name, sdk.NewCoins(coin))
}

func (k Keeper) SendCoinFromModuleToAccount(ctx sdk.Context, name string, address sdk.AccAddress, coin sdk.Coin) error {
	if coin.IsZero() {
		return rewardstypes.SendCoinsFromModuleToAccountInRewardsIsZero
	}

	return k.bank.SendCoinsFromModuleToAccount(ctx, name, address, sdk.NewCoins(coin))
}
func (k Keeper) SendCoinFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, coin sdk.Coins) error {
	if coin.IsZero() {
		return rewardstypes.SendCoinsFromModuleToModuleInRewardsIsZero
	}
	return k.bank.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, coin)
}

func (k Keeper) SpendableCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins {
	return k.bank.SpendableCoins(ctx, address)
}

func (k Keeper) GetNetFeeCollectedData(ctx sdk.Context, appID uint64) (netFeeData collectortypes.NetFeeCollectedData, found bool) {
	return k.collector.GetNetFeeCollectedData(ctx, appID)
}

func (k Keeper) SetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, fee sdk.Int) error {
	return k.collector.SetNetFeeCollectedData(ctx, appID, assetID, fee)
}

func (k Keeper) DecreaseNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, fee sdk.Int, netFeeCollectedData collectortypes.NetFeeCollectedData) error {
	return k.collector.DecreaseNetFeeCollectedData(ctx, appID, assetID, fee, netFeeCollectedData)
}

func (k Keeper) SetLockerTotalRewardsByAssetAppWise(ctx sdk.Context, lockerRewardsMapping types.LockerTotalRewardsByAssetAppWise) error {
	return k.locker.SetLockerTotalRewardsByAssetAppWise(ctx, lockerRewardsMapping)
}
func (k Keeper) GetLockerTotalRewardsByAssetAppWise(ctx sdk.Context, appID, assetID uint64) (lockerRewardsMapping types.LockerTotalRewardsByAssetAppWise, found bool) {
	return k.locker.GetLockerTotalRewardsByAssetAppWise(ctx, appID, assetID)
}

func (k Keeper) GetKillSwitchData(ctx sdk.Context, appID uint64) (esmtypes.KillSwitchParams, bool) {
	return k.esm.GetKillSwitchData(ctx, appID)
}

func (k Keeper) GetESMStatus(ctx sdk.Context, id uint64) (esmtypes.ESMStatus, bool) {
	return k.esm.GetESMStatus(ctx, id)
}

func (k Keeper) GetLockers(ctx sdk.Context) (locker []types.Locker) {
	return k.locker.GetLockers(ctx)
}
