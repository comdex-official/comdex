package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (k *Keeper) GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI {
	return k.account.GetModuleAccount(ctx, name)
}

func (k *Keeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return k.bank.GetBalance(ctx, addr, denom)
}

func (k *Keeper) MintCoin(ctx sdk.Context, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return nil
	}

	return k.bank.MintCoins(ctx, name, sdk.NewCoins(coin))
}

func (k *Keeper) SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error {
	return k.bank.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, amt)
}
func (k *Keeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	return k.bank.SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, amt)
}
func (k *Keeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	return k.bank.SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt)
}
