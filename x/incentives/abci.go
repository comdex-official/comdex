package incentives

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/incentives/keeper"
	"github.com/comdex-official/comdex/x/incentives/types"
)

// BeginBlocker for incentives module.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	k.TriggerAndUpdateEpochInfos(ctx)
}

// EndBlocker for incentives module.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {}
