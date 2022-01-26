package oracle

import (
	"fmt"
	"github.com/comdex-official/comdex/x/oracle/keeper"
	"github.com/comdex-official/comdex/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)


func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	fmt.Println("yes")
	if ctx.BlockHeight()>=59{
		fmt.Println("Setting Rates")
	k.SetRates(ctx, "ATOM")
	rates, _:= k.GetRates(ctx, "ATOM")

		var (
			market = types.Market{
				Symbol:   "ATOM",
				ScriptID: 37,
				Rates: rates,
			}
		)
		k.SetMarket(ctx, market)
	}
}
