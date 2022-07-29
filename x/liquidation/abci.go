package liquidation

import (
	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidation/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		err := k.LiquidateVaults(ctx)
		if err != nil {
			return err
		}
		err = k.UpdateLockedVaults(ctx)
		if err != nil {
			return err
		}
		err = k.LiquidateBorrows(ctx)
		if err != nil {
			return err
		}
		err = k.UpdateLockedBorrows(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
