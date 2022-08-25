package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
)

func (k Keeper) SendCoinFromAccountToModule(ctx sdk.Context, address sdk.AccAddress, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return lockertypes.SendCoinsFromAccountToModuleInLockerIsZero
	}

	return k.bank.SendCoinsFromAccountToModule(ctx, address, name, sdk.NewCoins(coin))
}

func (k Keeper) SendCoinFromModuleToAccount(ctx sdk.Context, name string, address sdk.AccAddress, coin sdk.Coin) error {
	if coin.IsZero() {
		return lockertypes.SendCoinsFromModuleToAccountInLockerIsZero
	}

	return k.bank.SendCoinsFromModuleToAccount(ctx, name, address, sdk.NewCoins(coin))
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

func (k Keeper) GetApps(ctx sdk.Context) (apps []assettypes.AppData, found bool) {
	return k.asset.GetApps(ctx)
}

func (k Keeper) UpdateCollector(ctx sdk.Context, appID, assetID uint64, collectedStabilityFee, collectedClosingFee, collectedOpeningFee, liquidationRewardsCollected sdk.Int) error {
	return k.collector.UpdateCollector(ctx, appID, assetID, collectedStabilityFee, collectedClosingFee, collectedOpeningFee, liquidationRewardsCollected)
}

func (k Keeper) GetCollectorLookupTable(ctx sdk.Context, appID, assetID uint64) (collectorLookup collectortypes.CollectorLookupTableData, found bool) {
	return k.collector.GetCollectorLookupTable(ctx, appID, assetID)
}

func (k Keeper) SendCoinFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, coin sdk.Coins) error {
	if coin.IsZero() {
		return lockertypes.SendCoinsFromModuleToModuleInLockerIsZero
	}
	return k.bank.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, coin)
}

func (k Keeper) GetKillSwitchData(ctx sdk.Context, appID uint64) (esmtypes.KillSwitchParams, bool) {
	return k.esm.GetKillSwitchData(ctx, appID)
}

func (k Keeper) GetESMStatus(ctx sdk.Context, id uint64) (esmtypes.ESMStatus, bool) {
	return k.esm.GetESMStatus(ctx, id)
}

func (k Keeper) CalculateLockerRewards(ctx sdk.Context, appID, assetID, lockerID uint64, Depositor string, NetBalance sdk.Int, blockHeight int64, lockerBlockTime int64) error {
	return k.rewards.CalculateLockerRewards(ctx, appID, assetID, lockerID, Depositor, NetBalance, blockHeight, lockerBlockTime)
}

func (k Keeper) DeleteLockerRewardTracker(ctx sdk.Context, rewards rewardstypes.LockerRewardsTracker) {
	k.rewards.DeleteLockerRewardTracker(ctx, rewards)
}
