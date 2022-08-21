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
	GetLockerProductAssetMapping(ctx sdk.Context, appMappingID, assetID uint64) (lockerProductMapping lockertypes.LockerProductAssetMapping, found bool)
	GetLocker(ctx sdk.Context, lockerID uint64) (locker lockertypes.Locker, found bool)
	GetLockers(ctx sdk.Context) (locker []lockertypes.Locker)
	GetLockerLookupTableByApp(ctx sdk.Context, appID uint64) (lockerLookupData []lockertypes.LockerLookupTableData, found bool)
	GetLockerLookupTable(ctx sdk.Context, appID, assetID uint64) (lockerLookupData lockertypes.LockerLookupTableData, found bool)
	SetLocker(ctx sdk.Context, locker lockertypes.Locker)
	SetLockerTotalRewardsByAssetAppWise(ctx sdk.Context, lockerRewardsMapping lockertypes.LockerTotalRewardsByAssetAppWise) error
	GetLockerTotalRewardsByAssetAppWise(ctx sdk.Context, appID, assetID uint64) (lockerRewardsMapping lockertypes.LockerTotalRewardsByAssetAppWise, found bool)
	SetLockerLookupTable(ctx sdk.Context, lockerLookupData lockertypes.LockerLookupTableData)
}

type CollectorKeeper interface {
	GetAppidToAssetCollectorMapping(ctx sdk.Context, appID, assetID uint64) (appAssetCollectorData collectortypes.AppToAssetIdCollectorMapping, found bool)
	GetCollectorLookupTable(ctx sdk.Context, appID, assetID uint64) (collectorLookup collectortypes.CollectorLookupTableData, found bool)
	GetAppToDenomsMapping(ctx sdk.Context, appID uint64) (appToDenom collectortypes.AppToDenomsMapping, found bool)
	GetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64) (netFeeData collectortypes.AppAssetIdToFeeCollectedData, found bool)
	SetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, fee sdk.Int) error
	DecreaseNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, amount sdk.Int, netFeeCollectedData collectortypes.AppAssetIdToFeeCollectedData) error
}

type VaultKeeper interface {
	GetAppMappingData(ctx sdk.Context, appMappingID uint64) (appExtendedPairVaultData []vaulttypes.AppExtendedPairVaultMappingData, found bool)
	CalculateCollaterlizationRatio(ctx sdk.Context, extendedPairVaultID uint64, amountIn sdk.Int, amountOut sdk.Int) (sdk.Dec, error)
	GetVault(ctx sdk.Context, id uint64) (vault vaulttypes.Vault, found bool)
	DeleteVault(ctx sdk.Context, id uint64)
	UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx sdk.Context, vaultData vaulttypes.Vault)
	UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, vaultLookupData vaulttypes.AppExtendedPairVaultMappingData, amount sdk.Int, changeType bool)
	UpdateTokenMintedAmountLockerMapping(ctx sdk.Context, vaultLookupData vaulttypes.AppExtendedPairVaultMappingData, amount sdk.Int, changeType bool)
	DeleteUserVaultExtendedPairMapping(ctx sdk.Context, address string, appID uint64, pairVaultID uint64)
	DeleteAddressFromAppExtendedPairVaultMapping(ctx sdk.Context, extendedPairID uint64, userVaultID uint64, appMappingID uint64)
	SetVault(ctx sdk.Context, vault vaulttypes.Vault)
	GetAppExtendedPairVaultMappingData(ctx sdk.Context, appMappingID uint64, pairVaultID uint64) (appExtendedPairVaultData vaulttypes.AppExtendedPairVaultMappingData, found bool)
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
	GetKillSwitchData(ctx sdk.Context, appID uint64) (esmtypes.KillSwitchParams, bool)
	GetESMStatus(ctx sdk.Context, id uint64) (esmStatus esmtypes.ESMStatus, found bool)
}
