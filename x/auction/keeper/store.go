package keeper

import (
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	protobuftypes "github.com/gogo/protobuf/types"
)

//Generic for all auctions.

func (k *Keeper) SetProtocolStatistics(ctx sdk.Context, appID, assetID uint64, amount sdk.Int) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.ProtocolStatisticsKey(appID, assetID)
	)
	stat, found := k.GetProtocolStat(ctx, appID, assetID)
	if found {
		stat.Loss = stat.Loss.Add(amount.ToDec())
		value := k.cdc.MustMarshal(&stat)
		store.Set(key, value)
	} else {
		var stats auctiontypes.ProtocolStatistics
		stats.AppId = appID
		stats.AssetId = assetID
		stats.Loss = amount.ToDec()
		value := k.cdc.MustMarshal(&stats)
		store.Set(key, value)
	}
}

func (k *Keeper) GetProtocolStat(ctx sdk.Context, appID, assetID uint64) (stats auctiontypes.ProtocolStatistics, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.ProtocolStatisticsKey(appID, assetID)
		value = store.Get(key)
	)
	if value == nil {
		return stats, false
	}
	k.cdc.MustUnmarshal(value, &stats)
	return stats, true
}

func (k *Keeper) GetAuctionID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionIDKey
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
		key   = auctiontypes.AuctionIDKey
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
		key   = auctiontypes.UserBiddingsIDKey
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
		key   = auctiontypes.UserBiddingsIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetAuctionType(ctx sdk.Context, auctionTypeID uint64, appID uint64) (string, error) {
	params, found := k.GetAuctionParams(ctx, appID)

	if !found {
		return "", auctiontypes.ErrorInvalidAuctionParams
	}
	if auctionTypeID == params.SurplusId {
		return auctiontypes.SurplusString, nil
	} else if auctionTypeID == params.DebtId {
		return auctiontypes.DebtString, nil
	} else if auctionTypeID == params.DutchId {
		return auctiontypes.DutchString, nil
	}

	return "", sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction mapping id %d not found", auctionTypeID)
}

func (k *Keeper) GetLendAuctionType(ctx sdk.Context, auctionTypeID uint64, appID uint64) (string, error) {
	params, found := k.GetAddAuctionParamsData(ctx, appID)

	if !found {
		return "", auctiontypes.ErrorInvalidAuctionParams
	}
	if auctionTypeID == params.SurplusId {
		return auctiontypes.SurplusString, nil
	} else if auctionTypeID == params.DebtId {
		return auctiontypes.DebtString, nil
	} else if auctionTypeID == params.DutchId {
		return auctiontypes.DutchString, nil
	}

	return "", sdkerrors.Wrapf(sdkerrors.ErrNotFound, "auction mapping id %d not found", auctionTypeID)
}

func (k *Keeper) GetAllAuctions(ctx sdk.Context) (auctions []auctiontypes.SurplusAuction) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, auctiontypes.AuctionKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.SurplusAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

//SURPLUS

