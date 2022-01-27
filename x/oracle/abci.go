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
	k.SetRates(ctx, "BTC")
	k.SetMarketForAsset(ctx, 1, "ATOM")
	k.SetMarketForAsset(ctx, 2, "BTC")

		rates1, _:= k.GetRates(ctx, "ATOM")
		rates2, _:= k.GetRates(ctx, "BTC")


		var (
			market1 = types.Market{
				Symbol:   "ATOM",
				ScriptID: 37,
				Rates: rates1,
			}
			market2 = types.Market{
				Symbol:   "BTC",
				ScriptID: 37,
				Rates: rates2,
			}
		)
		k.SetMarket(ctx, market1)
		k.SetMarket(ctx, market2)

	}
}
