package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/comdex-official/comdex/x/auctionsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetGenAuction(ctx sdk.Context, auction types.Auction) {

	var (
		store = k.Store(ctx)
		key   = types.AuctionKey(auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)

	store.Set(key, value)
}

func (k Keeper) SetGenAuctionLimitBidFeeData(ctx sdk.Context, feeData types.AuctionFeesCollectionFromLimitBidTx) {

	var (
		store = k.Store(ctx)
		key   = types.AuctionLimitBidFeeKey(feeData.AssetId)
		value = k.cdc.MustMarshal(&feeData)
	)

	store.Set(key, value)
}

func (k Keeper) GetGenAuctionLimitBidFeeData(ctx sdk.Context) (auctionFeesCollectionFromLimitBidTx []types.AuctionFeesCollectionFromLimitBidTx) {
	var (
		store = k.Store(ctx)
		iter  = storetypes.KVStorePrefixIterator(store, types.AuctionLimitBidFeeKeyPrefix)
	)

	defer func(iter storetypes.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var data types.AuctionFeesCollectionFromLimitBidTx
		k.cdc.MustUnmarshal(iter.Value(), &data)
		auctionFeesCollectionFromLimitBidTx = append(auctionFeesCollectionFromLimitBidTx, data)
	}

	return auctionFeesCollectionFromLimitBidTx
}
