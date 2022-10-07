package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

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
		iter  = sdk.KVStorePrefixIterator(store, types.BorrowPairKeyPrefix)
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
		key   = types.BorrowUserKey(ID)
	)

	store.Delete(key)
}

func (k Keeper) SetBorrowForAddressByPair(ctx sdk.Context, address sdk.AccAddress, pairID, ID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowForAddressByPair(address, pairID)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: ID,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) HasBorrowForAddressByPair(ctx sdk.Context, address sdk.AccAddress, pairID uint64) bool {
	var (
		store = k.Store(ctx)
		key   = types.BorrowForAddressByPair(address, pairID)
	)

	return store.Has(key)
}

func (k Keeper) DeleteBorrowForAddressByPair(ctx sdk.Context, address sdk.AccAddress, pairID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowForAddressByPair(address, pairID)
	)

	store.Delete(key)
}

// func (k Keeper) UpdateUserBorrowIDMapping(
// 	ctx sdk.Context,
// 	lendOwner string,
// 	borrowID uint64,
// 	isInsert bool,
// ) error {
// 	userBorrows, found := k.GetUserBorrows(ctx, lendOwner)

// 	if !found && isInsert {
// 		userBorrows = types.UserBorrowIdMapping{
// 			Owner:     lendOwner,
// 			BorrowIDs: nil,
// 		}
// 	} else if !found && !isInsert {
// 		return types.ErrorLendOwnerNotFound
// 	}

// 	if isInsert {
// 		userBorrows.BorrowIDs = append(userBorrows.BorrowIDs, borrowID)
// 	} else {
// 		for index, id := range userBorrows.BorrowIDs {
// 			if id == borrowID {
// 				userBorrows.BorrowIDs = append(userBorrows.BorrowIDs[:index], userBorrows.BorrowIDs[index+1:]...)
// 				break
// 			}
// 		}
// 	}

// 	k.SetUserBorrows(ctx, userBorrows)
// 	return nil
// }

// func (k Keeper) GetUserBorrows(ctx sdk.Context, address string) (userBorrows types.UserBorrowIdMapping, found bool) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.UserBorrowsForAddressKey(address)
// 		value = store.Get(key)
// 	)
// 	if value == nil {
// 		return userBorrows, false
// 	}
// 	k.cdc.MustUnmarshal(value, &userBorrows)

// 	return userBorrows, true
// }

// func (k Keeper) GetAllUserBorrows(ctx sdk.Context) (userBorrowIDMapping []types.UserBorrowIdMapping) {
// 	var (
// 		store = k.Store(ctx)
// 		iter  = sdk.KVStorePrefixIterator(store, types.UserBorrowsForAddressKeyPrefix)
// 	)

// 	defer func(iter sdk.Iterator) {
// 		err := iter.Close()
// 		if err != nil {
// 			return
// 		}
// 	}(iter)

// 	for ; iter.Valid(); iter.Next() {
// 		var asset types.UserBorrowIdMapping
// 		k.cdc.MustUnmarshal(iter.Value(), &asset)
// 		userBorrowIDMapping = append(userBorrowIDMapping, asset)
// 	}
// 	return userBorrowIDMapping
// }

// func (k Keeper) UserBorrows(ctx sdk.Context, address string) (userBorrows []types.BorrowAsset, found bool) {
// 	userBorrowID, _ := k.GetUserBorrows(ctx, address)
// 	for _, v := range userBorrowID.BorrowIDs {
// 		userBorrow, _ := k.GetBorrow(ctx, v)
// 		userBorrows = append(userBorrows, userBorrow)
// 	}
// 	return userBorrows, true
// }

// func (k Keeper) SetUserBorrows(ctx sdk.Context, userBorrows types.UserBorrowIdMapping) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.UserBorrowsForAddressKey(userBorrows.Owner)
// 		value = k.cdc.MustMarshal(&userBorrows)
// 	)
// 	store.Set(key, value)
// }

// func (k Keeper) UpdateBorrowIDByOwnerAndPoolMapping(
// 	ctx sdk.Context,
// 	borrowOwner string,
// 	borrowID uint64,
// 	poolID uint64,
// 	isInsert bool,
// ) error {
// 	userBorrows, found := k.GetBorrowIDByOwnerAndPool(ctx, borrowOwner, poolID)

// 	if !found && isInsert {
// 		userBorrows = types.BorrowIdByOwnerAndPoolMapping{
// 			Owner:     borrowOwner,
// 			PoolID:    poolID,
// 			BorrowIDs: nil,
// 		}
// 	} else if !found && !isInsert {
// 		return types.ErrorLendOwnerNotFound
// 	}

// 	if isInsert {
// 		userBorrows.BorrowIDs = append(userBorrows.BorrowIDs, borrowID)
// 	} else {
// 		for index, id := range userBorrows.BorrowIDs {
// 			if id == borrowID {
// 				userBorrows.BorrowIDs = append(userBorrows.BorrowIDs[:index], userBorrows.BorrowIDs[index+1:]...)
// 				break
// 			}
// 		}
// 	}

// 	k.SetBorrowIDByOwnerAndPool(ctx, userBorrows)
// 	return nil
// }

