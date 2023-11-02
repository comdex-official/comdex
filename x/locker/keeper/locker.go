package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/cosmos/gogoproto/types"

	esmtypes "github.com/comdex-official/comdex/x/esm/types"

	"github.com/comdex-official/comdex/x/locker/types"
)

// get locker lookup table.
func (k Keeper) SetLockerTotalRewardsByAssetAppWise(ctx sdk.Context, lockerRewardsMapping types.LockerTotalRewardsByAssetAppWise) error {
	var (
		store = k.Store(ctx)
		key   = types.LockerTotalRewardsByAssetAppWiseKey(lockerRewardsMapping.AppId, lockerRewardsMapping.AssetId)
		value = k.cdc.MustMarshal(&lockerRewardsMapping)
	)

	store.Set(key, value)
	return nil
}

func (k Keeper) GetLockerTotalRewardsByAssetAppWise(ctx sdk.Context, appID, assetID uint64) (lockerRewardsMapping types.LockerTotalRewardsByAssetAppWise, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LockerTotalRewardsByAssetAppWiseKey(appID, assetID)
		value = store.Get(key)
	)

	if value == nil {
		return lockerRewardsMapping, false
	}

	k.cdc.MustUnmarshal(value, &lockerRewardsMapping)
	return lockerRewardsMapping, true
}

func (k Keeper) GetLockerTotalRewardsByAppWise(ctx sdk.Context, appID uint64) (lockerProductMapping []types.LockerTotalRewardsByAssetAppWise, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LockerTotalRewardsByAppWiseKey(appID)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var mapData types.LockerTotalRewardsByAssetAppWise
		k.cdc.MustUnmarshal(iter.Value(), &mapData)
		lockerProductMapping = append(lockerProductMapping, mapData)
	}
	if lockerProductMapping == nil {
		return nil, false
	}
	return lockerProductMapping, true
}

func (k Keeper) GetAllLockerTotalRewardsByAssetAppWise(ctx sdk.Context) (lockerTotalRewardsByAssetAppWise []types.LockerTotalRewardsByAssetAppWise) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.LockerTotalRewardsByAssetAppWiseKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var lock types.LockerTotalRewardsByAssetAppWise
		k.cdc.MustUnmarshal(iter.Value(), &lock)
		lockerTotalRewardsByAssetAppWise = append(lockerTotalRewardsByAssetAppWise, lock)
	}
	return lockerTotalRewardsByAssetAppWise
}

func (k Keeper) SetLockerProductAssetMapping(ctx sdk.Context, lockerProductMapping types.LockerProductAssetMapping) {
	var (
		store = k.Store(ctx)
		key   = types.LockerProductAssetMappingKey(lockerProductMapping.AppId, lockerProductMapping.AssetId)
		value = k.cdc.MustMarshal(&lockerProductMapping)
	)

	store.Set(key, value)
}

func (k Keeper) GetLockerProductAssetMapping(ctx sdk.Context, appMappingID, assetID uint64) (lockerProductMapping types.LockerProductAssetMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LockerProductAssetMappingKey(appMappingID, assetID)
		value = store.Get(key)
	)

	if value == nil {
		return lockerProductMapping, false
	}

	k.cdc.MustUnmarshal(value, &lockerProductMapping)
	return lockerProductMapping, true
}

func (k Keeper) GetLockerProductAssetMappingByApp(ctx sdk.Context, appID uint64) (lockerProductMapping []types.LockerProductAssetMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LockerProductAssetMappingByAppKey(appID)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var mapData types.LockerProductAssetMapping
		k.cdc.MustUnmarshal(iter.Value(), &mapData)
		lockerProductMapping = append(lockerProductMapping, mapData)
	}
	if lockerProductMapping == nil {
		return nil, false
	}
	return lockerProductMapping, true
}

func (k Keeper) GetAllLockerProductAssetMapping(ctx sdk.Context) (lockerProductAssetMapping []types.LockerProductAssetMapping) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.LockerProductAssetMappingKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var lock types.LockerProductAssetMapping
		k.cdc.MustUnmarshal(iter.Value(), &lock)
		lockerProductAssetMapping = append(lockerProductAssetMapping, lock)
	}
	return lockerProductAssetMapping
}

