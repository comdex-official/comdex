package expected

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	"github.com/comdex-official/comdex/x/vault/types"
)

type AccountKeeper interface {
	GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI
}

type BankKeeper interface {
	BurnCoins(ctx sdk.Context, name string, coins sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, coins sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoinsFromModuleToAccount(ctx sdk.Context, name string, address sdk.AccAddress, coins sdk.Coins) error
}

type AssetKeeper interface {
	GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool)
	GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool)
	GetApps(ctx sdk.Context) (apps []assettypes.AppData, found bool)
	GetPairsVault(ctx sdk.Context, id uint64) (pairs assettypes.ExtendedPairVault, found bool)
}

type VaultKeeper interface {
	GetAppMappingData(ctx sdk.Context, appMappingID uint64) (appExtendedPairVaultData []types.AppExtendedPairVaultMappingData, found bool)
	CalculateCollateralizationRatio(ctx sdk.Context, extendedPairVaultID uint64, amountIn sdk.Int, amountOut sdk.Int) (sdk.Dec, error)
	GetVault(ctx sdk.Context, id uint64) (vault types.Vault, found bool)
	GetVaults(ctx sdk.Context) (vaults []types.Vault)
	GetIDForVault(ctx sdk.Context) uint64
	DeleteVault(ctx sdk.Context, id uint64)
	GetLengthOfVault(ctx sdk.Context) uint64
	SetLengthOfVault(ctx sdk.Context, length uint64)
	UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx sdk.Context, vaultData types.Vault)
	UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, appMappingID uint64, extendedPairID uint64, amount sdk.Int, changeType bool)
	UpdateTokenMintedAmountLockerMapping(ctx sdk.Context, appMappingID uint64, extendedPairID uint64, amount sdk.Int, changeType bool)
	DeleteUserVaultExtendedPairMapping(ctx sdk.Context, address string, appID uint64, pairVaultID uint64)
	DeleteAddressFromAppExtendedPairVaultMapping(ctx sdk.Context, extendedPairID uint64, userVaultID uint64, appMappingID uint64)
	SetVault(ctx sdk.Context, vault types.Vault)
}

type MarketKeeper interface {
	CalcAssetPrice(ctx sdk.Context, id uint64, amt sdk.Int) (price sdk.Dec, err error)
	GetTwa(ctx sdk.Context, id uint64) (twa markettypes.TimeWeightedAverage, found bool)
}

type AuctionKeeper interface {
	GetParams(ctx sdk.Context) auctiontypes.Params
	DutchActivator(ctx sdk.Context, lockedVault liquidationtypes.LockedVault) error
	LendDutchActivator(ctx sdk.Context, lockedVault liquidationtypes.LockedVault) error
}

type EsmKeeper interface {
	GetKillSwitchData(ctx sdk.Context, appID uint64) (esmtypes.KillSwitchParams, bool)
	GetESMStatus(ctx sdk.Context, id uint64) (esmStatus esmtypes.ESMStatus, found bool)
}

type LendKeeper interface {
	GetBorrows(ctx sdk.Context) (userBorrows []uint64, found bool)
	GetBorrow(ctx sdk.Context, id uint64) (borrow lendtypes.BorrowAsset, found bool)
	GetLendPair(ctx sdk.Context, id uint64) (pair lendtypes.Extended_Pair, found bool)
	GetAssetRatesParams(ctx sdk.Context, assetID uint64) (assetRatesStats lendtypes.AssetRatesParams, found bool)
	VerifyCollateralizationRatio(ctx sdk.Context, amountIn sdk.Int, assetIn assettypes.Asset, amountOut sdk.Int, assetOut assettypes.Asset, liquidationThreshold sdk.Dec) error
	CalculateCollateralizationRatio(ctx sdk.Context, amountIn sdk.Int, assetIn assettypes.Asset, amountOut sdk.Int, assetOut assettypes.Asset) (sdk.Dec, error)
	GetLend(ctx sdk.Context, id uint64) (lend lendtypes.LendAsset, found bool)
	CreteNewBorrow(ctx sdk.Context, liqBorrow liquidationtypes.LockedVault)
	GetPool(ctx sdk.Context, id uint64) (pool lendtypes.Pool, found bool)

	GetAssetStatsByPoolIDAndAssetID(ctx sdk.Context, assetID, poolID uint64) (AssetStats lendtypes.PoolAssetLBMapping, found bool)
	SetAssetStatsByPoolIDAndAssetID(ctx sdk.Context, AssetStats lendtypes.PoolAssetLBMapping)
	UpdateReserveBalances(ctx sdk.Context, assetID uint64, moduleName string, payment sdk.Coin, inc bool) error
	SetLend(ctx sdk.Context, lend lendtypes.LendAsset)
	SetBorrow(ctx sdk.Context, borrow lendtypes.BorrowAsset)
	CalculateBorrowInterestForLiquidation(ctx sdk.Context, borrowID uint64) (lendtypes.BorrowAsset, error)
	ReBalanceStableRates(ctx sdk.Context, borrowPos lendtypes.BorrowAsset) (lendtypes.BorrowAsset, error)
	DeleteBorrow(ctx sdk.Context, ID uint64)
}

type RewardsKeeper interface {
	CalculateVaultInterest(ctx sdk.Context, appID, assetID, lockerID uint64, NetBalance sdk.Int, blockHeight int64, lockerBlockTime int64) error
	DeleteVaultInterestTracker(ctx sdk.Context, vault rewardstypes.VaultInterestTracker)
}
