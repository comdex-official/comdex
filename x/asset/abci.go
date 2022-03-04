package asset

import (
	"fmt"
	"github.com/comdex-official/comdex/x/asset/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, k keeper.Keeper) {
	fmt.Println(k.GetParams(ctx))
}
