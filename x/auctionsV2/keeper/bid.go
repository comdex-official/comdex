package keeper

import (
	"github.com/comdex-official/comdex/x/auctionsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) PlaceDutchAuctionBid(ctx sdk.Context, auctionID uint64, bidder sdk.AccAddress, bid sdk.Coin, auctionData types.Auctions) error {
	



	return nil
}

func (k Keeper) PlaceEnglishAuctionBid(ctx sdk.Context, auctionID uint64, bidder sdk.AccAddress, bid sdk.Coin, auctionData types.Auctions) error {
	return nil
}
