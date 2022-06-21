package keeper

import (
	"github.com/comdex-official/comdex/x/tokenmint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
)

func (k *Keeper) BurnCoin(ctx sdk.Context, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return nil
	}

	return k.bank.BurnCoins(ctx, name, sdk.NewCoins(coin))
}

func (k *Keeper) MintCoin(ctx sdk.Context, name string, coin sdk.Coin) error {
	if coin.IsZero() {
		return types.ErrorMintingGenesisSupplyLessThanOne
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
func (k *Keeper) SendCoinFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, coin sdk.Coins) error {
	if coin.IsZero() {
		return nil
	}
	return k.bank.SendCoinsFromModuleToModule(ctx, senderModule, recipientModule, coin)
}

func (k *Keeper) SpendableCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins {
	return k.bank.SpendableCoins(ctx, address)
}

func (k *Keeper) GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool) {
	return k.asset.GetAsset(ctx, id)
}

func (k *Keeper) GetApp(ctx sdk.Context, id uint64) (assettypes.AppMapping, bool) {
	return k.asset.GetApp(ctx, id)
}

func (k *Keeper) GetApps(ctx sdk.Context) ([]assettypes.AppMapping, bool) {
	return k.asset.GetApps(ctx)
}

func (k *Keeper) GetMintGenesisTokenData(ctx sdk.Context, appId, assetId uint64) (assettypes.MintGenesisToken, bool) {
	return k.asset.GetMintGenesisTokenData(ctx, appId, assetId)
}

func (k *Keeper) GetAssetForDenom(ctx sdk.Context, denom string) (assettypes.Asset, bool){
	return k.asset.GetAssetForDenom(ctx,denom)
}