func (k *Keeper) SetSurplusAuction(ctx sdk.Context, auction auctiontypes.SurplusAuction) error {
	auctionType, err := k.GetAuctionType(ctx, auction.AuctionMappingId, auction.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionKey(auction.AppId, auctionType, auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)
	store.Set(key, value)
	return nil
}

func (k *Keeper) SetHistorySurplusAuction(ctx sdk.Context, auction auctiontypes.SurplusAuction) error {
	auctionType, err := k.GetAuctionType(ctx, auction.AuctionMappingId, auction.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryAuctionKey(auction.AppId, auctionType, auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)
	store.Set(key, value)
	return nil
}

func (k *Keeper) DeleteSurplusAuction(ctx sdk.Context, auction auctiontypes.SurplusAuction) error {
	auctionType, err := k.GetAuctionType(ctx, auction.AuctionMappingId, auction.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionKey(auction.AppId, auctionType, auction.AuctionId)
	)
	store.Delete(key)
	return nil
}

func (k *Keeper) GetSurplusAuction(ctx sdk.Context, appID, auctionMappingID, auctionID uint64) (auction auctiontypes.SurplusAuction, err error) {
	auctionType, err := k.GetAuctionType(ctx, auctionMappingID, appID)
	if err != nil {
		return auction, err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionKey(appID, auctionType, auctionID)
		value = store.Get(key)
	)
	if value == nil {
		return auction, sdkerrors.ErrNotFound
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, nil
}

func (k *Keeper) GetHistorySurplusAuction(ctx sdk.Context, appID, auctionMappingID, auctionID uint64) (auction auctiontypes.SurplusAuction, err error) {
	auctionType, err := k.GetAuctionType(ctx, auctionMappingID, appID)
	if err != nil {
		return auction, err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryAuctionKey(appID, auctionType, auctionID)
		value = store.Get(key)
	)
	if value == nil {
		return auction, sdkerrors.ErrNotFound
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, nil
}

func (k *Keeper) GetSurplusAuctions(ctx sdk.Context, appID uint64) (auctions []auctiontypes.SurplusAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionTypeKey(appID, auctiontypes.SurplusString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.SurplusAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) GetHistorySurplusAuctions(ctx sdk.Context, appID uint64) (auctions []auctiontypes.SurplusAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryAuctionTypeKey(appID, auctiontypes.SurplusString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.SurplusAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) SetSurplusUserBidding(ctx sdk.Context, userBiddings auctiontypes.SurplusBiddings) error {
	auctionType, err := k.GetAuctionType(ctx, userBiddings.AuctionMappingId, userBiddings.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
		value = k.cdc.MustMarshal(&userBiddings)
	)
	store.Set(key, value)
	return nil
}

func (k *Keeper) SetHistorySurplusUserBidding(ctx sdk.Context, userBiddings auctiontypes.SurplusBiddings) error {
	auctionType, err := k.GetAuctionType(ctx, userBiddings.AuctionMappingId, userBiddings.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryUserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
		value = k.cdc.MustMarshal(&userBiddings)
	)
	store.Set(key, value)
	return nil
}

func (k *Keeper) DeleteSurplusUserBidding(ctx sdk.Context, userBiddings auctiontypes.SurplusBiddings) error {
	auctionType, err := k.GetAuctionType(ctx, userBiddings.AuctionMappingId, userBiddings.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
	)
	store.Delete(key)
	return nil
}

func (k *Keeper) GetSurplusUserBidding(ctx sdk.Context, bidder string, appID, biddingID uint64) (userBidding auctiontypes.SurplusBiddings, err error) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserKey(bidder, appID, auctiontypes.SurplusString, biddingID)
		value = store.Get(key)
	)
	if value == nil {
		return userBidding, sdkerrors.ErrNotFound
	}

	k.cdc.MustUnmarshal(value, &userBidding)

	return userBidding, nil
}

func (k *Keeper) GetHistorySurplusUserBidding(ctx sdk.Context, bidder string, appID, biddingID uint64) (userBidding auctiontypes.SurplusBiddings, err error) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryUserKey(bidder, appID, auctiontypes.SurplusString, biddingID)
		value = store.Get(key)
	)
	if value == nil {
		return userBidding, sdkerrors.ErrNotFound
	}

	k.cdc.MustUnmarshal(value, &userBidding)

	return userBidding, nil
}

func (k *Keeper) GetSurplusUserBiddings(ctx sdk.Context, bidder string, appID uint64) (userBiddings []auctiontypes.SurplusBiddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserAuctionTypeKey(bidder, appID, auctiontypes.SurplusString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var userBidding auctiontypes.SurplusBiddings
		k.cdc.MustUnmarshal(iter.Value(), &userBidding)
		userBiddings = append(userBiddings, userBidding)
	}

	return userBiddings
}

func (k *Keeper) GetHistorySurplusUserBiddings(ctx sdk.Context, bidder string, appID uint64) (userBiddings []auctiontypes.SurplusBiddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryUserAuctionTypeKey(bidder, appID, auctiontypes.SurplusString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var userBidding auctiontypes.SurplusBiddings
		k.cdc.MustUnmarshal(iter.Value(), &userBidding)
		userBiddings = append(userBiddings, userBidding)
	}

	return userBiddings
}

//DEBT

func (k *Keeper) SetDebtAuction(ctx sdk.Context, auction auctiontypes.DebtAuction) error {
	auctionType, err := k.GetAuctionType(ctx, auction.AuctionMappingId, auction.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionKey(auction.AppId, auctionType, auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)
	store.Set(key, value)
	return nil
}

func (k *Keeper) SetHistoryDebtAuction(ctx sdk.Context, auction auctiontypes.DebtAuction) error {
	auctionType, err := k.GetAuctionType(ctx, auction.AuctionMappingId, auction.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryAuctionKey(auction.AppId, auctionType, auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)
	store.Set(key, value)
	return nil
}

func (k *Keeper) DeleteDebtAuction(ctx sdk.Context, auction auctiontypes.DebtAuction) error {
	auctionType, err := k.GetAuctionType(ctx, auction.AuctionMappingId, auction.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionKey(auction.AppId, auctionType, auction.AuctionId)
	)
	store.Delete(key)
	return nil
}

func (k *Keeper) GetDebtAuction(ctx sdk.Context, appID, auctionMappingID, auctionID uint64) (auction auctiontypes.DebtAuction, err error) {
	auctionType, err := k.GetAuctionType(ctx, auctionMappingID, appID)
	if err != nil {
		return auction, err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionKey(appID, auctionType, auctionID)
		value = store.Get(key)
	)
	if value == nil {
		return auction, sdkerrors.ErrNotFound
	}
	k.cdc.MustUnmarshal(value, &auction)
	return auction, nil
}

func (k *Keeper) GetHistoryDebtAuction(ctx sdk.Context, appID, auctionMappingID, auctionID uint64) (auction auctiontypes.DebtAuction, err error) {
	auctionType, err := k.GetAuctionType(ctx, auctionMappingID, appID)
	if err != nil {
		return auction, err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryAuctionKey(appID, auctionType, auctionID)
		value = store.Get(key)
	)
	if value == nil {
		return auction, sdkerrors.ErrNotFound
	}
	k.cdc.MustUnmarshal(value, &auction)
	return auction, nil
}

func (k *Keeper) GetDebtAuctions(ctx sdk.Context, appID uint64) (auctions []auctiontypes.DebtAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionTypeKey(appID, auctiontypes.DebtString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.DebtAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) GetHistoryDebtAuctions(ctx sdk.Context, appID uint64) (auctions []auctiontypes.DebtAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryAuctionTypeKey(appID, auctiontypes.DebtString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.DebtAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) SetDebtUserBidding(ctx sdk.Context, userBiddings auctiontypes.DebtBiddings) error {
	auctionType, err := k.GetAuctionType(ctx, userBiddings.AuctionMappingId, userBiddings.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
		value = k.cdc.MustMarshal(&userBiddings)
	)
	store.Set(key, value)
	return nil
}

func (k *Keeper) SetHistoryDebtUserBidding(ctx sdk.Context, userBiddings auctiontypes.DebtBiddings) error {
	auctionType, err := k.GetAuctionType(ctx, userBiddings.AuctionMappingId, userBiddings.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryUserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
		value = k.cdc.MustMarshal(&userBiddings)
	)
	store.Set(key, value)
	return nil
}

func (k *Keeper) DeleteDebtUserBidding(ctx sdk.Context, userBiddings auctiontypes.DebtBiddings) error {
	auctionType, err := k.GetAuctionType(ctx, userBiddings.AuctionMappingId, userBiddings.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
	)
	store.Delete(key)
	return nil
}

func (k *Keeper) GetDebtUserBidding(ctx sdk.Context, bidder string, appID, biddingID uint64) (userBidding auctiontypes.DebtBiddings, err error) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserKey(bidder, appID, auctiontypes.DebtString, biddingID)
		value = store.Get(key)
	)
	if value == nil {
		return userBidding, sdkerrors.ErrNotFound
	}
	k.cdc.MustUnmarshal(value, &userBidding)
	return userBidding, nil
}

func (k *Keeper) GetHistoryDebtUserBidding(ctx sdk.Context, bidder string, appID, biddingID uint64) (userBidding auctiontypes.DebtBiddings, err error) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryUserKey(bidder, appID, auctiontypes.DebtString, biddingID)
		value = store.Get(key)
	)
	if value == nil {
		return userBidding, sdkerrors.ErrNotFound
	}
	k.cdc.MustUnmarshal(value, &userBidding)
	return userBidding, nil
}

func (k *Keeper) GetDebtUserBiddings(ctx sdk.Context, bidder string, appID uint64) (userBiddings []auctiontypes.DebtBiddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserAuctionTypeKey(bidder, appID, auctiontypes.DebtString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var userBidding auctiontypes.DebtBiddings
		k.cdc.MustUnmarshal(iter.Value(), &userBidding)
		userBiddings = append(userBiddings, userBidding)
	}

	return userBiddings
}

func (k *Keeper) GetHistoryDebtUserBiddings(ctx sdk.Context, bidder string, appID uint64) (userBiddings []auctiontypes.DebtBiddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryUserAuctionTypeKey(bidder, appID, auctiontypes.DebtString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var userBidding auctiontypes.DebtBiddings
		k.cdc.MustUnmarshal(iter.Value(), &userBidding)
		userBiddings = append(userBiddings, userBidding)
	}

	return userBiddings
}

//DUTCH

func (k *Keeper) SetDutchAuction(ctx sdk.Context, auction auctiontypes.DutchAuction) error {
	auctionType, err := k.GetAuctionType(ctx, auction.AuctionMappingId, auction.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionKey(auction.AppId, auctionType, auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)

	store.Set(key, value)
	return nil
}

func (k *Keeper) SetHistoryDutchAuction(ctx sdk.Context, auction auctiontypes.DutchAuction) error {
	auctionType, err := k.GetAuctionType(ctx, auction.AuctionMappingId, auction.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryAuctionKey(auction.AppId, auctionType, auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)
	store.Set(key, value)
	return nil
}

func (k *Keeper) DeleteDutchAuction(ctx sdk.Context, auction auctiontypes.DutchAuction) error {
	auctionType, err := k.GetAuctionType(ctx, auction.AuctionMappingId, auction.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionKey(auction.AppId, auctionType, auction.AuctionId)
	)
	store.Delete(key)
	return nil
}

func (k *Keeper) GetDutchAuction(ctx sdk.Context, appID, auctionMappingID, auctionID uint64) (auction auctiontypes.DutchAuction, err error) {
	auctionType, err := k.GetAuctionType(ctx, auctionMappingID, appID)
	if err != nil {
		return auction, err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionKey(appID, auctionType, auctionID)
		value = store.Get(key)
	)

	if value == nil {
		return auction, sdkerrors.ErrNotFound
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, nil
}

func (k *Keeper) GetHistoryDutchAuction(ctx sdk.Context, appID, auctionMappingID, auctionID uint64) (auction auctiontypes.DutchAuction, err error) {
	auctionType, err := k.GetAuctionType(ctx, auctionMappingID, appID)
	if err != nil {
		return auction, err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryAuctionKey(appID, auctionType, auctionID)
		value = store.Get(key)
	)
	if value == nil {
		return auction, sdkerrors.ErrNotFound
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, nil
}

func (k *Keeper) GetDutchAuctions(ctx sdk.Context, appID uint64) (auctions []auctiontypes.DutchAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionTypeKey(appID, auctiontypes.DutchString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.DutchAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) GetHistoryDutchAuctions(ctx sdk.Context, appID uint64) (auctions []auctiontypes.DutchAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryAuctionTypeKey(appID, auctiontypes.DutchString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.DutchAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) SetDutchUserBidding(ctx sdk.Context, userBiddings auctiontypes.DutchBiddings) error {
	auctionType, err := k.GetAuctionType(ctx, userBiddings.AuctionMappingId, userBiddings.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
		value = k.cdc.MustMarshal(&userBiddings)
	)
	store.Set(key, value)
	return nil
}

func (k *Keeper) SetHistoryDutchUserBidding(ctx sdk.Context, userBiddings auctiontypes.DutchBiddings) error {
	auctionType, err := k.GetAuctionType(ctx, userBiddings.AuctionMappingId, userBiddings.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryUserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
		value = k.cdc.MustMarshal(&userBiddings)
	)
	store.Set(key, value)
	return nil
}

func (k *Keeper) DeleteDutchUserBidding(ctx sdk.Context, userBiddings auctiontypes.DutchBiddings) error {
	auctionType, err := k.GetAuctionType(ctx, userBiddings.AuctionMappingId, userBiddings.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
	)
	store.Delete(key)
	return nil
}

func (k *Keeper) GetDutchUserBidding(ctx sdk.Context, bidder string, appID, biddingID uint64) (userBidding auctiontypes.DutchBiddings, err error) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserKey(bidder, appID, auctiontypes.DutchString, biddingID)
		value = store.Get(key)
	)

	if value == nil {
		return userBidding, sdkerrors.ErrNotFound
	}

	k.cdc.MustUnmarshal(value, &userBidding)

	return userBidding, nil
}

func (k *Keeper) GetHistoryDutchUserBidding(ctx sdk.Context, bidder string, appID, biddingID uint64) (userBidding auctiontypes.DutchBiddings, err error) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryUserKey(bidder, appID, auctiontypes.DutchString, biddingID)
		value = store.Get(key)
	)

	if value == nil {
		return userBidding, sdkerrors.ErrNotFound
	}

	k.cdc.MustUnmarshal(value, &userBidding)

	return userBidding, nil
}

func (k *Keeper) GetDutchUserBiddings(ctx sdk.Context, bidder string, appID uint64) (userBiddings []auctiontypes.DutchBiddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.UserAuctionTypeKey(bidder, appID, auctiontypes.DutchString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var userBidding auctiontypes.DutchBiddings
		k.cdc.MustUnmarshal(iter.Value(), &userBidding)
		userBiddings = append(userBiddings, userBidding)
	}

	return userBiddings
}

func (k *Keeper) GetHistoryDutchUserBiddings(ctx sdk.Context, bidder string, appID uint64) (userBiddings []auctiontypes.DutchBiddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryUserAuctionTypeKey(bidder, appID, auctiontypes.DutchString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var userBidding auctiontypes.DutchBiddings
		k.cdc.MustUnmarshal(iter.Value(), &userBidding)
		userBiddings = append(userBiddings, userBidding)
	}

	return userBiddings
}

func (k *Keeper) SetAuctionParams(ctx sdk.Context, auctionParams auctiontypes.AuctionParams) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionParamsKey(auctionParams.AppId)
		value = k.cdc.MustMarshal(&auctionParams)
	)

	store.Set(key, value)
}

func (k *Keeper) GetAuctionParams(ctx sdk.Context, AppID uint64) (asset auctiontypes.AuctionParams, found bool) {
	key := auctiontypes.AuctionParamsKey(AppID)

	var (
		store = k.Store(ctx)

		value = store.Get(key)
	)

	if value == nil {
		return asset, false
	}

	k.cdc.MustUnmarshal(value, &asset)
	return asset, true
}

func (k *Keeper) GetLendAuctionID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.LendAuctionIDKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetLendAuctionID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.LendAuctionIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) SetDutchLendAuction(ctx sdk.Context, auction auctiontypes.DutchAuction) error {
	auctionType, err := k.GetLendAuctionType(ctx, auction.AuctionMappingId, auction.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.LendAuctionKey(auction.AppId, auctionType, auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)

	store.Set(key, value)
	return nil
}

func (k *Keeper) SetHistoryDutchLendAuction(ctx sdk.Context, auction auctiontypes.DutchAuction) error {
	auctionType, err := k.GetLendAuctionType(ctx, auction.AuctionMappingId, auction.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryLendAuctionKey(auction.AppId, auctionType, auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)
	store.Set(key, value)
	return nil
}

func (k *Keeper) DeleteDutchLendAuction(ctx sdk.Context, auction auctiontypes.DutchAuction) error {
	auctionType, err := k.GetLendAuctionType(ctx, auction.AuctionMappingId, auction.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.LendAuctionKey(auction.AppId, auctionType, auction.AuctionId)
	)
	store.Delete(key)
	return nil
}

func (k *Keeper) GetDutchLendAuction(ctx sdk.Context, appID, auctionMappingID, auctionID uint64) (auction auctiontypes.DutchAuction, err error) {
	auctionType, err := k.GetLendAuctionType(ctx, auctionMappingID, appID)
	if err != nil {
		return auction, err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.LendAuctionKey(appID, auctionType, auctionID)
		value = store.Get(key)
	)

	if value == nil {
		return auction, sdkerrors.ErrNotFound
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, nil
}

func (k *Keeper) GetHistoryDutchLendAuction(ctx sdk.Context, appID, auctionMappingID, auctionID uint64) (auction auctiontypes.DutchAuction, err error) {
	auctionType, err := k.GetLendAuctionType(ctx, auctionMappingID, appID)
	if err != nil {
		return auction, err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryLendAuctionKey(appID, auctionType, auctionID)
		value = store.Get(key)
	)
	if value == nil {
		return auction, sdkerrors.ErrNotFound
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, nil
}

func (k *Keeper) GetDutchLendAuctions(ctx sdk.Context, appID uint64) (auctions []auctiontypes.DutchAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.LendAuctionTypeKey(appID, auctiontypes.DutchString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.DutchAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) GetHistoryDutchLendAuctions(ctx sdk.Context, appID uint64) (auctions []auctiontypes.DutchAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryLendAuctionTypeKey(appID, auctiontypes.DutchString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.DutchAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) SetDutchUserLendBidding(ctx sdk.Context, userBiddings auctiontypes.DutchBiddings) error {
	auctionType, err := k.GetLendAuctionType(ctx, userBiddings.AuctionMappingId, userBiddings.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.LendUserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
		value = k.cdc.MustMarshal(&userBiddings)
	)
	store.Set(key, value)
	return nil
}

func (k *Keeper) SetHistoryDutchUserLendBidding(ctx sdk.Context, userBiddings auctiontypes.DutchBiddings) error {
	auctionType, err := k.GetLendAuctionType(ctx, userBiddings.AuctionMappingId, userBiddings.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryLendUserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
		value = k.cdc.MustMarshal(&userBiddings)
	)
	store.Set(key, value)
	return nil
}

func (k *Keeper) DeleteDutchLendUserBidding(ctx sdk.Context, userBiddings auctiontypes.DutchBiddings) error {
	auctionType, err := k.GetLendAuctionType(ctx, userBiddings.AuctionMappingId, userBiddings.AppId)
	if err != nil {
		return err
	}
	var (
		store = k.Store(ctx)
		key   = auctiontypes.LendUserKey(userBiddings.Bidder, userBiddings.AppId, auctionType, userBiddings.BiddingId)
	)
	store.Delete(key)
	return nil
}

func (k *Keeper) GetDutchLendUserBidding(ctx sdk.Context, bidder string, appID, biddingID uint64) (userBidding auctiontypes.DutchBiddings, err error) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.LendUserKey(bidder, appID, auctiontypes.DutchString, biddingID)
		value = store.Get(key)
	)

	if value == nil {
		return userBidding, sdkerrors.ErrNotFound
	}

	k.cdc.MustUnmarshal(value, &userBidding)

	return userBidding, nil
}

func (k *Keeper) GetHistoryDutchLendUserBidding(ctx sdk.Context, bidder string, appID, biddingID uint64) (userBidding auctiontypes.DutchBiddings, err error) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryLendUserKey(bidder, appID, auctiontypes.DutchString, biddingID)
		value = store.Get(key)
	)

	if value == nil {
		return userBidding, sdkerrors.ErrNotFound
	}

	k.cdc.MustUnmarshal(value, &userBidding)

	return userBidding, nil
}

func (k *Keeper) GetDutchLendUserBiddings(ctx sdk.Context, bidder string, appID uint64) (userBiddings []auctiontypes.DutchBiddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.LendUserAuctionTypeKey(bidder, appID, auctiontypes.DutchString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var userBidding auctiontypes.DutchBiddings
		k.cdc.MustUnmarshal(iter.Value(), &userBidding)
		userBiddings = append(userBiddings, userBidding)
	}

	return userBiddings
}

func (k *Keeper) GetHistoryDutchLendUserBiddings(ctx sdk.Context, bidder string, appID uint64) (userBiddings []auctiontypes.DutchBiddings) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.HistoryLendUserAuctionTypeKey(bidder, appID, auctiontypes.DutchString)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var userBidding auctiontypes.DutchBiddings
		k.cdc.MustUnmarshal(iter.Value(), &userBidding)
		userBiddings = append(userBiddings, userBidding)
	}

	return userBiddings
}
