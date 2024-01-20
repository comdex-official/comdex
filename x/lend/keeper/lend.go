package keeper

import (
	"sort"

	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/cosmos/gogoproto/types"

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
		iter  = storetypes.KVStorePrefixIterator(store, types.LendUserPrefix)
	)

	defer func(iter storetypes.Iterator) {
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

func (k Keeper) DeletePool(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.PoolKey(id)
	)

	store.Delete(key)
}

func (k Keeper) GetPools(ctx sdk.Context) (pools []types.Pool) {
	var (
		store = k.Store(ctx)
		iter  = storetypes.KVStorePrefixIterator(store, types.PoolKeyPrefix)
	)

	defer func(iter storetypes.Iterator) {
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
		iter  = storetypes.KVStorePrefixIterator(store, types.AssetToPairMappingKeyPrefix)
	)

	defer func(iter storetypes.Iterator) {
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

func (k Keeper) GetAllLendRewardTracker(ctx sdk.Context) (rewards []types.LendRewardsTracker) {
	var (
		store = k.Store(ctx)
		iter  = storetypes.KVStorePrefixIterator(store, types.LendRewardsTrackerKeyPrefix)
	)

	defer func(iter storetypes.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var tracker types.LendRewardsTracker
		k.cdc.MustUnmarshal(iter.Value(), &tracker)
		rewards = append(rewards, tracker)
	}
	return rewards
}

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
		key   = types.UserLendBorrowMappingKey(owner, lendID)
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
		iter  = storetypes.KVStorePrefixIterator(store, key)
	)

	defer func(iter storetypes.Iterator) {
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

func (k Keeper) GetAllUserTotalMappingData(ctx sdk.Context) (mappingData []types.UserAssetLendBorrowMapping) {
	var (
		store = k.Store(ctx)
		key   = types.UserLendBorrowMappingKeyPrefix
		iter  = storetypes.KVStorePrefixIterator(store, key)
	)

	defer func(iter storetypes.Iterator) {
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

func (k Keeper) HasLendForAddressByAsset(ctx sdk.Context, address string, assetID, poolID uint64) bool {
	mappingData := k.GetUserTotalMappingData(ctx, address)
	for _, data := range mappingData {
		if data.PoolId == poolID {
			lend, _ := k.GetLend(ctx, data.LendId)
			if lend.AssetID == assetID {
				return true
			}
		}
	}
	return false
}

func (k Keeper) DeleteLendForAddressByAsset(ctx sdk.Context, address string, lendingID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.UserLendBorrowMappingKey(address, lendingID)
	)

	store.Delete(key)
}

func (k Keeper) DeleteIDFromAssetStatsMapping(ctx sdk.Context, poolID, assetID, id uint64, typeOfID bool) {
	poolLBMappingData, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, poolID, assetID)
	if typeOfID {
		lengthOfIDs := len(poolLBMappingData.LendIds)

		dataIndex := sort.Search(lengthOfIDs, func(i int) bool { return poolLBMappingData.LendIds[i] >= id })

		if dataIndex < lengthOfIDs && poolLBMappingData.LendIds[dataIndex] == id {
			poolLBMappingData.LendIds = append(poolLBMappingData.LendIds[:dataIndex], poolLBMappingData.LendIds[dataIndex+1:]...)
			k.SetAssetStatsByPoolIDAndAssetID(ctx, poolLBMappingData)
		}
	} else {
		lengthOfIDs := len(poolLBMappingData.BorrowIds)

		dataIndex := sort.Search(lengthOfIDs, func(i int) bool { return poolLBMappingData.BorrowIds[i] >= id })

		if dataIndex < lengthOfIDs && poolLBMappingData.BorrowIds[dataIndex] == id {
			poolLBMappingData.BorrowIds = append(poolLBMappingData.BorrowIds[:dataIndex], poolLBMappingData.BorrowIds[dataIndex+1:]...)
			k.SetAssetStatsByPoolIDAndAssetID(ctx, poolLBMappingData)
		}
	}
}

func (k Keeper) SetReserveBuybackAssetData(ctx sdk.Context, reserve types.ReserveBuybackAssetData) {
	var (
		store = k.Store(ctx)
		key   = types.ReserveBuybackAssetDataKey(reserve.AssetID)
		value = k.cdc.MustMarshal(&reserve)
	)

	store.Set(key, value)
}

func (k Keeper) GetReserveBuybackAssetData(ctx sdk.Context, id uint64) (reserve types.ReserveBuybackAssetData, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.ReserveBuybackAssetDataKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return reserve, false
	}

	k.cdc.MustUnmarshal(value, &reserve)
	return reserve, true
}

func (k Keeper) GetAllReserveBuybackAssetData(ctx sdk.Context) (reserve []types.ReserveBuybackAssetData) {
	var (
		store = k.Store(ctx)
		iter  = storetypes.KVStorePrefixIterator(store, types.ReserveBuybackAssetDataKeyPrefix)
	)

	defer func(iter storetypes.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var tracker types.ReserveBuybackAssetData
		k.cdc.MustUnmarshal(iter.Value(), &tracker)
		reserve = append(reserve, tracker)
	}
	return reserve
}

func (k Keeper) DeleteBorrowIDFromUserMapping(ctx sdk.Context, owner string, lendID, borrowID uint64) {
	userData, _ := k.GetUserLendBorrowMapping(ctx, owner, lendID)
	lengthOfIDs := len(userData.BorrowId)

	dataIndex := sort.Search(lengthOfIDs, func(i int) bool { return userData.BorrowId[i] >= borrowID })

	if dataIndex < lengthOfIDs && userData.BorrowId[dataIndex] == borrowID {
		userData.BorrowId = append(userData.BorrowId[:dataIndex], userData.BorrowId[dataIndex+1:]...)
		k.SetUserLendBorrowMapping(ctx, userData)
	}
}

func (k Keeper) WasmHasBorrowForAddressAndAsset(ctx sdk.Context, assetID uint64, address string) bool {
	mappingData := k.GetUserTotalMappingData(ctx, address)
	for _, data := range mappingData {
		lend, _ := k.GetLend(ctx, data.LendId)
		if lend.AssetID == assetID {
			return true
		}
	}
	return false
}

func (k Keeper) GetModuleBalanceByPoolID(ctx sdk.Context, poolID uint64) (ModuleBalance types.ModuleBalance, found bool) {
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return ModuleBalance, false
	}
	for _, v := range pool.AssetData {
		asset, _ := k.Asset.GetAsset(ctx, v.AssetID)
		balance := k.ModuleBalance(ctx, pool.ModuleName, asset.Denom)
		tokenBal := sdk.NewCoin(asset.Denom, balance)
		modBalStats := types.ModuleBalanceStats{
			AssetID: asset.Id,
			Balance: tokenBal,
		}
		ModuleBalance.PoolID = poolID
		ModuleBalance.ModuleBalanceStats = append(ModuleBalance.ModuleBalanceStats, modBalStats)
	}
	return ModuleBalance, true
}

