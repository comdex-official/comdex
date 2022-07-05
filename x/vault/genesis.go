package vault

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/vault/keeper"
	"github.com/comdex-official/comdex/x/vault/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	// for _, item := range state.Vaults {
	// 	k.SetVault(ctx, item)
	// }

	// id := uint64(0)
	// for _, item := range state.Vaults {
	// 	if item.ID > id {
	// 		id = item.ID
	// 	}
	// }

	// k.SetID(ctx, id)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetVaults(ctx),
	)
}
