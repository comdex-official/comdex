package market

import (
	utils "github.com/comdex-official/comdex/types"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	bandkeeper "github.com/comdex-official/comdex/x/bandoracle/keeper"
	bandoraclemoduletypes "github.com/comdex-official/comdex/x/bandoracle/types"
	"github.com/comdex-official/comdex/x/market/keeper"
	"github.com/comdex-official/comdex/x/market/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper, bandKeeper bandkeeper.Keeper, assetKeeper assetkeeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		if bandKeeper.GetOracleValidationResult(ctx) {
			block := bandKeeper.GetLastBlockHeight(ctx)
			if block != types.Int64Zero {
				if ctx.BlockHeight()%types.Int64Twenty == types.Int64Zero && ctx.BlockHeight() != block && bandKeeper.GetCheckFlag(ctx) {
					assets := assetKeeper.GetAssets(ctx)
					id := bandKeeper.GetLastFetchPriceID(ctx)
					data, _ := bandKeeper.GetFetchPriceResult(ctx, bandoraclemoduletypes.OracleRequestID(id))
					scriptID := bandKeeper.GetFetchPriceMsg(ctx).OracleScriptID
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
			assets := assetKeeper.GetAssets(ctx)
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
