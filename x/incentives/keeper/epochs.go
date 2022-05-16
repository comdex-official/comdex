package keeper

import (
	"fmt"
	"time"

	"github.com/comdex-official/comdex/x/incentives/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) NewEpochInfo(ctx sdk.Context, duration time.Duration) types.EpochInfo {
	return types.EpochInfo{
		StartTime:               time.Time{},
		Duration:                duration,
		CurrentEpoch:            0,
		CurrentEpochStartTime:   ctx.BlockTime(),
		CurrentEpochStartHeight: 0,
	}
}

func (k Keeper) TriggerAndUpdateEpochInfos(ctx sdk.Context) {
	logger := k.Logger(ctx)

	epochInfos := k.GetAllEpochInfos(ctx)

	for _, epoch := range epochInfos {

		isFreshNewEpoch := epoch.StartTime == time.Time{} && epoch.CurrentEpoch == 0

		if isFreshNewEpoch {
			epoch.StartTime = ctx.BlockTime()
			epoch.CurrentEpochStartTime = epoch.CurrentEpochStartTime.Add(-epoch.Duration)
			k.SetEpochInfoByDuration(ctx, epoch)
			continue
		}

		// In case of chain halt/stop
		if epoch.CurrentEpochStartTime.Add(epoch.Duration * 2).Before(ctx.BlockTime()) {
			epoch.CurrentEpochStartTime = ctx.BlockTime()
			k.SetEpochInfoByDuration(ctx, epoch)
			continue
		}

		shouldTrigger := ctx.BlockTime().After(epoch.CurrentEpochStartTime.Add(epoch.Duration))
		if shouldTrigger {
			logger.Info(fmt.Sprintf("Starting new epoch with duration %d and epoch number %d", &epoch.Duration, epoch.CurrentEpoch))
			fmt.Println("Triggering epoch", ctx.BlockTime())
			epoch.CurrentEpoch += 1
			epoch.CurrentEpochStartTime = epoch.CurrentEpochStartTime.Add(epoch.Duration)
			epoch.CurrentEpochStartHeight = ctx.BlockHeight()

			k.SetEpochInfoByDuration(ctx, epoch)
		}
	}

}
