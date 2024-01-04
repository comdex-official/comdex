package rwa

import sdk "github.com/cosmos/cosmos-sdk/types"
import "github.com/comdex-official/comdex/x/rwa/keeper"
import "github.com/comdex-official/comdex/x/rwa/types"

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState()
}
