package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	"github.com/comdex-official/comdex/x/collector/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (k Keeper) GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI {
	return k.account.GetModuleAccount(ctx, name)
}

func (k Keeper) GetModuleAddress(_ sdk.Context, name string) sdk.AccAddress {
	return k.account.GetModuleAddress(name)
}
func (k Keeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return k.bank.GetBalance(ctx, addr, denom)
}

func (k Keeper) BurnCoins(ctx sdk.Context, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return auctiontypes.BurnCoinValueInCloseAuctionIsZero
	}

	return k.bank.BurnCoins(ctx, name, sdk.NewCoins(coin))
}

func (k Keeper) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, coin sdk.Coins) error {
	if coin.IsZero() {
		return auctiontypes.SendCoinsFromModuleToModuleInAuctionIsZero
	}

	return k.bank.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, coin)
}
func (k Keeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, coin sdk.Coins) error {
	if coin.IsZero() {
		return auctiontypes.SendCoinsFromModuleToAccountInAuctionIsZero
	}

	return k.bank.SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, coin)
}
func (k Keeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, coin sdk.Coins) error {
	if coin.IsZero() {
		return auctiontypes.SendCoinsFromAccountToModuleInAuctionIsZero
	}

	return k.bank.SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, coin)
}

func (k Keeper) GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool) {
	return k.market.GetPriceForAsset(ctx, id)
}

func (k Keeper) GetLockedVaults(ctx sdk.Context) (lockedVaults []liquidationtypes.LockedVault) {
	return k.liquidation.GetLockedVaults(ctx)
}

func (k Keeper) GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool) {
	return k.asset.GetPair(ctx, id)
}

func (k Keeper) GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool) {
	return k.asset.GetAsset(ctx, id)
}

func (k Keeper) SetFlagIsAuctionInProgress(ctx sdk.Context, appID, id uint64, flag bool) error {
	return k.liquidation.SetFlagIsAuctionInProgress(ctx, appID, id, flag)
}

func (k Keeper) SetFlagIsAuctionComplete(ctx sdk.Context, appID, id uint64, flag bool) error {
	return k.liquidation.SetFlagIsAuctionComplete(ctx, appID, id, flag)
}

func (k Keeper) GetAppidToAssetCollectorMapping(ctx sdk.Context, appID, assetID uint64) (appAssetCollectorData types.AppToAssetIdCollectorMapping, found bool) {
	return k.collector.GetAppidToAssetCollectorMapping(ctx, appID, assetID)
}

func (k Keeper) UpdateCollector(ctx sdk.Context, appID, assetID uint64, CollectedStabilityFee, CollectedClosingFee, CollectedOpeningFee, LiquidationRewardsCollected sdk.Int) error {
	return k.collector.UpdateCollector(ctx, appID, assetID, CollectedStabilityFee, CollectedClosingFee, CollectedOpeningFee, LiquidationRewardsCollected)
}

// func (k Keeper) SetCollectorLookupTable(ctx sdk.Context, records ...types.CollectorLookupTable) error {
// 	return k.collector.SetCollectorLookupTable(ctx, records...)
// }

func (k Keeper) GetCollectorLookupTable(ctx sdk.Context, appID, assetID uint64) (collectorLookup types.CollectorLookupTableData, found bool) {
	return k.collector.GetCollectorLookupTable(ctx, appID, assetID)
}

func (k Keeper) GetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64) (netFeeData types.AppAssetIdToFeeCollectedData, found bool) {
	return k.collector.GetNetFeeCollectedData(ctx, appID, assetID)
}
func (k Keeper) GetApps(ctx sdk.Context) (apps []assettypes.AppData, found bool) {
	return k.asset.GetApps(ctx)
}
func (k Keeper) GetApp(ctx sdk.Context, id uint64) (app assettypes.AppData, found bool) {
	return k.asset.GetApp(ctx, id)
}

func (k Keeper) MintNewTokensForApp(ctx sdk.Context, appMappingID uint64, assetID uint64, address string, amount sdk.Int) error {
	return k.tokenMint.MintNewTokensForApp(ctx, appMappingID, assetID, address, amount)
}

func (k Keeper) BurnTokensForApp(ctx sdk.Context, appMappingID uint64, assetID uint64, amount sdk.Int) error {
	return k.tokenMint.BurnTokensForApp(ctx, appMappingID, assetID, amount)
}

func (k Keeper) GetAmountFromCollector(ctx sdk.Context, appID, assetID uint64, amount sdk.Int) (sdk.Int, error) {
	return k.collector.GetAmountFromCollector(ctx, appID, assetID, amount)
}

func (k Keeper) SetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, fee sdk.Int) error {
	return k.collector.SetNetFeeCollectedData(ctx, appID, assetID, fee)
}

