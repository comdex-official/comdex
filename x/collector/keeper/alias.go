package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
)

func (k Keeper) HasAssetForDenom(ctx sdk.Context, id string) bool {
	return k.asset.HasAssetForDenom(ctx, id)
}

func (k Keeper) HasAsset(ctx sdk.Context, id uint64) bool {
	return k.asset.HasAsset(ctx, id)
}
func (k Keeper) GetAssetForDenom(ctx sdk.Context, id string) (types.Asset, bool) {
	return k.asset.GetAssetForDenom(ctx, id)
}

func (k Keeper) GetApp(ctx sdk.Context, id uint64) (types.AppData, bool) {
	return k.asset.GetApp(ctx, id)
}

func (k Keeper) GetAsset(ctx sdk.Context, id uint64) (types.Asset, bool) {
	return k.asset.GetAsset(ctx, id)
}

func (k Keeper) SendCoinFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, coin sdk.Coins) error {
	if coin.IsZero() {
		return collectortypes.SendCoinFromModuleToModuleIsZero
	}
	return k.bank.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, coin)
}

func (k Keeper) SendCoinsFromModuleToAccount(ctx sdk.Context, name string, address sdk.AccAddress, coins sdk.Coins) error {
	if coins.IsZero() {
		return collectortypes.SendCoinFromModuleToModuleIsZero
	}
	return k.bank.SendCoinsFromModuleToAccount(ctx, name, address, coins)
}

func (k Keeper) GetMintGenesisTokenData(ctx sdk.Context, appID, assetID uint64) (mintData types.MintGenesisToken, found bool) {
	return k.asset.GetMintGenesisTokenData(ctx, appID, assetID)
}

func (k Keeper) GetAuctionParams(ctx sdk.Context, AppID uint64) (asset auctiontypes.AuctionParams, found bool) {
	return k.auction.GetAuctionParams(ctx, AppID)
}

func (k Keeper) GetLockerLookupTable(ctx sdk.Context, appID, assetID uint64) (lockerLookupData lockertypes.LockerLookupTableData, found bool) {
	return k.locker.GetLockerLookupTable(ctx, appID, assetID)
}

func (k Keeper) GetReward(ctx sdk.Context, appId, assetID uint64) (rewards rewardstypes.InternalRewards, found bool) {
	return k.rewards.GetReward(ctx, appId, assetID)
}

func (k Keeper) GetLocker(ctx sdk.Context, lockerID uint64) (locker lockertypes.Locker, found bool) {
	return k.locker.GetLocker(ctx, lockerID)
}

func (k Keeper) CalculationOfRewards(ctx sdk.Context, amount sdk.Int, lsr sdk.Dec, bTime int64) (sdk.Dec, error) {
	return k.rewards.CalculationOfRewards(ctx, amount, lsr, bTime)
}

func (k Keeper) SetLocker(ctx sdk.Context, locker lockertypes.Locker) {
	k.locker.SetLocker(ctx, locker)
}

func (k Keeper) SetLockerLookupTable(ctx sdk.Context, lockerLookupData lockertypes.LockerLookupTableData) {
	k.locker.SetLockerLookupTable(ctx, lockerLookupData)
}

func (k Keeper) SetLockerRewardTracker(ctx sdk.Context, rewards rewardstypes.LockerRewardsTracker) {
	k.rewards.SetLockerRewardTracker(ctx, rewards)
}

func (k Keeper) GetLockerRewardTracker(ctx sdk.Context, id, appID uint64) (rewards rewardstypes.LockerRewardsTracker, found bool) {
	return k.rewards.GetLockerRewardTracker(ctx, id, appID)
}

func (k Keeper) SetLockerTotalRewardsByAssetAppWise(ctx sdk.Context, lockerRewardsMapping lockertypes.LockerTotalRewardsByAssetAppWise) error {
	return k.locker.SetLockerTotalRewardsByAssetAppWise(ctx, lockerRewardsMapping)
}
func (k Keeper) GetLockerTotalRewardsByAssetAppWise(ctx sdk.Context, appID, assetID uint64) (lockerRewardsMapping lockertypes.LockerTotalRewardsByAssetAppWise, found bool) {
	return k.locker.GetLockerTotalRewardsByAssetAppWise(ctx, appID, assetID)
}
