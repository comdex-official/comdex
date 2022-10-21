package keeper

import (
	migrationTypes "github.com/comdex-official/comdex/x/lend/migrations/types"
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) OldGetPools(ctx sdk.Context) (pools []migrationTypes.Pool) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.PoolKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var pool migrationTypes.Pool
		k.cdc.MustUnmarshal(iter.Value(), &pool)
		pools = append(pools, pool)
	}

	return pools
}

func (k Keeper) OldGetAllLend(ctx sdk.Context) (lendAsset []migrationTypes.LendAsset) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.LendUserPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var asset migrationTypes.LendAsset
		k.cdc.MustUnmarshal(iter.Value(), &asset)
		lendAsset = append(lendAsset, asset)
	}
	return lendAsset
}

func (k Keeper) OldGetAllBorrow(ctx sdk.Context) (borrowAsset []migrationTypes.BorrowAsset) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.BorrowPairKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var asset migrationTypes.BorrowAsset
		k.cdc.MustUnmarshal(iter.Value(), &asset)
		borrowAsset = append(borrowAsset, asset)
	}
	return borrowAsset
}
