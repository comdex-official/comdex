package auction

import (
	"github.com/comdex-official/comdex/x/auction/keeper"
	"github.com/comdex-official/comdex/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {

	k.SetParams(ctx, state.Params)

	for _, item := range state.SurplusAuction {
		k.SetSurplusAuction(ctx, item)
	}

	for _, item := range state.DebtAuction {
		k.SetDebtAuction(ctx, item)
	}

	for _, item := range state.DutchAuction {
		k.SetDutchAuction(ctx, item)
	}

	for _, item := range state.ProtocolStatistics {
		k.SetProtocolStatistics(ctx, item.AppId, item.AssetId, sdk.Int(item.Loss))
	}

	for _, item := range state.AuctionParams {
		k.SetAuctionParams(ctx, item)
	}

	// for _, item := range state.SurplusBiddings {
	// 	k.SetSurplusUserBidding(ctx, item)
	// }

	// for _, item := range state.DebtBiddings {
	// 	k.SetDebtUserBidding(ctx, item)
	// }

	// for _, item := range state.DutchBiddings {
	// 	k.SetDutchUserBidding(ctx, item)
	// }

}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {

	return types.NewGenesisState(
		k.GetAllSurplusAuctions(ctx),
		k.GetAllDebtAuctions(ctx),
		k.GetAllDutchAuctions(ctx),
		k.GetAllProtocolStat(ctx),
		k.GetAllAuctionParams(ctx),
		// k.GetAllSurplusUserBiddings(ctx),
		// k.GetAllDebtUserBidding(ctx),
		// k.GetAllDutchUserBiddings(ctx),
		k.GetParams(ctx),
	)
}

