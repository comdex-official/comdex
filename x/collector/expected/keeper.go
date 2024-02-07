package expected

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
	"github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
)

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx context.Context, address sdk.AccAddress, name string, coins sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, name string, address sdk.AccAddress, coins sdk.Coins) error

	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin 
}

type AssetKeeper interface {
	HasAssetForDenom(ctx sdk.Context, id string) bool
	HasAsset(ctx sdk.Context, id uint64) bool
	GetAssetForDenom(ctx sdk.Context, denom string) (types.Asset, bool)
	GetApp(ctx sdk.Context, id uint64) (types.AppData, bool)
	GetAsset(ctx sdk.Context, id uint64) (types.Asset, bool)
	GetMintGenesisTokenData(ctx sdk.Context, appID, assetID uint64) (mintData types.MintGenesisToken, found bool)
}

type AuctionKeeper interface {
	GetAuctionParams(ctx sdk.Context, appID uint64) (asset auctiontypes.AuctionParams, found bool)
}

type LockerKeeper interface {
	GetLockerLookupTable(ctx sdk.Context, appID, assetID uint64) (lockerLookupData lockertypes.LockerLookupTableData, found bool)
	GetLocker(ctx sdk.Context, lockerID uint64) (locker lockertypes.Locker, found bool)
	SetLocker(ctx sdk.Context, locker lockertypes.Locker)
	SetLockerLookupTable(ctx sdk.Context, lockerLookupData lockertypes.LockerLookupTableData)
	SetLockerTotalRewardsByAssetAppWise(ctx sdk.Context, lockerRewardsMapping lockertypes.LockerTotalRewardsByAssetAppWise) error
	GetLockerTotalRewardsByAssetAppWise(ctx sdk.Context, appID, assetID uint64) (lockerRewardsMapping lockertypes.LockerTotalRewardsByAssetAppWise, found bool)
}

type RewardsKeeper interface {
	GetReward(ctx sdk.Context, appID, assetID uint64) (rewards rewardstypes.InternalRewards, found bool)
	CalculationOfRewards(ctx sdk.Context, amount sdkmath.Int, lsr sdkmath.LegacyDec, bTime int64) (sdkmath.LegacyDec, error)
	SetLockerRewardTracker(ctx sdk.Context, rewards rewardstypes.LockerRewardsTracker)
	GetLockerRewardTracker(ctx sdk.Context, id, appID uint64) (rewards rewardstypes.LockerRewardsTracker, found bool)
}
