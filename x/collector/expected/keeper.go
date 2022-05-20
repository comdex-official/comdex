package expected

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AssetKeeper interface {
	//GetApps(ctx sdk.Context) (assettypes.AppMapping, bool)
	HasAssetForDenom(ctx sdk.Context, id string) bool
	HasAsset(ctx sdk.Context, id uint64) bool
}
