package expected

import (
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type LiquidationsV2Keeper interface {
	GetLiquidationWhiteListing(ctx sdk.Context, appId uint64) (liquidationWhiteListing types.LiquidationWhiteListing, found bool)
}
