package liquidationsV2

import (
	"github.com/comdex-official/comdex/x/liquidationsV2/keeper"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	var lockedVaultID uint64
	for _, item := range genState.LockedVault {
		k.SetLockedVault(ctx, item)
		lockedVaultID = lockedVaultID + 1
	}

	for _, item := range genState.LiquidationWhiteListing {
		k.SetLiquidationWhiteListing(ctx, item)
	}

	for _, item := range genState.AppReserveFunds {
		k.SetAppReserveFunds(ctx, item)
	}

	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetLockedVaults(ctx),
		k.GetGenLiquidationWhiteListing(ctx),
		k.GetGenAppReserveFunds(ctx),
		k.GetParams(ctx),
	)
}
