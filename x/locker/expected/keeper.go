package expected

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/petrichormoney/petri/x/asset/types"
	collectortypes "github.com/petrichormoney/petri/x/collector/types"
	esmtypes "github.com/petrichormoney/petri/x/esm/types"
	rewardstypes "github.com/petrichormoney/petri/x/rewards/types"
)

type BankKeeper interface {
	BurnCoins(ctx sdk.Context, name string, coins sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, coins sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, address sdk.AccAddress, name string, coins sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, name string, address sdk.AccAddress, coins sdk.Coins) error
	SpendableCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins
	SendCoinsFromModuleToModule(
		ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins,
	) error
}

type AssetKeeper interface {
	GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool)
	GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool)
	GetApp(ctx sdk.Context, id uint64) (assettypes.AppData, bool)
	GetApps(ctx sdk.Context) (apps []assettypes.AppData, found bool)
}

type CollectorKeeper interface {
	UpdateCollector(ctx sdk.Context, appID, assetID uint64, collectedStabilityFee, collectedClosingFee, collectedOpeningFee, liquidationRewardsCollected sdk.Int) error
	GetCollectorLookupTable(ctx sdk.Context, appID, assetID uint64) (collectorLookup collectortypes.CollectorLookupTableData, found bool)
}

type EsmKeeper interface {
	GetKillSwitchData(ctx sdk.Context, appID uint64) (esmtypes.KillSwitchParams, bool)
	GetESMStatus(ctx sdk.Context, id uint64) (esmStatus esmtypes.ESMStatus, found bool)
}

type RewardsKeeper interface {
	CalculateLockerRewards(ctx sdk.Context, appID, assetID, lockerID uint64, Depositor string, NetBalance sdk.Int, blockHeight int64, lockerBlockTime int64) error
	DeleteLockerRewardTracker(ctx sdk.Context, rewards rewardstypes.LockerRewardsTracker)
}
