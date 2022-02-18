package bandoracle

import (
	"github.com/comdex-official/comdex/x/bandoracle/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	id := k.GetLastFetchPriceID(ctx)
	if id != 0 {
		if ctx.BlockHeight()%20 == 0 {
			msg := k.GetFetchPriceMsg(ctx)
			_, err := k.FetchPrice(ctx, msg)
			if err != nil {
				return
			}
		}
	}
}
