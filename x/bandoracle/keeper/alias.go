package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k Keeper) GetAssets(ctx sdk.Context) (assets []types.Asset) {
	return k.assetKeeper.GetAssets(ctx)
}
