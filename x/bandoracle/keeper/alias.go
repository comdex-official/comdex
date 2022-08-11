package keeper

import (
	"github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAssets(ctx sdk.Context) (assets []types.Asset) {
	return k.assetKeeper.GetAssets(ctx)
}
