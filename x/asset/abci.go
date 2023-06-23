package asset

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/asset/keeper"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		return nil
	})
}
