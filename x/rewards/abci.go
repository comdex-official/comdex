package rewards

import (
	"time"

	"github.com/comdex-official/comdex/x/rewards/keeper"
	"github.com/comdex-official/comdex/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	k.TriggerAndUpdateEpochInfos(ctx)

	err := k.IterateLocker(ctx)
	if err != nil {
		return
	}

	appIDsVault := k.GetAppIDs(ctx).WhitelistedAppMappingIdsVaults
	for i := range appIDsVault {
		err := k.IterateVaults(ctx, appIDsVault[i])
		if err != nil {
			return
		}
	}

	err = k.DistributeExtRewardLocker(ctx)
	if err != nil {
		return
	}
	err = k.DistributeExtRewardVault(ctx)
	if err != nil {
		return
	}
	// nolint
	k.SetLastInterestTime(ctx, ctx.BlockTime().Unix())
}

// EndBlocker for incentives module.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {}
