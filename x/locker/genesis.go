package locker

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/petrichormoney/petri/x/locker/keeper"
	"github.com/petrichormoney/petri/x/locker/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	k.SetParams(ctx, state.Params)

	for _, item := range state.Lockers {
		k.SetLocker(ctx, item)
	}

	for _, item := range state.LockerProductAssetMapping {
		k.SetLockerProductAssetMapping(ctx, item)
	}

	for _, item := range state.LockerTotalRewardsByAssetAppWise {
		err := k.SetLockerTotalRewardsByAssetAppWise(ctx, item)
		if err != nil {
			return
		}
	}

	for _, item := range state.LockerLookupTable {
		k.SetLockerLookupTable(ctx, item)
	}

	for _, item := range state.UserLockerAssetMapping {
		k.SetUserLockerAssetMapping(ctx, item)
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetLockers(ctx),
		k.GetAllLockerProductAssetMapping(ctx),
		k.GetAllLockerTotalRewardsByAssetAppWise(ctx),
		k.GetAllLockerLookupTable(ctx),
		k.GetAllUserLockerAssetMapping(ctx),
		k.GetParams(ctx),
	)
}
