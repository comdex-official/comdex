package keeper

import (
	"time"

	"github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewEpochInfo return new EpochInfo object.
func (k Keeper) NewEpochInfo(ctx sdk.Context, duration time.Duration) types.EpochInfo {
	return types.EpochInfo{
		StartTime:               time.Time{},
		Duration:                duration,
		CurrentEpoch:            0,
		CurrentEpochStartTime:   ctx.BlockTime(),
		CurrentEpochStartHeight: 0,
	}
}

// TriggerAndUpdateEpochInfos updated the existing epoch and initiates the task if it is to be triggered.
/*func (k Keeper) TriggerAndUpdateEpochInfos(ctx sdk.Context) {
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
			missedEpochs := 2
			for {
				if epoch.CurrentEpochStartTime.Add(epoch.Duration * time.Duration(missedEpochs)).Before(ctx.BlockTime()) {
					missedEpochs++
				} else {
					epoch.CurrentEpochStartTime = epoch.CurrentEpochStartTime.Add(epoch.Duration * time.Duration(missedEpochs-1))
					break
				}
			}

			k.SetEpochInfoByDuration(ctx, epoch)
			continue
		}

		shouldTrigger := ctx.BlockTime().After(epoch.CurrentEpochStartTime.Add(epoch.Duration))
		if shouldTrigger {
			logger.Info(fmt.Sprintf("Starting new epoch with duration %s and epoch number %d", &epoch.Duration, epoch.CurrentEpoch))
			err := k.InitateGaugesForDuration(ctx, epoch.Duration)
			if err != nil {
				logger.Info(fmt.Sprintf("err occurred in epoch trigger : %v", err))
			}
			epoch.CurrentEpoch++
			epoch.CurrentEpochStartTime = epoch.CurrentEpochStartTime.Add(epoch.Duration)
			epoch.CurrentEpochStartHeight = ctx.BlockHeight()

			k.SetEpochInfoByDuration(ctx, epoch)
		}
	}
}*/
