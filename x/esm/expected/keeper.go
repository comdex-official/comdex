package expected

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
)

type AssetKeeper interface {
	GetApp(ctx sdk.Context, id uint64) (app assettypes.AppData, found bool)
	GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool)
	GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool)
	GetApps(ctx sdk.Context) (apps []assettypes.AppData, found bool)
	GetAssets(ctx sdk.Context) (assets []assettypes.Asset)
	GetPairsVault(ctx sdk.Context, id uint64) (pairs assettypes.ExtendedPairVault, found bool)
	GetAssetForDenom(ctx sdk.Context, denom string) (asset assettypes.Asset, found bool)
}

type VaultKeeper interface {
	GetVaults(ctx sdk.Context) (vaults []vaulttypes.Vault)
	DeleteVault(ctx sdk.Context, id uint64)
	DeleteAddressFromAppExtendedPairVaultMapping(ctx sdk.Context, extendedPairID uint64, userVaultID uint64, appMappingID uint64)
	GetStableMintVaults(ctx sdk.Context) (stableVaults []vaulttypes.StableMintVault)
	UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, appMappingID uint64, extendedPairID uint64, amount sdkmath.Int, changeType bool)
	UpdateTokenMintedAmountLockerMapping(ctx sdk.Context, appMappingID uint64, extendedPairID uint64, amount sdkmath.Int, changeType bool)
	DeleteUserVaultExtendedPairMapping(ctx sdk.Context, address string, appID uint64, pairVaultID uint64)
	GetLengthOfVault(ctx sdk.Context) uint64
	SetLengthOfVault(ctx sdk.Context, length uint64)
}

type BankKeeper interface {
	BurnCoins(ctx context.Context, name string, coins sdk.Coins) error
	MintCoins(ctx context.Context, name string, coins sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, address sdk.AccAddress, name string, coins sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, name string, address sdk.AccAddress, coins sdk.Coins) error
	SpendableCoins(ctx context.Context, address sdk.AccAddress) sdk.Coins
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromModuleToModule(
		ctx context.Context, senderModule, recipientModule string, amt sdk.Coins,
	) error
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin 
}

type MarketKeeper interface {
	GetTwa(ctx sdk.Context, id uint64) (twa markettypes.TimeWeightedAverage, found bool)
}

type Tokenmint interface {
	BurnTokensForApp(ctx sdk.Context, appMappingID uint64, assetID uint64, amount sdkmath.Int) error
}

type Collector interface {
	GetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64) (netFeeData collectortypes.AppAssetIdToFeeCollectedData, found bool)
	GetAppNetFeeCollectedData(ctx sdk.Context, appID uint64) (netFeeData []collectortypes.AppAssetIdToFeeCollectedData, found bool)
	DecreaseNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, amount sdkmath.Int) error
}
