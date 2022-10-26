package keeper

import (
	v5types "github.com/comdex-official/comdex/app/upgrades/testnet/v5_0_0/types"
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) OldGetPools(ctx sdk.Context) (pools []v5types.PoolOld) {
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
		var pool v5types.PoolOld
		k.cdc.MustUnmarshal(iter.Value(), &pool)
		pools = append(pools, pool)
	}

	return pools
}

func (k Keeper) OldGetAllLend(ctx sdk.Context) (lendAsset []v5types.LendAssetOld) {
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
		var asset v5types.LendAssetOld
		k.cdc.MustUnmarshal(iter.Value(), &asset)
		lendAsset = append(lendAsset, asset)
	}
	return lendAsset
}

func (k Keeper) OldGetAllBorrow(ctx sdk.Context) (borrowAsset []v5types.BorrowAssetOld) {
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
		var asset v5types.BorrowAssetOld
		k.cdc.MustUnmarshal(iter.Value(), &asset)
		borrowAsset = append(borrowAsset, asset)
	}
	return borrowAsset
}
