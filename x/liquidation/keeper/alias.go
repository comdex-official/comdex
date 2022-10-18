package keeper

// import (
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

// 	assettypes "github.com/comdex-official/comdex/x/asset/types"
// 	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
// 	esmtypes "github.com/comdex-official/comdex/x/esm/types"
// 	lendtypes "github.com/comdex-official/comdex/x/lend/types"
// 	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
// 	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
// 	"github.com/comdex-official/comdex/x/vault/types"
// )

// func (k Keeper) GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI {
// 	return k.account.GetModuleAccount(ctx, name)
// }

// func (k Keeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
// 	return k.bank.GetBalance(ctx, addr, denom)
// }

// func (k Keeper) GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool) {
// 	return k.asset.GetPair(ctx, id)
// }

// func (k Keeper) GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool) {
// 	return k.asset.GetAsset(ctx, id)
// }

// func (k Keeper) GetApps(ctx sdk.Context) (apps []assettypes.AppData, found bool) {
// 	return k.asset.GetApps(ctx)
// }

// func (k Keeper) SendCoinFromModuleToAccount(ctx sdk.Context, name string, address sdk.AccAddress, coin sdk.Coin) error {
// 	if coin.IsZero() {
// 		return liquidationtypes.SendCoinsFromModuleToAccountInLiquidationIsZero
// 	}

// 	return k.bank.SendCoinsFromModuleToAccount(ctx, name, address, sdk.NewCoins(coin))
// }

// func (k Keeper) GetAppMappingData(ctx sdk.Context, appMappingID uint64) (appExtendedPairVaultData []types.AppExtendedPairVaultMappingData, found bool) {
// 	return k.vault.GetAppMappingData(ctx, appMappingID)
// }

// func (k Keeper) CalculateCollateralizationRatio(ctx sdk.Context, extendedPairVaultID uint64, amountIn sdk.Int, amountOut sdk.Int) (sdk.Dec, error) {
// 	return k.vault.CalculateCollateralizationRatio(ctx, extendedPairVaultID, amountIn, amountOut)
// }

// func (k Keeper) GetVault(ctx sdk.Context, id uint64) (vault types.Vault, found bool) {
// 	return k.vault.GetVault(ctx, id)
// }

// func (k Keeper) GetVaults(ctx sdk.Context) (vaults []types.Vault) {
// 	return k.vault.GetVaults(ctx)
// }

// func (k Keeper) GetIDForVault(ctx sdk.Context) uint64 {
// 	return k.vault.GetIDForVault(ctx)
// }

// func (k Keeper) GetLengthOfVault(ctx sdk.Context) uint64 {
// 	return k.vault.GetLengthOfVault(ctx)
// }

// func (k Keeper) SetLengthOfVault(ctx sdk.Context, length uint64) {
// 	k.vault.SetLengthOfVault(ctx, length)
// }

// func (k Keeper) DeleteVault(ctx sdk.Context, id uint64) {
// 	k.vault.DeleteVault(ctx, id)
// }

// func (k Keeper) UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx sdk.Context, vaultData types.Vault) {
// 	k.vault.UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx, vaultData)
// }

// func (k Keeper) UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, appMappingID uint64, extendedPairID uint64, amount sdk.Int, changeType bool) {
// 	k.vault.UpdateCollateralLockedAmountLockerMapping(ctx, appMappingID, extendedPairID, amount, changeType)
// }

// func (k Keeper) DeleteUserVaultExtendedPairMapping(ctx sdk.Context, address string, appID uint64, pairVaultID uint64) {
// 	k.vault.DeleteUserVaultExtendedPairMapping(ctx, address, appID, pairVaultID)
// }

// func (k Keeper) DeleteAddressFromAppExtendedPairVaultMapping(ctx sdk.Context, extendedPairID uint64, userVaultID uint64, appMappingID uint64) {
// 	k.vault.DeleteAddressFromAppExtendedPairVaultMapping(ctx, extendedPairID, userVaultID, appMappingID)
// }

// func (k Keeper) GetPairsVault(ctx sdk.Context, id uint64) (pairs assettypes.ExtendedPairVault, found bool) {
// 	return k.asset.GetPairsVault(ctx, id)
// }

// func (k Keeper) GetAuctionParams(ctx sdk.Context) auctiontypes.Params {
// 	return k.auction.GetParams(ctx)
// }

// func (k Keeper) SetVault(ctx sdk.Context, vault types.Vault) {
// 	k.vault.SetVault(ctx, vault)
// }

// func (k Keeper) GetKillSwitchData(ctx sdk.Context, appID uint64) (esmtypes.KillSwitchParams, bool) {
// 	return k.esm.GetKillSwitchData(ctx, appID)
// }

// func (k Keeper) GetESMStatus(ctx sdk.Context, id uint64) (esmtypes.ESMStatus, bool) {
// 	return k.esm.GetESMStatus(ctx, id)
// }

// func (k Keeper) GetBorrows(ctx sdk.Context) (userBorrows []uint64, found bool) {
// 	return k.lend.GetBorrows(ctx)
// }

