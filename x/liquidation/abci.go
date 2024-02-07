package liquidation

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/liquidation/keeper"
	"github.com/comdex-official/comdex/x/liquidation/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	err := k.LiquidateVaults(ctx)
	if err != nil {
		ctx.Logger().Error("error in LiquidateVaults")
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeLiquidateVaultsErr,
				sdk.NewAttribute(types.Error, fmt.Sprintf("%s", err)),
			),
		)
	}
	err = k.LiquidateBorrows(ctx)
	if err != nil {
		ctx.Logger().Error("error in LiquidateBorrows")
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeLiquidateBorrowsErr,
				sdk.NewAttribute(types.Error, fmt.Sprintf("%s", err)),
			),
		)
	}
}
