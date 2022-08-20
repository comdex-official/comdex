package expected

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
)

type BankKeeper interface {
	BurnCoins(ctx sdk.Context, name string, coins sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, coins sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, address sdk.AccAddress, name string, coins sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, name string, address sdk.AccAddress, coins sdk.Coins) error

	SendCoinsFromModuleToModule(
		ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins,
	) error

	SpendableCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins
}

type AssetKeeper interface {
	GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool)
	GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool)
	GetApp(ctx sdk.Context, id uint64) (assettypes.AppData, bool)
	GetPairsVault(ctx sdk.Context, pairID uint64) (assettypes.ExtendedPairVault, bool)
}

type MarketKeeper interface {
	GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool)
}

type CollectorKeeper interface {
	UpdateCollector(ctx sdk.Context, appID, assetID uint64, CollectedStabilityFee, CollectedClosingFee, CollectedOpeningFee, LiquidationRewardsCollected sdk.Int) error
}

type EsmKeeper interface {
	GetKillSwitchData(ctx sdk.Context, appID uint64) (esmtypes.KillSwitchParams, bool)
	GetESMStatus(ctx sdk.Context, id uint64) (esmStatus esmtypes.ESMStatus, found bool)
	GetSnapshotOfPrices(ctx sdk.Context, appID, assetID uint64) (price uint64, found bool)
	GetESMTriggerParams(ctx sdk.Context, id uint64) (esmTriggerParams esmtypes.ESMTriggerParams, found bool)
}

type TokenMintKeeper interface {
	UpdateAssetDataInTokenMintByApp(ctx sdk.Context, appMappingID uint64, assetID uint64, changeType bool, amount sdk.Int)
}

type RewardsKeeper interface {
	CalculateVaultInterest(ctx sdk.Context, appID, assetID, lockerID uint64, NetBalance sdk.Int, blockHeight int64, lockerBlockTime int64) error
}
