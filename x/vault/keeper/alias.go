package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
)

func (k Keeper) BurnCoin(ctx sdk.Context, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return vaulttypes.BurnCoinValueInVaultIsZero
	}

	return k.bank.BurnCoins(ctx, name, sdk.NewCoins(coin))
}

func (k Keeper) MintCoin(ctx sdk.Context, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return vaulttypes.MintCoinValueInVaultIsZero
	}

	return k.bank.MintCoins(ctx, name, sdk.NewCoins(coin))
}

func (k Keeper) SendCoinFromAccountToModule(ctx sdk.Context, address sdk.AccAddress, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return vaulttypes.SendCoinsFromAccountToModuleInVaultIsZero
	}

	return k.bank.SendCoinsFromAccountToModule(ctx, address, name, sdk.NewCoins(coin))
}

func (k Keeper) SendCoinFromModuleToAccount(ctx sdk.Context, name string, address sdk.AccAddress, coin sdk.Coin) error {
	if coin.IsZero() {
		return vaulttypes.SendCoinsFromModuleToAccountInVaultIsZero
	}

	return k.bank.SendCoinsFromModuleToAccount(ctx, name, address, sdk.NewCoins(coin))
}

func (k Keeper) SendCoinFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, coin sdk.Coins) error {
	if coin.IsZero() {
		return vaulttypes.SendCoinsFromModuleToModuleInVaultIsZero
	}
	return k.bank.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, coin)
}

func (k Keeper) SpendableCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins {
	return k.bank.SpendableCoins(ctx, address)
}

func (k Keeper) GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool) {
	return k.asset.GetAsset(ctx, id)
}

func (k Keeper) GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool) {
	return k.asset.GetPair(ctx, id)
}

func (k Keeper) GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool) {
	return k.oracle.GetPriceForAsset(ctx, id)
}

func (k Keeper) GetApp(ctx sdk.Context, id uint64) (assettypes.AppData, bool) {
	return k.asset.GetApp(ctx, id)
}

func (k Keeper) GetPairsVault(ctx sdk.Context, pairID uint64) (assettypes.ExtendedPairVault, bool) {
	return k.asset.GetPairsVault(ctx, pairID)
}

func (k Keeper) UpdateCollector(ctx sdk.Context, appID, assetID uint64, CollectedStabilityFee, CollectedClosingFee, CollectedOpeningFee, LiquidationRewardsCollected sdk.Int) error {
	return k.collector.UpdateCollector(ctx, appID, assetID, CollectedStabilityFee, CollectedClosingFee, CollectedOpeningFee, LiquidationRewardsCollected)
}

func (k Keeper) GetKillSwitchData(ctx sdk.Context, appID uint64) (esmtypes.KillSwitchParams, bool) {
	return k.esm.GetKillSwitchData(ctx, appID)
}

func (k Keeper) GetESMStatus(ctx sdk.Context, id uint64) (esmtypes.ESMStatus, bool) {
	return k.esm.GetESMStatus(ctx, id)
}

func (k Keeper) GetSnapshotOfPrices(ctx sdk.Context, appID, assetID uint64) (price uint64, found bool) {
	return k.esm.GetSnapshotOfPrices(ctx, appID, assetID)
}

func (k Keeper) GetESMTriggerParams(ctx sdk.Context, id uint64) (esmTriggerParams esmtypes.ESMTriggerParams, found bool) {
	return k.esm.GetESMTriggerParams(ctx, id)
}

func (k Keeper) UpdateAssetDataInTokenMintByApp(ctx sdk.Context, appMappingID uint64, assetID uint64, changeType bool, amount sdk.Int) {
	k.tokenmint.UpdateAssetDataInTokenMintByApp(ctx, appMappingID, assetID, changeType, amount)
}

func (k Keeper) CalculateVaultInterest(ctx sdk.Context, appID, assetID, lockerID uint64, NetBalance sdk.Int, blockHeight int64, lockerBlockTime int64) error {
	return k.rewards.CalculateVaultInterest(ctx, appID, assetID, lockerID, NetBalance, blockHeight, lockerBlockTime)
}

func (k Keeper) DeleteVaultInterestTracker(ctx sdk.Context, vault rewardstypes.VaultInterestTracker) {
	k.rewards.DeleteVaultInterestTracker(ctx, vault)
}

func (k Keeper) CalcAssetPrice(ctx sdk.Context, id uint64, amt sdk.Int) (price sdk.Int, err error) {
	return k.oracle.CalcAssetPrice(ctx, id, amt)
}
