package esm

import (
	"github.com/comdex-official/comdex/x/esm/keeper"
	"github.com/comdex-official/comdex/x/esm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		k.SetKillSwitchData(ctx, item)
	}

	for _, item := range state.UsersDepositMapping {
		k.SetUserDepositByApp(ctx, item)
	}

	for _, item := range state.ESMMarketPrice {
		k.SetESMMarketForAsset(ctx, item)
	}

	for _, item := range state.DataAfterCoolOff {
		k.SetDataAfterCoolOff(ctx, item)
	}

	for _, item := range state.AssetToAmountValue {
		k.SetAssetToAmountValue(ctx, item)
	}

	for _, item := range state.AppToAmountValue {
		k.SetAppToAmtValue(ctx, item)
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {

	return types.NewGenesisState(
		k.GetAllESMTriggerParams(ctx),
		k.GetAllCurrentDepositStats(ctx),
		k.GetAllESMStatus(ctx),
		k.GetAllKillSwitchData(ctx),
		k.GetAllUserDepositByApp(ctx),
		k.GetAllESMMarketForAsset(ctx),
		k.GetAllDataAfterCoolOff(ctx),
		k.GetAllAssetToAmountValue(ctx),
		k.GetAllAppToAmtValue(ctx),
		k.GetParams(ctx),
	)
}
