package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k *Keeper) SetUserBorrowIDHistory(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowHistoryIdPrefix
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k *Keeper) GetUserBorrowIDHistory(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.BorrowHistoryIdPrefix
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetBorrow(ctx sdk.Context, lend types.BorrowAsset) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowUserKey(lend.ID)
		value = k.cdc.MustMarshal(&lend)
	)

	store.Set(key, value)
}

func (k *Keeper) GetBorrow(ctx sdk.Context, id uint64) (lend types.BorrowAsset, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowUserKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return lend, false
	}

	k.cdc.MustUnmarshal(value, &lend)
	return lend, true
}

func (k *Keeper) DeleteBorrow(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowUserKey(id)
	)

	store.Delete(key)
}

func (k *Keeper) SetBorrowForAddressByPair(ctx sdk.Context, address sdk.AccAddress, pairID, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowForAddressByPair(address, pairID)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) HasBorrowForAddressByPair(ctx sdk.Context, address sdk.AccAddress, pairID uint64) bool {
	var (
		store = k.Store(ctx)
		key   = types.BorrowForAddressByPair(address, pairID)
	)

	return store.Has(key)
}

func (k *Keeper) DeleteBorrowForAddressByPair(ctx sdk.Context, address sdk.AccAddress, pairID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowForAddressByPair(address, pairID)
	)

	store.Delete(key)
}

func (k *Keeper) UpdateUserBorrowIdMapping(
	ctx sdk.Context,
	lendOwner string,
	borrowId uint64,
	isInsert bool,
) error {

	userVaults, found := k.GetUserBorrows(ctx, lendOwner)

	if !found && isInsert {
		userVaults = types.UserBorrowIdMapping{
			Owner:     lendOwner,
			BorrowIds: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorLendOwnerNotFound
	}

	if isInsert {
		userVaults.BorrowIds = append(userVaults.BorrowIds, borrowId)
	} else {
		for index, id := range userVaults.BorrowIds {
			if id == borrowId {
				userVaults.BorrowIds = append(userVaults.BorrowIds[:index], userVaults.BorrowIds[index+1:]...)
				break
			}
		}
	}

	k.SetUserBorrows(ctx, userVaults)
	return nil
}

func (k *Keeper) GetUserBorrows(ctx sdk.Context, address string) (userBorrows types.UserBorrowIdMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserBorrowsForAddressKey(address)
		value = store.Get(key)
	)
	if value == nil {
		return userBorrows, false
	}
	k.cdc.MustUnmarshal(value, &userBorrows)

	return userBorrows, true
}

func (k *Keeper) SetUserBorrows(ctx sdk.Context, userBorrows types.UserBorrowIdMapping) {
	var (
		store = k.Store(ctx)
		key   = types.UserBorrowsForAddressKey(userBorrows.Owner)
		value = k.cdc.MustMarshal(&userBorrows)
	)
	store.Set(key, value)
}

func (k *Keeper) UpdateBorrowIdByOwnerAndPoolMapping(
	ctx sdk.Context,
	borrowOwner string,
	borrowId uint64,
	poolId uint64,
	isInsert bool,
) error {

	userLends, found := k.GetBorrowIdByOwnerAndPool(ctx, borrowOwner, poolId)

	if !found && isInsert {
		userLends = types.BorrowIdByOwnerAndPoolMapping{
			Owner:     borrowOwner,
			PoolId:    poolId,
			BorrowIds: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorLendOwnerNotFound
	}

	if isInsert {
		userLends.BorrowIds = append(userLends.BorrowIds, borrowId)
	} else {
		for index, id := range userLends.BorrowIds {
			if id == borrowId {
				userLends.BorrowIds = append(userLends.BorrowIds[:index], userLends.BorrowIds[index+1:]...)
				break
			}
		}
	}

	k.SetBorrowIdByOwnerAndPool(ctx, userLends)
	return nil
}

func (k *Keeper) GetBorrowIdByOwnerAndPool(ctx sdk.Context, address string, poolId uint64) (userBorrows types.BorrowIdByOwnerAndPoolMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowByUserAndPoolKey(address, poolId)
		value = store.Get(key)
	)
	if value == nil {
		return userBorrows, false
	}
	k.cdc.MustUnmarshal(value, &userBorrows)

	return userBorrows, true
}

func (k *Keeper) SetBorrowIdByOwnerAndPool(ctx sdk.Context, userBorrows types.BorrowIdByOwnerAndPoolMapping) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowByUserAndPoolKey(userBorrows.Owner, userBorrows.PoolId)
		value = k.cdc.MustMarshal(&userBorrows)
	)
	store.Set(key, value)
}

func (k *Keeper) UpdateBorrowIdsMapping(
	ctx sdk.Context,
	borrowId uint64,
	isInsert bool,
) error {

	userVaults, found := k.GetBorrows(ctx)

	if !found && isInsert {
		userVaults = types.BorrowMapping{
			BorrowIds: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorLendOwnerNotFound
	}

	if isInsert {
		userVaults.BorrowIds = append(userVaults.BorrowIds, borrowId)
	} else {
		for index, id := range userVaults.BorrowIds {
			if id == borrowId {
				userVaults.BorrowIds = append(userVaults.BorrowIds[:index], userVaults.BorrowIds[index+1:]...)
				break
			}
		}
	}

	k.SetBorrows(ctx, userVaults)
	return nil
}

func (k *Keeper) GetBorrows(ctx sdk.Context) (userBorrows types.BorrowMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowsKey
		value = store.Get(key)
	)
	if value == nil {
		return userBorrows, false
	}
	k.cdc.MustUnmarshal(value, &userBorrows)

	return userBorrows, true
}

func (k *Keeper) SetBorrows(ctx sdk.Context, userBorrows types.BorrowMapping) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowsKey
		value = k.cdc.MustMarshal(&userBorrows)
	)
	store.Set(key, value)
}
