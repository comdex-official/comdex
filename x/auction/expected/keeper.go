package expected

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/collector/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type AccountKeeper interface {
	GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI
}

type BankKeeper interface {
	MintCoins(ctx sdk.Context, name string, coins sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, coins sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type MarketKeeper interface {
	GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool)
}

type LiquidationKeeper interface {
	SetFlagIsAuctionInProgress(ctx sdk.Context, id uint64, flag bool) error
	SetFlagIsAuctionComplete(ctx sdk.Context, id uint64, flag bool) error
	GetLockedVaults(ctx sdk.Context) (locked_vaults []liquidationtypes.LockedVault)
	//UpdateAssetQuantitiesInLockedVault(ctx sdk.Context, collateral_auction auctiontypes.CollateralAuction, amountIn sdk.Int, assetIn assettypes.Asset, amountOut sdk.Int, assetOut assettypes.Asset) error
}

type AssetKeeper interface {
	GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool)
	GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool)
	GetApps(ctx sdk.Context) (apps []assettypes.AppMapping, found bool)
}

type VaultKeeper interface {
	//CalculateCollaterlizationRatio(
	//	ctx sdk.Context,
	//	amountIn sdk.Int,
	//	assetIn assettypes.Asset,
	//	amountOut sdk.Int,
	//	assetOut assettypes.Asset,
	//) (sdk.Dec, error)
	//BurnCAssets(ctx sdk.Context, moduleName string, collateralDenom string, denom string, amount sdk.Int) error
}

type CollectorKeeper interface {
	GetAppidToAssetCollectorMapping(ctx sdk.Context, app_id uint64) (appAssetCollectorData types.AppIdToAssetCollectorMapping, found bool)
	UpdateCollector(ctx sdk.Context, appId, asset_id uint64, CollectedStabilityFee, CollectedClosingFee, CollectedOpeningFee, LiquidationRewardsCollected sdk.Int) error
	SetCollectorLookupTable(ctx sdk.Context, records ...types.CollectorLookupTable) error
	GetCollectorLookupTable(ctx sdk.Context, app_id uint64) (collectorLookup types.CollectorLookup, found bool)
	SetCollectorAuctionLookupTable(ctx sdk.Context, records ...types.CollectorAuctionLookupTable) error
	GetCollectorAuctionLookupTable(ctx sdk.Context, app_id uint64) (appIdToAuctionData types.CollectorAuctionLookupTable, found bool)
	GetNetFeeCollectedData(ctx sdk.Context, app_id uint64) (netFeeData types.NetFeeCollectedData, found bool)
	GetAmountFromCollector(ctx sdk.Context, appId, asset_id uint64, amount sdk.Int) (sdk.Int, error)
	SetNetFeeCollectedData(ctx sdk.Context, app_id, asset_id uint64, fee sdk.Int) error
}

type TokenMintKeeper interface {
	MintNewTokensForApp(ctx sdk.Context, appMappingId uint64, assetId uint64, address string, amount sdk.Int) error
	BurnTokensForApp(ctx sdk.Context, appMappingId uint64, assetId uint64, amount sdk.Int) error
}
