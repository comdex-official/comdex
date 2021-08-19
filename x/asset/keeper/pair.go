package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k *Keeper) SetPairID(ctx sdk.Context, id uint64) {
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

func (k *Keeper) GetPairID(ctx sdk.Context) uint64 {
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

func (k *Keeper) SetPair(ctx sdk.Context, pair types.Pair) {
	var (
		store = k.Store(ctx)
		key   = types.PairKey(pair.ID)
		value = k.cdc.MustMarshal(&pair)
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

	k.cdc.MustUnmarshal(value, &pair)
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
		k.cdc.MustUnmarshal(iter.Value(), &pair)
		pairs = append(pairs, pair)
	}

	return pairs
}
