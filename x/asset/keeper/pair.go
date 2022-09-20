package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k Keeper) SetPairID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.PairIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetPairID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.PairIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var count protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &count)

	return count.GetValue()
}

func (k Keeper) SetPair(ctx sdk.Context, pair types.Pair) {
	var (
		store = k.Store(ctx)
		key   = types.PairKey(pair.Id)
		value = k.cdc.MustMarshal(&pair)
	)

	store.Set(key, value)
}

func (k Keeper) GetPair(ctx sdk.Context, id uint64) (pair types.Pair, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.PairKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return pair, false
	}

	k.cdc.MustUnmarshal(value, &pair)
	return pair, true
}

func (k Keeper) GetPairs(ctx sdk.Context) (pairs []types.Pair) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.PairKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var pair types.Pair
		k.cdc.MustUnmarshal(iter.Value(), &pair)
		pairs = append(pairs, pair)
	}

	return pairs
}

func (k Keeper) NewAddPair(ctx sdk.Context, msg *types.MsgAddPairRequest) (*types.MsgAddPairResponse, error) {
	if !k.HasAsset(ctx, msg.AssetIn) {
		return nil, types.ErrorAssetDoesNotExist
	}
	if !k.HasAsset(ctx, msg.AssetOut) {
		return nil, types.ErrorAssetDoesNotExist
	}
	var (
		id   = k.GetPairID(ctx)
		pair = types.Pair{
			Id:       id + 1,
			AssetIn:  msg.AssetIn,
			AssetOut: msg.AssetOut,
		}
	)
	k.SetPairID(ctx, pair.Id)
	k.SetPair(ctx, pair)
	return &types.MsgAddPairResponse{}, nil
}

func (k *Keeper) UpdatePairRecords(ctx sdk.Context, msg types.Pair) error {
	pair, found := k.GetPair(ctx, msg.Id)
	if !found {
		return types.ErrorPairDoesNotExist
	}
	if !k.HasAsset(ctx, msg.AssetIn) {
		return types.ErrorAssetDoesNotExist
	}
	if !k.HasAsset(ctx, msg.AssetOut) {
		return types.ErrorAssetDoesNotExist
	}
	pair.AssetIn = msg.AssetIn
	pair.AssetOut = msg.AssetOut
	k.SetPair(ctx, pair)
	return nil
}