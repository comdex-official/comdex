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

func (k Keeper) SetAuction(ctx sdk.Context, auction types.Auctions) error {

	var (
		store = k.Store(ctx)
		key   = types.AuctionKey(auction.AppId, auction.AuctionId)
		value = k.cdc.MustMarshal(&auction)
	)

	store.Set(key, value)
	return nil
}

func (k Keeper) DeleteAuction(ctx sdk.Context, auction types.Auctions) error {

	var (
		store = k.Store(ctx)
		key   = types.AuctionKey(auction.AppId, auction.AuctionId)
	)
	store.Delete(key)
	return nil
}

func (k Keeper) GetAuction(ctx sdk.Context, appID, auctionMappingID, auctionID uint64) (auction types.Auctions, err error) {
	var (
		store = k.Store(ctx)
		key   = types.AuctionKey(appID, auctionID)
		value = store.Get(key)
	)

	if value == nil {
		return auction, sdkerrors.ErrNotFound
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, nil
}
