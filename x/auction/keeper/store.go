package keeper

import (
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k *Keeper) GetSurplusAuctionID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.CollateralAuctionIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetSurplusAuctionID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.CollateralAuctionIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetDebtAuctionID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtAuctionIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetDebtAuctionID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtAuctionIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetDutchAuctionID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchAuctionIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetDutchAuctionID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchAuctionIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) SetSurplusAuction(ctx sdk.Context, auction auctiontypes.SurplusAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.CollateralAuctionKey(auction.Id)
		value = k.cdc.MustMarshal(&auction)
	)
	store.Set(key, value)
}

func (k *Keeper) DeleteSurplusAuction(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.CollateralAuctionKey(id)
	)
	store.Delete(key)
}

func (k *Keeper) SetDebtAuction(ctx sdk.Context, auction auctiontypes.DebtAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtAuctionKey(auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)
	store.Set(key, value)
}

func (k *Keeper) DeleteDebtAuction(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtAuctionKey(id)
	)
	store.Delete(key)
}

func (k *Keeper) SetDutchAuction(ctx sdk.Context, auction auctiontypes.DutchAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchAuctionKey(auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)
	store.Set(key, value)
}

func (k *Keeper) DeleteDutchAuction(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchAuctionKey(id)
	)
	store.Delete(key)
}

func (k *Keeper) GetSurplusAuction(ctx sdk.Context, id uint64) (auction auctiontypes.SurplusAuction, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.CollateralAuctionKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return auction, false
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, true
}

func (k *Keeper) GetSurplusAuctions(ctx sdk.Context) (auctions []auctiontypes.SurplusAuction) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, auctiontypes.CollateralAuctionKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.SurplusAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) GetDebtAuction(ctx sdk.Context, id uint64) (auction auctiontypes.DebtAuction, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtAuctionKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return auction, false
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, true
}

func (k *Keeper) GetDebtAuctions(ctx sdk.Context) (auctions []auctiontypes.DebtAuction) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, auctiontypes.DebtAuctionKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.DebtAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) GetDutchAuction(ctx sdk.Context, id uint64) (auction auctiontypes.DutchAuction, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchAuctionKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return auction, false
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, true
}

func (k *Keeper) GetDutchAuctions(ctx sdk.Context) (auctions []auctiontypes.DutchAuction) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, auctiontypes.DutchAuctionKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.DutchAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) GetSurplusBiddingID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.BiddingsIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetSurplusBiddingID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.BiddingsIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetDebtBiddingID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtBiddingsIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetDebtBiddingID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtBiddingsIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetDutchBiddingID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchBiddingsIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetDutchBiddingID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchBiddingsIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}
func (k *Keeper) SetSurplusBidding(ctx sdk.Context, bidding auctiontypes.Biddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.BiddingsKey(bidding.Id)
		value = k.cdc.MustMarshal(&bidding)
	)
	store.Set(key, value)
}

func (k *Keeper) GetSurplusBidding(ctx sdk.Context, id uint64) (bidding auctiontypes.Biddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.BiddingsKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return bidding, false
	}

	k.cdc.MustUnmarshal(value, &bidding)
	return bidding, true
}

func (k *Keeper) SetDebtBidding(ctx sdk.Context, bidding auctiontypes.Biddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtBiddingsKey(bidding.Id)
		value = k.cdc.MustMarshal(&bidding)
	)
	store.Set(key, value)
}

func (k *Keeper) GetDebtBidding(ctx sdk.Context, id uint64) (bidding auctiontypes.Biddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtBiddingsKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return bidding, false
	}

	k.cdc.MustUnmarshal(value, &bidding)
	return bidding, true
}

func (k *Keeper) SetDutchBidding(ctx sdk.Context, bidding auctiontypes.DutchBiddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchBiddingsKey(bidding.BiddingId)
		value = k.cdc.MustMarshal(&bidding)
	)
	store.Set(key, value)
}

func (k *Keeper) GetDutchBidding(ctx sdk.Context, id uint64) (bidding auctiontypes.DutchBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchBiddingsKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return bidding, false
	}

	k.cdc.MustUnmarshal(value, &bidding)
	return bidding, true
}

func (k *Keeper) GetSurplusBiddings(ctx sdk.Context) (biddings []auctiontypes.Biddings) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, auctiontypes.BiddingsKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var bidding auctiontypes.Biddings
		k.cdc.MustUnmarshal(iter.Value(), &bidding)
		biddings = append(biddings, bidding)
	}

	return biddings
}

func (k *Keeper) GetDebtBiddings(ctx sdk.Context) (biddings []auctiontypes.Biddings) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, auctiontypes.DebtAuctionKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var bidding auctiontypes.Biddings
		k.cdc.MustUnmarshal(iter.Value(), &bidding)
		biddings = append(biddings, bidding)
	}

	return biddings
}

func (k *Keeper) GetDutchBiddings(ctx sdk.Context) (biddings []auctiontypes.DutchBiddings) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, auctiontypes.DutchAuctionKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var bidding auctiontypes.DutchBiddings
		k.cdc.MustUnmarshal(iter.Value(), &bidding)
		biddings = append(biddings, bidding)
	}

	return biddings
}

func (k *Keeper) GetSurplusUserBiddingID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserBiddingsIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetSurplusUserBiddingID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserBiddingsIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetDebtUserBiddingID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtUserBiddingsIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetDebtUserBiddingID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtUserBiddingsIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetDutchUserBiddingID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchUserBiddingsIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetDutchUserBiddingID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchUserBiddingsIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) SetSurplusUserBidding(ctx sdk.Context, userBiddings auctiontypes.UserBiddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserBiddingsKey(userBiddings.Bidder)
		value = k.cdc.MustMarshal(&userBiddings)
	)
	store.Set(key, value)
}

func (k *Keeper) GetSurplusUserBiddings(ctx sdk.Context, bidder string) (userBiddings auctiontypes.UserBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserBiddingsKey(bidder)
		value = store.Get(key)
	)

	if value == nil {
		return userBiddings, false
	}

	k.cdc.MustUnmarshal(value, &userBiddings)
	return userBiddings, true
}

func (k *Keeper) SetDebtUserBidding(ctx sdk.Context, userBiddings auctiontypes.UserBiddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtUserBiddingsKey(userBiddings.Bidder)
		value = k.cdc.MustMarshal(&userBiddings)
	)
	store.Set(key, value)
}

func (k *Keeper) GetDebtUserBiddings(ctx sdk.Context, bidder string) (userBiddings auctiontypes.UserBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DebtUserBiddingsKey(bidder)
		value = store.Get(key)
	)

	if value == nil {
		return userBiddings, false
	}

	k.cdc.MustUnmarshal(value, &userBiddings)
	return userBiddings, true
}

func (k *Keeper) SetDutchUserBidding(ctx sdk.Context, userBiddings auctiontypes.UserBiddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchUserBiddingsKey(userBiddings.Bidder)
		value = k.cdc.MustMarshal(&userBiddings)
	)
	store.Set(key, value)
}

func (k *Keeper) GetDutchUserBiddings(ctx sdk.Context, bidder string) (userBiddings auctiontypes.UserBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.DutchUserBiddingsKey(bidder)
		value = store.Get(key)
	)

	if value == nil {
		return userBiddings, false
	}

	k.cdc.MustUnmarshal(value, &userBiddings)
	return userBiddings, true
}
