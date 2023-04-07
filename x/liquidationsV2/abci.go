package liquidationsV2

import (
	"fmt"
	"github.com/comdex-official/comdex/x/liquidationsV2/keeper"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	err := k.Liquidate(ctx)
	if err != nil {
		ctx.Logger().Error("error in Liquidate function")
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeLiquidateVaultsErr,
				sdk.NewAttribute(types.Error, fmt.Sprintf("%s", err)),
			),
		)
	}
}
