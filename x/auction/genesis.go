package auction

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/petrichormoney/petri/x/auction/keeper"
	"github.com/petrichormoney/petri/x/auction/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	k.SetParams(ctx, state.Params)
	var auctionID uint64
	var lendAuctionID uint64

	for _, item := range state.SurplusAuction {
		k.SetGenSurplusAuction(ctx, item)
	}

	for _, item := range state.DebtAuction {
		k.SetGenDebtAuction(ctx, item)
	}

	for _, item := range state.DutchAuction {
		k.SetGenDutchAuction(ctx, item)
		auctionID = item.AuctionId
	}

	k.SetAuctionID(ctx, auctionID)
	k.SetUserBiddingID(ctx, state.UserBiddingID)

	for _, item := range state.ProtocolStatistics {
		k.SetGenProtocolStatistics(ctx, item.AppId, item.AssetId, item.Loss)
	}

	for _, item := range state.AuctionParams {
		k.SetAuctionParams(ctx, item)
	}

	for _, item := range state.DutchAuction {
		k.SetGenLendDutchLendAuction(ctx, item)
		lendAuctionID = item.AuctionId
	}
	k.SetLendAuctionID(ctx, lendAuctionID)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetAllSurplusAuctions(ctx),
		k.GetAllDebtAuctions(ctx),
		k.GetAllDutchAuctions(ctx),
		k.GetAllProtocolStat(ctx),
		k.GetAllAuctionParams(ctx),
		k.GetDutchLendAuctions(ctx, 3),
		k.GetParams(ctx),
		k.GetUserBiddingID(ctx),
	)
}
