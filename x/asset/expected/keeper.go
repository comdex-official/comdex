package expected

import (
	"github.com/comdex-official/comdex/x/market/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MarketKeeper interface {
	GetMarketForAsset(ctx sdk.Context, id uint64) (types.Market, bool)
	GetPriceForMarket(ctx sdk.Context, symbol string) (uint64, bool)
}
