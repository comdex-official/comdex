package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/lend/types"
)

func (k Keeper) SetUserLendIDCounter(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendCounterIDPrefix
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k Keeper) GetUserLendIDCounter(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LendCounterIDPrefix
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) SetLend(ctx sdk.Context, lend types.LendAsset) {
	var (
		store = k.Store(ctx)
		key   = types.LendUserKey(lend.ID)
		value = k.cdc.MustMarshal(&lend)
	)

	store.Set(key, value)
}

func (k Keeper) GetLend(ctx sdk.Context, id uint64) (lend types.LendAsset, found bool) {
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

func (k Keeper) GetAllLend(ctx sdk.Context) (lendAsset []types.LendAsset) {
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
		var asset types.LendAsset
		k.cdc.MustUnmarshal(iter.Value(), &asset)
		lendAsset = append(lendAsset, asset)
	}
	return lendAsset
}

func (k Keeper) DeleteLend(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendUserKey(id)
	)

	store.Delete(key)
}

func (k Keeper) SetPoolID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.PoolIDPrefix
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k Keeper) GetPoolID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.PoolIDPrefix
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	var (
		store = k.Store(ctx)
		key   = types.PoolKey(pool.PoolID)
		value = k.cdc.MustMarshal(&pool)
	)

	store.Set(key, value)
}

func (k Keeper) GetPool(ctx sdk.Context, id uint64) (pool types.Pool, found bool) {
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

func (k Keeper) GetPools(ctx sdk.Context) (pools []types.Pool) {
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
		var pool types.Pool
		k.cdc.MustUnmarshal(iter.Value(), &pool)
		pools = append(pools, pool)
	}

	return pools
}

func (k Keeper) SetAssetToPair(ctx sdk.Context, assetToPair types.AssetToPairMapping) {
	var (
		store = k.Store(ctx)
		key   = types.AssetToPairMappingKey(assetToPair.AssetID, assetToPair.PoolID)
		value = k.cdc.MustMarshal(&assetToPair)
	)

	store.Set(key, value)
}

func (k Keeper) GetAssetToPair(ctx sdk.Context, assetID, poolID uint64) (assetToPair types.AssetToPairMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AssetToPairMappingKey(assetID, poolID)
		value = store.Get(key)
	)

	if value == nil {
		return assetToPair, false
	}

	k.cdc.MustUnmarshal(value, &assetToPair)
	return assetToPair, true
}

func (k Keeper) GetAllAssetToPair(ctx sdk.Context) (assetToPairMapping []types.AssetToPairMapping) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.AssetToPairMappingKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var asset types.AssetToPairMapping
		k.cdc.MustUnmarshal(iter.Value(), &asset)
		assetToPairMapping = append(assetToPairMapping, asset)
	}
	return assetToPairMapping
}

func (k Keeper) SetLendForAddressByAsset(ctx sdk.Context, address sdk.AccAddress, assetID, id, poolID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendForAddressByAsset(address, assetID, poolID)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) HasLendForAddressByAsset(ctx sdk.Context, address sdk.AccAddress, assetID, poolID uint64) bool {
	var (
		store = k.Store(ctx)
		key   = types.LendForAddressByAsset(address, assetID, poolID)
	)

	return store.Has(key)
}

func (k Keeper) DeleteLendForAddressByAsset(ctx sdk.Context, address sdk.AccAddress, assetID, poolID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendForAddressByAsset(address, assetID, poolID)
	)

	store.Delete(key)
}

// func (k Keeper) UpdateUserLendIDMapping(
// 	ctx sdk.Context,
// 	lendOwner string,
// 	lendID uint64,
// 	isInsert bool,
// ) error {
// 	userVaults, found := k.GetUserLends(ctx, lendOwner)

// 	if !found && isInsert {
// 		userVaults = types.UserLendIdMapping{
// 			Owner:   lendOwner,
// 			LendIDs: nil,
// 		}
// 	} else if !found && !isInsert {
// 		return types.ErrorLendOwnerNotFound
// 	}

// 	if isInsert {
// 		userVaults.LendIDs = append(userVaults.LendIDs, lendID)
// 	} else {
// 		for index, id := range userVaults.LendIDs {
// 			if id == lendID {
// 				userVaults.LendIDs = append(userVaults.LendIDs[:index], userVaults.LendIDs[index+1:]...)
// 				break
// 			}
// 		}
// 	}

