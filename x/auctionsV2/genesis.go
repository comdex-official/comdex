package auctionsV2

import (
	"github.com/comdex-official/comdex/x/auctionsV2/keeper"
	"github.com/comdex-official/comdex/x/auctionsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	var auctionID, userBidID uint64
	for _, item := range genState.Auction {
		k.SetGenAuction(ctx, item)
	}

	k.SetAuctionParams(ctx, genState.AuctionParams)

	for _, item := range genState.AuctionFeesCollectionFromLimitBidTx {
		k.SetGenAuctionLimitBidFeeData(ctx, item)
	}

	k.SetParams(ctx, genState.Params)
	k.SetAuctionID(ctx, auctionID)
	k.SetUserBidID(ctx, userBidID)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	auctionParams, _ := k.GetAuctionParams(ctx)
	return types.NewGenesisState(
		k.GetAuctions(ctx),
		auctionParams,
		k.GetGenAuctionLimitBidFeeData(ctx),
		k.GetParams(ctx),
		k.GetAuctionID(ctx),
		k.GetUserBidID(ctx),
	)
}
