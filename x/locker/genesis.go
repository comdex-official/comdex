package locker

import (
	"github.com/comdex-official/comdex/x/locker/keeper"
	"github.com/comdex-official/comdex/x/locker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		k.SetLockerTotalRewardsByAssetAppWise(ctx, item)
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
