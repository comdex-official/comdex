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

	// pools := k.GetAllPools(ctx)

	// for _, pool := range pools {
	// 	data, found := k.GetPoolLiquidityProvidersData(ctx, pool.Id)
	// 	if found {
	// 		fmt.Println("QUEUED........")
	// 		sort.Slice(data.QueuedLiquidityProviders, func(i, j int) bool {
	// 			return data.QueuedLiquidityProviders[i].CreatedAt.After(data.QueuedLiquidityProviders[j].CreatedAt)
	// 		})
	// 		for _, queued := range data.QueuedLiquidityProviders {
	// 			fmt.Println(queued.CreatedAt, " , ", queued.Address, " , ", queued.SupplyProvided)
	// 		}

	// 		fmt.Println("\nValid Liquidity Providers .........")
	// 		for k, v := range data.LiquidityProviders {
	// 			fmt.Println(k, " : ", v)
	// 		}
	// 	}
	// }

	k.DeleteOutdatedRequests(ctx)
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	params := k.GetParams(ctx)
	if ctx.BlockHeight()%int64(params.BatchSize) == 0 {
		k.ExecuteRequests(ctx)
		k.ProcessQueuedLiquidityProviders(ctx)
	}
}
