package expected

import (
	"context"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctionsV2types "github.com/comdex-official/comdex/x/auctionsV2/types"
	"github.com/comdex-official/comdex/x/collector/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	liquidationsV2types "github.com/comdex-official/comdex/x/liquidationsV2/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"

	// vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type LiquidationsV2Keeper interface {
	GetLiquidationWhiteListing(ctx sdk.Context, appId uint64) (liquidationWhiteListing liquidationsV2types.LiquidationWhiteListing, found bool)
	GetLockedVault(ctx sdk.Context, appID, id uint64) (lockedVault liquidationsV2types.LockedVault, found bool)
	DeleteLockedVault(ctx sdk.Context, appID, id uint64)
	WithdrawAppReserveFundsFn(ctx sdk.Context, appId, assetId uint64, tokenQuantity sdk.Coin) error
	MsgCloseDutchAuctionForBorrow(ctx sdk.Context, liquidationData liquidationsV2types.LockedVault, auctionData auctionsV2types.Auction) error
}

type MarketKeeper interface {
	CalcAssetPrice(ctx sdk.Context, id uint64, amt sdkmath.Int) (price sdkmath.LegacyDec, err error)
	GetTwa(ctx sdk.Context, id uint64) (twa markettypes.TimeWeightedAverage, found bool)
}

type AssetKeeper interface {
	GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool)
	GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool)
	GetApps(ctx sdk.Context) (apps []assettypes.AppData, found bool)
	GetApp(ctx sdk.Context, id uint64) (app assettypes.AppData, found bool)
	GetPairsVault(ctx sdk.Context, id uint64) (pairs assettypes.ExtendedPairVault, found bool)
	GetAssetForDenom(ctx sdk.Context, denom string) (asset assettypes.Asset, found bool)
}

type EsmKeeper interface {
	GetKillSwitchData(ctx sdk.Context, appID uint64) (esmtypes.KillSwitchParams, bool)
	GetESMStatus(ctx sdk.Context, id uint64) (esmStatus esmtypes.ESMStatus, found bool)
	CalcDollarValueOfToken(ctx sdk.Context, rate uint64, amt sdkmath.Int, decimals sdkmath.Int) (price sdkmath.LegacyDec)
	SetAssetToAmount(ctx sdk.Context, assetToAmount esmtypes.AssetToAmount)
	GetDataAfterCoolOff(ctx sdk.Context, id uint64) (esmDataAfterCoolOff esmtypes.DataAfterCoolOff, found bool)
	SetDataAfterCoolOff(ctx sdk.Context, esmDataAfterCoolOff esmtypes.DataAfterCoolOff)
	GetSnapshotOfPrices(ctx sdk.Context, appID, assetID uint64) (price uint64, found bool)
}

type BankKeeper interface {
	MintCoins(ctx context.Context, name string, coins sdk.Coins) error
	BurnCoins(ctx context.Context, name string, coins sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin 
}
type VaultKeeper interface {
	GetAmountOfOtherToken(ctx sdk.Context, id1 uint64, rate1 sdkmath.LegacyDec, amt1 sdkmath.Int, id2 uint64, rate2 sdkmath.LegacyDec) (sdkmath.LegacyDec, sdkmath.Int, error)
	UpdateTokenMintedAmountLockerMapping(ctx sdk.Context, appMappingID uint64, extendedPairID uint64, amount sdkmath.Int, changeType bool)
	UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, appMappingID uint64, extendedPairID uint64, amount sdkmath.Int, changeType bool)
	CreateNewVault(ctx sdk.Context, From string, AppID uint64, ExtendedPairVaultID uint64, AmountIn sdkmath.Int, AmountOut sdkmath.Int) error
}
type CollectorKeeper interface {
	SetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, fee sdkmath.Int) error
	GetAuctionMappingForApp(ctx sdk.Context, appID, assetID uint64) (collectorAuctionLookupTable types.AppAssetIdToAuctionLookupTable, found bool)
	SetAuctionMappingForApp(ctx sdk.Context, records types.AppAssetIdToAuctionLookupTable) error
}

type TokenMintKeeper interface {
	MintNewTokensForApp(ctx sdk.Context, appMappingID uint64, assetID uint64, address string, amount sdkmath.Int) error
	BurnTokensForApp(ctx sdk.Context, appMappingID uint64, assetID uint64, amount sdkmath.Int) error
}
