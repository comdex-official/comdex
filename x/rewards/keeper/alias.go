package keeper

import (
	collecortypes "github.com/comdex-official/comdex/x/collector/types"
	"github.com/comdex-official/comdex/x/locker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetLockerProductAssetMapping(ctx sdk.Context, appMappingId uint64) (lockerProductMapping types.LockerProductAssetMapping, found bool) {
	return k.locker.GetLockerProductAssetMapping(ctx, appMappingId)
}

func (k Keeper) GetAppidToAssetCollectorMapping(ctx sdk.Context, app_id uint64) (appAssetCollectorData collecortypes.AppIdToAssetCollectorMapping, found bool) {
	return k.collector.GetAppidToAssetCollectorMapping(ctx, app_id)
}

func (k Keeper) GetCollectorLookupTable(ctx sdk.Context, app_id uint64) (collectorLookup collecortypes.CollectorLookup, found bool) {
	return k.collector.GetCollectorLookupTable(ctx, app_id)
}

func (k Keeper) GetAppToDenomsMapping(ctx sdk.Context, AppId uint64) (appToDenom collecortypes.AppToDenomsMapping, found bool) {
	return k.collector.GetAppToDenomsMapping(ctx, AppId)
}

func (k Keeper) GetLocker(ctx sdk.Context, lockerId string) (locker types.Locker, found bool) {
	return k.locker.GetLocker(ctx, lockerId)
}

func (k Keeper) GetLockerLookupTable(ctx sdk.Context, appMappingId uint64) (lockerLookupData types.LockerLookupTable, found bool) {
	return k.locker.GetLockerLookupTable(ctx, appMappingId)
}

func (k Keeper) GetCollectorLookupByAsset(ctx sdk.Context, app_id, asset_id uint64) (collectorLookup collecortypes.CollectorLookup, found bool) {
	return k.collector.GetCollectorLookupByAsset(ctx, app_id, asset_id)
}
