package liquidation

import (
	"github.com/comdex-official/comdex/x/liquidation/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	err := k.LiquidateVaults(ctx)
	if err != nil {
		return
	}
	err = k.UpdateLockedVaults(ctx)
	if err != nil {
		return
	}
}
