package expected

import (
	"github.com/comdex-official/comdex/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OracleKeeper interface {
	GetMarketForAsset(ctx sdk.Context, id uint64) (types.Market, bool)
	GetPriceForMarket(ctx sdk.Context, symbol string) (uint64, bool)
}
