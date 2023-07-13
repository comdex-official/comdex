package keeper

import (
	"github.com/comdex-official/comdex/x/auctionsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k Keeper) SetAuctionID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.AuctionIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}
func (k Keeper) GetAuctionID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.AuctionIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) SetUserBidID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.UserBidIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k Keeper) GetUserBidID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.UserBidIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) SetAuction(ctx sdk.Context, auction types.Auction) error {

	var (
		store = k.Store(ctx)
		key   = types.AuctionKey(auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)

	store.Set(key, value)
	return nil
}

func (k Keeper) SetAuctionLimitBidFeeData(ctx sdk.Context, feeData types.AuctionFeesCollectionFromLimitBidTx) error {

	var (
		store = k.Store(ctx)
		key   = types.AuctionLimitBidFeeKey(feeData.AssetId)
		value = k.cdc.MustMarshal(&feeData)
	)

	store.Set(key, value)
	return nil
}

func (k Keeper) GetAuctionLimitBidFeeData(ctx sdk.Context, assetId uint64) (feeData types.AuctionFeesCollectionFromLimitBidTx, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AuctionLimitBidFeeKey(assetId)
		value = store.Get(key)
	)

	if value == nil {
		return feeData, false
	}

	k.cdc.MustUnmarshal(value, &feeData)
	return feeData, true
}

func (k Keeper) SetAuctionHistorical(ctx sdk.Context, auction types.AuctionHistorical) error {

	var (
		store = k.Store(ctx)
		key   = types.AuctionHistoricalKey(auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)

	store.Set(key, value)
	return nil
}

func (k Keeper) SetBidHistorical(ctx sdk.Context, bid types.Bid) error {

	var (
		store = k.Store(ctx)
		key   = types.BidHistoricalKey(bid.BiddingId, bid.BidderAddress)
		value = k.cdc.MustMarshal(&bid)
	)

	store.Set(key, value)
	return nil
}

func (k Keeper) SetIndividualUserBid(ctx sdk.Context, userBid types.Bid) error {

	var (
		store = k.Store(ctx)
		key   = types.UserBidHistoricalKey(userBid.BiddingId, userBid.BidderAddress)
		value = k.cdc.MustMarshal(&userBid)
	)

	store.Set(key, value)
	return nil
}

func (k Keeper) SetUserBid(ctx sdk.Context, userBid types.Bid) error {

	var (
		store = k.Store(ctx)
		key   = types.UserBidKey(userBid.BiddingId)
		value = k.cdc.MustMarshal(&userBid)
	)

	store.Set(key, value)
	return nil
}

func (k Keeper) DeleteIndividualUserBid(ctx sdk.Context, userBid types.Bid) error {

	var (
		store = k.Store(ctx)
		key   = types.UserBidHistoricalKey(userBid.BiddingId, userBid.BidderAddress)
	)
	store.Delete(key)
	return nil
}

func (k Keeper) DeleteAuction(ctx sdk.Context, auction types.Auction) error {

	var (
		store = k.Store(ctx)
		key   = types.AuctionKey(auction.AuctionId)
	)
	store.Delete(key)
	return nil
}

func (k Keeper) GetUserBid(ctx sdk.Context, userBidId uint64) (userBid types.Bid, err error) {
	var (
		store = k.Store(ctx)
		key   = types.UserBidKey(userBidId)
		value = store.Get(key)
	)

	if value == nil {
		return userBid, sdkerrors.ErrNotFound
	}

	k.cdc.MustUnmarshal(value, &userBid)
	return userBid, nil
}

func (k Keeper) GetAuction(ctx sdk.Context, auctionID uint64) (auction types.Auction, err error) {
	var (
		store = k.Store(ctx)
		key   = types.AuctionKey(auctionID)
		value = store.Get(key)
	)

	if value == nil {
		return auction, sdkerrors.ErrNotFound
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, nil
}

func (k Keeper) GetAuctionHistorical(ctx sdk.Context, auctionID uint64) (auction types.AuctionHistorical, err error) {
	var (
		store = k.Store(ctx)
		key   = types.AuctionHistoricalKey(auctionID)
		value = store.Get(key)
	)

	if value == nil {
		return auction, sdkerrors.ErrNotFound
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, nil
}

func (k Keeper) GetAuctions(ctx sdk.Context) (auctions []types.Auction) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.AuctionKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var auction types.Auction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k Keeper) GetAuctionHistoricals(ctx sdk.Context) (auctions []types.AuctionHistorical) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.AuctionHistoricalKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var auction types.AuctionHistorical
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k Keeper) GetUserBids(ctx sdk.Context) (userBids []types.Bid) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.UserBidKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var userBid types.Bid
		k.cdc.MustUnmarshal(iter.Value(), &userBid)
		userBids = append(userBids, userBid)
	}

	return userBids
}

func (k Keeper) SetAuctionLimitBidFeeDataExternal(ctx sdk.Context, feeData types.AuctionFeesCollectionFromLimitBidTx) error {

	var (
		store = k.Store(ctx)
		key   = types.ExternalAuctionLimitBidFeeKey(feeData.AssetId)
		value = k.cdc.MustMarshal(&feeData)
	)

	store.Set(key, value)
	return nil
}
func (k Keeper) GetAuctionLimitBidFeeDataExternal(ctx sdk.Context, assetId uint64) (feeData types.AuctionFeesCollectionFromLimitBidTx, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.ExternalAuctionLimitBidFeeKey(assetId)
		value = store.Get(key)
	)

	if value == nil {
		return feeData, false
	}

	k.cdc.MustUnmarshal(value, &feeData)
	return feeData, true
}
