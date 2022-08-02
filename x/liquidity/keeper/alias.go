package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
)

func (k Keeper) GetApps(ctx sdk.Context) (apps []assettypes.AppData, found bool) {
	return k.assetKeeper.GetApps(ctx)
}
