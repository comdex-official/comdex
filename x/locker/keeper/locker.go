package keeper

import (

	// "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// "github.com/tendermint/tendermint/libs/log"

	// "github.com/comdex-official/comdex/x/locker/expected"
	"github.com/comdex-official/comdex/x/locker/types"
)

//get locker lookup table

func (k *Keeper) SetLockerProductAssetMapping(ctx sdk.Context, lockerProductMapping types.LockerProductAssetMapping) {

	var (
		store = k.Store(ctx)
		key   = types.LockerProductAssetMappingKey(lockerProductMapping.AppMappingId)
		value = k.cdc.MustMarshal(&lockerProductMapping)
	)

	store.Set(key, value)

}

func (k *Keeper) GetLockerProductAssetMapping(ctx sdk.Context, appMappingId uint64) (lockerProductMapping types.LockerProductAssetMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LockerProductAssetMappingKey(appMappingId)
		value = store.Get(key)
	)

	if value == nil {
		return lockerProductMapping, false
	}

	k.cdc.MustUnmarshal(value, &lockerProductMapping)
	return lockerProductMapping, true
}

func (k *Keeper) SetLockerLookupTable(ctx sdk.Context, lockerLookupData types.LockerLookupTable) {

	var (
		store = k.Store(ctx)
		key   = types.LockerLookupTableKey(lockerLookupData.AppMappingId)
		value = k.cdc.MustMarshal(&lockerLookupData)
	)

	store.Set(key, value)

}

func (k *Keeper) GetLockerLookupTable(ctx sdk.Context, appMappingId uint64) (lockerLookupData types.LockerLookupTable, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LockerLookupTableKey(appMappingId)
		value = store.Get(key)
	)

	if value == nil {
		return lockerLookupData, false
	}

	k.cdc.MustUnmarshal(value, &lockerLookupData)
	return lockerLookupData, true
}

func (k *Keeper) CheckLockerProductAssetMapping(ctx sdk.Context, assetId uint64, lockerProductMapping types.LockerProductAssetMapping) (found bool) {

	for _, id := range lockerProductMapping.AssetIds {

		if id == assetId {
			return true
		} else {
			continue
		}
	}
	return false

}

//For updating token locker mappping in lookup table
func (k *Keeper) UpdateTokenLockerMapping(ctx sdk.Context, lockerLookupData types.LockerLookupTable, counter uint64, userLockerData types.Locker) {

	for _, lockerData := range lockerLookupData.Lockers {
		if lockerData.AssetId == userLockerData.AssetDepositId {

			lockerData.DepositedAmount = lockerData.DepositedAmount.Add(userLockerData.NetBalance)
			lockerData.LockerIds = append(lockerData.LockerIds, userLockerData.LockerId)

		}

	}
	lockerLookupData.Counter = counter
	k.SetLockerLookupTable(ctx, lockerLookupData)

}

//For updating token locker mappping in lookup table
func (k *Keeper) UpdateAmountLockerMapping(ctx sdk.Context, lockerLookupData types.LockerLookupTable, assetId uint64, amount sdk.Int, changeType bool) {

	//if Change type true = Add to deposits
	//If change type false = Substract from the deposits

	for _, lockerData := range lockerLookupData.Lockers {
		if lockerData.AssetId == assetId {
			if changeType {
				lockerData.DepositedAmount = lockerData.DepositedAmount.Add(amount)
			} else {
				lockerData.DepositedAmount = lockerData.DepositedAmount.Sub(amount)
			}

		}

	}
	k.SetLockerLookupTable(ctx, lockerLookupData)

}

//User Locker Functions:
func (k *Keeper) SetUserLockerAssetMapping(ctx sdk.Context, userLockerAssetData types.UserLockerAssetMapping) {

	var (
		store = k.Store(ctx)
		key   = types.UserLockerAssetMappingKey(userLockerAssetData.Owner)
		value = k.cdc.MustMarshal(&userLockerAssetData)
	)

	store.Set(key, value)

}

func (k *Keeper) GetUserLockerAssetMapping(ctx sdk.Context, address string) (userLockerAssetData types.UserLockerAssetMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserLockerAssetMappingKey(address)
		value = store.Get(key)
	)

	if value == nil {
		return userLockerAssetData, false
	}

	k.cdc.MustUnmarshal(value, &userLockerAssetData)
	return userLockerAssetData, true
}

//Checking if for a certain user for the app type , whether there exists a certain asset or not and if it contains a locker id or not
func (k *Keeper) CheckUserAppToAssetMapping(ctx sdk.Context, userLockerAssetData types.UserLockerAssetMapping, assetId uint64, appMappingId uint64) (locker_id string, found bool) {

	for _, locker_app_mapping := range userLockerAssetData.LockerAppMapping {

		if locker_app_mapping.AppMappingId == appMappingId {
			for _, asset_to_lockerid_mapping := range locker_app_mapping.UserAssetLocker {

				if asset_to_lockerid_mapping.AssetId == assetId && len(asset_to_lockerid_mapping.LockerId) > 0 {

					locker_id = asset_to_lockerid_mapping.LockerId
					return locker_id, true

				}

			}

		}

	}
	return locker_id, false

}

func (k *Keeper) SetLocker(ctx sdk.Context, locker types.Locker) {

	var (
		store = k.Store(ctx)
		key   = types.LockerKey(locker.LockerId)
		value = k.cdc.MustMarshal(&locker)
	)

	store.Set(key, value)

}

func (k *Keeper) GetLocker(ctx sdk.Context, lockerId string) (locker types.Locker, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LockerKey(lockerId)
		value = store.Get(key)
	)

	if value == nil {
		return locker, false
	}

	k.cdc.MustUnmarshal(value, &locker)
	return locker, true
}

//Target
//user sends create request
//it comdes to the function and check if user data exists or not. if not create locker
//if user data exists- check app mapping , from app mapping check asset id . if it does then fail tx.
// else user locker id  exists use that to create this struct and set it.
