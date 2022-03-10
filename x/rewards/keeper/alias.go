package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) GetAssets(ctx sdk.Context) (assets []assettypes.Asset) {
	return k.asset.GetAssets(ctx)
}

// func (k *Keeper) GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI {
// 	return k.account.GetModuleAccount(ctx, name)
// }

// func (k *Keeper) GetModuleAddress(name string) sdk.AccAddress {
// 	return k.account.GetModuleAddress(name)
// }
