package oracle

import (
	"github.com/comdex-official/comdex/x/oracle/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)


func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	if ctx.BlockHeight()>=59 && ctx.BlockHeight()%20==0{
	k.SetRates(ctx, "symbol")
	}
}
