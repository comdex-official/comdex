package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	collecortypes "github.com/comdex-official/comdex/x/collector/types"
	"github.com/comdex-official/comdex/x/locker/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetLockerProductAssetMapping(ctx sdk.Context, appMappingId uint64) (lockerProductMapping types.LockerProductAssetMapping, found bool) {
	return k.locker.GetLockerProductAssetMapping(ctx, appMappingId)
}

func (k Keeper) GetAppidToAssetCollectorMapping(ctx sdk.Context, app_id uint64) (appAssetCollectorData collecortypes.AppIdToAssetCollectorMapping, found bool) {
	return k.collector.GetAppidToAssetCollectorMapping(ctx, app_id)
}

func (k Keeper) GetCollectorLookupTable(ctx sdk.Context, app_id uint64) (collectorLookup collecortypes.CollectorLookup, found bool) {
	return k.collector.GetCollectorLookupTable(ctx, app_id)
}

func (k Keeper) GetAppToDenomsMapping(ctx sdk.Context, AppId uint64) (appToDenom collecortypes.AppToDenomsMapping, found bool) {
	return k.collector.GetAppToDenomsMapping(ctx, AppId)
}

func (k Keeper) GetLocker(ctx sdk.Context, lockerId string) (locker types.Locker, found bool) {
	return k.locker.GetLocker(ctx, lockerId)
}

func (k Keeper) GetSupply(ctx sdk.Context, denom string) sdk.Coin {
	return k.bank.GetSupply(ctx, denom)
}

func (k Keeper) GetLockerLookupTable(ctx sdk.Context, appMappingId uint64) (lockerLookupData types.LockerLookupTable, found bool) {
	return k.locker.GetLockerLookupTable(ctx, appMappingId)
}

func (k Keeper) GetCollectorLookupByAsset(ctx sdk.Context, app_id, asset_id uint64) (collectorLookup collecortypes.CollectorLookup, found bool) {
	return k.collector.GetCollectorLookupByAsset(ctx, app_id, asset_id)
}

func (k Keeper) UpdateLocker(ctx sdk.Context, locker types.Locker) {
	k.locker.UpdateLocker(ctx, locker)
}

func (k *Keeper) GetAppExtendedPairVaultMapping(ctx sdk.Context, appMappingId uint64) (appExtendedPairVaultData vaulttypes.AppExtendedPairVaultMapping, found bool) {
	return k.vault.GetAppExtendedPairVaultMapping(ctx, appMappingId)
}

func (k *Keeper) CalculateCollaterlizationRatio(ctx sdk.Context, extendedPairVaultId uint64, amountIn sdk.Int, amountOut sdk.Int) (sdk.Dec, error) {
	return k.vault.CalculateCollaterlizationRatio(ctx, extendedPairVaultId, amountIn, amountOut)
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

func (k *Keeper) UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, valutLookupData vaulttypes.AppExtendedPairVaultMapping, extendedPairId uint64, amount sdk.Int, changeType bool) {
	k.vault.UpdateCollateralLockedAmountLockerMapping(ctx, valutLookupData, extendedPairId, amount, changeType)
}

func (k *Keeper) UpdateUserVaultExtendedPairMapping(ctx sdk.Context, extendedPairId uint64, userAddress string, appMappingId uint64) {
	k.vault.UpdateUserVaultExtendedPairMapping(ctx, extendedPairId, userAddress, appMappingId)
}

func (k *Keeper) DeleteAddressFromAppExtendedPairVaultMapping(ctx sdk.Context, extendedPairId uint64, userVaultId string, appMappingId uint64) {
	k.vault.DeleteAddressFromAppExtendedPairVaultMapping(ctx, extendedPairId, userVaultId, appMappingId)
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

func (k *Keeper) GetNetFeeCollectedData(ctx sdk.Context, app_id uint64) (netFeeData collecortypes.NetFeeCollectedData, found bool) {
	return k.collector.GetNetFeeCollectedData(ctx, app_id)
}

func (k *Keeper) SetNetFeeCollectedData(ctx sdk.Context, app_id, asset_id uint64, fee sdk.Int) error {
	return k.collector.SetNetFeeCollectedData(ctx, app_id, asset_id, fee)
}
