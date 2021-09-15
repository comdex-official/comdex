package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/oracle/keeper"
	"github.com/comdex-official/comdex/x/oracle/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {

	for _, item := range state.Markets {
		k.SetMarket(ctx, item)
	}

}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetMarkets(ctx),
		k.GetParams(ctx),
	)
}