func (k Keeper) GetLockedVault(ctx sdk.Context, appID, id uint64) (lockedVault liquidationtypes.LockedVault, found bool) {
	return k.liquidation.GetLockedVault(ctx, appID, id)
}

func (k Keeper) SetLockedVault(ctx sdk.Context, lockedVault liquidationtypes.LockedVault) {
	k.liquidation.SetLockedVault(ctx, lockedVault)
}

func (k Keeper) GetPairsVault(ctx sdk.Context, id uint64) (pairs assettypes.ExtendedPairVault, found bool) {
	return k.asset.GetPairsVault(ctx, id)
}

func (k Keeper) GetAppExtendedPairVaultMappingData(ctx sdk.Context, appMappingID uint64, pairVaultsID uint64) (appExtendedPairVaultData vaulttypes.AppExtendedPairVaultMappingData, found bool) {
	return k.vault.GetAppExtendedPairVaultMappingData(ctx, appMappingID, pairVaultsID)
}

func (k Keeper) SetAppExtendedPairVaultMappingData(ctx sdk.Context, appExtendedPairVaultData vaulttypes.AppExtendedPairVaultMappingData)  {
	k.vault.SetAppExtendedPairVaultMappingData(ctx, appExtendedPairVaultData)
}

func (k Keeper) GetAuctionMappingForApp(ctx sdk.Context, appID, assetID uint64) (collectorAuctionLookupTable types.AppAssetIdToAuctionLookupTable, found bool) {
	return k.collector.GetAuctionMappingForApp(ctx, appID, assetID)
}
func (k Keeper) SetAuctionMappingForApp(ctx sdk.Context, records types.AppAssetIdToAuctionLookupTable) error {
	return k.collector.SetAuctionMappingForApp(ctx, records)
}

func (k Keeper) UpdateTokenMintedAmountLockerMapping(ctx sdk.Context, appMappingID uint64, extendedPairID uint64, amount sdk.Int, changeType bool) {
	k.vault.UpdateTokenMintedAmountLockerMapping(ctx, appMappingID, extendedPairID, amount, changeType)
}
func (k Keeper) UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, appMappingID uint64, extendedPairID uint64, amount sdk.Int, changeType bool) {
	k.vault.UpdateCollateralLockedAmountLockerMapping(ctx, appMappingID, extendedPairID, amount, changeType)
}

func (k Keeper) GetAllAuctionMappingForApp(ctx sdk.Context) (collectorAuctionLookupTable []types.AppAssetIdToAuctionLookupTable, found bool) {
	return k.collector.GetAllAuctionMappingForApp(ctx)
}
func (k Keeper) DeleteLockedVault(ctx sdk.Context, appId, id uint64) {
	k.liquidation.DeleteLockedVault(ctx, appId, id)
}

func (k Keeper) DeleteUserVaultExtendedPairMapping(ctx sdk.Context, address string, appID uint64, pairVaultID uint64) {
	k.vault.DeleteUserVaultExtendedPairMapping(ctx, address, appID, pairVaultID)
}

func (k Keeper) CreateLockedVaultHistory(ctx sdk.Context, lockedVault liquidationtypes.LockedVault) error {
	return k.liquidation.CreateLockedVaultHistory(ctx, lockedVault)
}

func (k Keeper) GetKillSwitchData(ctx sdk.Context, appID uint64) (esmtypes.KillSwitchParams, bool) {
	return k.esm.GetKillSwitchData(ctx, appID)
}

func (k Keeper) GetESMStatus(ctx sdk.Context, id uint64) (esmtypes.ESMStatus, bool) {
	return k.esm.GetESMStatus(ctx, id)
}

func (k Keeper) CreateNewVault(ctx sdk.Context, From string, AppID uint64, ExtendedPairVaultID uint64, AmountIn sdk.Int, AmountOut sdk.Int) error {
	return k.vault.CreateNewVault(ctx, From, AppID, ExtendedPairVaultID, AmountIn, AmountOut)
}

func (k Keeper) GetUserAppExtendedPairMappingData(ctx sdk.Context, from string, appMapping uint64, extendedPairVault uint64) (userVaultAssetData vaulttypes.OwnerAppExtendedPairVaultMappingData, found bool) {
	return k.vault.GetUserAppExtendedPairMappingData(ctx, from, appMapping, extendedPairVault)
}

func (k Keeper) GetUserAppMappingData(ctx sdk.Context, address string, appID uint64) (mappingData []vaulttypes.OwnerAppExtendedPairVaultMappingData, found bool) {
	return k.vault.GetUserAppMappingData(ctx, address, appID)
}

// func (k Keeper) CheckUserAppToExtendedPairMapping(ctx sdk.Context, userVaultAssetData vaulttypes.UserVaultAssetMapping, extendedPairVaultID uint64, appMappingID uint64) (vaultID uint64, found bool) {
// 	return k.vault.CheckUserAppToExtendedPairMapping(ctx, userVaultAssetData, extendedPairVaultID, appMappingID)
// }