// 	k.SetUserLends(ctx, userVaults)
// 	return nil
// }

// func (k Keeper) GetUserLends(ctx sdk.Context, address string) (userVaults types.UserLendIdMapping, found bool) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.UserLendsForAddressKey(address)
// 		value = store.Get(key)
// 	)
// 	if value == nil {
// 		return userVaults, false
// 	}
// 	k.cdc.MustUnmarshal(value, &userVaults)

// 	return userVaults, true
// }

// func (k Keeper) GetAllUserLends(ctx sdk.Context) (userLendIDMapping []types.UserLendIdMapping) {
// 	var (
// 		store = k.Store(ctx)
// 		iter  = sdk.KVStorePrefixIterator(store, types.UserLendsForAddressKeyPrefix)
// 	)

// 	defer func(iter sdk.Iterator) {
// 		err := iter.Close()
// 		if err != nil {
// 			return
// 		}
// 	}(iter)

// 	for ; iter.Valid(); iter.Next() {
// 		var asset types.UserLendIdMapping
// 		k.cdc.MustUnmarshal(iter.Value(), &asset)
// 		userLendIDMapping = append(userLendIDMapping, asset)
// 	}
// 	return userLendIDMapping
// }

// func (k Keeper) UserLends(ctx sdk.Context, address string) (userLends []types.LendAsset, found bool) {
// 	userLendID, _ := k.GetUserLends(ctx, address)
// 	for _, v := range userLendID.LendIDs {
// 		userLend, _ := k.GetLend(ctx, v)
// 		userLends = append(userLends, userLend)
// 	}
// 	return userLends, true
// }

// func (k Keeper) SetUserLends(ctx sdk.Context, userVaults types.UserLendIdMapping) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.UserLendsForAddressKey(userVaults.Owner)
// 		value = k.cdc.MustMarshal(&userVaults)
// 	)
// 	store.Set(key, value)
// }

// func (k Keeper) UpdateLendIDByOwnerAndPoolMapping(
// 	ctx sdk.Context,
// 	lendOwner string,
// 	lendID uint64,
// 	poolID uint64,
// 	isInsert bool,
// ) error {
// 	userLends, found := k.GetLendIDByOwnerAndPool(ctx, lendOwner, poolID)

// 	if !found && isInsert {
// 		userLends = types.LendIdByOwnerAndPoolMapping{
// 			Owner:   lendOwner,
// 			PoolID:  poolID,
// 			LendIDs: nil,
// 		}
// 	} else if !found && !isInsert {
// 		return types.ErrorLendOwnerNotFound
// 	}

// 	if isInsert {
// 		userLends.LendIDs = append(userLends.LendIDs, lendID)
// 	} else {
// 		for index, id := range userLends.LendIDs {
// 			if id == lendID {
// 				userLends.LendIDs = append(userLends.LendIDs[:index], userLends.LendIDs[index+1:]...)
// 				break
// 			}
// 		}
// 	}

// 	k.SetLendIDByOwnerAndPool(ctx, userLends)
// 	return nil
// }

// func (k Keeper) GetLendIDByOwnerAndPool(ctx sdk.Context, address string, poolID uint64) (userLends types.LendIdByOwnerAndPoolMapping, found bool) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.LendByUserAndPoolKey(address, poolID)
// 		value = store.Get(key)
// 	)
// 	if value == nil {
// 		return userLends, false
// 	}
// 	k.cdc.MustUnmarshal(value, &userLends)

// 	return userLends, true
// }

// func (k Keeper) GetAllLendIDByOwnerAndPool(ctx sdk.Context) (lendIDByOwnerAndPoolMapping []types.LendIdByOwnerAndPoolMapping) {
// 	var (
// 		store = k.Store(ctx)
// 		iter  = sdk.KVStorePrefixIterator(store, types.LendByUserAndPoolPrefix)
// 	)

// 	defer func(iter sdk.Iterator) {
// 		err := iter.Close()
// 		if err != nil {
// 			return
// 		}
// 	}(iter)

// 	for ; iter.Valid(); iter.Next() {
// 		var asset types.LendIdByOwnerAndPoolMapping
// 		k.cdc.MustUnmarshal(iter.Value(), &asset)
// 		lendIDByOwnerAndPoolMapping = append(lendIDByOwnerAndPoolMapping, asset)
// 	}
// 	return lendIDByOwnerAndPoolMapping
// }

