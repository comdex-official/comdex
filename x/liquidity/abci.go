package liquidity

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/liquidity/keeper"
	"github.com/comdex-official/comdex/x/liquidity/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	allApps, found := k.GetApps(ctx)
	if found {

		for _, app := range allApps {

			k.DeleteOutdatedRequests(ctx, app.Id)
			if ctx.BlockHeight()%150 == 0 {
				k.ConvertAccumulatedSwapFeesWithSwapDistrToken(ctx, app.Id)
			}
		}
	}
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	allApps, found := k.GetApps(ctx)
	if found {

		for _, app := range allApps {
			params, err := k.GetGenericParams(ctx, app.Id)
			if err != nil {
				continue
			}
			if ctx.BlockHeight()%int64(params.BatchSize) == 0 {
				k.ExecuteRequests(ctx, app.Id)
				k.ProcessQueuedLiquidityProviders(ctx, app.Id)
			}
		}
	}
}
