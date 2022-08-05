package liquidation

import (
	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidation/keeper"
	"github.com/comdex-official/comdex/x/liquidation/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {

	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		err := k.LiquidateVaults(ctx)
		if err != nil {
			ctx.Logger().Error("error in LiquidateVaults")
		}
		err = k.UpdateLockedVaults(ctx)
		if err != nil {
			ctx.Logger().Error("error in UpdateLockedVaults")
		}
		err = k.LiquidateBorrows(ctx)
		if err != nil {
			ctx.Logger().Error("error in LiquidateBorrows")
		}
		err = k.UpdateLockedBorrows(ctx)
		if err != nil {
			ctx.Logger().Error("error in UpdateLockedBorrows")
		}
		return nil
	})
}
