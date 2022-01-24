package keeper

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k Keeper) StartCollateralAuction(
	ctx sdk.Context,
	vault vaulttypes.Vault,
	collateralizationRatio sdk.Dec,
	assetIn assettypes.Asset,
	assetOut assettypes.Asset,
) {
	auction := auctiontypes.NewCollateralAuction(
		vault,
		collateralizationRatio,
		assetIn,
		assetOut,
	)
	auction.ID = k.GetAuctionID(ctx)
	k.SetAuctionID(ctx, auction.ID)
	k.SetAuction(ctx, auction)
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
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

	return id.GetValue() + 1
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

func (k *Keeper) SetAuction(ctx sdk.Context, auction auctiontypes.CollateralAuction) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionKey(auction.ID)
		value = k.cdc.MustMarshal(&auction)
	)
	store.Set(key, value)
}

func (k *Keeper) GetAuctions(ctx sdk.Context) (auctions []auctiontypes.CollateralAuction) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, auctiontypes.AuctionKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var auction auctiontypes.CollateralAuction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		auctions = append(auctions, auction)
	}

	return auctions
}

func (k *Keeper) GetAuction(ctx sdk.Context, id uint64) (auction auctiontypes.CollateralAuction, found bool) {
	var (
		store = k.Store(ctx)
		key   = auctiontypes.AuctionKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return auction, false
	}

	k.cdc.MustUnmarshal(value, &auction)
	return auction, true
}
