package expected

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/comdex-official/comdex/x/asset/types"

)

type AssetKeeper interface {
	// GetApps(ctx sdk.Context) (assettypes.AppMapping, bool)
	HasAssetForDenom(ctx sdk.Context, id string) (bool)
	HasAsset(ctx sdk.Context, id uint64) (bool)
	GetAssetForDenom(ctx sdk.Context, denom string) (types.Asset, bool)
	GetApp(ctx sdk.Context, id uint64) (types.AppMapping, bool)
}