func (k Keeper) SetLockerLookupTable(ctx sdk.Context, lockerLookupData types.LockerLookupTableData) {
	var (
		store = k.Store(ctx)
		key   = types.LockerLookupTableKey(lockerLookupData.AppId, lockerLookupData.AssetId)
		value = k.cdc.MustMarshal(&lockerLookupData)
	)

	store.Set(key, value)
}

func (k Keeper) GetLockerLookupTable(ctx sdk.Context, appID, assetID uint64) (lockerLookupData types.LockerLookupTableData, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LockerLookupTableKey(appID, assetID)
		value = store.Get(key)
	)

	if value == nil {
		return lockerLookupData, false
	}

	k.cdc.MustUnmarshal(value, &lockerLookupData)
	return lockerLookupData, true
}

func (k Keeper) GetLockerLookupTableByApp(ctx sdk.Context, appID uint64) (lockerLookupData []types.LockerLookupTableData, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LockerLookupTableByAppKey(appID)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var mapData types.LockerLookupTableData
		k.cdc.MustUnmarshal(iter.Value(), &mapData)
		lockerLookupData = append(lockerLookupData, mapData)
	}
	if lockerLookupData == nil {
		return nil, false
	}

	return lockerLookupData, true
}

func (k Keeper) GetAllLockerLookupTable(ctx sdk.Context) (lockerLookupTable []types.LockerLookupTableData) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.LockerLookupTableKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var lock types.LockerLookupTableData
		k.cdc.MustUnmarshal(iter.Value(), &lock)
		lockerLookupTable = append(lockerLookupTable, lock)
	}
	return lockerLookupTable
}

// UpdateAmountLockerMapping For updating token locker mapping in lookup table.
func (k Keeper) UpdateAmountLockerMapping(ctx sdk.Context, appID uint64, assetID uint64, amount sdk.Int, changeType bool) {
	// if Change type true = Add to deposits

	// If change type false = Subtract from the deposits
	lookupTableData, exists := k.GetLockerLookupTable(ctx, appID, assetID)
	if exists {
		if changeType {
			lookupTableData.DepositedAmount = lookupTableData.DepositedAmount.Add(amount)
		} else {
			lookupTableData.DepositedAmount = lookupTableData.DepositedAmount.Sub(amount)
		}

		k.SetLockerLookupTable(ctx, lookupTableData)
	}
}

// SetUserLockerAssetMapping User Locker Functions.
func (k Keeper) SetUserLockerAssetMapping(ctx sdk.Context, userLockerAssetData types.UserAppAssetLockerMapping) {
	var (
		store = k.Store(ctx)
		key   = types.UserAppAssetLockerMappingKey(userLockerAssetData.Owner, userLockerAssetData.AppId, userLockerAssetData.AssetId)
		value = k.cdc.MustMarshal(&userLockerAssetData)
	)

	store.Set(key, value)
}

func (k Keeper) GetUserLockerAssetMapping(ctx sdk.Context, address string, appID, assetID uint64) (userLockerAssetData types.UserAppAssetLockerMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserAppAssetLockerMappingKey(address, appID, assetID)
		value = store.Get(key)
	)

	if value == nil {
		return userLockerAssetData, false
	}

	k.cdc.MustUnmarshal(value, &userLockerAssetData)
	return userLockerAssetData, true
}

func (k Keeper) GetUserLockerAppMapping(ctx sdk.Context, address string, appID uint64) (userLockerAssetData []types.UserAppAssetLockerMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserAppLockerMappingKey(address, appID)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var mapData types.UserAppAssetLockerMapping
		k.cdc.MustUnmarshal(iter.Value(), &mapData)
		userLockerAssetData = append(userLockerAssetData, mapData)
	}
	if userLockerAssetData == nil {
		return nil, false
	}

	return userLockerAssetData, true
}

