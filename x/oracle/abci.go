package oracle

import (
	"github.com/comdex-official/comdex/x/oracle/keeper"
	"github.com/comdex-official/comdex/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	block := k.GetLastBlockheight(ctx)

	if ctx.BlockHeight()%20 == 0 && ctx.BlockHeight() > block+21 {
		assets := k.GetAssets(ctx)
		for _, asset := range assets {
			k.SetRates(ctx, asset.Name)
			k.SetMarketForAsset(ctx, asset.Id, asset.Name)
			rate, _ := k.GetRates(ctx, asset.Name)
			var (
				market = types.Market{
					Symbol:   asset.Name,
					ScriptID: 112,
					Rates:    rate,
				}
			)
			k.SetMarket(ctx, market)
		}
	}
}
