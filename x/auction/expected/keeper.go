package expected

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	"github.com/comdex-official/comdex/x/collector/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
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
	CalcAssetPrice(ctx sdk.Context, id uint64, amt sdk.Int) (price sdk.Dec, err error)
	GetTwa(ctx sdk.Context, id uint64) (twa markettypes.TimeWeightedAverage, found bool)
}

type LiquidationKeeper interface {
	SetFlagIsAuctionInProgress(ctx sdk.Context, appID, id uint64, flag bool) error
	SetFlagIsAuctionComplete(ctx sdk.Context, appID, id uint64, flag bool) error
	GetLockedVaults(ctx sdk.Context) (lockedVaults []liquidationtypes.LockedVault)
	GetLockedVault(ctx sdk.Context, appID, id uint64) (lockedVault liquidationtypes.LockedVault, found bool)
	SetLockedVault(ctx sdk.Context, lockedVault liquidationtypes.LockedVault)
	DeleteLockedVault(ctx sdk.Context, appID, id uint64)
	CreateLockedVaultHistory(ctx sdk.Context, lockedVault liquidationtypes.LockedVault) error
	UnLiquidateLockedBorrows(ctx sdk.Context, appID, id uint64, dutchAuction auctiontypes.DutchAuction) error
}

type AssetKeeper interface {
	GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool)
	GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool)
	GetApps(ctx sdk.Context) (apps []assettypes.AppData, found bool)
	GetApp(ctx sdk.Context, id uint64) (app assettypes.AppData, found bool)
	GetPairsVault(ctx sdk.Context, id uint64) (pairs assettypes.ExtendedPairVault, found bool)
}

type VaultKeeper interface {
	GetAppExtendedPairVaultMappingData(ctx sdk.Context, appMappingID uint64, pairVaultsID uint64) (appExtendedPairVaultData vaulttypes.AppExtendedPairVaultMappingData, found bool)
	SetAppExtendedPairVaultMappingData(ctx sdk.Context, appExtendedPairVaultData vaulttypes.AppExtendedPairVaultMappingData)
	UpdateTokenMintedAmountLockerMapping(ctx sdk.Context, appMappingID uint64, extendedPairID uint64, amount sdk.Int, changeType bool)
	UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, appMappingID uint64, extendedPairID uint64, amount sdk.Int, changeType bool)
	DeleteUserVaultExtendedPairMapping(ctx sdk.Context, from string, appMapping uint64, extendedPairVault uint64)
	CreateNewVault(ctx sdk.Context, From string, AppID uint64, ExtendedPairVaultID uint64, AmountIn sdk.Int, AmountOut sdk.Int) error
	GetUserAppExtendedPairMappingData(ctx sdk.Context, from string, appMapping uint64, extendedPairVault uint64) (userVaultAssetData vaulttypes.OwnerAppExtendedPairVaultMappingData, found bool)
	GetUserAppMappingData(ctx sdk.Context, from string, appMapping uint64) (userVaultAssetData []vaulttypes.OwnerAppExtendedPairVaultMappingData, found bool)
	// CheckUserAppToExtendedPairMapping(ctx sdk.Context, userVaultAssetData vaulttypes.UserVaultAssetMapping, extendedPairVaultID uint64, appMappingID uint64) (vaultID uint64, found bool)
	SetVault(ctx sdk.Context, vault vaulttypes.Vault)
	GetVault(ctx sdk.Context, id uint64) (vault vaulttypes.Vault, found bool)
	GetAmountOfOtherToken(ctx sdk.Context, id1 uint64, rate1 sdk.Dec, amt1 sdk.Int, id2 uint64, rate2 sdk.Dec) (sdk.Dec, sdk.Int, error)
}

