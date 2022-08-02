package lend

import (
	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/lend/keeper"
	"github.com/comdex-official/comdex/x/lend/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {

	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		err := k.IterateLends(ctx)
		if err != nil {
			return err
		}
		err = k.IterateBorrows(ctx)
		if err != nil {
			return err
		}
		err = k.RebalanceStableRates(ctx)
		if err != nil {
			return err
		}
		err = k.SetLastInterestTime(ctx, ctx.BlockTime().Unix())
		if err != nil {
			return err
		}
		return nil
	})
}