func (k Keeper) SetVault(ctx sdk.Context, vault vaulttypes.Vault) {
	k.vault.SetVault(ctx, vault)
}

func (k Keeper) GetVault(ctx sdk.Context, id uint64) (vault vaulttypes.Vault, found bool) {
	return k.vault.GetVault(ctx, id)
}

func (k Keeper) GetBorrows(ctx sdk.Context) (userBorrows lendtypes.BorrowMapping, found bool) {
	return k.lend.GetBorrows(ctx)
}

func (k Keeper) GetBorrow(ctx sdk.Context, id uint64) (borrow lendtypes.BorrowAsset, found bool) {
	return k.lend.GetBorrow(ctx, id)
}

func (k Keeper) GetLendPair(ctx sdk.Context, id uint64) (pair lendtypes.Extended_Pair, found bool) {
	return k.lend.GetLendPair(ctx, id)
}

func (k Keeper) GetAssetRatesStats(ctx sdk.Context, assetID uint64) (assetRatesStats lendtypes.AssetRatesStats, found bool) {
	return k.lend.GetAssetRatesStats(ctx, assetID)
}

func (k Keeper) VerifyCollaterlizationRatio(ctx sdk.Context, amountIn sdk.Int, assetIn assettypes.Asset, amountOut sdk.Int, assetOut assettypes.Asset, liquidationThreshold sdk.Dec) error {
	return k.lend.VerifyCollaterlizationRatio(ctx, amountIn, assetIn, amountOut, assetOut, liquidationThreshold)
}

func (k Keeper) CalculateLendCollaterlizationRatio(ctx sdk.Context, amountIn sdk.Int, assetIn assettypes.Asset, amountOut sdk.Int, assetOut assettypes.Asset) (sdk.Dec, error) {
	return k.lend.CalculateCollaterlizationRatio(ctx, amountIn, assetIn, amountOut, assetOut)
}

func (k Keeper) GetLend(ctx sdk.Context, id uint64) (lend lendtypes.LendAsset, found bool) {
	return k.lend.GetLend(ctx, id)
}

func (k Keeper) DeleteBorrow(ctx sdk.Context, id uint64) {
	k.lend.DeleteBorrow(ctx, id)
}

func (k Keeper) DeleteBorrowForAddressByPair(ctx sdk.Context, address sdk.AccAddress, pairID uint64) {
	k.lend.DeleteBorrowForAddressByPair(ctx, address, pairID)
}

func (k Keeper) UpdateUserBorrowIDMapping(ctx sdk.Context, borrowOwner string, borrowID uint64, isInsert bool) error {
	return k.lend.UpdateUserBorrowIDMapping(ctx, borrowOwner, borrowID, isInsert)
}

func (k Keeper) UpdateBorrowIDByOwnerAndPoolMapping(ctx sdk.Context, borrowOwner string, borrowID uint64, poolID uint64, isInsert bool) error {
	return k.lend.UpdateBorrowIDByOwnerAndPoolMapping(ctx, borrowOwner, borrowID, poolID, isInsert)
}

func (k Keeper) UpdateBorrowIdsMapping(ctx sdk.Context, borrowID uint64, isInsert bool) error {
	return k.lend.UpdateBorrowIdsMapping(ctx, borrowID, isInsert)
}

func (k Keeper) CreteNewBorrow(ctx sdk.Context, liqBorrow liquidationtypes.LockedVault) {
	k.lend.CreteNewBorrow(ctx, liqBorrow)
}

func (k Keeper) GetPool(ctx sdk.Context, id uint64) (pool lendtypes.Pool, found bool) {
	return k.lend.GetPool(ctx, id)
}

func (k Keeper) GetAddAuctionParamsData(ctx sdk.Context, appID uint64) (auctionParams lendtypes.AuctionParams, found bool) {
	return k.lend.GetAddAuctionParamsData(ctx, appID)
}

func (k Keeper) GetReserveDepositStats(ctx sdk.Context) (depositStats lendtypes.DepositStats, found bool) {
	return k.lend.GetReserveDepositStats(ctx)
}

func (k Keeper) ModuleBalance(ctx sdk.Context, moduleName string, denom string) sdk.Int {
	return k.lend.ModuleBalance(ctx, moduleName, denom)
}

func (k Keeper) UpdateReserveBalances(ctx sdk.Context, assetID uint64, moduleName string, payment sdk.Coin, inc bool) error {
	return k.lend.UpdateReserveBalances(ctx, assetID, moduleName, payment, inc)
}
func (k Keeper) UnLiquidateLockedBorrows(ctx sdk.Context, appID, id uint64) error {
	return k.liquidation.UnLiquidateLockedBorrows(ctx, appID, id)
}

func (k Keeper) SetLend(ctx sdk.Context, lend lendtypes.LendAsset) {
	k.lend.SetLend(ctx, lend)
}
