package liquidation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/petrichormoney/petri/x/liquidation/keeper"
	"github.com/petrichormoney/petri/x/liquidation/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	k.SetParams(ctx, state.Params)
	var lockedVaultID uint64

	for _, item := range state.LockedVault {
		k.SetLockedVault(ctx, item)
		lockedVaultID = lockedVaultID + 1
	}

	for _, item := range state.WhitelistedApps {
		k.SetAppIDForLiquidation(ctx, item)
	}

	k.SetLockedVaultID(ctx, lockedVaultID)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetLockedVaults(ctx),
		k.GetAppIdsForLiquidation(ctx),
		k.GetParams(ctx),
	)
}
