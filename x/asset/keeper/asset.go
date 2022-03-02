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

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetAsset(ctx sdk.Context, asset types.Asset) {
	var (
		store = k.Store(ctx)
		key   = types.AssetKey(asset.Id)
		value = k.cdc.MustMarshal(&asset)
	)

	store.Set(key, value)
}

func (k *Keeper) HasAsset(ctx sdk.Context, id uint64) bool {
	var (
		store = k.Store(ctx)
		key   = types.AssetKey(id)
	)

	return store.Has(key)
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

func (k *Keeper) SetAssetForDenom(ctx sdk.Context, denom string, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.AssetForDenomKey(denom)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) HasAssetForDenom(ctx sdk.Context, denom string) bool {
	var (
		store = k.Store(ctx)
		key   = types.AssetForDenomKey(denom)
	)

	return store.Has(key)
}

func (k *Keeper) GetAssetForDenom(ctx sdk.Context, denom string) (asset types.Asset, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AssetForDenomKey(denom)
		value = store.Get(key)
	)

	if value == nil {
		return asset, false
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return k.GetAsset(ctx, id.GetValue())
}

func (k *Keeper) DeleteAssetForDenom(ctx sdk.Context, denom string) {
	var (
		store = k.Store(ctx)
		key   = types.AssetForDenomKey(denom)
	)

	store.Delete(key)
}

func (k *Keeper) GetPriceForAsset(ctx sdk.Context, id uint64) (uint64, bool) {
	market, found := k.oracle.GetMarketForAsset(ctx, id)
	if !found {
		return 0, false
	}

	return k.oracle.GetPriceForMarket(ctx, market.Symbol)
}

func (k *Keeper) UpdateLiquidationRatio(ctx sdk.Context, prop *types.UpdateLiquidationRatio) error {

	params := k.GetParams(ctx)

	params.LiquidationRatio = prop.LiquidationRatio
	k.SetParams(ctx, params)

	return nil
}
