package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
)

func (k Keeper) getModAccountBalances(ctx sdk.Context, accountName string, denom string) sdk.Int {
	macc := k.GetModuleAccount(ctx, accountName)
	return k.GetBalance(ctx, macc.GetAddress(), denom).Amount
}

func (k Keeper) SeizeCollateral(
	ctx sdk.Context,
	vault vaulttypes.Vault,
	assetIn assettypes.Asset,
	assetOut assettypes.Asset,
	collateralizationRatio sdk.Dec,
) error {

	collateralAvailableInVaultModule := k.getModAccountBalances(ctx, vaulttypes.ModuleName, assetIn.Denom)
	collateralCoins := sdk.NewCoin(assetIn.Denom, sdk.MinInt(vault.AmountIn, collateralAvailableInVaultModule))

	err := k.bank.SendCoinsFromModuleToModule(ctx, vaulttypes.ModuleName, auctiontypes.ModuleName, sdk.NewCoins(collateralCoins))
	if err != nil {
		return err
	}
	vaultOwner, _ := sdk.AccAddressFromBech32(vault.Owner)
	k.DeleteVault(ctx, vault.ID)
	k.DeleteVaultForAddressByPair(ctx, vaultOwner, vault.PairID)
	k.StartCollateralAuction(
		ctx,
		vault,
		collateralizationRatio,
		assetIn,
		assetOut,
	)
	return nil
}
