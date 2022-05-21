package expected

import (
	collecortypes "github.com/comdex-official/comdex/x/collector/types"
	"github.com/comdex-official/comdex/x/locker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type LockerKeeper interface {
	GetLockerProductAssetMapping(ctx sdk.Context, appMappingId uint64) (lockerProductMapping types.LockerProductAssetMapping, found bool)
	GetLocker(ctx sdk.Context, lockerId string) (locker types.Locker, found bool)
	GetLockerLookupTable(ctx sdk.Context, appMappingId uint64) (lockerLookupData types.LockerLookupTable, found bool)
}

type CollectorKeeper interface {
	GetAppidToAssetCollectorMapping(ctx sdk.Context, app_id uint64) (appAssetCollectorData collecortypes.AppIdToAssetCollectorMapping, found bool)
	GetCollectorLookupTable(ctx sdk.Context, app_id uint64) (collectorLookup collecortypes.CollectorLookup, found bool)
	GetAppToDenomsMapping(ctx sdk.Context, AppId uint64) (appToDenom collecortypes.AppToDenomsMapping, found bool)
	GetCollectorLookupByAsset(ctx sdk.Context, app_id, asset_id uint64) (collectorLookup collecortypes.CollectorLookup, found bool)
}
