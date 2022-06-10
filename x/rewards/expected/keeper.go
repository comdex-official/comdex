package expected

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	collecortypes "github.com/comdex-official/comdex/x/collector/types"
	liquiditytypes "github.com/comdex-official/comdex/x/liquidity/types"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	"github.com/comdex-official/comdex/x/rewards/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias).
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	// Methods imported from account should be defined here
}

type LiquidityKeeper interface {
	GetPair(ctx sdk.Context, id uint64) (pair liquiditytypes.Pair, found bool)
	GetPool(ctx sdk.Context, id uint64) (pool liquiditytypes.Pool, found bool)
	GetFarmingRewardsData(ctx sdk.Context, coinToDistribute sdk.Coin, liquidityGaugeData types.LiquidtyGaugeMetaData) ([]types.RewardDistributionDataCollector, error)
	TransferFundsForSwapFeeDistribution(ctx sdk.Context, poolId uint64) (sdk.Coin, error)
}

type AssetKeeper interface {
	GetPairsVault(ctx sdk.Context, id uint64) (pairs assettypes.ExtendedPairVault, found bool)
	HasAssetForDenom(ctx sdk.Context, denom string) bool
	GetAssetForDenom(ctx sdk.Context, denom string) (asset assettypes.Asset, found bool)
	GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool)
}

type MarketKeeper interface {
	GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool)
}

type LockerKeeper interface {
	GetLockerProductAssetMapping(ctx sdk.Context, appMappingId uint64) (lockerProductMapping lockertypes.LockerProductAssetMapping, found bool)
	GetLocker(ctx sdk.Context, lockerId string) (locker lockertypes.Locker, found bool)
	GetLockerLookupTable(ctx sdk.Context, appMappingId uint64) (lockerLookupData lockertypes.LockerLookupTable, found bool)
	UpdateLocker(ctx sdk.Context, locker lockertypes.Locker)
	SetLockerTotalRewardsByAssetAppWise(ctx sdk.Context, lockerRewardsMapping lockertypes.LockerTotalRewardsByAssetAppWise) error
	GetLockerTotalRewardsByAssetAppWise(ctx sdk.Context, app_id, asset_id uint64) (lockerRewardsMapping lockertypes.LockerTotalRewardsByAssetAppWise, found bool)
}

type CollectorKeeper interface {
	GetAppidToAssetCollectorMapping(ctx sdk.Context, app_id uint64) (appAssetCollectorData collecortypes.AppIdToAssetCollectorMapping, found bool)
	GetCollectorLookupTable(ctx sdk.Context, app_id uint64) (collectorLookup collecortypes.CollectorLookup, found bool)
	GetAppToDenomsMapping(ctx sdk.Context, AppId uint64) (appToDenom collecortypes.AppToDenomsMapping, found bool)
	GetCollectorLookupByAsset(ctx sdk.Context, app_id, asset_id uint64) (collectorLookup collecortypes.CollectorLookupTable, found bool)
	GetNetFeeCollectedData(ctx sdk.Context, app_id uint64) (netFeeData collecortypes.NetFeeCollectedData, found bool)
	SetNetFeeCollectedData(ctx sdk.Context, app_id, asset_id uint64, fee sdk.Int) error
	DecreaseNetFeeCollectedData(ctx sdk.Context, appId, assetId uint64, amount sdk.Int) error
}

type VaultKeeper interface {
	GetAppExtendedPairVaultMapping(ctx sdk.Context, appMappingId uint64) (appExtendedPairVaultData vaulttypes.AppExtendedPairVaultMapping, found bool)
	CalculateCollaterlizationRatio(ctx sdk.Context, extendedPairVaultId uint64, amountIn sdk.Int, amountOut sdk.Int) (sdk.Dec, error)
	GetVault(ctx sdk.Context, id string) (vault vaulttypes.Vault, found bool)
	DeleteVault(ctx sdk.Context, id string)
	UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx sdk.Context, counter uint64, vaultData vaulttypes.Vault)
	UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, valutLookupData vaulttypes.AppExtendedPairVaultMapping, extendedPairId uint64, amount sdk.Int, changeType bool)
	UpdateTokenMintedAmountLockerMapping(ctx sdk.Context, valutLookupData vaulttypes.AppExtendedPairVaultMapping, extendedPairId uint64, amount sdk.Int, changeType bool)
	UpdateUserVaultExtendedPairMapping(ctx sdk.Context, extendedPairId uint64, userAddress string, appMappingId uint64)
	DeleteAddressFromAppExtendedPairVaultMapping(ctx sdk.Context, extendedPairId uint64, userVaultId string, appMappingId uint64)
	SetVault(ctx sdk.Context, vault vaulttypes.Vault)
}

type BankKeeper interface {
	BurnCoins(ctx sdk.Context, name string, coins sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, coins sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, address sdk.AccAddress, name string, coins sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, name string, address sdk.AccAddress, coins sdk.Coins) error

	SendCoinsFromModuleToModule(
		ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins,
	) error

	SpendableCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
}