// func (k Keeper) LendIDByOwnerAndPool(ctx sdk.Context, address string, poolID uint64) (userLends []types.LendAsset, found bool) {
// 	userLendID, _ := k.GetLendIDByOwnerAndPool(ctx, address, poolID)
// 	for _, v := range userLendID.LendIDs {
// 		userLend, _ := k.GetLend(ctx, v)
// 		userLends = append(userLends, userLend)
// 	}
// 	return userLends, true
// }

// func (k Keeper) SetLendIDByOwnerAndPool(ctx sdk.Context, userLends types.LendIdByOwnerAndPoolMapping) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.LendByUserAndPoolKey(userLends.Owner, userLends.PoolID)
// 		value = k.cdc.MustMarshal(&userLends)
// 	)
// 	store.Set(key, value)
// }

// func (k Keeper) SetLendIDToBorrowIDMapping(ctx sdk.Context, lendIDToBorrowIDMapping types.LendIdToBorrowIdMapping) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.LendIDToBorrowIDMappingKey(lendIDToBorrowIDMapping.LendingID)
// 		value = k.cdc.MustMarshal(&lendIDToBorrowIDMapping)
// 	)

// 	store.Set(key, value)
// }

// func (k Keeper) UpdateLendIDToBorrowIDMapping(
// 	ctx sdk.Context,
// 	lendID uint64,
// 	borrowID uint64,
// 	isInsert bool,
// ) error {
// 	lendIDToBorrowIDMapping, found := k.GetLendIDToBorrowIDMapping(ctx, lendID)

// 	if !found && isInsert {
// 		lendIDToBorrowIDMapping = types.LendIdToBorrowIdMapping{
// 			LendingID:   lendID,
// 			BorrowingID: nil,
// 		}
// 	} else if !found && !isInsert {
// 		return types.ErrorLendOwnerNotFound
// 	}

// 	if isInsert {
// 		lendIDToBorrowIDMapping.BorrowingID = append(lendIDToBorrowIDMapping.BorrowingID, borrowID)
// 	} else {
// 		for index, id := range lendIDToBorrowIDMapping.BorrowingID {
// 			if id == borrowID {
// 				lendIDToBorrowIDMapping.BorrowingID = append(lendIDToBorrowIDMapping.BorrowingID[:index], lendIDToBorrowIDMapping.BorrowingID[index+1:]...)
// 				break
// 			}
// 		}
// 	}

// 	k.SetLendIDToBorrowIDMapping(ctx, lendIDToBorrowIDMapping)
// 	return nil
// }

// func (k Keeper) GetLendIDToBorrowIDMapping(ctx sdk.Context, id uint64) (lendIDToBorrowIDMapping types.LendIdToBorrowIdMapping, found bool) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.LendIDToBorrowIDMappingKey(id)
// 		value = store.Get(key)
// 	)

// 	if value == nil {
// 		return lendIDToBorrowIDMapping, false
// 	}

// 	k.cdc.MustUnmarshal(value, &lendIDToBorrowIDMapping)
// 	return lendIDToBorrowIDMapping, true
// }

// func (k Keeper) GetAllLendIDToBorrowIDMapping(ctx sdk.Context) (lendIDToBorrowIdMapping []types.LendIdToBorrowIdMapping) {
// 	var (
// 		store = k.Store(ctx)
// 		iter  = sdk.KVStorePrefixIterator(store, types.LendIDToBorrowIDMappingKeyPrefix)
// 	)

// 	defer func(iter sdk.Iterator) {
// 		err := iter.Close()
// 		if err != nil {
// 			return
// 		}
// 	}(iter)

// 	for ; iter.Valid(); iter.Next() {
// 		var asset types.LendIdToBorrowIdMapping
// 		k.cdc.MustUnmarshal(iter.Value(), &asset)
// 		lendIDToBorrowIdMapping = append(lendIDToBorrowIdMapping, asset)
// 	}
// 	return lendIDToBorrowIdMapping
// }

// func (k Keeper) DeleteLendIDToBorrowIDMapping(ctx sdk.Context, lendingID uint64) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.LendIDToBorrowIDMappingKey(lendingID)
// 	)

