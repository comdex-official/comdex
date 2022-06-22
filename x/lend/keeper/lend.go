package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k Keeper) GetCollateralAmount(ctx sdk.Context, borrowerAddr sdk.AccAddress, denom string) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	collateral := sdk.NewCoin(denom, sdk.ZeroInt())
	key := types.CreateCollateralAmountKey(borrowerAddr, denom)

	if bz := store.Get(key); bz != nil {
		err := collateral.Amount.Unmarshal(bz)
		if err != nil {
			panic(err)
		}
	}

	return collateral
}

func (k Keeper) setCollateralAmount(ctx sdk.Context, borrowerAddr sdk.AccAddress, collateral sdk.Coin) error {
	if !collateral.IsValid() {
		return sdkerrors.Wrap(types.ErrInvalidAsset, collateral.String())
	}

	if borrowerAddr.Empty() {
		return types.ErrEmptyAddress
	}

	bz, err := collateral.Amount.Marshal()
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	key := types.CreateCollateralAmountKey(borrowerAddr, collateral.Denom)

	if collateral.Amount.IsZero() {
		store.Delete(key)
	} else {
		store.Set(key, bz)
	}
	return nil
}

func (k *Keeper) SetUserLendIDHistory(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendHistoryIdPrefix
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k *Keeper) GetUserLendIDHistory(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LendHistoryIdPrefix
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetLend(ctx sdk.Context, lend types.LendAsset) {
	var (
		store = k.Store(ctx)
		key   = types.LendUserKey(lend.ID)
		value = k.cdc.MustMarshal(&lend)
	)

	store.Set(key, value)
}

func (k *Keeper) GetLend(ctx sdk.Context, id uint64) (lend types.LendAsset, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LendUserKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return lend, false
	}

	k.cdc.MustUnmarshal(value, &lend)
	return lend, true
}

func (k *Keeper) DeleteLend(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendUserKey(id)
	)

	store.Delete(key)
}

func (k *Keeper) SetPoolId(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.PoolIdPrefix
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k *Keeper) GetPoolId(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.PoolIdPrefix
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	var (
		store = k.Store(ctx)
		key   = types.PoolKey(pool.PoolId)
		value = k.cdc.MustMarshal(&pool)
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

	k.cdc.MustUnmarshal(value, &pool)
	return pool, true
}

func (k *Keeper) SetAssetToPair(ctx sdk.Context, assetToPair types.AssetToPairMapping) {
	var (
		store = k.Store(ctx)
		key   = types.AssetToPairMappingKey(assetToPair.AssetId)
		value = k.cdc.MustMarshal(&assetToPair)
	)

	store.Set(key, value)
}

func (k *Keeper) GetAssetToPair(ctx sdk.Context, id uint64) (assetToPair types.AssetToPairMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AssetToPairMappingKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return assetToPair, false
	}

	k.cdc.MustUnmarshal(value, &assetToPair)
	return assetToPair, true
}

func (k *Keeper) SetLendForAddressByAsset(ctx sdk.Context, address sdk.AccAddress, assetID, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendForAddressByAsset(address, assetID)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) HasLendForAddressByAsset(ctx sdk.Context, address sdk.AccAddress, assetID uint64) bool {
	var (
		store = k.Store(ctx)
		key   = types.LendForAddressByAsset(address, assetID)
	)

	return store.Has(key)
}

func (k *Keeper) DeleteLendForAddressByAsset(ctx sdk.Context, address sdk.AccAddress, assetID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendForAddressByAsset(address, assetID)
	)

	store.Delete(key)
}

func (k *Keeper) UpdateUserLendIdMapping(
	ctx sdk.Context,
	lendOwner string,
	lend types.LendAsset,
	isInsert bool,
) error {

	userVaults, found := k.GetUserLends(ctx, lendOwner)

	if !found && isInsert {
		userVaults = types.UserLendIdMapping{
			Owner: lendOwner,
			Lends: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorLendOwnerNotFound
	}

	if isInsert {
		userVaults.Lends = append(userVaults.Lends, lend)
	} else {

		for index, newLend := range userVaults.Lends {

			if newLend.ID == lend.ID {
				userVaults.Lends = append(userVaults.Lends[:index], userVaults.Lends[index+1:]...)
				break
			}
		}
	}

	k.SetUserLends(ctx, userVaults)
	return nil
}

func (k *Keeper) GetUserLends(ctx sdk.Context, address string) (userVaults types.UserLendIdMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserLendsForAddressKey(address)
		value = store.Get(key)
	)
	if value == nil {
		return userVaults, false
	}
	k.cdc.MustUnmarshal(value, &userVaults)

	return userVaults, true
}