func (k Keeper) GetUserLockerMapping(ctx sdk.Context, address string) (userLockerAssetData []types.UserAppAssetLockerMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserLockerMappingKey(address)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var mapData types.UserAppAssetLockerMapping
		k.cdc.MustUnmarshal(iter.Value(), &mapData)
		userLockerAssetData = append(userLockerAssetData, mapData)
	}
	if userLockerAssetData == nil {
		return nil, false
	}

	return userLockerAssetData, true
}

func (k Keeper) GetAllUserLockerAssetMapping(ctx sdk.Context) (userLockerAssetMapping []types.UserAppAssetLockerMapping) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.UserLockerAssetMappingKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var lock types.UserAppAssetLockerMapping
		k.cdc.MustUnmarshal(iter.Value(), &lock)
		userLockerAssetMapping = append(userLockerAssetMapping, lock)
	}
	return userLockerAssetMapping
}

func (k Keeper) SetIDForLocker(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LockerIDPrefix
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetIDForLocker(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LockerIDPrefix
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) SetLocker(ctx sdk.Context, locker types.Locker) {
	var (
		store = k.Store(ctx)
		key   = types.LockerKey(locker.LockerId)
		value = k.cdc.MustMarshal(&locker)
	)

	store.Set(key, value)
}

func (k Keeper) DeleteLocker(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LockerKey(id)
	)

	store.Delete(key)
}

func (k Keeper) GetLocker(ctx sdk.Context, lockerID uint64) (locker types.Locker, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LockerKey(lockerID)
		value = store.Get(key)
	)

	if value == nil {
		return locker, false
	}

	k.cdc.MustUnmarshal(value, &locker)
	return locker, true
}

func (k Keeper) GetLockers(ctx sdk.Context) (locker []types.Locker) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.LockerKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var lock types.Locker
		k.cdc.MustUnmarshal(iter.Value(), &lock)
		locker = append(locker, lock)
	}
	return locker
}

func (k Keeper) WasmAddWhiteListedAssetQuery(ctx sdk.Context, appMappingID, AssetID uint64) (bool, string) {
	_, found := k.asset.GetApp(ctx, appMappingID)
	if !found {
		return false, types.ErrorAppMappingDoesNotExist.Error()
	}
	_, found = k.asset.GetAsset(ctx, AssetID)
	if !found {
		return false, types.ErrorAssetDoesNotExist.Error()
	}
	_, found1 := k.GetLockerProductAssetMapping(ctx, appMappingID, AssetID)

	if found1 {
		return false, types.ErrorLockerProductAssetMappingExists.Error()
	}
	return true, ""
}

func (k Keeper) AddWhiteListedAsset(ctx sdk.Context, msg *types.MsgAddWhiteListedAssetRequest) (*types.MsgAddWhiteListedAssetResponse, error) {
	esmStatus, found := k.esm.GetESMStatus(ctx, msg.AppId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}
	klwsParams, _ := k.esm.GetKillSwitchData(ctx, msg.AppId)
	if klwsParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	appMapping, found := k.asset.GetApp(ctx, msg.AppId)
	if !found {
		return nil, types.ErrorAppMappingDoesNotExist
	}
	asset, found := k.asset.GetAsset(ctx, msg.AssetId)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}
	_, found1 := k.GetLockerProductAssetMapping(ctx, msg.AppId, msg.AssetId)

	if !found1 {
		// Set a new instance of Locker Product Asset  Mapping

		var locker types.LockerProductAssetMapping
		locker.AppId = appMapping.Id
		locker.AssetId = msg.AssetId
		k.SetLockerProductAssetMapping(ctx, locker)

		// Also Create a LockerLookup table Instance and set it with the new asset id
		var lockerLookupData types.LockerLookupTableData

		lockerLookupData.AssetId = asset.Id
		lockerLookupData.AppId = appMapping.Id
		k.SetLockerLookupTable(ctx, lockerLookupData)
		return &types.MsgAddWhiteListedAssetResponse{}, nil
	}
	// Check if the asset from msg exists or not ,
	return nil, types.ErrorLockerProductAssetMappingExists
}
