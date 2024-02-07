package expected

import (
	"context"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
)

type BankKeeper interface {
	BurnCoins(ctx context.Context, name string, coins sdk.Coins) error
	MintCoins(ctx context.Context, name string, coins sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, address sdk.AccAddress, name string, coins sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, name string, address sdk.AccAddress, coins sdk.Coins) error

	SendCoinsFromModuleToModule(
		ctx context.Context, senderModule, recipientModule string, amt sdk.Coins,
	) error

	SpendableCoins(ctx context.Context, address sdk.AccAddress) sdk.Coins
}

type AssetKeeper interface {
	GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool)
	GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool)
	GetApp(ctx sdk.Context, id uint64) (assettypes.AppData, bool)
	GetPairsVault(ctx sdk.Context, pairID uint64) (assettypes.ExtendedPairVault, bool)
}

type MarketKeeper interface {
	CalcAssetPrice(ctx sdk.Context, id uint64, amt sdkmath.Int) (price sdkmath.LegacyDec, err error)
}

type CollectorKeeper interface {
	UpdateCollector(ctx sdk.Context, appID, assetID uint64, CollectedStabilityFee, CollectedClosingFee, CollectedOpeningFee, LiquidationRewardsCollected sdkmath.Int) error
}

type EsmKeeper interface {
	GetKillSwitchData(ctx sdk.Context, appID uint64) (esmtypes.KillSwitchParams, bool)
	GetESMStatus(ctx sdk.Context, id uint64) (esmStatus esmtypes.ESMStatus, found bool)
	GetSnapshotOfPrices(ctx sdk.Context, appID, assetID uint64) (price uint64, found bool)
	GetESMTriggerParams(ctx sdk.Context, id uint64) (esmTriggerParams esmtypes.ESMTriggerParams, found bool)
}

type TokenMintKeeper interface {
	UpdateAssetDataInTokenMintByApp(ctx sdk.Context, appMappingID uint64, assetID uint64, changeType bool, amount sdkmath.Int)
}

type RewardsKeeper interface {
	CalculateVaultInterest(ctx sdk.Context, appID, assetID, lockerID uint64, NetBalance sdkmath.Int, blockHeight int64, lockerBlockTime int64) error
	DeleteVaultInterestTracker(ctx sdk.Context, vault rewardstypes.VaultInterestTracker)
	GetExternalRewardStableVaultByApp(ctx sdk.Context, appID uint64) (VaultExternalRewards rewardstypes.StableVaultExternalRewards, found bool)
	VerifyAppIDInRewards(ctx sdk.Context, appID uint64) bool
}
