package keeper

import (
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

//Generic for all auctions
func (k *Keeper) GetAuctionID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetAuctionID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetUserBiddingID(ctx sdk.Context) uint64 {
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

func (k *Keeper) SetUserBiddingID(ctx sdk.Context, id uint64) {
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

func (k *Keeper) GetAuctionType(ctx sdk.Context, auctionTypeId uint64) string {
	params := k.GetParams(ctx)
	if auctionTypeId == params.SurplusId {
		return auctiontypes.SurplusString
	} else if auctionTypeId == params.DebtId {
		return auctiontypes.DebtString
	} else if auctionTypeId == params.DutchId {
		return auctiontypes.DutchString
	}

	return ""
}

func (k *Keeper) GetAllAuctions(ctx sdk.Context) (auctions []auctiontypes.SurplusAuction) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, auctiontypes.AuctionKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.SurplusAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

//SURPLUS

func (k *Keeper) SetSurplusAuction(ctx sdk.Context, auction auctiontypes.SurplusAuction) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, auction.AuctionMappingId)
		key         = auctiontypes.AuctionKey(auction.AppId, auctionType, auction.AuctionId)
		value       = k.cdc.MustMarshal(&auction)
	)
	if auctionType == "" {
		return
	}
	store.Set(key, value)
}

func (k *Keeper) SetHistorySurplusAuction(ctx sdk.Context, auction auctiontypes.SurplusAuction) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, auction.AuctionMappingId)
		key         = auctiontypes.HistoryAuctionKey(auction.AppId, auctionType, auction.AuctionId)
		value       = k.cdc.MustMarshal(&auction)
	)
	if auctionType == "" {
		return
	}
	store.Set(key, value)
}

func (k *Keeper) DeleteSurplusAuction(ctx sdk.Context, auction auctiontypes.SurplusAuction) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, auction.AuctionMappingId)
		key         = auctiontypes.AuctionKey(auction.AppId, auctionType, auction.AuctionId)
	)
	if auctionType == "" {
		return
	}
	store.Delete(key)
}

func (k *Keeper) GetSurplusAuction(ctx sdk.Context, appId, auctionMappingId, auctionId uint64) (auction auctiontypes.SurplusAuction, found bool) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, auctionMappingId)
		key         = auctiontypes.AuctionKey(appId, auctionType, auctionId)
		value       = store.Get(key)
	)
	if auctionType == "" {
		return auction, false
	}
	if value == nil {
		return auction, false
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, true
}

func (k *Keeper) GetHistorySurplusAuction(ctx sdk.Context, appId, auctionMappingId, auctionId uint64) (auction auctiontypes.SurplusAuction, found bool) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, auctionMappingId)
		key         = auctiontypes.HistoryAuctionKey(appId, auctionType, auctionId)
		value       = store.Get(key)
	)
	if auctionType == "" {
		return auction, false
	}
	if value == nil {
		return auction, false
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, true
}

func (k *Keeper) GetSurplusAuctions(ctx sdk.Context, appId uint64) (auctions []auctiontypes.SurplusAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionTypeKey(appId, auctiontypes.SurplusString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.SurplusAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) GetHistorySurplusAuctions(ctx sdk.Context, appId uint64) (auctions []auctiontypes.SurplusAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryAuctionTypeKey(appId, auctiontypes.SurplusString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.SurplusAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) SetSurplusUserBidding(ctx sdk.Context, userBiddings auctiontypes.SurplusBiddings) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, userBiddings.AuctionMappingId)
		key         = auctiontypes.UserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
		value       = k.cdc.MustMarshal(&userBiddings)
	)
	if auctionType == "" {
		return
	}
	store.Set(key, value)
}

func (k *Keeper) SetHistorySurplusUserBidding(ctx sdk.Context, userBiddings auctiontypes.SurplusBiddings) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, userBiddings.AuctionMappingId)
		key         = auctiontypes.HistoryUserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
		value       = k.cdc.MustMarshal(&userBiddings)
	)
	if auctionType == "" {
		return
	}
	store.Set(key, value)
}

