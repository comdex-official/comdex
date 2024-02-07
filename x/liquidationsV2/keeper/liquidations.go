package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/cosmos/gogoproto/types"
)

// locked vaults kvs

func (k Keeper) SetLockedVaultID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k Keeper) GetLockedVaultID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) SetLockedVault(ctx sdk.Context, lockedVault types.LockedVault) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKey(lockedVault.AppId, lockedVault.LockedVaultId)
		value = k.cdc.MustMarshal(&lockedVault)
	)
	store.Set(key, value)
}

func (k Keeper) DeleteLockedVault(ctx sdk.Context, appID, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKey(appID, id)
	)
	store.Delete(key)
}

func (k Keeper) GetLockedVault(ctx sdk.Context, appID, id uint64) (lockedVault types.LockedVault, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKey(appID, id)
		value = store.Get(key)
	)

	if value == nil {
		return lockedVault, false
	}

	k.cdc.MustUnmarshal(value, &lockedVault)
	return lockedVault, true
}

func (k Keeper) GetLockedVaultByApp(ctx sdk.Context, appID uint64) (lockedVault []types.LockedVault) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKeyByApp(appID)
		iter  = storetypes.KVStorePrefixIterator(store, key)
	)

	defer func(iter storetypes.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var mapData types.LockedVault
		k.cdc.MustUnmarshal(iter.Value(), &mapData)
		lockedVault = append(lockedVault, mapData)
	}
	return lockedVault
}

func (k Keeper) GetLockedVaults(ctx sdk.Context) (lockedVaults []types.LockedVault) {
	var (
		store = k.Store(ctx)
		iter  = storetypes.KVStorePrefixIterator(store, types.LockedVaultKeyPrefix)
	)

	defer func(iter storetypes.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var lockedVault types.LockedVault
		k.cdc.MustUnmarshal(iter.Value(), &lockedVault)
		lockedVaults = append(lockedVaults, lockedVault)
	}

	return lockedVaults
}
