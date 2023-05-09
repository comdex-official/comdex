package keeper

import (
	"github.com/comdex-official/comdex/x/auctionsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k Keeper) PlaceDutchAuctionBid(ctx sdk.Context, auctionID uint64, bidder sdk.AccAddress, bid sdk.Coin, auctionData types.Auctions) error {
	



	return nil
}

func (k Keeper) PlaceEnglishAuctionBid(ctx sdk.Context, auctionID uint64, bidder sdk.AccAddress, bid sdk.Coin, auctionData types.Auctions) error {
	return nil
}

func (k Keeper) SetLimitAuctionBidID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LimitAuctionBidIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetLimitAuctionBidID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LimitAuctionBidIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) PlaceLimitAuctionBid(ctx sdk.Context, bidder string, CollateralTokenId, DebtTokenId uint64, PremiumDiscount string, amount sdk.Coin) error {
	id := k.GetLimitAuctionBidID(ctx)
	bidderAddr, err := sdk.AccAddressFromBech32(bidder)
	if err != nil {
		return nil
	}
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, bidderAddr, types.ModuleName, sdk.NewCoins(amount))
	if err != nil {
		return err
	}
	premiumDiscount, err := sdk.NewDecFromStr(PremiumDiscount)
	if err != nil {
		return err
	}

	limitBid := types.LimitOrderBid{
		LimitOrderBiddingId: id + 1,
		BidderAddress:       bidder,
		CollateralToken:     sdk.Coin{},
		DebtToken:           sdk.Coin{},
		BiddingId:           nil,
		PremiumDiscount:     premiumDiscount,
	}
	k.SetLimitAuctionBidID(ctx, limitBid.LimitOrderBiddingId)
	return nil
}

func (k Keeper) CancelLimitAuctionBid(ctx sdk.Context, bidder sdk.AccAddress, CollateralTokenId, DebtTokenId uint64) error {
	return nil
}

func (k Keeper) WithdrawLimitAuctionBid(ctx sdk.Context, bidder sdk.AccAddress, CollateralTokenId, DebtTokenId uint64, PremiumDiscount string) error {
	return nil
}
