package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k Keeper) IsWhitelistedAsset(ctx sdk.Context, tokenDenom string) bool {
	//store := ctx.KVStore(k.storeKey)
	//key := types.CreateRegisteredTokenKey(tokenDenom)
	//
	//return store.Has(key)
	return true
}

func (k *Keeper) SetAssetID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.WhitelistedAssetIDKey
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
		key   = types.WhitelistedAssetIDKey
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
		iter  = sdk.KVStorePrefixIterator(store, types.WhitelistedAssetKeyPrefix)
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

func (k *Keeper) GetBalancesForModule(ctx sdk.Context, module string) (balances string) {
	if module != "osmo" && module != "cmdx" && module != "atom" {
		return "0"
	}
	moduleBalances := k.bank.GetAllBalances(ctx, k.account.GetModuleAddress(module))
	result := moduleBalances.String()
	if result == "" {
		return "0"
	}
	return result
}
