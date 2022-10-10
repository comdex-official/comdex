package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	"github.com/comdex-official/comdex/x/vault/types"
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

func (k Keeper) GetAppMappingData(ctx sdk.Context, appMappingID uint64) (appExtendedPairVaultData []types.AppExtendedPairVaultMappingData, found bool) {
	return k.vault.GetAppMappingData(ctx, appMappingID)
}

func (k Keeper) CalculateCollaterlizationRatio(ctx sdk.Context, extendedPairVaultID uint64, amountIn sdk.Int, amountOut sdk.Int) (sdk.Dec, error) {
	return k.vault.CalculateCollaterlizationRatio(ctx, extendedPairVaultID, amountIn, amountOut)
}

func (k Keeper) GetVault(ctx sdk.Context, id uint64) (vault types.Vault, found bool) {
	return k.vault.GetVault(ctx, id)
}

func (k Keeper) DeleteVault(ctx sdk.Context, id uint64) {
	k.vault.DeleteVault(ctx, id)
}

func (k Keeper) UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx sdk.Context, vaultData types.Vault) {
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

func (k Keeper) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, coin sdk.Coins) error {
	if coin.IsZero() {
		return auctiontypes.SendCoinsFromModuleToModuleInAuctionIsZero
	}

	return k.bank.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, coin)
}

func (k Keeper) CalculateVaultInterest(ctx sdk.Context, appID, assetID, lockerID uint64, NetBalance sdk.Int, blockHeight int64, lockerBlockTime int64) error {
	return k.rewards.CalculateVaultInterest(ctx, appID, assetID, lockerID, NetBalance, blockHeight, lockerBlockTime)
}

func (k Keeper) DeleteVaultInterestTracker(ctx sdk.Context, vault rewardstypes.VaultInterestTracker) {
	k.rewards.DeleteVaultInterestTracker(ctx, vault)
}

func (k Keeper) DutchActivator(ctx sdk.Context, lockedVault liquidationtypes.LockedVault) error {
	return k.auction.DutchActivator(ctx, lockedVault)
}
