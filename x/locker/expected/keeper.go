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

type OracleKeeper interface {
	GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool)
}

type CollectorKeeper interface {
	UpdateCollector(ctx sdk.Context, appID, assetID uint64, collectedStabilityFee, collectedClosingFee, collectedOpeningFee, liquidationRewardsCollected sdk.Int) error
}

type EsmKeeper interface {
	GetKillSwitchData(ctx sdk.Context, app_id uint64) (esmtypes.KillSwitchParams, bool)
}