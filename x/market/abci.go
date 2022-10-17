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
				if ctx.BlockHeight()%types.Int64Twenty == types.Int64Zero && ctx.BlockHeight() != block && k.GetCheckFlag(ctx) {
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
						}
					}
				}
			}
		} else {
			assets := k.GetAssets(ctx)
			for _, asset := range assets {
				twa, found := k.GetTwa(ctx, asset.Id)
				if !found {
					continue
				}
				twa.IsPriceActive = false
				k.SetTwa(ctx, twa)
			}
		}
		return nil
	})
}
