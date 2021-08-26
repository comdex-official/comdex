package cdp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/cdp/keeper"
	"github.com/comdex-official/comdex/x/cdp/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	for _, item := range state.CDPs {
		k.SetCDP(ctx, item)
	}

	id := uint64(0)
	for _, item := range state.CDPs {
		if item.ID > id {
			id = item.ID
		}
	}

	k.SetID(ctx, id)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetCDPs(ctx),
	)
}
