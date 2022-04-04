package market

import (
	"github.com/comdex-official/comdex/x/market/keeper"
	"github.com/comdex-official/comdex/x/market/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	block := k.GetLastBlockheight(ctx)
	if block !=0 {
		if ctx.BlockHeight()%20-1 == 0 && ctx.BlockHeight() > block+21 {
			assets := k.GetAssets(ctx)
			for _, asset := range assets {
				k.SetRates(ctx, asset.Name)
				k.SetMarketForAsset(ctx, asset.Id, asset.Name)
				rate, _ := k.GetRates(ctx, asset.Name)
				scriptID := k.GetFetchPriceMsg(ctx).OracleScriptID
				var (
					market = types.Market{
						Symbol:   asset.Name,
						ScriptID: scriptID,
						Rates:    rate,
					}
				)
				k.SetMarket(ctx, market)
			}
		}
	}
}
