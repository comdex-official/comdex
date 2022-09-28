package liquidation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/liquidation/keeper"
	"github.com/comdex-official/comdex/x/liquidation/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	k.SetParams(ctx, state.Params)

	for _, item := range state.LockedVault {
		k.SetLockedVault(ctx, item)
	}

	for _, item := range state.WhitelistedApps {
		k.SetAppIDForLiquidation(ctx, item)
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetLockedVaults(ctx),
		k.GetAppIdsForLiquidation(ctx),
		k.GetParams(ctx),
	)
}
