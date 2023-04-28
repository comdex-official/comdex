package expected

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	liquidationsV2types "github.com/comdex-official/comdex/x/liquidationsV2/types"
	markettypes "github.com/comdex-official/comdex/x/market/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type LiquidationsV2Keeper interface {
	GetLiquidationWhiteListing(ctx sdk.Context, appId uint64) (liquidationWhiteListing liquidationsV2types.LiquidationWhiteListing, found bool)
}

type MarketKeeper interface {
	CalcAssetPrice(ctx sdk.Context, id uint64, amt sdk.Int) (price sdk.Dec, err error)
	GetTwa(ctx sdk.Context, id uint64) (twa markettypes.TimeWeightedAverage, found bool)
}

type AssetKeeper interface {
	GetAsset(ctx sdk.Context, id uint64) (assettypes.Asset, bool)
	GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool)
	GetApps(ctx sdk.Context) (apps []assettypes.AppData, found bool)
	GetApp(ctx sdk.Context, id uint64) (app assettypes.AppData, found bool)
	GetPairsVault(ctx sdk.Context, id uint64) (pairs assettypes.ExtendedPairVault, found bool)
}
