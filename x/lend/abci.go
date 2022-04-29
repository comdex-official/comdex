package lend

import (
	"fmt"
	"github.com/comdex-official/comdex/x/lend/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	fmt.Println(k.GetAsset(ctx, 3))
	fmt.Println(k.GetPairs(ctx))
	fmt.Println(k.GetUserLendHistory(ctx,1))
}
