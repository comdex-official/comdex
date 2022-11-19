package keeper

import (
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/petrichormoney/petri/x/asset/types"
)

func (k Keeper) SetAssetID(ctx sdk.Context, id uint64) {
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

func (k Keeper) GetAssetID(ctx sdk.Context) uint64 {
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

func (k Keeper) SetAsset(ctx sdk.Context, asset types.Asset) {
	var (
		store = k.Store(ctx)
		key   = types.AssetKey(asset.Id)
		value = k.cdc.MustMarshal(&asset)
	)

	store.Set(key, value)
}

func (k Keeper) HasAsset(ctx sdk.Context, id uint64) bool {
	var (
		store = k.Store(ctx)
		key   = types.AssetKey(id)
	)

	return store.Has(key)
}

func (k Keeper) GetAsset(ctx sdk.Context, id uint64) (asset types.Asset, found bool) {
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

func (k Keeper) GetAssetDenom(ctx sdk.Context, id uint64) string {
	asset, _ := k.GetAsset(ctx, id)
	return asset.Denom
}

func (k Keeper) GetAssets(ctx sdk.Context) (assets []types.Asset) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.AssetKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var asset types.Asset
		k.cdc.MustUnmarshal(iter.Value(), &asset)
		assets = append(assets, asset)
	}

	return assets
}

func (k Keeper) SetAssetForDenom(ctx sdk.Context, denom string, id uint64) {
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

func (k Keeper) HasAssetForDenom(ctx sdk.Context, denom string) bool {
	var (
		store = k.Store(ctx)
		key   = types.AssetForDenomKey(denom)
	)

	return store.Has(key)
}

func (k Keeper) GetAssetForDenom(ctx sdk.Context, denom string) (asset types.Asset, found bool) {
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

func (k Keeper) DeleteAssetForDenom(ctx sdk.Context, denom string) {
	var (
		store = k.Store(ctx)
		key   = types.AssetForDenomKey(denom)
	)

	store.Delete(key)
}

func (k Keeper) SetAssetForName(ctx sdk.Context, name string, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.AssetForNameKey(name)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) HasAssetForName(ctx sdk.Context, name string) bool {
	var (
		store = k.Store(ctx)
		key   = types.AssetForNameKey(name)
	)

	return store.Has(key)
}

func (k Keeper) DeleteAssetForName(ctx sdk.Context, name string) {
	var (
		store = k.Store(ctx)
		key   = types.AssetForNameKey(name)
	)

	store.Delete(key)
}

func (k Keeper) AddAssetRecords(ctx sdk.Context, msg types.Asset) error {
	if k.HasAssetForDenom(ctx, msg.Denom) || k.HasAssetForName(ctx, msg.Name) {
		return types.ErrorDuplicateAsset
	}

	IsLetter := regexp.MustCompile(`^[A-Z]+$`).MatchString

	if !IsLetter(msg.Name) || len(msg.Name) > 10 {
		return types.ErrorNameDidNotMeetCriterion
	}
	if !msg.IsOnChain && msg.IsCdpMintable {
		return types.ErrorOffChainAssetCannotBeMintable
	}

	var (
		id    = k.GetAssetID(ctx)
		asset = types.Asset{
			Id:                    id + 1,
			Name:                  msg.Name,
			Denom:                 msg.Denom,
			Decimals:              msg.Decimals,
			IsOnChain:             msg.IsOnChain,
			IsOraclePriceRequired: msg.IsOraclePriceRequired,
			IsCdpMintable:         msg.IsCdpMintable,
		}
	)
	if msg.IsOraclePriceRequired {
		k.bandoracle.SetCheckFlag(ctx, false)
	}

	k.SetAssetID(ctx, asset.Id)
	k.SetAsset(ctx, asset)
	k.SetAssetForDenom(ctx, asset.Denom, asset.Id)
	k.SetAssetForName(ctx, asset.Name, asset.Id)

	return nil
}

func (k Keeper) UpdateAssetRecords(ctx sdk.Context, msg types.Asset) error {
	asset, found := k.GetAsset(ctx, msg.Id)
	if !found {
		return types.ErrorAssetDoesNotExist
	}
	IsLetter := regexp.MustCompile(`^[A-Z]+$`).MatchString

	if !IsLetter(msg.Name) || len(msg.Name) > 10 {
		return types.ErrorNameDidNotMeetCriterion
	}

	if k.HasAssetForName(ctx, msg.Name) && asset.Name != msg.Name {
		return types.ErrorDuplicateAsset
	}
	k.DeleteAssetForName(ctx, asset.Name)
	asset.Name = msg.Name
	k.SetAssetForName(ctx, asset.Name, asset.Id)

	if k.HasAssetForDenom(ctx, msg.Denom) && asset.Denom != msg.Denom {
		return types.ErrorDuplicateAsset
	}

	k.DeleteAssetForDenom(ctx, asset.Denom)
	asset.Denom = msg.Denom
	k.SetAssetForDenom(ctx, asset.Denom, asset.Id)

	if msg.Decimals.GT(sdk.ZeroInt()) {
		asset.Decimals = msg.Decimals
	}
	asset.IsOraclePriceRequired = msg.IsOraclePriceRequired
	if msg.IsOraclePriceRequired {
		k.bandoracle.SetCheckFlag(ctx, false)
	}

	k.SetAsset(ctx, asset)
	return nil
}
