package auction

import (
	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/auction/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		err1 := k.SurplusActivator(ctx)
		if err1 != nil {
			return err1
		}
		err2 := k.DebtActivator(ctx)
		if err2 != nil {
			return err2
		}
		err3 := k.DutchActivator(ctx)
		if err3 != nil {
			return err3
		}
		err4 := k.RestartDutch(ctx)
		if err4 != nil {
			return err4
		}
		err5 := k.LendDutchActivator(ctx)
		if err5 != nil {
			return err5
		}
		err6 := k.RestartLendDutch(ctx)
		if err6 != nil {
			return err6
		}
		return nil
	})
}
