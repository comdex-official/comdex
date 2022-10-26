package keeper

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/lend/types"
)

func (k Keeper) SetUserBorrowIDCounter(ctx sdk.Context, ID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.NewBorrowCounterIDPrefix
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
		key   = types.NewBorrowCounterIDPrefix
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
		key   = types.NewBorrowUserKey(borrow.ID)
		value = k.cdc.MustMarshal(&borrow)
	)

	store.Set(key, value)
}

func (k Keeper) GetBorrow(ctx sdk.Context, ID uint64) (borrow types.BorrowAsset, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.NewBorrowUserKey(ID)
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
		iter  = sdk.KVStorePrefixIterator(store, types.NewBorrowPairKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
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
		key   = types.NewBorrowUserKey(ID)
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
	assetStats := k.GetAllAssetStatsByPoolIDAndAssetID(ctx)
	for _, data := range assetStats {
		borrowIds = append(borrowIds, data.BorrowIds...)
	}
	return borrowIds, true
}

func (k Keeper) SetBorrowInterestTracker(ctx sdk.Context, interest types.BorrowInterestTracker) {
	var (
		store = k.Store(ctx)
		key   = types.NewBorrowInterestTrackerKey(interest.BorrowingId)
		value = k.cdc.MustMarshal(&interest)
	)

	store.Set(key, value)
}

func (k Keeper) GetBorrowInterestTracker(ctx sdk.Context, ID uint64) (interest types.BorrowInterestTracker, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.NewBorrowInterestTrackerKey(ID)
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
		iter  = sdk.KVStorePrefixIterator(store, types.NewBorrowInterestTrackerKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
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
		key   = types.NewBorrowInterestTrackerKey(ID)
	)

	store.Delete(key)
}

func (k Keeper) SetStableBorrowIds(ctx sdk.Context, borrow types.StableBorrowMapping) {
	var (
		store = k.Store(ctx)
		key   = types.NewStableBorrowIDsKeyPrefix
		value = k.cdc.MustMarshal(&borrow)
	)

	store.Set(key, value)
}

func (k Keeper) GetStableBorrowIds(ctx sdk.Context) (borrow types.StableBorrowMapping) {
	var (
		store = k.Store(ctx)
		key   = types.NewStableBorrowIDsKeyPrefix
		value = store.Get(key)
	)

	k.cdc.MustUnmarshal(value, &borrow)
	return borrow
}

func (k Keeper) DeleteIDFromStableBorrowMapping(ctx sdk.Context, id uint64) {
	stableIds := k.GetStableBorrowIds(ctx)
	lengthOfIDs := len(stableIds.StableBorrowIDs)
	dataIndex := sort.Search(lengthOfIDs, func(i int) bool { return stableIds.StableBorrowIDs[i] >= id })

	if dataIndex < lengthOfIDs && stableIds.StableBorrowIDs[dataIndex] == id {
		stableIds.StableBorrowIDs = append(stableIds.StableBorrowIDs[:dataIndex], stableIds.StableBorrowIDs[dataIndex+1:]...)
		k.SetStableBorrowIds(ctx, stableIds)
	}
}
