package liquidation

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/petrichormoney/petri/x/liquidation/keeper"
	"github.com/petrichormoney/petri/x/liquidation/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
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
