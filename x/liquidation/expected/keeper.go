package expected

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type AccountKeeper interface {
	GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI
}

type BankKeeper interface {
	MintCoins(ctx sdk.Context, name string, coins sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type AssetKeeper interface {
	GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool)
	GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool)
}

type VaultKeeper interface {
	GetVaults(ctx sdk.Context) (vaults []vaulttypes.Vault)
	CalculateCollaterlizationRatio(
		ctx sdk.Context,
		amountIn sdk.Int,
		assetIn assettypes.Asset,
		amountOut sdk.Int,
		assetOut assettypes.Asset,
	) (sdk.Dec, error)
	CreteNewVault(ctx sdk.Context, pairdId uint64, from string, assetIn assettypes.Asset, amountIn sdk.Int, assetOut assettypes.Asset, amountOut sdk.Int) error
	DeleteVault(ctx sdk.Context, id uint64)
	DeleteVaultForAddressByPair(ctx sdk.Context, address sdk.AccAddress, pairID uint64)
	GetID(ctx sdk.Context) uint64
	SetID(ctx sdk.Context, id uint64)
	SetVault(ctx sdk.Context, vault vaulttypes.Vault)
	SetVaultForAddressByPair(ctx sdk.Context, address sdk.AccAddress, pairID uint64, id uint64)
	UpdateUserVaultIdMapping(ctx sdk.Context, vaultOwner string, vaultId uint64, isInsert bool) error
	UpdateCollateralVaultIdMapping(ctx sdk.Context, assetInDenom string, assetOutDenom string, vaultId uint64, isInsert bool) error
}

type OracleKeeper interface {
	GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool)
}
