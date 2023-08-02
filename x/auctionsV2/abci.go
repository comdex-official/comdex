package auctionsV2

import (
	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/auctionsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/auctionsV2/keeper"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		err := k.AuctionIterator(ctx)
		if err != nil {
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeAuctionIteratorErr,
				),
			)
			ctx.Logger().Error("error in Auction iterator")
		}
		return nil
	})

	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		err := k.LimitOrderBid(ctx)
		if err != nil {
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeLimitOrderBidErr,
				),
			)
			ctx.Logger().Error("error in Auction Limit order bid")
		}
		return nil
	})
}