func (k *Keeper) DeleteSurplusUserBidding(ctx sdk.Context, userBiddings auctiontypes.SurplusBiddings) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, userBiddings.AuctionMappingId)
		key         = auctiontypes.UserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
	)
	if auctionType == "" {
		return
	}
	store.Delete(key)
}

func (k *Keeper) GetSurplusUserBidding(ctx sdk.Context, bidder string, appId, biddingId uint64) (userBidding auctiontypes.SurplusBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserKey(bidder, appId, auctiontypes.SurplusString, biddingId)
		value = store.Get(key)
	)
	if value == nil {
		return userBidding, false
	}

	k.cdc.MustUnmarshal(value, &userBidding)

	return userBidding, true
}

func (k *Keeper) GetHistorySurplusUserBidding(ctx sdk.Context, bidder string, appId, biddingId uint64) (userBidding auctiontypes.SurplusBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryUserKey(bidder, appId, auctiontypes.SurplusString, biddingId)
		value = store.Get(key)
	)
	if value == nil {
		return userBidding, false
	}

	k.cdc.MustUnmarshal(value, &userBidding)

	return userBidding, true
}

func (k *Keeper) GetSurplusUserBiddings(ctx sdk.Context, bidder string, appId uint64) (userBiddings []auctiontypes.SurplusBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserAuctionTypeKey(bidder, appId, auctiontypes.SurplusString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var userBidding auctiontypes.SurplusBiddings
		k.cdc.MustUnmarshal(iter.Value(), &userBidding)
		userBiddings = append(userBiddings, userBidding)
	}

	return userBiddings, true
}

func (k *Keeper) GetHistorySurplusUserBiddings(ctx sdk.Context, bidder string, appId uint64) (userBiddings []auctiontypes.SurplusBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryUserAuctionTypeKey(bidder, appId, auctiontypes.SurplusString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var userBidding auctiontypes.SurplusBiddings
		k.cdc.MustUnmarshal(iter.Value(), &userBidding)
		userBiddings = append(userBiddings, userBidding)
	}

	return userBiddings, true
}

//DEBT

func (k *Keeper) SetDebtAuction(ctx sdk.Context, auction auctiontypes.DebtAuction) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, auction.AuctionMappingId)
		key         = auctiontypes.AuctionKey(auction.AppId, auctionType, auction.AuctionId)
		value       = k.cdc.MustMarshal(&auction)
	)
	if auctionType == "" {
		return
	}
	store.Set(key, value)
}

func (k *Keeper) SetHistoryDebtAuction(ctx sdk.Context, auction auctiontypes.DebtAuction) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, auction.AuctionMappingId)
		key         = auctiontypes.HistoryAuctionKey(auction.AppId, auctionType, auction.AuctionId)
		value       = k.cdc.MustMarshal(&auction)
	)
	if auctionType == "" {
		return
	}
	store.Set(key, value)
}

func (k *Keeper) DeleteDebtAuction(ctx sdk.Context, auction auctiontypes.DebtAuction) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, auction.AuctionMappingId)
		key         = auctiontypes.AuctionKey(auction.AppId, auctionType, auction.AuctionId)
	)
	if auctionType == "" {
		return
	}
	store.Delete(key)
}

func (k *Keeper) GetDebtAuction(ctx sdk.Context, appId, auctionMappingId, auctionId uint64) (auction auctiontypes.DebtAuction, found bool) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, auctionMappingId)
		key         = auctiontypes.AuctionKey(appId, auctionType, auctionId)
		value       = store.Get(key)
	)
	if auctionType == "" {
		return auction, false
	}
	if value == nil {
		return auction, false
	}
	k.cdc.MustUnmarshal(value, &auction)
	return auction, true
}

func (k *Keeper) GetHistoryDebtAuction(ctx sdk.Context, appId, auctionMappingId, auctionId uint64) (auction auctiontypes.DebtAuction, found bool) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, auctionMappingId)
		key         = auctiontypes.HistoryAuctionKey(appId, auctionType, auctionId)
		value       = store.Get(key)
	)
	if auctionType == "" {
		return auction, false
	}
	if value == nil {
		return auction, false
	}
	k.cdc.MustUnmarshal(value, &auction)
	return auction, true
}

