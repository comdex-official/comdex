package auction

import (
	"github.com/comdex-official/comdex/x/auction/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	err1 := k.SurplusActivator(ctx)
	if err1 != nil {
		return
	}
	err2 := k.DebtActivator(ctx)
	if err2 != nil {
		return
	}
	err3 := k.DutchActivator(ctx)
	if err3 != nil {
		return
	}
	err4 := k.RestartDutch(ctx)
	if err4 != nil {
		return
	}
	err5 := k.LendDutchActivator(ctx)
	if err5 != nil {
		return
	}
	err6 := k.RestartLendDutch(ctx)
	if err6 != nil {
		return
	}
}