// func (k Keeper) GetBorrowIDByOwnerAndPool(ctx sdk.Context, address string, poolID uint64) (userBorrows types.BorrowIdByOwnerAndPoolMapping, found bool) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.BorrowByUserAndPoolKey(address, poolID)
// 		value = store.Get(key)
// 	)
// 	if value == nil {
// 		return userBorrows, false
// 	}
// 	k.cdc.MustUnmarshal(value, &userBorrows)

// 	return userBorrows, true
// }

// func (k Keeper) GetAllBorrowIDByOwnerAndPool(ctx sdk.Context) (borrowIDByOwnerAndPoolMapping []types.BorrowIdByOwnerAndPoolMapping) {
// 	var (
// 		store = k.Store(ctx)
// 		iter  = sdk.KVStorePrefixIterator(store, types.BorrowByUserAndPoolPrefix)
// 	)

// 	defer func(iter sdk.Iterator) {
// 		err := iter.Close()
// 		if err != nil {
// 			return
// 		}
// 	}(iter)

// 	for ; iter.Valid(); iter.Next() {
// 		var asset types.BorrowIdByOwnerAndPoolMapping
// 		k.cdc.MustUnmarshal(iter.Value(), &asset)
// 		borrowIDByOwnerAndPoolMapping = append(borrowIDByOwnerAndPoolMapping, asset)
// 	}
// 	return borrowIDByOwnerAndPoolMapping
// }

// func (k Keeper) BorrowIDByOwnerAndPool(ctx sdk.Context, address string, poolID uint64) (userBorrows []types.BorrowAsset, found bool) {
// 	userLendID, _ := k.GetBorrowIDByOwnerAndPool(ctx, address, poolID)
// 	for _, v := range userLendID.BorrowIDs {
// 		userBorrow, _ := k.GetBorrow(ctx, v)
// 		userBorrows = append(userBorrows, userBorrow)
// 	}
// 	return userBorrows, true
// }

// func (k Keeper) SetBorrowIDByOwnerAndPool(ctx sdk.Context, userBorrows types.BorrowIdByOwnerAndPoolMapping) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.BorrowByUserAndPoolKey(userBorrows.Owner, userBorrows.PoolID)
// 		value = k.cdc.MustMarshal(&userBorrows)
// 	)
// 	store.Set(key, value)
// }

// func (k Keeper) UpdateBorrowIdsMapping(
// 	ctx sdk.Context,
// 	borrowID uint64,
// 	isInsert bool,
// ) error {
// 	userBorrows, found := k.GetBorrows(ctx)

// 	if !found && isInsert {
// 		userBorrows = types.BorrowMapping{
// 			BorrowIDs: nil,
// 		}
// 	} else if !found && !isInsert {
// 		return types.ErrorLendOwnerNotFound
// 	}

// 	if isInsert {
// 		userBorrows.BorrowIDs = append(userBorrows.BorrowIDs, borrowID)
// 	} else {
// 		for index, id := range userBorrows.BorrowIDs {
// 			if id == borrowID {
// 				userBorrows.BorrowIDs = append(userBorrows.BorrowIDs[:index], userBorrows.BorrowIDs[index+1:]...)
// 				break
// 			}
// 		}
// 	}

// 	k.SetBorrows(ctx, userBorrows)
// 	return nil
// }

// func (k Keeper) GetBorrows(ctx sdk.Context) (userBorrows types.BorrowMapping, found bool) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.BorrowsKey
// 		value = store.Get(key)
// 	)
// 	if value == nil {
// 		return userBorrows, false
// 	}
// 	k.cdc.MustUnmarshal(value, &userBorrows)

// 	return userBorrows, true
// }

// func (k Keeper) SetBorrows(ctx sdk.Context, userBorrows types.BorrowMapping) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.BorrowsKey
// 		value = k.cdc.MustMarshal(&userBorrows)
// 	)
// 	store.Set(key, value)
// }

// func (k Keeper) SetReservePoolRecordsForBorrow(ctx sdk.Context, borrow types.ReservePoolRecordsForBorrow) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.ReservePoolRecordsForBorrowKey(borrow.ID)
// 		value = k.cdc.MustMarshal(&borrow)
// 	)

// 	store.Set(key, value)
// }

// func (k Keeper) GetReservePoolRecordsForBorrow(ctx sdk.Context, ID uint64) (borrow types.ReservePoolRecordsForBorrow, found bool) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.ReservePoolRecordsForBorrowKey(ID)
// 		value = store.Get(key)
// 	)

// 	if value == nil {
// 		return borrow, false
// 	}

// 	k.cdc.MustUnmarshal(value, &borrow)
// 	return borrow, true
// }

// func (k Keeper) DeleteReservePoolRecordsForBorrow(ctx sdk.Context, ID uint64) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.ReservePoolRecordsForBorrowKey(ID)
// 	)

// 	store.Delete(key)
// }

func (k Keeper) SetBorrowInterestTracker(ctx sdk.Context, interest types.BorrowInterestTracker) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowInterestTrackerKey(interest.BorrowingId)
		value = k.cdc.MustMarshal(&interest)
	)

	store.Set(key, value)
}

func (k Keeper) GetBorrowInterestTracker(ctx sdk.Context, id uint64) (interest types.BorrowInterestTracker, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowInterestTrackerKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return interest, false
	}

	k.cdc.MustUnmarshal(value, &interest)
	return interest, true
}