func (k *Keeper) GetDebtAuctions(ctx sdk.Context, appId uint64) (auctions []auctiontypes.DebtAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionTypeKey(appId, auctiontypes.DebtString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.DebtAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) GetHistoryDebtAuctions(ctx sdk.Context, appId uint64) (auctions []auctiontypes.DebtAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryAuctionTypeKey(appId, auctiontypes.DebtString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.DebtAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) SetDebtUserBidding(ctx sdk.Context, userBiddings auctiontypes.DebtBiddings) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, userBiddings.AuctionMappingId)
		key         = auctiontypes.UserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
		value       = k.cdc.MustMarshal(&userBiddings)
	)
	if auctionType == "" {
		return
	}
	store.Set(key, value)
}

func (k *Keeper) SetHistoryDebtUserBidding(ctx sdk.Context, userBiddings auctiontypes.DebtBiddings) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, userBiddings.AuctionMappingId)
		key         = auctiontypes.HistoryUserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
		value       = k.cdc.MustMarshal(&userBiddings)
	)
	if auctionType == "" {
		return
	}
	store.Set(key, value)
}

func (k *Keeper) DeleteDebtUserBidding(ctx sdk.Context, userBiddings auctiontypes.DebtBiddings) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, userBiddings.AuctionMappingId)
		key         = auctiontypes.UserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
	)
	if auctionType == "" {
		return
	}
	store.Delete(key)
}

func (k *Keeper) GetDebtUserBidding(ctx sdk.Context, bidder string, appId, biddingId uint64) (userBidding auctiontypes.DebtBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserKey(bidder, appId, auctiontypes.DebtString, biddingId)
		value = store.Get(key)
	)
	if value == nil {
		return userBidding, false
	}
	k.cdc.MustUnmarshal(value, &userBidding)
	return userBidding, true
}

func (k *Keeper) GetHistoryDebtUserBidding(ctx sdk.Context, bidder string, appId, biddingId uint64) (userBidding auctiontypes.DebtBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryUserKey(bidder, appId, auctiontypes.DebtString, biddingId)
		value = store.Get(key)
	)
	if value == nil {
		return userBidding, false
	}
	k.cdc.MustUnmarshal(value, &userBidding)
	return userBidding, true
}

func (k *Keeper) GetDebtUserBiddings(ctx sdk.Context, bidder string, appId uint64) (userBiddings []auctiontypes.DebtBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserAuctionTypeKey(bidder, appId, auctiontypes.DebtString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var userBidding auctiontypes.DebtBiddings
		k.cdc.MustUnmarshal(iter.Value(), &userBidding)
		userBiddings = append(userBiddings, userBidding)
	}

	return userBiddings, true
}

func (k *Keeper) GetHistoryDebtUserBiddings(ctx sdk.Context, bidder string, appId uint64) (userBiddings []auctiontypes.DebtBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryUserAuctionTypeKey(bidder, appId, auctiontypes.DebtString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var userBidding auctiontypes.DebtBiddings
		k.cdc.MustUnmarshal(iter.Value(), &userBidding)
		userBiddings = append(userBiddings, userBidding)
	}

	return userBiddings, true
}

//DUTCH

func (k *Keeper) SetDutchAuction(ctx sdk.Context, auction auctiontypes.DutchAuction) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, auction.AuctionMappingId)
		key         = auctiontypes.AuctionKey(auction.AppId, auctionType, auction.AuctionId)
		value       = k.cdc.MustMarshal(&auction)
	)
	if auctionType == "" {
		return
	}
	store.Set(key, value)
}

func (k *Keeper) SetHistoryDutchAuction(ctx sdk.Context, auction auctiontypes.DutchAuction) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, auction.AuctionMappingId)
		key         = auctiontypes.HistoryAuctionKey(auction.AppId, auctionType, auction.AuctionId)
		value       = k.cdc.MustMarshal(&auction)
	)
	if auctionType == "" {
		return
	}
	store.Set(key, value)
}

func (k *Keeper) DeleteDutchAuction(ctx sdk.Context, auction auctiontypes.DutchAuction) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, auction.AuctionMappingId)
		key         = auctiontypes.AuctionKey(auction.AppId, auctionType, auction.AuctionId)
	)
	if auctionType == "" {
		return
	}
	store.Delete(key)
}

