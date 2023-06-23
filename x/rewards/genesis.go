package rewards

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/rewards/keeper"
	"github.com/comdex-official/comdex/x/rewards/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	var (
		gaugeID       uint64
		lendRewardsID uint64
	)

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

	for _, item := range state.LendExternalRewards {
		if item.Id > lendRewardsID {
			lendRewardsID = item.Id
		}
		k.SetExternalRewardLend(ctx, item)
	}

	k.SetGaugeID(ctx, gaugeID)
	k.SetExternalRewardsLendID(ctx, lendRewardsID)
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
		k.GetExternalRewardLends(ctx),
	)
}
