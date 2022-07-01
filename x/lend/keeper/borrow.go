package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k *Keeper) SetUserBorrowIDHistory(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowHistoryIDPrefix
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
		key   = types.BorrowHistoryIDPrefix
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetBorrow(ctx sdk.Context, borrow types.BorrowAsset) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowUserKey(borrow.ID)
		value = k.cdc.MustMarshal(&borrow)
	)

	store.Set(key, value)
}

func (k *Keeper) GetBorrow(ctx sdk.Context, id uint64) (borrow types.BorrowAsset, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowUserKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return borrow, false
	}

	k.cdc.MustUnmarshal(value, &borrow)
	return borrow, true
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

func (k *Keeper) UpdateUserBorrowIDMapping(
	ctx sdk.Context,
	lendOwner string,
	borrowID uint64,
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
		userVaults.BorrowIds = append(userVaults.BorrowIds, borrowID)
	} else {
		for index, id := range userVaults.BorrowIds {
			if id == borrowID {
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

func (k *Keeper) UserBorrows(ctx sdk.Context, address string) (userBorrows []types.BorrowAsset, found bool) {
	userLendID, _ := k.GetUserBorrows(ctx, address)
	for _, v := range userLendID.BorrowIds {
		userBorrow, _ := k.GetBorrow(ctx, v)
		userBorrows = append(userBorrows, userBorrow)
	}
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

func (k *Keeper) UpdateBorrowIDByOwnerAndPoolMapping(
	ctx sdk.Context,
	borrowOwner string,
	borrowID uint64,
	poolID uint64,
	isInsert bool,
) error {

	userLends, found := k.GetBorrowIDByOwnerAndPool(ctx, borrowOwner, poolID)

	if !found && isInsert {
		userLends = types.BorrowIdByOwnerAndPoolMapping{
			Owner:     borrowOwner,
			PoolId:    poolID,
			BorrowIds: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorLendOwnerNotFound
	}

	if isInsert {
		userLends.BorrowIds = append(userLends.BorrowIds, borrowID)
	} else {
		for index, id := range userLends.BorrowIds {
			if id == borrowID {
				userLends.BorrowIds = append(userLends.BorrowIds[:index], userLends.BorrowIds[index+1:]...)
				break
			}
		}
	}

	k.SetBorrowIDByOwnerAndPool(ctx, userLends)
	return nil
}

func (k *Keeper) GetBorrowIDByOwnerAndPool(ctx sdk.Context, address string, poolID uint64) (userBorrows types.BorrowIdByOwnerAndPoolMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowByUserAndPoolKey(address, poolID)
		value = store.Get(key)
	)
	if value == nil {
		return userBorrows, false
	}
	k.cdc.MustUnmarshal(value, &userBorrows)

	return userBorrows, true
}

func (k *Keeper) BorrowIDByOwnerAndPool(ctx sdk.Context, address string, poolID uint64) (userBorrows []types.BorrowAsset, found bool) {
	userLendId, _ := k.GetBorrowIDByOwnerAndPool(ctx, address, poolID)
	for _, v := range userLendId.BorrowIds {
		userBorrow, _ := k.GetBorrow(ctx, v)
		userBorrows = append(userBorrows, userBorrow)
	}
	return userBorrows, true
}

func (k *Keeper) SetBorrowIDByOwnerAndPool(ctx sdk.Context, userBorrows types.BorrowIdByOwnerAndPoolMapping) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowByUserAndPoolKey(userBorrows.Owner, userBorrows.PoolId)
		value = k.cdc.MustMarshal(&userBorrows)
	)
	store.Set(key, value)
}

func (k *Keeper) UpdateBorrowIdsMapping(
	ctx sdk.Context,
	borrowID uint64,
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
		userVaults.BorrowIds = append(userVaults.BorrowIds, borrowID)
	} else {
		for index, id := range userVaults.BorrowIds {
			if id == borrowID {
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

func (k *Keeper) UpdateStableBorrowIdsMapping(
	ctx sdk.Context,
	borrowID uint64,
	isInsert bool,
) error {

	userVaults, found := k.GetStableBorrows(ctx)

	if !found && isInsert {
		userVaults = types.StableBorrowMapping{
			StableBorrowIds: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorLendOwnerNotFound
	}

	if isInsert {
		userVaults.StableBorrowIds = append(userVaults.StableBorrowIds, borrowID)
	} else {
		for index, id := range userVaults.StableBorrowIds {
			if id == borrowID {
				userVaults.StableBorrowIds = append(userVaults.StableBorrowIds[:index], userVaults.StableBorrowIds[index+1:]...)
				break
			}
		}
	}

	k.SetStableBorrows(ctx, userVaults)
	return nil
}

func (k *Keeper) GetStableBorrows(ctx sdk.Context) (userBorrows types.StableBorrowMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.StableBorrowsKey
		value = store.Get(key)
	)
	if value == nil {
		return userBorrows, false
	}
	k.cdc.MustUnmarshal(value, &userBorrows)

	return userBorrows, true
}

func (k *Keeper) SetStableBorrows(ctx sdk.Context, userBorrows types.StableBorrowMapping) {
	var (
		store = k.Store(ctx)
		key   = types.StableBorrowsKey
		value = k.cdc.MustMarshal(&userBorrows)
	)
	store.Set(key, value)
}
