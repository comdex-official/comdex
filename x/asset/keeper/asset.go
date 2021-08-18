package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k *Keeper) SetAssetID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.AssetIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetAssetID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.AssetIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var count protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &count)

	return count.GetValue()
}

func (k *Keeper) SetAsset(ctx sdk.Context, asset types.Asset) {
	var (
		store = k.Store(ctx)
		key   = types.AssetKey(asset.ID)
		value = k.cdc.MustMarshal(&asset)
	)

	store.Set(key, value)
}

func (k *Keeper) GetAsset(ctx sdk.Context, id uint64) (asset types.Asset, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AssetKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return asset, false
	}

	k.cdc.MustUnmarshal(value, &asset)
	return asset, true
}

func (k *Keeper) GetAssets(ctx sdk.Context) (assets []types.Asset) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.AssetKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var asset types.Asset
		k.cdc.MustUnmarshal(iter.Value(), &asset)
		assets = append(assets, asset)
	}

	return assets
}

func (k *Keeper) SetAssetIDForMarket(ctx sdk.Context, id uint64, symbol string) {
	var (
		store = k.Store(ctx)
		key   = types.AssetForMarketKey(symbol)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetAssetIDForMarket(ctx sdk.Context, symbol string) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.AssetForMarketKey(symbol)
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) GetAssetForMarket(ctx sdk.Context, symbol string) (asset types.Asset, found bool) {
	id := k.GetAssetIDForMarket(ctx, symbol)
	if id == 0 {
		return asset, false
	}

	return k.GetAsset(ctx, id)
}
