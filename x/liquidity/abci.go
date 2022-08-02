package liquidity

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidity/keeper"
	"github.com/comdex-official/comdex/x/liquidity/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		allApps, found := k.GetApps(ctx)
		if found {
			for _, app := range allApps {
				k.DeleteOutdatedRequests(ctx, app.Id)
				if ctx.BlockHeight()%150 == 0 {
					k.ConvertAccumulatedSwapFeesWithSwapDistrToken(ctx, app.Id)
				}
			}
		}
		return nil
	})
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyEndBlocker)

	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		allApps, found := k.GetApps(ctx)
		if found {
			for _, app := range allApps {
				params, err := k.GetGenericParams(ctx, app.Id)
				if err != nil {
					continue
				}
				if ctx.BlockHeight()%int64(params.BatchSize) == 0 {
					k.ExecuteRequests(ctx, app.Id)
					k.ProcessQueuedFarmers(ctx, app.Id)
				}
			}
		}
		return nil
	})
}
