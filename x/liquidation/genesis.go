package liquidation

import (
	"github.com/comdex-official/comdex/x/liquidation/keeper"
	"github.com/comdex-official/comdex/x/liquidation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {

	k.SetParams(ctx, state.Params)

	for _, item := range state.LockedVault {

		k.SetLockedVault(ctx, item)
	}

	for _, item := range state.LockedVaultToAppMapping {

		k.SetLockedVaultByAppID(ctx, item)
	}

	k.SetAppID(ctx, state.WhitelistedAppIds)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetLockedVaults(ctx),
		k.GetAllLockedVaultByAppID(ctx),
		k.GetAppIds(ctx),
		k.GetParams(ctx),
	)
}
