package keeper

import (
	"github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k *Keeper) SetReward(ctx sdk.Context, rewards types.InternalRewards) {
	var (
		store = k.Store(ctx)
		key   = types.RewardsKey(rewards.App_mapping_ID)
		value = k.cdc.MustMarshal(&rewards)
	)

	store.Set(key, value)
}

func (k *Keeper) GetReward(ctx sdk.Context, id uint64) (asset types.InternalRewards, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.RewardsKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return asset, false
	}

	k.cdc.MustUnmarshal(value, &asset)
	return asset, true
}

func (k *Keeper) GetRewards(ctx sdk.Context) (lends []types.InternalRewards) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.RewardsKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var rewards types.InternalRewards
		k.cdc.MustUnmarshal(iter.Value(), &rewards)
		lends = append(lends, rewards)
	}

	return lends
}

func (k *Keeper) SetAppId(ctx sdk.Context, AppIds types.WhitelistedAppIdsVault) {
	var (
		store = k.Store(ctx)
		key   = types.AppIdsVaultKeyPrefix
		value = k.cdc.MustMarshal(&AppIds)
	)

	store.Set(key, value)
}

func (k *Keeper) GetAppIds(ctx sdk.Context) (appIds types.WhitelistedAppIdsVault) {
	var (
		store = k.Store(ctx)
		key   = types.AppIdsVaultKeyPrefix
		value = store.Get(key)
	)

	if value == nil {
		return appIds
	}

	k.cdc.MustUnmarshal(value, &appIds)
	return appIds
}

func (k *Keeper) SetExternalRewardsLockers(ctx sdk.Context, LockerExternalRewards types.LockerExternalRewards) {
	var (
		store = k.Store(ctx)
		key   = types.ExternalRewardsLockerMappingKey(LockerExternalRewards.Id)
		value = k.cdc.MustMarshal(&LockerExternalRewards)
	)
	store.Set(key, value)
}

func (k *Keeper) GetExternalRewardsLocker(ctx sdk.Context, Id uint64) (LockerExternalRewards types.LockerExternalRewards) {
	var (
		store = k.Store(ctx)
		key   = types.ExternalRewardsLockerMappingKey(Id)
		value = store.Get(key)
	)
	if value == nil {
		return LockerExternalRewards
	}
	k.cdc.MustUnmarshal(value, &LockerExternalRewards)
	return LockerExternalRewards
}

func (k *Keeper) GetExternalRewardsLockers(ctx sdk.Context) (LockerExternalRewards []types.LockerExternalRewards) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.ExternalRewardsLockerKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var LockerExternalReward types.LockerExternalRewards
		k.cdc.MustUnmarshal(iter.Value(), &LockerExternalReward)
		LockerExternalRewards = append(LockerExternalRewards, LockerExternalReward)
	}

	return LockerExternalRewards
}

func (k *Keeper) SetExternalRewardsLockersId(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.ExtRewardsLockerIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetExternalRewardsLockersId(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.ExtRewardsLockerIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetExternalRewardsLockersCounter(ctx sdk.Context, appId uint64, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.AssetForDenomKey(appId)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) GetExternalRewardsLockersCounter(ctx sdk.Context, appId uint64) (asset uint64, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AssetForDenomKey(appId)
		value = store.Get(key)
	)

	if value == nil {
		return asset, false
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue(), true
}

func (k *Keeper) SetEpochTime(ctx sdk.Context, appId uint64, time int64) {
	var (
		store = k.Store(ctx)
		key   = types.AssetForDenomKey(appId)
		value = k.cdc.MustMarshal(
			&protobuftypes.Int64Value{
				Value: time,
			})
	)

	store.Set(key, value)
}

func (k *Keeper) GetEpochTime(ctx sdk.Context, appId uint64) (time int64, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AssetForDenomKey(appId)
		value = store.Get(key)
	)

	if value == nil {
		return time, false
	}
	var id protobuftypes.Int64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue(), true
}
