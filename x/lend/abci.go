package lend

import (
	"github.com/comdex-official/comdex/x/lend/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	err := k.IterateLends(ctx)
	if err != nil {
		return
	}
	err = k.IterateBorrows(ctx)
	if err != nil {
		return
	}
	err = k.RebalanceStableRates(ctx)
	if err != nil {
		return
	}
	err = k.SetLastInterestTime(ctx, ctx.BlockTime().Unix())
	if err != nil {
		return
	}
}
