package auctionsV2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/auctionsV2/keeper"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	err := k.AuctionIterator(ctx)
	if err != nil {
		return
	}
	err = k.LimitOrderBid(ctx)
	if err != nil {
		return
	}
}
