package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) PlaceDutchAuctionBid(ctx sdk.Context, appID, auctionID uint64, bidder sdk.AccAddress, bid sdk.Coin) error {
	return nil
}