func (k *Keeper) GetDutchAuction(ctx sdk.Context, appId, auctionMappingId, auctionId uint64) (auction auctiontypes.DutchAuction, found bool) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, auctionMappingId)
		key         = auctiontypes.AuctionKey(appId, auctionType, auctionId)
		value       = store.Get(key)
	)
	if auctionType == "" {
		return auction, false
	}
	if value == nil {
		return auction, false
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, true
}

func (k *Keeper) GetHistoryDutchAuction(ctx sdk.Context, appId, auctionMappingId, auctionId uint64) (auction auctiontypes.DutchAuction, found bool) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, auctionMappingId)
		key         = auctiontypes.HistoryAuctionKey(appId, auctionType, auctionId)
		value       = store.Get(key)
	)
	if auctionType == "" {
		return auction, false
	}
	if value == nil {
		return auction, false
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, true
}

func (k *Keeper) GetDutchAuctions(ctx sdk.Context, appId uint64) (auctions []auctiontypes.DutchAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionTypeKey(appId, auctiontypes.DutchString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.DutchAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) GetHistoryDutchAuctions(ctx sdk.Context, appId uint64) (auctions []auctiontypes.DutchAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryAuctionTypeKey(appId, auctiontypes.DutchString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.DutchAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) SetDutchUserBidding(ctx sdk.Context, userBiddings auctiontypes.DutchBiddings) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, userBiddings.AuctionMappingId)
		key         = auctiontypes.UserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
		value       = k.cdc.MustMarshal(&userBiddings)
	)
	if auctionType == "" {
		return
	}
	store.Set(key, value)
}

func (k *Keeper) SetHistoryDutchUserBidding(ctx sdk.Context, userBiddings auctiontypes.DutchBiddings) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, userBiddings.AuctionMappingId)
		key         = auctiontypes.HistoryUserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
		value       = k.cdc.MustMarshal(&userBiddings)
	)
	if auctionType == "" {
		return
	}
	store.Set(key, value)
}

func (k *Keeper) DeleteDutchUserBidding(ctx sdk.Context, userBiddings auctiontypes.DutchBiddings) {
	var (
		store       = k.Store(ctx)
		auctionType = k.GetAuctionType(ctx, userBiddings.AuctionMappingId)
		key         = auctiontypes.UserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
	)
	if auctionType == "" {
		return
	}
	store.Delete(key)
}

func (k *Keeper) GetDutchUserBidding(ctx sdk.Context, bidder string, appId, biddingId uint64) (userBidding auctiontypes.DutchBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserKey(bidder, appId, auctiontypes.DutchString, biddingId)
		value = store.Get(key)
	)

	if value == nil {
		return userBidding, false
	}

	k.cdc.MustUnmarshal(value, &userBidding)

	return userBidding, true
}

func (k *Keeper) GetHistoryDutchUserBidding(ctx sdk.Context, bidder string, appId, biddingId uint64) (userBidding auctiontypes.DutchBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryUserKey(bidder, appId, auctiontypes.DutchString, biddingId)
		value = store.Get(key)
	)

	if value == nil {
		return userBidding, false
	}

	k.cdc.MustUnmarshal(value, &userBidding)

	return userBidding, true
}

func (k *Keeper) GetDutchUserBiddings(ctx sdk.Context, bidder string, appId uint64) (userBiddings []auctiontypes.DutchBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserAuctionTypeKey(bidder, appId, auctiontypes.DutchString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var userBidding auctiontypes.DutchBiddings
		k.cdc.MustUnmarshal(iter.Value(), &userBidding)
		userBiddings = append(userBiddings, userBidding)
	}

	return userBiddings, true
}

func (k *Keeper) GetHistoryDutchUserBiddings(ctx sdk.Context, bidder string, appId uint64) (userBiddings []auctiontypes.DutchBiddings, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryUserAuctionTypeKey(bidder, appId, auctiontypes.DutchString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var userBidding auctiontypes.DutchBiddings
		k.cdc.MustUnmarshal(iter.Value(), &userBidding)
		userBiddings = append(userBiddings, userBidding)
	}

	return userBiddings, true
}
