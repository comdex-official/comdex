package market

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/market/keeper"
	"github.com/comdex-official/comdex/x/market/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		block := k.GetLastBlockheight(ctx)
		if block != types.Int64Zero {
			if ctx.BlockHeight()%types.Int64Twenty-types.Int64One == types.Int64Zero && ctx.BlockHeight() > block+types.Int64TwentyOne {
				assets := k.GetAssets(ctx)
				for _, asset := range assets {
					if asset.IsOraclePriceRequired {
						k.SetRates(ctx, asset.Name)
						k.SetMarketForAsset(ctx, asset.Id, asset.Name)
						rate, _ := k.GetRates(ctx, asset.Name)
						scriptID := k.GetFetchPriceMsg(ctx).OracleScriptID
						market := types.Market{
							Symbol:   asset.Name,
							ScriptID: scriptID,
							Rates:    rate,
						}
						k.SetMarket(ctx, market)
					}
				}
			}
		}
		return nil
	})
}
