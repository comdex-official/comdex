package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
)

func (k Keeper) BurnCoin(ctx sdk.Context, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return lendtypes.BurnCoinValueInLendIsZero
	}

	return k.bank.BurnCoins(ctx, name, sdk.NewCoins(coin))
}

func (k Keeper) MintCoin(ctx sdk.Context, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return lendtypes.MintCoinValueInLendIsZero
	}

	return k.bank.MintCoins(ctx, name, sdk.NewCoins(coin))
}

func (k Keeper) SendCoinFromAccountToModule(ctx sdk.Context, address sdk.AccAddress, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return lendtypes.SendCoinsFromAccountToModuleInLendIsZero
	}

	return k.bank.SendCoinsFromAccountToModule(ctx, address, name, sdk.NewCoins(coin))
}

func (k Keeper) SendCoinFromModuleToAccount(ctx sdk.Context, name string, address sdk.AccAddress, coin sdk.Coin) error {
	if coin.IsZero() {
		return lendtypes.SendCoinsFromModuleToAccountInLendIsZero
	}

	return k.bank.SendCoinsFromModuleToAccount(ctx, name, address, sdk.NewCoins(coin))
}

func (k Keeper) SpendableCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins {
	return k.bank.SpendableCoins(ctx, address)
}

func (k Keeper) GetModuleAddress(name string) sdk.AccAddress {
	return k.account.GetModuleAddress(name)
}

func (k Keeper) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return k.bank.GetAllBalances(ctx, addr)
}

func (k Keeper) GetAsset(ctx sdk.Context, ID uint64) (assettypes.Asset, bool) {
	return k.asset.GetAsset(ctx, ID)
}

func (k Keeper) GetPriceForAsset(ctx sdk.Context, ID uint64) (uint64, bool) {
	return k.market.GetPriceForAsset(ctx, ID)
}

func (k Keeper) SendCoinFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, coin sdk.Coins) error {
	if coin.IsZero() {
		return lendtypes.SendCoinsFromModuleToModuleInLendIsZero
	}
	return k.bank.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, coin)
}

func (k Keeper) GetApp(ctx sdk.Context, id uint64) (assettypes.AppData, bool) {
	return k.asset.GetApp(ctx, id)
}

func (k Keeper) GetKillSwitchData(ctx sdk.Context, appID uint64) (esmtypes.KillSwitchParams, bool) {
	return k.esm.GetKillSwitchData(ctx, appID)
}
