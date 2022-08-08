package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	"github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (k Keeper) GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI {
	return k.account.GetModuleAccount(ctx, name)
}

func (k Keeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return k.bank.GetBalance(ctx, addr, denom)
}

func (k Keeper) GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool) {
	return k.asset.GetPair(ctx, id)
}

func (k Keeper) GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool) {
	return k.asset.GetAsset(ctx, id)
}

func (k Keeper) GetApps(ctx sdk.Context) (apps []assettypes.AppData, found bool) {
	return k.asset.GetApps(ctx)
}

func (k Keeper) SendCoinFromModuleToAccount(ctx sdk.Context, name string, address sdk.AccAddress, coin sdk.Coin) error {
	if coin.IsZero() {
		return liquidationtypes.SendCoinsFromModuleToAccountInLiquidationIsZero
	}

	return k.bank.SendCoinsFromModuleToAccount(ctx, name, address, sdk.NewCoins(coin))
}

func (k Keeper) GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool) {
	return k.market.GetPriceForAsset(ctx, id)
}

func (k Keeper) GetAppExtendedPairVaultMapping(ctx sdk.Context, appMappingID uint64) (appExtendedPairVaultData types.AppExtendedPairVaultMapping, found bool) {
	return k.vault.GetAppExtendedPairVaultMapping(ctx, appMappingID)
}

func (k Keeper) CalculateCollaterlizationRatio(ctx sdk.Context, extendedPairVaultID uint64, amountIn sdk.Int, amountOut sdk.Int) (sdk.Dec, error) {
	return k.vault.CalculateCollaterlizationRatio(ctx, extendedPairVaultID, amountIn, amountOut)
}

func (k Keeper) GetVault(ctx sdk.Context, id string) (vault types.Vault, found bool) {
	return k.vault.GetVault(ctx, id)
}

func (k Keeper) DeleteVault(ctx sdk.Context, id string) {
	k.vault.DeleteVault(ctx, id)
}

func (k Keeper) UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx sdk.Context, counter uint64, vaultData types.Vault) {
	k.vault.UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx, counter, vaultData)
}

func (k Keeper) UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, vaultLookupData types.AppExtendedPairVaultMapping, extendedPairID uint64, amount sdk.Int, changeType bool) {
	k.vault.UpdateCollateralLockedAmountLockerMapping(ctx, vaultLookupData, extendedPairID, amount, changeType)
}

func (k Keeper) UpdateUserVaultExtendedPairMapping(ctx sdk.Context, extendedPairID uint64, userAddress string, appMappingID uint64) {
	k.vault.UpdateUserVaultExtendedPairMapping(ctx, extendedPairID, userAddress, appMappingID)
}

func (k Keeper) DeleteAddressFromAppExtendedPairVaultMapping(ctx sdk.Context, extendedPairID uint64, userVaultID string, appMappingID uint64) {
	k.vault.DeleteAddressFromAppExtendedPairVaultMapping(ctx, extendedPairID, userVaultID, appMappingID)
}

func (k Keeper) GetPairsVault(ctx sdk.Context, id uint64) (pairs assettypes.ExtendedPairVault, found bool) {
	return k.asset.GetPairsVault(ctx, id)
}

func (k Keeper) GetAuctionParams(ctx sdk.Context) auctiontypes.Params {
	return k.auction.GetParams(ctx)
}

func (k Keeper) SetVault(ctx sdk.Context, vault types.Vault) {
	k.vault.SetVault(ctx, vault)
}

func (k Keeper) GetKillSwitchData(ctx sdk.Context, appID uint64) (esmtypes.KillSwitchParams, bool) {
	return k.esm.GetKillSwitchData(ctx, appID)
}

func (k Keeper) GetESMStatus(ctx sdk.Context, id uint64) (esmtypes.ESMStatus, bool) {
	return k.esm.GetESMStatus(ctx, id)
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

func (k Keeper) GetBorrowStats(ctx sdk.Context) (borrowStats lendtypes.DepositStats, found bool) {
	return k.lend.GetBorrowStats(ctx)
}

func (k Keeper) SetBorrowStats(ctx sdk.Context, borrowStats lendtypes.DepositStats) {
	k.lend.SetBorrowStats(ctx, borrowStats)
}

func (k Keeper) GetAssetStatsByPoolIDAndAssetID(ctx sdk.Context, assetID, poolID uint64) (AssetStats lendtypes.AssetStats, found bool) {
	return k.lend.GetAssetStatsByPoolIDAndAssetID(ctx, assetID, poolID)
}

func (k Keeper) SetAssetStatsByPoolIDAndAssetID(ctx sdk.Context, AssetStats lendtypes.AssetStats) {
	k.lend.SetAssetStatsByPoolIDAndAssetID(ctx, AssetStats)
}

func (k Keeper) UpdateBorrowStats(ctx sdk.Context, pair lendtypes.Extended_Pair, borrowPos lendtypes.BorrowAsset, amount sdk.Int, inc bool) {
	k.lend.UpdateBorrowStats(ctx, pair, borrowPos, amount, inc)
}
