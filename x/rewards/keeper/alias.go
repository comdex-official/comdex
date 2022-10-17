package keeper

import (
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	"github.com/comdex-official/comdex/x/liquidity/amm"
	liquiditytypes "github.com/comdex-official/comdex/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	"github.com/comdex-official/comdex/x/locker/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
)

func (k Keeper) GetLockerProductAssetMapping(ctx sdk.Context, appID, assetID uint64) (lockerProductMapping types.LockerProductAssetMapping, found bool) {
	return k.locker.GetLockerProductAssetMapping(ctx, appID, assetID)
}

func (k Keeper) GetAppidToAssetCollectorMapping(ctx sdk.Context, appID, assetID uint64) (appAssetCollectorData collectortypes.AppToAssetIdCollectorMapping, found bool) {
	return k.collector.GetAppidToAssetCollectorMapping(ctx, appID, assetID)
}

func (k Keeper) GetCollectorLookupTable(ctx sdk.Context, appID, assetID uint64) (collectorLookup collectortypes.CollectorLookupTableData, found bool) {
	return k.collector.GetCollectorLookupTable(ctx, appID, assetID)
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

func (k Keeper) GetLockerLookupTableByApp(ctx sdk.Context, appID uint64) (lockerLookupData []types.LockerLookupTableData, found bool) {
	return k.locker.GetLockerLookupTableByApp(ctx, appID)
}

// func (k Keeper) GetCollectorLookupByAsset(ctx sdk.Context, appID, assetID uint64) (collectorLookup collectortypes.CollectorLookupTable, found bool) {
// 	return k.collector.GetCollectorLookupByAsset(ctx, appID, assetID)
// }

func (k Keeper) SetLocker(ctx sdk.Context, locker types.Locker) {
	k.locker.SetLocker(ctx, locker)
}

func (k Keeper) GetAppMappingData(ctx sdk.Context, appMappingID uint64) (appExtendedPairVaultData []vaulttypes.AppExtendedPairVaultMappingData, found bool) {
	return k.vault.GetAppMappingData(ctx, appMappingID)
}

func (k Keeper) CalculateCollateralizationRatio(ctx sdk.Context, extendedPairVaultID uint64, amountIn sdk.Int, amountOut sdk.Int) (sdk.Dec, error) {
	return k.vault.CalculateCollateralizationRatio(ctx, extendedPairVaultID, amountIn, amountOut)
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

func (k Keeper) UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx sdk.Context, vaultData vaulttypes.Vault) {
	k.vault.UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx, vaultData)
}

func (k Keeper) UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, appMappingID uint64, extendedPairID uint64, amount sdk.Int, changeType bool) {
	k.vault.UpdateCollateralLockedAmountLockerMapping(ctx, appMappingID, extendedPairID, amount, changeType)
}

func (k Keeper) DeleteUserVaultExtendedPairMapping(ctx sdk.Context, address string, appID uint64, pairVaultID uint64) {
	k.vault.DeleteUserVaultExtendedPairMapping(ctx, address, appID, pairVaultID)
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

func (k Keeper) GetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64) (netFeeData collectortypes.AppAssetIdToFeeCollectedData, found bool) {
	return k.collector.GetNetFeeCollectedData(ctx, appID, assetID)
}

func (k Keeper) SetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, fee sdk.Int) error {
	return k.collector.SetNetFeeCollectedData(ctx, appID, assetID, fee)
}

func (k Keeper) DecreaseNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, fee sdk.Int, netFeeCollectedData collectortypes.AppAssetIdToFeeCollectedData) error {
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

func (k Keeper) SetLockerLookupTable(ctx sdk.Context, lockerLookupData types.LockerLookupTableData) {
	k.locker.SetLockerLookupTable(ctx, lockerLookupData)
}

func (k Keeper) GetLockerLookupTable(ctx sdk.Context, appID, assetID uint64) (lockerLookupData types.LockerLookupTableData, found bool) {
	return k.locker.GetLockerLookupTable(ctx, appID, assetID)
}

func (k Keeper) GetAppExtendedPairVaultMappingData(ctx sdk.Context, appMappingID uint64, pairVaultID uint64) (appExtendedPairVaultData vaulttypes.AppExtendedPairVaultMappingData, found bool) {
	return k.vault.GetAppExtendedPairVaultMappingData(ctx, appMappingID, pairVaultID)
}

func (k Keeper) CalcAssetPrice(ctx sdk.Context, id uint64, amt sdk.Int) (price sdk.Dec, err error) {
	return k.marketKeeper.CalcAssetPrice(ctx, id, amt)
}

func (k Keeper) GetTwa(ctx sdk.Context, id uint64) (twa markettypes.TimeWeightedAverage, found bool) {
	return k.marketKeeper.GetTwa(ctx, id)
}

func (k Keeper) GetBorrow(ctx sdk.Context, id uint64) (borrow lendtypes.BorrowAsset, found bool) {
	return k.lend.GetBorrow(ctx, id)
}

func (k Keeper) GetLend(ctx sdk.Context, id uint64) (lend lendtypes.LendAsset, found bool) {
	return k.lend.GetLend(ctx, id)
}

func (k Keeper) GetAssetStatsByPoolIDAndAssetID(ctx sdk.Context, poolID, assetID uint64) (PoolAssetLBMapping lendtypes.PoolAssetLBMapping, found bool) {
	return k.lend.GetAssetStatsByPoolIDAndAssetID(ctx, poolID, assetID)
}

func (k Keeper) GetActiveFarmer(ctx sdk.Context, appID, poolID uint64, farmer sdk.AccAddress) (activeFarmer liquiditytypes.ActiveFarmer, found bool) {
	return k.liquidityKeeper.GetActiveFarmer(ctx, appID, poolID, farmer)
}

func (k Keeper) GetAMMPoolInterfaceObject(ctx sdk.Context, appID, poolID uint64) (*liquiditytypes.Pool, *liquiditytypes.Pair, *amm.BasicPool, error) {
	return k.liquidityKeeper.GetAMMPoolInterfaceObject(ctx, appID, poolID)
}

func (k Keeper) CalculateXYFromPoolCoin(ctx sdk.Context, ammPool *amm.BasicPool, poolCoin sdk.Coin) (sdk.Int, sdk.Int, error) {
	return k.liquidityKeeper.CalculateXYFromPoolCoin(ctx, ammPool, poolCoin)
}

func (k Keeper) GetAssetForDenom(ctx sdk.Context, denom string) (asset assettypes.Asset, found bool) {
	return k.asset.GetAssetForDenom(ctx, denom)
}