func (k *Keeper) SetUserLends(ctx sdk.Context, userVaults types.UserLendIdMapping) {
	var (
		store = k.Store(ctx)
		key   = types.UserLendsForAddressKey(userVaults.Owner)
		value = k.cdc.MustMarshal(&userVaults)
	)
	store.Set(key, value)
}

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
	borrow types.BorrowAsset,
	isInsert bool,
) error {

	userVaults, found := k.GetUserBorrows(ctx, lendOwner)

	if !found && isInsert {
		userVaults = types.UserBorrowIdMapping{
			Owner:   lendOwner,
			Borrows: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorLendOwnerNotFound
	}

	if isInsert {
		userVaults.Borrows = append(userVaults.Borrows, borrow)
	} else {
		for index, id := range userVaults.Borrows {
			if id.ID == borrow.ID {
				userVaults.Borrows = append(userVaults.Borrows[:index], userVaults.Borrows[index+1:]...)
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

func (k *Keeper) SetLendIdToBorrowIdMapping(ctx sdk.Context, LendIdToBorrowIdMapping types.LendIdToBorrowIdMapping) {
	var (
		store = k.Store(ctx)
		key   = types.LendIdToBorrowIdMappingKey(LendIdToBorrowIdMapping.LendingID)
		value = k.cdc.MustMarshal(&LendIdToBorrowIdMapping)
	)

	store.Set(key, value)
}

func (k *Keeper) UpdateLendIdToBorrowIdMapping(
	ctx sdk.Context,
	lendId uint64,
	borrowId uint64,
	isInsert bool,
) error {

	LendIdToBorrowIdMapping, found := k.GetLendIdToBorrowIdMapping(ctx, lendId)

	if !found && isInsert {
		LendIdToBorrowIdMapping = types.LendIdToBorrowIdMapping{
			LendingID:   lendId,
			BorrowingID: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorLendOwnerNotFound
	}

	if isInsert {
		LendIdToBorrowIdMapping.BorrowingID = append(LendIdToBorrowIdMapping.BorrowingID, borrowId)
	} else {
		for index, id := range LendIdToBorrowIdMapping.BorrowingID {
			if id == borrowId {
				LendIdToBorrowIdMapping.BorrowingID = append(LendIdToBorrowIdMapping.BorrowingID[:index], LendIdToBorrowIdMapping.BorrowingID[index+1:]...)
				break
			}
		}
	}

	k.SetLendIdToBorrowIdMapping(ctx, LendIdToBorrowIdMapping)
	return nil
}

func (k *Keeper) GetLendIdToBorrowIdMapping(ctx sdk.Context, id uint64) (LendIdToBorrowIdMapping types.LendIdToBorrowIdMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LendIdToBorrowIdMappingKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return LendIdToBorrowIdMapping, false
	}

	k.cdc.MustUnmarshal(value, &LendIdToBorrowIdMapping)
	return LendIdToBorrowIdMapping, true
}

func (k *Keeper) DeleteLendIdToBorrowIdMapping(ctx sdk.Context, lendingId uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendIdToBorrowIdMappingKey(lendingId)
	)

	store.Delete(key)
}

func (k *Keeper) SetAssetStatsByPoolIdAndAssetId(ctx sdk.Context, assetID, poolId uint64, AssetStats types.AssetStats) {
	var (
		store = k.Store(ctx)
		key   = types.SetAssetStatsByPoolIdAndAssetId(assetID, poolId)
		value = k.cdc.MustMarshal(&AssetStats)
	)

	store.Set(key, value)
}

func (k *Keeper) GetAssetStatsByPoolIdAndAssetId(ctx sdk.Context, assetID, poolId uint64) (AssetStats types.AssetStats, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.SetAssetStatsByPoolIdAndAssetId(assetID, poolId)
		value = store.Get(key)
	)

	if value == nil {
		return AssetStats, false
	}

	k.cdc.MustUnmarshal(value, &AssetStats)
	return AssetStats, true
}

func (k *Keeper) SetLastInterestTime(ctx sdk.Context, interestTime int64) error {
	store := ctx.KVStore(k.storeKey)
	timeKey := types.CreateLastInterestTimeKey()

	bz, err := k.cdc.Marshal(&protobuftypes.Int64Value{Value: interestTime})
	if err != nil {
		return err
	}

	store.Set(timeKey, bz)
	return nil
}

// GetLastInterestTime gets last time at which interest was accrued.
func (k Keeper) GetLastInterestTime(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.storeKey)
	timeKey := types.CreateLastInterestTimeKey()
	bz := store.Get(timeKey)

	val := protobuftypes.Int64Value{}

	if err := k.cdc.Unmarshal(bz, &val); err != nil {
		panic(err)
	}

	return val.Value
}

func (k *Keeper) UpdateLendIdsMapping(
	ctx sdk.Context,
	lendId uint64,
	isInsert bool,
) error {

	userVaults, found := k.GetLends(ctx)

	if !found && isInsert {
		userVaults = types.LendMapping{
			LendIds: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorLendOwnerNotFound
	}

	if isInsert {
		userVaults.LendIds = append(userVaults.LendIds, lendId)
	} else {
		for index, id := range userVaults.LendIds {
			if id == lendId {
				userVaults.LendIds = append(userVaults.LendIds[:index], userVaults.LendIds[index+1:]...)
				break
			}
		}
	}

	k.SetLends(ctx, userVaults)
	return nil
}

func (k *Keeper) GetLends(ctx sdk.Context) (userVaults types.LendMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LendsKey
		value = store.Get(key)
	)
	if value == nil {
		return userVaults, false
	}
	k.cdc.MustUnmarshal(value, &userVaults)

	return userVaults, true
}

func (k *Keeper) SetLends(ctx sdk.Context, userLends types.LendMapping) {
	var (
		store = k.Store(ctx)
		key   = types.LendsKey
		value = k.cdc.MustMarshal(&userLends)
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
