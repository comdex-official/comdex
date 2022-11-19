package esm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/petrichormoney/petri/x/esm/keeper"
	"github.com/petrichormoney/petri/x/esm/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	k.SetParams(ctx, state.Params)

	for _, item := range state.ESMTriggerParams {
		k.SetESMTriggerParams(ctx, item)
	}

	for _, item := range state.CurrentDepositStats {
		k.SetCurrentDepositStats(ctx, item)
	}

	for _, item := range state.ESMStatus {
		k.SetESMStatus(ctx, item)
	}

	for _, item := range state.KillSwitchParams {
		err := k.SetKillSwitchData(ctx, item)
		if err != nil {
			return
		}
	}

	for _, item := range state.UsersDepositMapping {
		k.SetUserDepositByApp(ctx, item)
	}

	for _, item := range state.DataAfterCoolOff {
		k.SetDataAfterCoolOff(ctx, item)
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetAllESMTriggerParams(ctx),
		k.GetAllCurrentDepositStats(ctx),
		k.GetAllESMStatus(ctx),
		k.GetAllKillSwitchData(ctx),
		k.GetAllUserDepositByApp(ctx),
		k.GetAllDataAfterCoolOff(ctx),
		k.GetParams(ctx),
	)
}
