package bandoracle

import (
	"fmt"
	"github.com/comdex-official/comdex/x/bandoracle/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	fmt.Println(k.GetOracleValidationResult(ctx))
	block := k.GetLastBlockheight(ctx)
	if ctx.BlockHeight()%20 == 0 {
		req := k.GetTempFetchPriceID(ctx)
		res := k.OraclePriceValidationByRequestId(ctx, req)
		k.SetOracleValidationResult(ctx, res)
	} else if ctx.BlockHeight()%20-1 == 0 && ctx.BlockHeight() > block+11 {
		id := k.GetLastFetchPriceID(ctx)
		k.SetTempFetchPriceID(ctx, id)
		msg := k.GetFetchPriceMsg(ctx)
		_, err := k.FetchPrice(ctx, msg)
		if err != nil {
			return
		}
	}
}
