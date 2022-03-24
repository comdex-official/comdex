package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) GetAssets(ctx sdk.Context) (assets []assettypes.Asset) {
	return k.asset.GetAssets(ctx)
}

func (k *Keeper) GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool) {
	return k.asset.GetAsset(ctx, id)
}

func (k *Keeper) GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool) {
	return k.asset.GetPair(ctx, id)
}

func (k *Keeper) GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool) {
	return k.oracle.GetPriceForAsset(ctx, id)
}

func (k *Keeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	return k.bank.SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt)
}

func (k *Keeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	return k.bank.SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt)
}

func (k *Keeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return k.bank.GetBalance(ctx, addr, denom)
}

func (k *Keeper) GetCollateralBasedVaults(ctx sdk.Context, collateral_denom string) (collateralBasedVaults vaulttypes.CollateralVaultIdMapping, found bool) {
	return k.vault.GetCollateralBasedVaults(ctx, collateral_denom)
}

func (k *Keeper) GetVault(ctx sdk.Context, id uint64) (vault vaulttypes.Vault, found bool) {
	return k.vault.GetVault(ctx, id)
}

func (k *Keeper) GetCAssetTotalValueMintedForCollateral(ctx sdk.Context, collateralType assettypes.Asset) sdk.Dec {
	return k.vault.GetCAssetTotalValueMintedForCollateral(ctx, collateralType)
}
