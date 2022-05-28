package auction

import (
	"github.com/comdex-official/comdex/x/auction/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	//k.CreateNewAuctions(ctx)
	err := k.CloseAndRestartAuctions(ctx)
	if err != nil {
		return
	}
}
