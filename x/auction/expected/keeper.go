package expected

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/collector/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type AccountKeeper interface {
	GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI
	GetModuleAddress(name string) sdk.AccAddress
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
	GetLockedVaults(ctx sdk.Context) (lockedVaults []liquidationtypes.LockedVault)
	GetLockedVault(ctx sdk.Context, id uint64) (lockedVault liquidationtypes.LockedVault, found bool)
	SetLockedVault(ctx sdk.Context, lockedVault liquidationtypes.LockedVault)
	DeleteLockedVault(ctx sdk.Context, id uint64)
	CreateLockedVaultHistory(ctx sdk.Context, lockedVault liquidationtypes.LockedVault) error
	//UpdateAssetQuantitiesInLockedVault(ctx sdk.Context, collateral_auction auctiontypes.CollateralAuction, amountIn sdk.Int, assetIn assettypes.Asset, amountOut sdk.Int, assetOut assettypes.Asset) error
}

type AssetKeeper interface {
	GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool)
	GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool)
	GetApps(ctx sdk.Context) (apps []assettypes.AppData, found bool)
	GetApp(ctx sdk.Context, id uint64) (app assettypes.AppData, found bool)
	GetPairsVault(ctx sdk.Context, id uint64) (pairs assettypes.ExtendedPairVault, found bool)
}

type VaultKeeper interface {
	GetAppExtendedPairVaultMapping(ctx sdk.Context, appMappingID uint64) (appExtendedPairVaultData vaulttypes.AppExtendedPairVaultMapping, found bool)
	SetAppExtendedPairVaultMapping(ctx sdk.Context, appExtendedPairVaultData vaulttypes.AppExtendedPairVaultMapping) error
	UpdateTokenMintedAmountLockerMapping(ctx sdk.Context, vaultLookupData vaulttypes.AppExtendedPairVaultMapping, extendedPairID uint64, amount sdk.Int, changeType bool)
	UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, vaultLookupData vaulttypes.AppExtendedPairVaultMapping, extendedPairID uint64, amount sdk.Int, changeType bool)
	UpdateUserVaultExtendedPairMapping(ctx sdk.Context, extendedPairID uint64, userAddress string, appMappingID uint64)
}

type CollectorKeeper interface {
	GetAppidToAssetCollectorMapping(ctx sdk.Context, appID uint64) (appAssetCollectorData types.AppIdToAssetCollectorMapping, found bool)
	UpdateCollector(ctx sdk.Context, appID, assetID uint64, CollectedStabilityFee, CollectedClosingFee, CollectedOpeningFee, LiquidationRewardsCollected sdk.Int) error
	SetCollectorLookupTable(ctx sdk.Context, records ...types.CollectorLookupTable) error
	GetCollectorLookupTable(ctx sdk.Context, appID uint64) (collectorLookup types.CollectorLookup, found bool)
	GetAuctionMappingForApp(ctx sdk.Context, appID uint64) (collectorAuctionLookupTable types.CollectorAuctionLookupTable, found bool)
	GetNetFeeCollectedData(ctx sdk.Context, appID uint64) (netFeeData types.NetFeeCollectedData, found bool)
	GetAmountFromCollector(ctx sdk.Context, appID, assetID uint64, amount sdk.Int) (sdk.Int, error)
	SetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, fee sdk.Int) error
	SetAuctionMappingForApp(ctx sdk.Context, records ...types.CollectorAuctionLookupTable) error
	GetAllAuctionMappingForApp(ctx sdk.Context) (collectorAuctionLookupTable []types.CollectorAuctionLookupTable, found bool)
}

type TokenMintKeeper interface {
	MintNewTokensForApp(ctx sdk.Context, appMappingID uint64, assetID uint64, address string, amount sdk.Int) error
	BurnTokensForApp(ctx sdk.Context, appMappingID uint64, assetID uint64, amount sdk.Int) error
}
