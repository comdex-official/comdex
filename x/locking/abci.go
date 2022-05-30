package locking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/locking/keeper"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	k.DeleteMaturedUnlocks(ctx)
}
