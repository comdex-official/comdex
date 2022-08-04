package expected

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
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
	GetPair(ctx sdk.Context, appID, id uint64) (pair liquiditytypes.Pair, found bool)
	GetPool(ctx sdk.Context, appID, id uint64) (pool liquiditytypes.Pool, found bool)
	GetFarmingRewardsData(ctx sdk.Context, appID uint64, coinToDistribute sdk.Coin, liquidityGaugeData types.LiquidtyGaugeMetaData) ([]types.RewardDistributionDataCollector, error)
	TransferFundsForSwapFeeDistribution(ctx sdk.Context, appID, poolID uint64) (sdk.Coin, error)
}

type AssetKeeper interface {
	GetPairsVault(ctx sdk.Context, id uint64) (pairs assettypes.ExtendedPairVault, found bool)
	HasAssetForDenom(ctx sdk.Context, denom string) bool
	GetAssetForDenom(ctx sdk.Context, denom string) (asset assettypes.Asset, found bool)
	GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool)
	GetApp(ctx sdk.Context, id uint64) (app assettypes.AppData, found bool)
	GetApps(ctx sdk.Context) (apps []assettypes.AppData, found bool)
}

type MarketKeeper interface {
	GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool)
}

type LockerKeeper interface {
	GetLockerProductAssetMapping(ctx sdk.Context, appMappingID uint64) (lockerProductMapping lockertypes.LockerProductAssetMapping, found bool)
	GetLocker(ctx sdk.Context, lockerID string) (locker lockertypes.Locker, found bool)
	GetLockerLookupTable(ctx sdk.Context, appMappingID uint64) (lockerLookupData lockertypes.LockerLookupTable, found bool)
	UpdateLocker(ctx sdk.Context, locker lockertypes.Locker)
	SetLockerTotalRewardsByAssetAppWise(ctx sdk.Context, lockerRewardsMapping lockertypes.LockerTotalRewardsByAssetAppWise) error
	GetLockerTotalRewardsByAssetAppWise(ctx sdk.Context, appID, assetID uint64) (lockerRewardsMapping lockertypes.LockerTotalRewardsByAssetAppWise, found bool)
}

type CollectorKeeper interface {
	GetAppidToAssetCollectorMapping(ctx sdk.Context, appID uint64) (appAssetCollectorData collectortypes.AppIdToAssetCollectorMapping, found bool)
	GetCollectorLookupTable(ctx sdk.Context, appID uint64) (collectorLookup collectortypes.CollectorLookup, found bool)
	GetAppToDenomsMapping(ctx sdk.Context, appID uint64) (appToDenom collectortypes.AppToDenomsMapping, found bool)
	GetCollectorLookupByAsset(ctx sdk.Context, appID, assetID uint64) (collectorLookup collectortypes.CollectorLookupTable, found bool)
	GetNetFeeCollectedData(ctx sdk.Context, appID uint64) (netFeeData collectortypes.NetFeeCollectedData, found bool)
	SetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, fee sdk.Int) error
	DecreaseNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, amount sdk.Int, netFeeCollectedData collectortypes.NetFeeCollectedData) error
}

type VaultKeeper interface {
	GetAppExtendedPairVaultMapping(ctx sdk.Context, appMappingID uint64) (appExtendedPairVaultData vaulttypes.AppExtendedPairVaultMapping, found bool)
	CalculateCollaterlizationRatio(ctx sdk.Context, extendedPairVaultID uint64, amountIn sdk.Int, amountOut sdk.Int) (sdk.Dec, error)
	GetVault(ctx sdk.Context, id string) (vault vaulttypes.Vault, found bool)
	DeleteVault(ctx sdk.Context, id string)
	UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx sdk.Context, counter uint64, vaultData vaulttypes.Vault)
	UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, vaultLookupData vaulttypes.AppExtendedPairVaultMapping, extendedPairID uint64, amount sdk.Int, changeType bool)
	UpdateTokenMintedAmountLockerMapping(ctx sdk.Context, vaultLookupData vaulttypes.AppExtendedPairVaultMapping, extendedPairID uint64, amount sdk.Int, changeType bool)
	UpdateUserVaultExtendedPairMapping(ctx sdk.Context, extendedPairID uint64, userAddress string, appMappingID uint64)
	DeleteAddressFromAppExtendedPairVaultMapping(ctx sdk.Context, extendedPairID uint64, userVaultID string, appMappingID uint64)
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

type EsmKeeper interface {
	GetKillSwitchData(ctx sdk.Context, app_id uint64) (esmtypes.KillSwitchParams, bool)
	GetESMStatus(ctx sdk.Context, id uint64) (esmStatus esmtypes.ESMStatus, found bool)
}
