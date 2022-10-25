package lend

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/lend/keeper"
	"github.com/comdex-official/comdex/x/lend/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		// err := k.ReBalanceStableRates(ctx)
		// if err != nil {
		// 	ctx.Logger().Error("error in ReBalance Stable Rates")
		// }
		if ctx.BlockHeight() == 153 {
			err2 := k.MigrateData(ctx)
			if err2 != nil {
				ctx.Logger().Error("error in Migrate Data")
			}
		}
		return nil
	})
}
