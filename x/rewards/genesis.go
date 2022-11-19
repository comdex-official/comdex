package rewards

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/petrichormoney/petri/x/rewards/keeper"
	"github.com/petrichormoney/petri/x/rewards/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	var gaugeID uint64

	k.SetParams(ctx, state.Params)

	for _, item := range state.InternalRewards {
		k.SetReward(ctx, item)
	}

	for _, item := range state.LockerRewardsTracker {
		k.SetLockerRewardTracker(ctx, item)
	}

	for _, item := range state.VaultInterestTracker {
		k.SetVaultInterestTracker(ctx, item)
	}

	for _, item := range state.LockerExternalRewards {
		k.SetExternalRewardsLockers(ctx, item)
	}

	for _, item := range state.VaultExternalRewards {
		k.SetExternalRewardVault(ctx, item)
	}

	for _, item := range state.AppIDs {
		k.SetAppByAppID(ctx, item)
	}

	for _, item := range state.EpochInfo {
		k.SetEpochInfoByDuration(ctx, item)
	}

	for _, item := range state.Gauge {
		if item.Id > gaugeID {
			gaugeID = item.Id
		}
		k.SetGauge(ctx, item)
	}

	for _, item := range state.GaugeByTriggerDuration {
		k.SetGaugeIdsByTriggerDuration(ctx, item)
	}

	k.SetGaugeID(ctx, gaugeID)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetRewards(ctx),
		k.GetAllLockerRewardTracker(ctx),
		k.GetAllVaultInterestTracker(ctx),
		k.GetExternalRewardsLockers(ctx),
		k.GetExternalRewardVaults(ctx),
		k.GetAppIDs(ctx),
		k.GetAllEpochInfos(ctx),
		k.GetAllGauges(ctx),
		k.GetAllGaugeIdsByTriggerDuration(ctx),
		k.GetParams(ctx),
	)
}
