package liquidity

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	utils "github.com/comdex-official/comdex/types"
	expected "github.com/comdex-official/comdex/x/liquidity/expected"
	"github.com/comdex-official/comdex/x/liquidity/keeper"
	"github.com/comdex-official/comdex/x/liquidity/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper, assetKeeper expected.AssetKeeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	allApps, found := assetKeeper.GetApps(ctx)
	if found {
		for _, app := range allApps {
			_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
				k.DeleteOutdatedRequests(ctx, app.Id)
				if ctx.BlockHeight()%150 == 0 {
					k.ConvertAccumulatedSwapFeesWithSwapDistrToken(ctx, app.Id)
				}
				return nil
			})
		}
	}
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper, assetKeeper expected.AssetKeeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyEndBlocker)

	allApps, found := assetKeeper.GetApps(ctx)
	if found {
		for _, app := range allApps {
			_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
				params, err := k.GetGenericParams(ctx, app.Id)
				if err != nil {
					return err
				}
				if ctx.BlockHeight()%int64(params.BatchSize) == 0 {
					k.ExecuteRequests(ctx, app.Id)
					k.ProcessQueuedFarmers(ctx, app.Id)
				}
				k.BlockAllFarmCoinTransfers(ctx, app.Id)
				return nil
			})
		}
	}
}
