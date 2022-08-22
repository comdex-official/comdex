package rewards

import (
	"github.com/comdex-official/comdex/x/rewards/keeper"
	"github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {

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
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetRewards(ctx),
		k.GetAllLockerRewardTracker(ctx),
		k.GetAllVaultInterestTracker(ctx),
		k.GetExternalRewardsLockers(ctx),
		k.GetExternalRewardVaults(ctx),
		k.GetParams(ctx),
	)
}
