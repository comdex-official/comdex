package keeper

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/cosmos/gogoproto/types"

	"github.com/comdex-official/comdex/x/lend/types"
)

func (k Keeper) SetUserBorrowIDCounter(ctx sdk.Context, ID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowCounterIDPrefix
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: ID,
			},
		)
	)
	store.Set(key, value)
}

func (k Keeper) GetUserBorrowIDCounter(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.BorrowCounterIDPrefix
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var ID protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &ID)

	return ID.GetValue()
}

func (k Keeper) SetBorrow(ctx sdk.Context, borrow types.BorrowAsset) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowUserKey(borrow.ID)
		value = k.cdc.MustMarshal(&borrow)
	)

	store.Set(key, value)
}

func (k Keeper) GetBorrow(ctx sdk.Context, ID uint64) (borrow types.BorrowAsset, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowUserKey(ID)
		value = store.Get(key)
	)

	if value == nil {
		return borrow, false
	}

	k.cdc.MustUnmarshal(value, &borrow)
	return borrow, true
}

func (k Keeper) GetAllBorrow(ctx sdk.Context) (borrowAsset []types.BorrowAsset) {
	var (
		store = k.Store(ctx)
		iter  = storetypes.KVStorePrefixIterator(store, types.BorrowPairKeyPrefix)
	)

	defer func(iter storetypes.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var asset types.BorrowAsset
		k.cdc.MustUnmarshal(iter.Value(), &asset)
		borrowAsset = append(borrowAsset, asset)
	}
	return borrowAsset
}

func (k Keeper) DeleteBorrow(ctx sdk.Context, ID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowUserKey(ID)
	)

	store.Delete(key)
}

func (k Keeper) HasBorrowForAddressByPair(ctx sdk.Context, address string, pairID uint64) bool {
	mappingData := k.GetUserTotalMappingData(ctx, address)
	for _, data := range mappingData {
		for _, indata := range data.BorrowId {
			borrowData, _ := k.GetBorrow(ctx, indata)
			if borrowData.PairID == pairID {
				return true
			}
		}
	}
	return false
}

func (k Keeper) GetBorrows(ctx sdk.Context) (borrowIds []uint64, found bool) {
	var (
		store = k.Store(ctx)
		iter  = storetypes.KVStorePrefixIterator(store, types.AssetStatsByPoolIDAndAssetIDKeyPrefix)
	)

	defer func(iter storetypes.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var asset types.PoolAssetLBMapping
		k.cdc.MustUnmarshal(iter.Value(), &asset)
		borrowIds = append(borrowIds, asset.BorrowIds...)
	}
	return borrowIds, true
}

func (k Keeper) SetBorrowInterestTracker(ctx sdk.Context, interest types.BorrowInterestTracker) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowInterestTrackerKey(interest.BorrowingId)
		value = k.cdc.MustMarshal(&interest)
	)

	store.Set(key, value)
}

func (k Keeper) GetBorrowInterestTracker(ctx sdk.Context, ID uint64) (interest types.BorrowInterestTracker, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowInterestTrackerKey(ID)
		value = store.Get(key)
	)

	if value == nil {
		return interest, false
	}

	k.cdc.MustUnmarshal(value, &interest)
	return interest, true
}

func (k Keeper) GetAllBorrowInterestTracker(ctx sdk.Context) (interest []types.BorrowInterestTracker) {
	var (
		store = k.Store(ctx)
		iter  = storetypes.KVStorePrefixIterator(store, types.BorrowInterestTrackerKeyPrefix)
	)

	defer func(iter storetypes.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var tracker types.BorrowInterestTracker
		k.cdc.MustUnmarshal(iter.Value(), &tracker)
		interest = append(interest, tracker)
	}
	return interest
}

func (k Keeper) DeleteBorrowInterestTracker(ctx sdk.Context, ID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowInterestTrackerKey(ID)
	)

	store.Delete(key)
}

func (k Keeper) GetBorrowByUserAndAssetID(ctx sdk.Context, owner, debtDenom string, assetIn uint64) types.BorrowAsset {
	var borrowAsset types.BorrowAsset
	mappingData := k.GetUserTotalMappingData(ctx, owner)
	for _, data := range mappingData {
		lend, _ := k.GetLend(ctx, data.LendId)
		if lend.AssetID == assetIn {
			for _, borrowID := range data.BorrowId {
				borrow, _ := k.GetBorrow(ctx, borrowID)
				if borrow.AmountOut.Denom == debtDenom {
					borrowAsset = borrow
				}
			}
		}
	}
	return borrowAsset
}