// 	store.Delete(key)
// }

func (k Keeper) SetAssetStatsByPoolIDAndAssetID(ctx sdk.Context, PoolAssetLBMapping types.PoolAssetLBMapping) {
	var (
		store = k.Store(ctx)
		key   = types.SetAssetStatsByPoolIDAndAssetID(PoolAssetLBMapping.PoolID, PoolAssetLBMapping.AssetID)
		value = k.cdc.MustMarshal(&PoolAssetLBMapping)
	)

	store.Set(key, value)
}

func (k Keeper) GetAssetStatsByPoolIDAndAssetID(ctx sdk.Context, poolID, assetID uint64) (PoolAssetLBMapping types.PoolAssetLBMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.SetAssetStatsByPoolIDAndAssetID(poolID, assetID)
		value = store.Get(key)
	)

	if value == nil {
		return PoolAssetLBMapping, false
	}

	k.cdc.MustUnmarshal(value, &PoolAssetLBMapping)
	return PoolAssetLBMapping, true
}

func (k Keeper) GetAllAssetStatsByPoolIDAndAssetID(ctx sdk.Context) (assetStats []types.PoolAssetLBMapping) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.AssetStatsByPoolIDAndAssetIDKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var asset types.PoolAssetLBMapping
		k.cdc.MustUnmarshal(iter.Value(), &asset)
		assetStats = append(assetStats, asset)
	}
	return assetStats
}

func (k Keeper) AssetStatsByPoolIDAndAssetID(ctx sdk.Context, poolID, assetID uint64) (PoolAssetLBMapping types.PoolAssetLBMapping, found bool) {
	PoolAssetLBMapping, found = k.UpdateAPR(ctx, poolID, assetID)
	if !found {
		return PoolAssetLBMapping, false
	}
	return PoolAssetLBMapping, true
}

// func (k Keeper) UpdateLendIDsMapping(
// 	ctx sdk.Context,
// 	lendID uint64,
// 	isInsert bool,
// ) error {
// 	userVaults, found := k.GetLends(ctx)

// 	if !found && isInsert {
// 		userVaults = types.LendMapping{
// 			LendIDs: nil,
// 		}
// 	} else if !found && !isInsert {
// 		return types.ErrorLendOwnerNotFound
// 	}

// 	if isInsert {
// 		userVaults.LendIDs = append(userVaults.LendIDs, lendID)
// 	} else {
// 		for index, id := range userVaults.LendIDs {
// 			if id == lendID {
// 				userVaults.LendIDs = append(userVaults.LendIDs[:index], userVaults.LendIDs[index+1:]...)
// 				break
// 			}
// 		}
// 	}

// 	k.SetLends(ctx, userVaults)
// 	return nil
// }

// func (k Keeper) GetLends(ctx sdk.Context) (userVaults types.LendMapping, found bool) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.LendsKey
// 		value = store.Get(key)
// 	)
// 	if value == nil {
// 		return userVaults, false
// 	}
// 	k.cdc.MustUnmarshal(value, &userVaults)

// 	return userVaults, true
// }

// func (k Keeper) SetLends(ctx sdk.Context, userLends types.LendMapping) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.LendsKey
// 		value = k.cdc.MustMarshal(&userLends)
// 	)
// 	store.Set(key, value)
// }

// func (k Keeper) GetModuleBalanceByPoolID(ctx sdk.Context, poolID uint64) (ModuleBalance types.ModuleBalance, found bool) {
// 	pool, found := k.GetPool(ctx, poolID)
// 	if !found {
// 		return ModuleBalance, false
// 	}
// 	for _, v := range pool.AssetData {
// 		asset, _ := k.GetAsset(ctx, v.AssetID)
// 		balance := k.ModuleBalance(ctx, pool.ModuleName, asset.Denom)
// 		tokenBal := sdk.NewCoin(asset.Denom, balance)
// 		modBalStats := types.ModuleBalanceStats{
// 			AssetID: asset.Id,
// 			Balance: tokenBal,
// 		}
// 		ModuleBalance.PoolID = poolID
// 		ModuleBalance.ModuleBalanceStats = append(ModuleBalance.ModuleBalanceStats, &modBalStats)
// 	}
// 	return ModuleBalance, true
// }

