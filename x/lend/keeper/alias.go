package keeper

import (
	"github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) BurnCoin(ctx sdk.Context, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return nil
	}

	return k.bank.BurnCoins(ctx, name, sdk.NewCoins(coin))
}

func (k *Keeper) MintCoin(ctx sdk.Context, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return nil
	}

	return k.bank.MintCoins(ctx, name, sdk.NewCoins(coin))
}

func (k *Keeper) SendCoinFromAccountToModule(ctx sdk.Context, address sdk.AccAddress, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return nil
	}

	return k.bank.SendCoinsFromAccountToModule(ctx, address, name, sdk.NewCoins(coin))
}

func (k *Keeper) SendCoinFromModuleToAccount(ctx sdk.Context, name string, address sdk.AccAddress, coin sdk.Coin) error {
	if coin.IsZero() {
		return nil
	}

	return k.bank.SendCoinsFromModuleToAccount(ctx, name, address, sdk.NewCoins(coin))
}

func (k *Keeper) SpendableCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins {
	return k.bank.SpendableCoins(ctx, address)
}

func (k *Keeper) GetModuleAddress(name string) sdk.AccAddress {
	return k.account.GetModuleAddress(name)
}

func (k *Keeper) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return k.bank.GetAllBalances(ctx, addr)
}

func (k *Keeper) GetWhitelistAsset(ctx sdk.Context, id uint64) (asset types.ExtendedAsset, found bool) {
	return k.asset.GetWhitelistAsset(ctx, id)
}

func (k *Keeper) GetWhitelistPair(ctx sdk.Context, id uint64) (pair types.ExtendedPairLend, found bool) {
	return k.asset.GetWhitelistPair(ctx, id)
}

func (k *Keeper) GetPair(ctx sdk.Context, id uint64) (pair types.Pair, found bool) {
	return k.asset.GetPair(ctx, id)
}

func (k *Keeper) GetAsset(ctx sdk.Context, id uint64) (asset types.Asset, found bool) {
	return k.asset.GetAsset(ctx, id)
}

func (k *Keeper) GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool) {
	return k.oracle.GetPriceForAsset(ctx, id)
}
