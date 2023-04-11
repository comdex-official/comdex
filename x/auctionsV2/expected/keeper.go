package expected

import (
	liquidationsV2types "github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type LiquidationsV2Keeper interface {
	GetLiquidationWhiteListing(ctx sdk.Context, appId uint64) (liquidationWhiteListing liquidationsV2types.LiquidationWhiteListing, found bool)
}
