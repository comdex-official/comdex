package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
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

func (k Keeper) SetUserLimitBidData(ctx sdk.Context, mappingData types.LimitOrderBid, collateralTokenID uint64, debtTokenID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.UserLimitBidKey(mappingData.BidderAddress, collateralTokenID, debtTokenID)
		value = k.cdc.MustMarshal(&mappingData)
	)

	store.Set(key, value)
}

func (k Keeper) GetUserLimitBidData(ctx sdk.Context, address string, collateralTokenID uint64, debtTokenID uint64) (mappingData types.LimitOrderBid, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserLimitBidKey(address, collateralTokenID, debtTokenID)
		value = store.Get(key)
	)

	if value == nil {
		return mappingData, false
	}

	k.cdc.MustUnmarshal(value, &mappingData)
	return mappingData, true
}

func (k Keeper) DepositLimitAuctionBid(ctx sdk.Context, bidder string, CollateralTokenId, DebtTokenId uint64, PremiumDiscount string, amount sdk.Coin) error {
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
	collateralAsset, found := k.asset.GetAsset(ctx, CollateralTokenId)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	collateralAssetToken := sdk.NewCoin(collateralAsset.Denom, sdk.NewInt(0))
	userLimitBid, found := k.GetUserLimitBidData(ctx, bidder, CollateralTokenId, DebtTokenId)
	if !found {
		userLimitBid = types.LimitOrderBid{
			LimitOrderBiddingId: id + 1,
			BidderAddress:       bidder,
			CollateralToken:     collateralAssetToken, // zero
			DebtToken:           amount,               // user's balance
			BiddingId:           nil,
			PremiumDiscount:     premiumDiscount,
		}
	} else {
		userLimitBid.CollateralToken = userLimitBid.CollateralToken.Add(amount)
	}

	k.SetLimitAuctionBidID(ctx, userLimitBid.LimitOrderBiddingId)
	return nil
}

func (k Keeper) CancelLimitAuctionBid(ctx sdk.Context, bidder sdk.AccAddress, CollateralTokenId, DebtTokenId uint64) error {
	return nil
}

func (k Keeper) WithdrawLimitAuctionBid(ctx sdk.Context, bidder string, CollateralTokenId, DebtTokenId uint64) error {
	return nil
}