// func (k Keeper) SetUserDepositStats(ctx sdk.Context, depositStats types.DepositStats) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.UserDepositStatsPrefix
// 		value = k.cdc.MustMarshal(&depositStats)
// 	)

// 	store.Set(key, value)
// }

// func (k Keeper) GetUserDepositStats(ctx sdk.Context) (depositStats types.DepositStats, found bool) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.UserDepositStatsPrefix
// 		value = store.Get(key)
// 	)

// 	if value == nil {
// 		return depositStats, false
// 	}

// 	k.cdc.MustUnmarshal(value, &depositStats)
// 	return depositStats, true
// }

// func (k Keeper) SetReserveDepositStats(ctx sdk.Context, depositStats types.DepositStats) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.ReserveDepositStatsPrefix
// 		value = k.cdc.MustMarshal(&depositStats)
// 	)

// 	store.Set(key, value)
// }

// func (k Keeper) GetReserveDepositStats(ctx sdk.Context) (depositStats types.DepositStats, found bool) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.ReserveDepositStatsPrefix
// 		value = store.Get(key)
// 	)

// 	if value == nil {
// 		return depositStats, false
// 	}

// 	k.cdc.MustUnmarshal(value, &depositStats)
// 	return depositStats, true
// }

// func (k Keeper) SetBuyBackDepositStats(ctx sdk.Context, depositStats types.DepositStats) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.BuyBackDepositStatsPrefix
// 		value = k.cdc.MustMarshal(&depositStats)
// 	)

// 	store.Set(key, value)
// }

// func (k Keeper) GetBuyBackDepositStats(ctx sdk.Context) (depositStats types.DepositStats, found bool) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.BuyBackDepositStatsPrefix
// 		value = store.Get(key)
// 	)

// 	if value == nil {
// 		return depositStats, false
// 	}

// 	k.cdc.MustUnmarshal(value, &depositStats)
// 	return depositStats, true
// }

// func (k Keeper) SetBorrowStats(ctx sdk.Context, borrowStats types.DepositStats) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.BorrowStatsPrefix
// 		value = k.cdc.MustMarshal(&borrowStats)
// 	)

// 	store.Set(key, value)
// }

// func (k Keeper) GetBorrowStats(ctx sdk.Context) (borrowStats types.DepositStats, found bool) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.BorrowStatsPrefix
// 		value = store.Get(key)
// 	)

// 	if value == nil {
// 		return borrowStats, false
// 	}

// 	k.cdc.MustUnmarshal(value, &borrowStats)
// 	return borrowStats, true
// }

func (k Keeper) SetLendRewardTracker(ctx sdk.Context, rewards types.LendRewardsTracker) {
	var (
		store = k.Store(ctx)
		key   = types.LendRewardsTrackerKey(rewards.LendingId)
		value = k.cdc.MustMarshal(&rewards)
	)

	store.Set(key, value)
}

func (k Keeper) GetLendRewardTracker(ctx sdk.Context, id uint64) (rewards types.LendRewardsTracker, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LendRewardsTrackerKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return rewards, false
	}

	k.cdc.MustUnmarshal(value, &rewards)
	return rewards, true
}

// only called while borrowing
func (k Keeper) SetUserLendBorrowMapping(ctx sdk.Context, userMapping types.UserAssetLendBorrowMapping) {
	var (
		store = k.Store(ctx)
		key   = types.UserLendBorrowMappingKey(userMapping.Owner, userMapping.LendId)
		value = k.cdc.MustMarshal(&userMapping)
	)

	store.Set(key, value)
}

func (k Keeper) GetUserLendBorrowMapping(ctx sdk.Context, owner string, lendID uint64) (userMapping types.UserAssetLendBorrowMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserLendBorrowMappingKey(userMapping.Owner, userMapping.LendId)
		value = store.Get(key)
	)

	if value == nil {
		return userMapping, false
	}

	k.cdc.MustUnmarshal(value, &userMapping)
	return userMapping, true
}

func (k Keeper) GetUserTotalMappingData(ctx sdk.Context, address string) (mappingData []types.UserAssetLendBorrowMapping) {
	var (
		store = k.Store(ctx)
		key   = types.UserLendBorrowKey(address)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var mapData types.UserAssetLendBorrowMapping
		k.cdc.MustUnmarshal(iter.Value(), &mapData)
		mappingData = append(mappingData, mapData)
	}

	return mappingData
}
