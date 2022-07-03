package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/collector/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
)

func (k *Keeper) GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI {
	return k.account.GetModuleAccount(ctx, name)
}

func (k *Keeper) GetModuleAddress(ctx sdk.Context, name string) sdk.AccAddress {
	return k.account.GetModuleAddress(name)
}
func (k *Keeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return k.bank.GetBalance(ctx, addr, denom)
}

func (k *Keeper) MintCoins(ctx sdk.Context, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return nil
	}

	return k.bank.MintCoins(ctx, name, sdk.NewCoins(coin))
}

func (k *Keeper) BurnCoins(ctx sdk.Context, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return nil
	}

	return k.bank.BurnCoins(ctx, name, sdk.NewCoins(coin))
}

func (k *Keeper) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error {
	return k.bank.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, amt)
}
func (k *Keeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	return k.bank.SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt)
}
func (k *Keeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	return k.bank.SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt)
}

func (k *Keeper) GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool) {
	return k.market.GetPriceForAsset(ctx, id)
}

func (k *Keeper) GetLockedVaults(ctx sdk.Context) (lockedVaults []liquidationtypes.LockedVault) {
	return k.liquidation.GetLockedVaults(ctx)
}

func (k *Keeper) GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool) {
	return k.asset.GetPair(ctx, id)
}

func (k *Keeper) GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool) {
	return k.asset.GetAsset(ctx, id)
}

func (k *Keeper) SetFlagIsAuctionInProgress(ctx sdk.Context, id uint64, flag bool) error {
	return k.liquidation.SetFlagIsAuctionInProgress(ctx, id, flag)
}

func (k *Keeper) SetFlagIsAuctionComplete(ctx sdk.Context, id uint64, flag bool) error {
	return k.liquidation.SetFlagIsAuctionComplete(ctx, id, flag)
}

func (k *Keeper) GetAppidToAssetCollectorMapping(ctx sdk.Context, appID uint64) (appAssetCollectorData types.AppIdToAssetCollectorMapping, found bool) {
	return k.collector.GetAppidToAssetCollectorMapping(ctx, appID)
}

func (k *Keeper) UpdateCollector(ctx sdk.Context, appID, assetID uint64, CollectedStabilityFee, CollectedClosingFee, CollectedOpeningFee, LiquidationRewardsCollected sdk.Int) error {
	return k.collector.UpdateCollector(ctx, appID, assetID, CollectedStabilityFee, CollectedClosingFee, CollectedOpeningFee, LiquidationRewardsCollected)
}

func (k *Keeper) SetCollectorLookupTable(ctx sdk.Context, records ...types.CollectorLookupTable) error {
	return k.collector.SetCollectorLookupTable(ctx, records...)
}

func (k *Keeper) GetCollectorLookupTable(ctx sdk.Context, appID uint64) (collectorLookup types.CollectorLookup, found bool) {
	return k.collector.GetCollectorLookupTable(ctx, appID)
}

func (k *Keeper) GetNetFeeCollectedData(ctx sdk.Context, appID uint64) (netFeeData types.NetFeeCollectedData, found bool) {
	return k.collector.GetNetFeeCollectedData(ctx, appID)
}
func (k *Keeper) GetApps(ctx sdk.Context) (apps []assettypes.AppData, found bool) {
	return k.asset.GetApps(ctx)
}
func (k *Keeper) GetApp(ctx sdk.Context, id uint64) (app assettypes.AppData, found bool) {
	return k.asset.GetApp(ctx, id)
}

func (k *Keeper) MintNewTokensForApp(ctx sdk.Context, appMappingID uint64, assetID uint64, address string, amount sdk.Int) error {
	return k.tokenMint.MintNewTokensForApp(ctx, appMappingID, assetID, address, amount)
}

func (k *Keeper) BurnTokensForApp(ctx sdk.Context, appMappingID uint64, assetID uint64, amount sdk.Int) error {
	return k.tokenMint.BurnTokensForApp(ctx, appMappingID, assetID, amount)
}

func (k *Keeper) GetAmountFromCollector(ctx sdk.Context, appID, assetID uint64, amount sdk.Int) (sdk.Int, error) {
	return k.collector.GetAmountFromCollector(ctx, appID, assetID, amount)
}

func (k *Keeper) SetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, fee sdk.Int) error {
	return k.collector.SetNetFeeCollectedData(ctx, appID, assetID, fee)
}

func (k *Keeper) GetLockedVault(ctx sdk.Context, id uint64) (lockedVault liquidationtypes.LockedVault, found bool) {
	return k.liquidation.GetLockedVault(ctx, id)
}

