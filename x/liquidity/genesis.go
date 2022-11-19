package liquidity

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/petrichormoney/petri/x/liquidity/keeper"
	"github.com/petrichormoney/petri/x/liquidity/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.InitGenesis(ctx, genState)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return k.ExportGenesis(ctx)
}
