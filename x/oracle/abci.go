package oracle

import (
	"github.com/comdex-official/comdex/x/oracle/keeper"
	"github.com/comdex-official/comdex/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)


func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {

	if ctx.BlockHeight()>=59{
	k.SetRates(ctx, "ATOM")
	k.SetRates(ctx, "XAU")
	k.SetMarketForAsset(ctx, 1, "ATOM")
	k.SetMarketForAsset(ctx, 2, "XAU")

		rates1, _:= k.GetRates(ctx, "ATOM")
		rates2, _:= k.GetRates(ctx, "XAU")


		var (
			market1 = types.Market{
				Symbol:   "ATOM",
				ScriptID: 112,
				Rates: rates1,
			}
			market2 = types.Market{
				Symbol:   "XAU",
				ScriptID: 112,
				Rates: rates2,
			}
		)
		k.SetMarket(ctx, market1)
		k.SetMarket(ctx, market2)

	}
}
