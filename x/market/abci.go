package market

import (
	utils "github.com/comdex-official/comdex/types"
	bandoraclemoduletypes "github.com/comdex-official/comdex/x/bandoracle/types"
	"github.com/comdex-official/comdex/x/market/keeper"
	"github.com/comdex-official/comdex/x/market/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		if k.GetOracleValidationResult(ctx) {
			block := k.GetLastBlockheight(ctx)
			if block != types.Int64Zero {
				if ctx.BlockHeight()%types.Int64Twenty == types.Int64Zero && ctx.BlockHeight() != block {
					assets := k.GetAssets(ctx)
					id := k.GetLastFetchPriceID(ctx)
					data, _ := k.GetFetchPriceResult(ctx, bandoraclemoduletypes.OracleRequestID(id))
					scriptID := k.GetFetchPriceMsg(ctx).OracleScriptID
					index := -1
					length := len(data.Rates)
					for _, asset := range assets {
						if asset.IsOraclePriceRequired && data.Rates != nil {
							index = index + 1
							if length > index {
								k.UpdatePriceList(ctx, asset.Id, scriptID, data.Rates[index])
							}

							// below code to be removed
							// k.SetRates(ctx, asset.Name)
							// k.SetMarketForAsset(ctx, asset.Id, asset.Name)
							// rate, _ := k.GetRates(ctx, asset.Name)
							// scriptID := k.GetFetchPriceMsg(ctx).OracleScriptID
							// var (
							// 	market = types.Market{
							// 		Symbol:   asset.Name,
							// 		ScriptID: scriptID,
							// 		Rates:    rate,
							// 	}
							// )
							// k.SetMarket(ctx, market)
							// above code to be removed

						}
					}
				}
			}
		}
		return nil
	})
}

// set twa
// get 1 twa
//get all twa
// get current price for asset. - if price not active then error
//calculate asset price  - if price not active then error
