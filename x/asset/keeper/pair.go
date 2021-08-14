package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k *Keeper) SetCount(ctx sdk.Context, count uint64) {
	var (
		store = k.Store(ctx)
		key   = types.CountKey
		value = k.cdc.MustMarshalBinaryBare(
			&protobuftypes.UInt64Value{
				Value: count,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetCount(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.CountKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var count protobuftypes.UInt64Value
	k.cdc.MustUnmarshalBinaryBare(value, &count)

	return count.GetValue()
}

func (k *Keeper) SetPair(ctx sdk.Context, pair types.Pair) {
	var (
		store = k.Store(ctx)
		key   = types.PairKey(pair.Id)
		value = k.cdc.MustMarshalBinaryBare(&pair)
	)

	store.Set(key, value)
}

func (k *Keeper) GetPair(ctx sdk.Context, id uint64) (pair types.Pair, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.PairKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return pair, false
	}

	k.cdc.MustUnmarshalBinaryBare(value, &pair)
	return pair, true
}

func (k *Keeper) GetPairs(ctx sdk.Context) (pairs []types.Pair) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.PairKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var pair types.Pair
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &pair)
		pairs = append(pairs, pair)
	}

	return pairs
}