func (k Keeper) SetFundModBal(ctx sdk.Context, modBal types.ModBal) {
	var (
		store = k.Store(ctx)
		key   = types.KeyFundModBal
		value = k.cdc.MustMarshal(&modBal)
	)

	store.Set(key, value)
}

func (k Keeper) GetFundModBal(ctx sdk.Context) (modBal types.ModBal, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.KeyFundModBal
		value = store.Get(key)
	)

	if value == nil {
		return modBal, false
	}

	k.cdc.MustUnmarshal(value, &modBal)
	return modBal, true
}

func (k Keeper) GetAllFundModBal(ctx sdk.Context) (modBal types.ModBal) {
	var (
		store = k.Store(ctx)
		iter  = storetypes.KVStorePrefixIterator(store, types.KeyFundModBal)
	)

	defer func(iter storetypes.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	k.cdc.MustUnmarshal(iter.Value(), &modBal)
	return modBal
}

func (k Keeper) SetFundReserveBal(ctx sdk.Context, resBal types.ReserveBal) {
	var (
		store = k.Store(ctx)
		key   = types.KeyFundReserveBal
		value = k.cdc.MustMarshal(&resBal)
	)

	store.Set(key, value)
}

func (k Keeper) GetFundReserveBal(ctx sdk.Context) (resBal types.ReserveBal, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.KeyFundReserveBal
		value = store.Get(key)
	)

	if value == nil {
		return resBal, false
	}

	k.cdc.MustUnmarshal(value, &resBal)
	return resBal, true
}

func (k Keeper) GetAllFundReserveBal(ctx sdk.Context) (resBal types.ReserveBal) {
	var (
		store = k.Store(ctx)
		iter  = storetypes.KVStorePrefixIterator(store, types.KeyFundReserveBal)
	)

	defer func(iter storetypes.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	k.cdc.MustUnmarshal(iter.Value(), &resBal)
	return resBal
}

func (k Keeper) SetAllReserveStatsByAssetID(ctx sdk.Context, allReserveStats types.AllReserveStats) {
	var (
		store = k.Store(ctx)
		key   = types.AllReserveStatsKey(allReserveStats.AssetID)
		value = k.cdc.MustMarshal(&allReserveStats)
	)

	store.Set(key, value)
}

func (k Keeper) GetAllReserveStatsByAssetID(ctx sdk.Context, id uint64) (allReserveStats types.AllReserveStats, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AllReserveStatsKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return allReserveStats, false
	}

	k.cdc.MustUnmarshal(value, &allReserveStats)
	return allReserveStats, true
}

func (k Keeper) GetTotalReserveStatsByAssetID(ctx sdk.Context) (reserve []types.AllReserveStats) {
	var (
		store = k.Store(ctx)
		iter  = storetypes.KVStorePrefixIterator(store, types.AllReserveStatsPrefix)
	)

	defer func(iter storetypes.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var tracker types.AllReserveStats
		k.cdc.MustUnmarshal(iter.Value(), &tracker)
		reserve = append(reserve, tracker)
	}
	return reserve
}

func (k Keeper) SetFundModBalByAssetPool(ctx sdk.Context, assetID, poolID uint64, amt sdk.Coin) {
	var (
		store = k.Store(ctx)
		key   = types.FundModBalanceKey(assetID, poolID)
		value = k.cdc.MustMarshal(&amt)
	)

	store.Set(key, value)
}

func (k Keeper) GetFundModBalByAssetPool(ctx sdk.Context, assetID, poolID uint64) (amt sdk.Coin, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.FundModBalanceKey(assetID, poolID)
		value = store.Get(key)
	)

	if value == nil {
		return amt, false
	}

	k.cdc.MustUnmarshal(value, &amt)
	return amt, true
}
