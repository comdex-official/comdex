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

func (k *Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	var (
		store = k.Store(ctx)
		key   = types.PoolKey(pool.Id)
		value = k.cdc.MustMarshalBinaryBare(&pool)
	)

	store.Set(key, value)
}

func (k *Keeper) GetPool(ctx sdk.Context, id uint64) (pool types.Pool, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.PoolKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return pool, false
	}

	k.cdc.MustUnmarshalBinaryBare(value, &pool)
	return pool, true
}

func (k *Keeper) GetPools(ctx sdk.Context) (pools []types.Pool) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.PoolKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var pool types.Pool
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &pool)
		pools = append(pools, pool)
	}

	return pools
}
