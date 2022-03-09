package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) GetAssets(ctx sdk.Context) (assets []assettypes.Asset) {
	return k.asset.GetAssets(ctx)
}
