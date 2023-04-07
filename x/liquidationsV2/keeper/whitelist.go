package keeper

import (
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

// whitelisted appIds kvs

func (k Keeper) SetAppIDForLiquidation(ctx sdk.Context, appID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.WhitelistAppKeyByApp(appID)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: appID,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetAppIDByAppForLiquidation(ctx sdk.Context, appID uint64) (uint64, bool) {
	var (
		store = k.Store(ctx)
		key   = types.WhitelistAppKeyByApp(appID)
		value = store.Get(key)
	)

	if value == nil {
		return 0, false
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue(), true
}

func (k Keeper) GetAppIdsForLiquidation(ctx sdk.Context) (appIds []uint64) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.AppIdsKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var app protobuftypes.UInt64Value
		k.cdc.MustUnmarshal(iter.Value(), &app)
		appIds = append(appIds, app.Value)
	}
	return appIds
}

func (k Keeper) DeleteAppID(ctx sdk.Context, appID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.WhitelistAppKeyByApp(appID)
	)

	store.Delete(key)
}
