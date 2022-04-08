package expected

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
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

type AssetKeeper interface {
	GetAssets(ctx sdk.Context) (assets []assettypes.Asset)
	GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool)
	GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool)
}

type VaultKeeper interface {
	GetCollateralBasedVaults(ctx sdk.Context, collateral_denom string) (collateralBasedVaults vaulttypes.CollateralVaultIdMapping, found bool)
	GetVault(ctx sdk.Context, id uint64) (vault vaulttypes.Vault, found bool)
	SetVault(ctx sdk.Context, vault vaulttypes.Vault)
	GetCAssetTotalValueMintedForCollateral(ctx sdk.Context, collateralType assettypes.Asset) sdk.Dec
}

type MarketKeeper interface {
	GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool)
}
