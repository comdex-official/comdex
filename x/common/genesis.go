package common

import (
	"github.com/comdex-official/comdex/x/common/keeper"
	"github.com/comdex-official/comdex/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState *types.GenesisState) {
	var (
		GameID uint64
	)
	// this line is used by starport scaffolding # genesis/module/init
	for _, item := range genState.WhitelistedContracts {
		if item.GameId > GameID {
			GameID = item.GameId
		}
		k.SetContract(ctx, item)
	}
	k.SetGameID(ctx, GameID)
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetAllContract(ctx),
		k.GetParams(ctx),
	)
}
