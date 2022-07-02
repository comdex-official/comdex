package expected

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	"github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
)

type AccountKeeper interface {
	GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI
}

type BankKeeper interface {
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
	GetAppExtendedPairVaultMapping(ctx sdk.Context, appMappingID uint64) (appExtendedPairVaultData types.AppExtendedPairVaultMapping, found bool)
	CalculateCollaterlizationRatio(ctx sdk.Context, extendedPairVaultID uint64, amountIn sdk.Int, amountOut sdk.Int) (sdk.Dec, error)
	GetVault(ctx sdk.Context, id string) (vault types.Vault, found bool)
	DeleteVault(ctx sdk.Context, id string)
	UpdateAppExtendedPairVaultMappingDataOnMsgCreate(ctx sdk.Context, counter uint64, vaultData types.Vault)
	UpdateCollateralLockedAmountLockerMapping(ctx sdk.Context, vaultLookupData types.AppExtendedPairVaultMapping, extendedPairID uint64, amount sdk.Int, changeType bool)
	UpdateTokenMintedAmountLockerMapping(ctx sdk.Context, vaultLookupData types.AppExtendedPairVaultMapping, extendedPairID uint64, amount sdk.Int, changeType bool)
	UpdateUserVaultExtendedPairMapping(ctx sdk.Context, extendedPairID uint64, userAddress string, appMappingID uint64)
	DeleteAddressFromAppExtendedPairVaultMapping(ctx sdk.Context, extendedPairID uint64, userVaultID string, appMappingID uint64)
	SetVault(ctx sdk.Context, vault types.Vault)
}

type MarketKeeper interface {
	GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool)
}

type AuctionKeeper interface {
	GetParams(ctx sdk.Context) auctiontypes.Params
}

type EsmKeeper interface {
	GetKillSwitchData(ctx sdk.Context, app_id uint64) (esmtypes.KillSwitchParams, bool)
}