// func (k Keeper) GetBorrow(ctx sdk.Context, id uint64) (borrow lendtypes.BorrowAsset, found bool) {
// 	return k.lend.GetBorrow(ctx, id)
// }

// func (k Keeper) GetLendPair(ctx sdk.Context, id uint64) (pair lendtypes.Extended_Pair, found bool) {
// 	return k.lend.GetLendPair(ctx, id)
// }

// func (k Keeper) GetAssetRatesParams(ctx sdk.Context, assetID uint64) (assetRatesStats lendtypes.AssetRatesParams, found bool) {
// 	return k.lend.GetAssetRatesParams(ctx, assetID)
// }

// func (k Keeper) VerifyCollateralizationRatio(ctx sdk.Context, amountIn sdk.Int, assetIn assettypes.Asset, amountOut sdk.Int, assetOut assettypes.Asset, liquidationThreshold sdk.Dec) error {
// 	return k.lend.VerifyCollateralizationRatio(ctx, amountIn, assetIn, amountOut, assetOut, liquidationThreshold)
// }

// func (k Keeper) CalculateLendCollateralizationRatio(ctx sdk.Context, amountIn sdk.Int, assetIn assettypes.Asset, amountOut sdk.Int, assetOut assettypes.Asset) (sdk.Dec, error) {
// 	return k.lend.CalculateCollateralizationRatio(ctx, amountIn, assetIn, amountOut, assetOut)
// }

// func (k Keeper) GetLend(ctx sdk.Context, id uint64) (lend lendtypes.LendAsset, found bool) {
// 	return k.lend.GetLend(ctx, id)
// }

// func (k Keeper) CreteNewBorrow(ctx sdk.Context, liqBorrow liquidationtypes.LockedVault) {
// 	k.lend.CreteNewBorrow(ctx, liqBorrow)
// }

// func (k Keeper) GetPool(ctx sdk.Context, id uint64) (pool lendtypes.Pool, found bool) {
// 	return k.lend.GetPool(ctx, id)
// }

// func (k Keeper) GetAssetStatsByPoolIDAndAssetID(ctx sdk.Context, poolID, assetID uint64) (AssetStats lendtypes.PoolAssetLBMapping, found bool) {
// 	return k.lend.GetAssetStatsByPoolIDAndAssetID(ctx, poolID, assetID)
// }

// func (k Keeper) SetAssetStatsByPoolIDAndAssetID(ctx sdk.Context, AssetStats lendtypes.PoolAssetLBMapping) {
// 	k.lend.SetAssetStatsByPoolIDAndAssetID(ctx, AssetStats)
// }

// //func (k Keeper) UpdateBorrowStats(ctx sdk.Context, pair lendtypes.Extended_Pair, borrowPos lendtypes.BorrowAsset, amount sdk.Int, inc bool) {
// //	k.lend.UpdateBorrowStats(ctx, pair, borrowPos, amount, inc)
// //}

// func (k Keeper) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, coin sdk.Coins) error {
// 	if coin.IsZero() {
// 		return auctiontypes.SendCoinsFromModuleToModuleInAuctionIsZero
// 	}

// 	return k.bank.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, coin)
// }

// func (k Keeper) UpdateReserveBalances(ctx sdk.Context, assetID uint64, moduleName string, payment sdk.Coin, inc bool) error {
// 	return k.lend.UpdateReserveBalances(ctx, assetID, moduleName, payment, inc)
// }

// func (k Keeper) SetLend(ctx sdk.Context, lend lendtypes.LendAsset) {
// 	k.lend.SetLend(ctx, lend)
// }

// func (k Keeper) BurnCoin(ctx sdk.Context, name string, coin sdk.Coin) error {
// 	if coin.IsZero() {
// 		return lendtypes.BurnCoinValueInLendIsZero
// 	}

// 	return k.bank.BurnCoins(ctx, name, sdk.NewCoins(coin))
// }

// func (k Keeper) CalculateVaultInterest(ctx sdk.Context, appID, assetID, lockerID uint64, NetBalance sdk.Int, blockHeight int64, lockerBlockTime int64) error {
// 	return k.rewards.CalculateVaultInterest(ctx, appID, assetID, lockerID, NetBalance, blockHeight, lockerBlockTime)
// }

// func (k Keeper) DeleteVaultInterestTracker(ctx sdk.Context, vault rewardstypes.VaultInterestTracker) {
// 	k.rewards.DeleteVaultInterestTracker(ctx, vault)
// }

// func (k Keeper) DutchActivator(ctx sdk.Context, lockedVault liquidationtypes.LockedVault) error {
// 	return k.auction.DutchActivator(ctx, lockedVault)
// }

// func (k Keeper) LendDutchActivator(ctx sdk.Context, lockedVault liquidationtypes.LockedVault) error {
// 	return k.auction.LendDutchActivator(ctx, lockedVault)
// }

// func (k Keeper) SetBorrow(ctx sdk.Context, borrow lendtypes.BorrowAsset) {
// 	k.lend.SetBorrow(ctx, borrow)
// }

// func (k Keeper) CalcAssetPrice(ctx sdk.Context, id uint64, amt sdk.Int) (price sdk.Dec, err error) {
// 	return k.market.CalcAssetPrice(ctx, id, amt)
// }
