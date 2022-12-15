package rewards

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/rewards/keeper"
	"github.com/comdex-official/comdex/x/rewards/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		k.TriggerAndUpdateEpochInfos(ctx)

		err := k.DistributeExtRewardLocker(ctx)
		if err != nil {
			ctx.Logger().Error("error in DistributeExtRewardLocker")
		}
		err = k.DistributeExtRewardVault(ctx)
		if err != nil {
			ctx.Logger().Error("error in DistributeExtRewardVault")
		}
		err = k.DistributeExtRewardLend(ctx)
		if err != nil {
			ctx.Logger().Error("error in DistributeExtRewardLend")
		}
		err = k.CombinePSMUserPositions(ctx)
		if err != nil {
			ctx.Logger().Error("error in CombinePSMUserPositions")
		}
		return nil
	})
}

// EndBlocker for incentives module.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {}
