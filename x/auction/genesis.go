package auction

import (
	"github.com/comdex-official/comdex/x/auction/keeper"
	"github.com/comdex-official/comdex/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {

	k.SetParams(ctx, state.Params)

	for _, item := range state.SurplusAuction {
		err := k.SetSurplusAuction(ctx, item)
		if err != nil {
			return
		}
	}

	for _, item := range state.DebtAuction {
		err := k.SetDebtAuction(ctx, item)
		if err != nil {
			return
		}
	}

	for _, item := range state.DutchAuction {
		err := k.SetDutchAuction(ctx, item)
		if err != nil {
			return
		}
	}

	for _, item := range state.ProtocolStatistics {
		k.SetProtocolStatistics(ctx, item.AppId, item.AssetId, sdk.Int(item.Loss))
	}

	for _, item := range state.AuctionParams {
		k.SetAuctionParams(ctx, item)
	}

}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {

	return types.NewGenesisState(
		k.GetAllSurplusAuctions(ctx),
		k.GetAllDebtAuctions(ctx),
		k.GetAllDutchAuctions(ctx),
		k.GetAllProtocolStat(ctx),
		k.GetAllAuctionParams(ctx),
		k.GetParams(ctx),
	)
}