type CollectorKeeper interface {
	GetAppidToAssetCollectorMapping(ctx sdk.Context, appID, assetID uint64) (appAssetCollectorData types.AppToAssetIdCollectorMapping, found bool)
	UpdateCollector(ctx sdk.Context, appID, assetID uint64, CollectedStabilityFee, CollectedClosingFee, CollectedOpeningFee, LiquidationRewardsCollected sdk.Int) error
	// SetCollectorLookupTable(ctx sdk.Context, records ...types.CollectorLookupTable) error
	GetCollectorLookupTable(ctx sdk.Context, appID, assetID uint64) (collectorLookup types.CollectorLookupTableData, found bool)
	GetAuctionMappingForApp(ctx sdk.Context, appID, assetID uint64) (collectorAuctionLookupTable types.AppAssetIdToAuctionLookupTable, found bool)
	GetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64) (netFeeData types.AppAssetIdToFeeCollectedData, found bool)
	GetAmountFromCollector(ctx sdk.Context, appID, assetID uint64, amount sdk.Int) (sdk.Int, error)
	SetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, fee sdk.Int) error
	SetAuctionMappingForApp(ctx sdk.Context, records types.AppAssetIdToAuctionLookupTable) error
	GetAllAuctionMappingForApp(ctx sdk.Context) (collectorAuctionLookupTable []types.AppAssetIdToAuctionLookupTable, found bool)
}

type TokenMintKeeper interface {
	MintNewTokensForApp(ctx sdk.Context, appMappingID uint64, assetID uint64, address string, amount sdk.Int) error
	BurnTokensForApp(ctx sdk.Context, appMappingID uint64, assetID uint64, amount sdk.Int) error
}

type EsmKeeper interface {
	GetKillSwitchData(ctx sdk.Context, appID uint64) (esmtypes.KillSwitchParams, bool)
	GetESMStatus(ctx sdk.Context, id uint64) (esmStatus esmtypes.ESMStatus, found bool)
}

type LendKeeper interface {
	GetBorrows(ctx sdk.Context) (userBorrows lendtypes.BorrowMapping, found bool)
	GetBorrow(ctx sdk.Context, id uint64) (borrow lendtypes.BorrowAsset, found bool)
	GetLendPair(ctx sdk.Context, id uint64) (pair lendtypes.Extended_Pair, found bool)
	GetAssetRatesStats(ctx sdk.Context, assetID uint64) (assetRatesStats lendtypes.AssetRatesStats, found bool)
	VerifyCollateralizationRatio(ctx sdk.Context, amountIn sdk.Int, assetIn assettypes.Asset, amountOut sdk.Int, assetOut assettypes.Asset, liquidationThreshold sdk.Dec) error
	CalculateCollateralizationRatio(ctx sdk.Context, amountIn sdk.Int, assetIn assettypes.Asset, amountOut sdk.Int, assetOut assettypes.Asset) (sdk.Dec, error)
	GetLend(ctx sdk.Context, id uint64) (lend lendtypes.LendAsset, found bool)
	DeleteBorrow(ctx sdk.Context, id uint64)

	DeleteBorrowForAddressByPair(ctx sdk.Context, address sdk.AccAddress, pairID uint64)
	UpdateUserBorrowIDMapping(ctx sdk.Context, lendOwner string, borrowID uint64, isInsert bool) error
	UpdateBorrowIDByOwnerAndPoolMapping(ctx sdk.Context, borrowOwner string, borrowID uint64, poolID uint64, isInsert bool) error
	UpdateBorrowIdsMapping(ctx sdk.Context, borrowID uint64, isInsert bool) error
	CreteNewBorrow(ctx sdk.Context, liqBorrow liquidationtypes.LockedVault)
	GetPool(ctx sdk.Context, id uint64) (pool lendtypes.Pool, found bool)
	GetAddAuctionParamsData(ctx sdk.Context, appID uint64) (auctionParams lendtypes.AuctionParams, found bool)
	GetReserveDepositStats(ctx sdk.Context) (depositStats lendtypes.DepositStats, found bool)
	ModuleBalance(ctx sdk.Context, moduleName string, denom string) sdk.Int
	UpdateReserveBalances(ctx sdk.Context, assetID uint64, moduleName string, payment sdk.Coin, inc bool) error
	SetLend(ctx sdk.Context, lend lendtypes.LendAsset)
}
