package keeper

import (
	"github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) GetAssetsForOracle(ctx sdk.Context) (assets []types.Asset) {
	return k.assetKeeper.GetAssetsForOracle(ctx)
}
