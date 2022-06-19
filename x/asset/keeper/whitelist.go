package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k *Keeper) GetWhitelistAssetID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.WhitelistAssetIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetWhitelistAssetID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.WhitelistAssetIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) SetWhitelistAsset(ctx sdk.Context, asset types.ExtendedAsset) {
	var (
		store = k.Store(ctx)
		key   = types.WhitelistAssetKey(asset.Id)
		value = k.cdc.MustMarshal(&asset)
	)

	store.Set(key, value)
}

func (k *Keeper) GetWhitelistAsset(ctx sdk.Context, id uint64) (asset types.ExtendedAsset, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.WhitelistAssetKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return asset, false
	}

	k.cdc.MustUnmarshal(value, &asset)
	return asset, true
}

func (k *Keeper) GetWhiteListPairID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.WhitelistPairIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var count protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &count)

	return count.GetValue()
}

func (k *Keeper) HasWhitelistAsset(ctx sdk.Context, id uint64) bool {
	var (
		store = k.Store(ctx)
		key   = types.WhitelistAssetKey(id)
	)

	return store.Has(key)
}

func (k *Keeper) SetWhitelistPairID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.WhitelistPairIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) SetWhitelistPair(ctx sdk.Context, pair types.ExtendedPairLend) {
	var (
		store = k.Store(ctx)
		key   = types.WhitelistPairKey(pair.Id)
		value = k.cdc.MustMarshal(&pair)
	)

	store.Set(key, value)
}

func (k *Keeper) GetWhitelistPair(ctx sdk.Context, id uint64) (pair types.ExtendedPairLend, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.WhitelistPairKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return pair, false
	}

	k.cdc.MustUnmarshal(value, &pair)
	return pair, true
}

func (k *Keeper) AddWhitelistedAssetRecords(ctx sdk.Context, records ...types.ExtendedAsset) error {
	for _, msg := range records {
		if !k.HasAsset(ctx, msg.AssetId) {
			return types.ErrorAssetDoesNotExist
		}
		if k.HasWhitelistAsset(ctx, msg.Id) {
			return types.ErrorDuplicateAsset
		}

		var (
			id    = k.GetWhitelistAssetID(ctx)
			asset = types.ExtendedAsset{
				Id:                   id + 1,
				AssetId:              msg.AssetId,
				CollateralWeight:     msg.CollateralWeight,
				LiquidationThreshold: msg.LiquidationThreshold,
				IsBridgedAsset:       msg.IsBridgedAsset,
			}
		)

		k.SetWhitelistAssetID(ctx, asset.Id)
		k.SetWhitelistAsset(ctx, asset)

	}

	return nil
}

func (k *Keeper) UpdateWhitelistedAssetRecords(ctx sdk.Context, msg types.ExtendedAsset) error {
	asset, found := k.GetWhitelistAsset(ctx, msg.Id)
	if !found {
		return types.ErrorAssetDoesNotExist
	}

	if !msg.CollateralWeight.IsZero() {
		asset.CollateralWeight = msg.CollateralWeight
	}
	if !msg.LiquidationThreshold.IsZero() {
		asset.LiquidationThreshold = msg.LiquidationThreshold
	}
	if msg.IsBridgedAsset || !msg.IsBridgedAsset {
		asset.IsBridgedAsset = msg.IsBridgedAsset

	}

	k.SetWhitelistAsset(ctx, asset)
	return nil
}

func (k *Keeper) AddWhitelistedPairsRecords(ctx sdk.Context, records ...types.ExtendedPairLend) error {
	for _, msg := range records {
		_, found := k.GetPair(ctx, msg.PairId)
		if !found {
			return types.ErrorPairDoesNotExist
		}
		_, got := k.GetWhitelistPair(ctx, msg.Id)
		if got {
			return types.ErrorDuplicatePair
		}

		var (
			id   = k.GetWhiteListPairID(ctx)
			pair = types.ExtendedPairLend{
				Id:                    id + 1,
				PairId:                msg.PairId,
				ModuleAcc:             msg.ModuleAcc,
				BaseBorrowRateAsset_1: msg.BaseBorrowRateAsset_1,
				BaseLendRateAsset_1:   msg.BaseLendRateAsset_1,
				BaseBorrowRateAsset_2: msg.BaseBorrowRateAsset_2,
				BaseLendRateAsset_2:   msg.BaseLendRateAsset_2,
			}
		)

		k.SetWhitelistPairID(ctx, pair.Id)
		k.SetWhitelistPair(ctx, pair)
	}
	return nil
}

func (k *Keeper) UpdateWhitelistedPairRecords(ctx sdk.Context, msg types.ExtendedPairLend) error {

	pair, found := k.GetWhitelistPair(ctx, msg.Id)
	if !found {
		return types.ErrorPairDoesNotExist
	}

	if len(msg.ModuleAcc) > 0 {
		pair.ModuleAcc = msg.ModuleAcc
	}
	if !msg.BaseBorrowRateAsset_1.IsZero() {
		pair.BaseBorrowRateAsset_1 = msg.BaseBorrowRateAsset_1
	}
	if !msg.BaseBorrowRateAsset_2.IsZero() {
		pair.BaseBorrowRateAsset_2 = msg.BaseBorrowRateAsset_2
	}
	if !msg.BaseLendRateAsset_1.IsZero() {
		pair.BaseLendRateAsset_1 = msg.BaseLendRateAsset_1
	}
	if !msg.BaseLendRateAsset_2.IsZero() {
		pair.BaseLendRateAsset_2 = msg.BaseLendRateAsset_2
	}

	k.SetWhitelistPair(ctx, pair)
	return nil
}