func (k *Keeper) SetLockedVault(ctx sdk.Context, lockedVault liquidationtypes.LockedVault) {
	k.liquidation.SetLockedVault(ctx, lockedVault)
}

func (k *Keeper) GetPairsVault(ctx sdk.Context, id uint64) (pairs assettypes.ExtendedPairVault, found bool) {
	return k.asset.GetPairsVault(ctx, id)
}

func (k *Keeper) GetAppExtendedPairVaultMapping(ctx sdk.Context, appMappingID uint64) (appExtendedPairVaultData vaulttypes.AppExtendedPairVaultMapping, found bool) {
	return k.vault.GetAppExtendedPairVaultMapping(ctx, appMappingID)
}

func (k *Keeper) SetAppExtendedPairVaultMapping(ctx sdk.Context, appExtendedPairVaultData vaulttypes.AppExtendedPairVaultMapping) error {
	return k.vault.SetAppExtendedPairVaultMapping(ctx, appExtendedPairVaultData)
}

func (k *Keeper) GetAuctionMappingForApp(ctx sdk.Context, appID uint64) (collectorAuctionLookupTable types.CollectorAuctionLookupTable, found bool) {
	return k.collector.GetAuctionMappingForApp(ctx, appID)
}
func (k *Keeper) SetAuctionMappingForApp(ctx sdk.Context, records ...types.CollectorAuctionLookupTable) error {
	return k.collector.SetAuctionMappingForApp(ctx, records...)
}

func (k *Keeper) UpdateTokenMintedAmountLockerMapping(ctx sdk.Context, vaultLookupData vaulttypes.AppExtendedPairVaultMapping, extendedPairID uint64, amount sdk.Int, changeType bool) {
	k.vault.UpdateTokenMintedAmountLockerMapping(ctx, vaultLookupData, extendedPairID, amount, changeType)
}
func (k *Keeper) UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, vaultLookupData vaulttypes.AppExtendedPairVaultMapping, extendedPairID uint64, amount sdk.Int, changeType bool) {
	k.vault.UpdateCollateralLockedAmountLockerMapping(ctx, vaultLookupData, extendedPairID, amount, changeType)
}

func (k *Keeper) GetAllAuctionMappingForApp(ctx sdk.Context) (collectorAuctionLookupTable []types.CollectorAuctionLookupTable, found bool) {
	return k.collector.GetAllAuctionMappingForApp(ctx)
}
func (k *Keeper) DeleteLockedVault(ctx sdk.Context, id uint64) {
	k.liquidation.DeleteLockedVault(ctx, id)
}

func (k *Keeper) UpdateUserVaultExtendedPairMapping(ctx sdk.Context, extendedPairID uint64, userAddress string, appMappingID uint64) {
	k.vault.UpdateUserVaultExtendedPairMapping(ctx, extendedPairID, userAddress, appMappingID)
}

func (k Keeper) CreateLockedVaultHistory(ctx sdk.Context, lockedVault liquidationtypes.LockedVault) error {
	return k.liquidation.CreateLockedVaultHistory(ctx, lockedVault)
}

func (k *Keeper) GetKillSwitchData(ctx sdk.Context, app_id uint64) (esmtypes.KillSwitchParams, bool) {
	return k.esm.GetKillSwitchData(ctx, app_id)
}

func (k *Keeper) GetESMStatus(ctx sdk.Context, id uint64) (esmtypes.ESMStatus, bool) {
	return k.esm.GetESMStatus(ctx,id)
}

func (k *Keeper) CreateNewVault(ctx sdk.Context, From string, AppId uint64, ExtendedPairVaultID uint64, AmountIn sdk.Int, AmountOut sdk.Int) error {
	return k.vault.CreateNewVault(ctx, From, AppId, ExtendedPairVaultID, AmountIn, AmountOut)
}

func (k *Keeper) GetUserVaultExtendedPairMapping(ctx sdk.Context, address string) (userVaultAssetData vaulttypes.UserVaultAssetMapping, found bool){
	return k.vault.GetUserVaultExtendedPairMapping(ctx, address)
}

func (k *Keeper) CheckUserAppToExtendedPairMapping(ctx sdk.Context, userVaultAssetData vaulttypes.UserVaultAssetMapping, extendedPairVaultID uint64, appMappingID uint64) (vaultID string, found bool) {
	return k.vault.CheckUserAppToExtendedPairMapping(ctx, userVaultAssetData, extendedPairVaultID, appMappingID)
}

func (k *Keeper) SetVault(ctx sdk.Context, vault vaulttypes.Vault) {
	k.vault.SetVault(ctx, vault)
}

func (k *Keeper) GetVault(ctx sdk.Context, id string) (vault vaulttypes.Vault, found bool){
	return k.vault.GetVault(ctx, id)
}