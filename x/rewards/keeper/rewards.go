package keeper

import (
	"context"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	"github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k Keeper) SetReward(ctx sdk.Context, rewards types.InternalRewards) {
	var (
		store = k.Store(ctx)
		key   = types.RewardsKey(rewards.App_mapping_ID, rewards.Asset_ID)
		value = k.cdc.MustMarshal(&rewards)
	)

	store.Set(key, value)
}

func (k Keeper) GetReward(ctx sdk.Context, appId, assetID uint64) (rewards types.InternalRewards, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.RewardsKey(appId, assetID)
		value = store.Get(key)
	)

	if value == nil {
		return rewards, false
	}

	k.cdc.MustUnmarshal(value, &rewards)
	return rewards, true
}

func (k Keeper) DeleteReward(ctx sdk.Context, appID, assetID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.RewardsKey(appID, assetID)
	)
	store.Delete(key)
}

func (k Keeper) GetRewardByApp(ctx sdk.Context, appId uint64) (rewards []types.InternalRewards, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.RewardsKeyByApp(appId)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var mapData types.InternalRewards
		k.cdc.MustUnmarshal(iter.Value(), &mapData)
		rewards = append(rewards, mapData)
	}
	if rewards == nil {
		return nil, false
	}
	return rewards, true
}

func (k Keeper) GetRewards(ctx sdk.Context) (lends []types.InternalRewards) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.RewardsKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var rewards types.InternalRewards
		k.cdc.MustUnmarshal(iter.Value(), &rewards)
		lends = append(lends, rewards)
	}

	return lends
}

func (k Keeper) SetAppID(ctx sdk.Context, AppIds types.WhitelistedAppIdsVault) {
	var (
		store = k.Store(ctx)
		key   = types.AppIdsVaultKeyPrefix
		value = k.cdc.MustMarshal(&AppIds)
	)

	store.Set(key, value)
}

func (k Keeper) GetAppIDs(ctx sdk.Context) (appIds types.WhitelistedAppIdsVault) {
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

func (k Keeper) SetExternalRewardsLockers(ctx sdk.Context, LockerExternalRewards types.LockerExternalRewards) {
	var (
		store = k.Store(ctx)
		key   = types.ExternalRewardsLockerMappingKey(LockerExternalRewards.Id)
		value = k.cdc.MustMarshal(&LockerExternalRewards)
	)
	store.Set(key, value)
}

func (k Keeper) GetExternalRewardsLocker(ctx sdk.Context, id uint64) (LockerExternalRewards types.LockerExternalRewards) {
	var (
		store = k.Store(ctx)
		key   = types.ExternalRewardsLockerMappingKey(id)
		value = store.Get(key)
	)
	if value == nil {
		return LockerExternalRewards
	}
	k.cdc.MustUnmarshal(value, &LockerExternalRewards)
	return LockerExternalRewards
}

func (k Keeper) GetExternalRewardsLockers(ctx sdk.Context) (LockerExternalRewards []types.LockerExternalRewards) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.ExternalRewardsLockerKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var LockerExternalReward types.LockerExternalRewards
		k.cdc.MustUnmarshal(iter.Value(), &LockerExternalReward)
		LockerExternalRewards = append(LockerExternalRewards, LockerExternalReward)
	}

	return LockerExternalRewards
}

func (k Keeper) SetExternalRewardsLockersID(ctx sdk.Context, id uint64) {
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

func (k Keeper) GetExternalRewardsLockersID(ctx sdk.Context) uint64 {
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

func (k Keeper) SetEpochTimeID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.EpochTimeIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetEpochTimeID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.EpochTimeIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) SetEpochTime(ctx sdk.Context, epoch types.EpochTime) {
	var (
		store = k.Store(ctx)
		key   = types.EpochForLockerKey(epoch.Id)
		value = k.cdc.MustMarshal(&epoch)
	)

	store.Set(key, value)
}

func (k Keeper) GetEpochTime(ctx sdk.Context, id uint64) (epoch types.EpochTime, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.EpochForLockerKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return epoch, false
	}
	k.cdc.MustUnmarshal(value, &epoch)

	return epoch, true
}

func (k Keeper) SetExternalRewardsVaultID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.ExtRewardsVaultIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetExternalRewardsVaultID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.ExtRewardsVaultIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) GetExternalRewardVaults(ctx sdk.Context) (VaultExternalRewards []types.VaultExternalRewards) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.ExternalRewardsVaultKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var VaultExternalReward types.VaultExternalRewards
		k.cdc.MustUnmarshal(iter.Value(), &VaultExternalReward)
		VaultExternalRewards = append(VaultExternalRewards, VaultExternalReward)
	}

	return VaultExternalRewards
}

func (k Keeper) SetExternalRewardVault(ctx sdk.Context, VaultExternalRewards types.VaultExternalRewards) {
	var (
		store = k.Store(ctx)
		key   = types.ExternalRewardsVaultMappingKey(VaultExternalRewards.Id)
		value = k.cdc.MustMarshal(&VaultExternalRewards)
	)
	store.Set(key, value)
}

// Wasm query checks

func (k Keeper) GetRemoveWhitelistAppIDLockerRewardsCheck(ctx sdk.Context, appMappingID uint64, assetIDs []uint64) (found bool, err string) {
	_, found = k.GetRewardByApp(ctx, appMappingID)
	if !found {
		return false, "not found"
	}
	return true, ""
}

func (k Keeper) GetWhitelistAppIDVaultInterestCheck(ctx sdk.Context, appMappingID uint64) (found bool, err string) {
	found = uint64InSlice(appMappingID, k.GetAppIDs(ctx).WhitelistedAppMappingIdsVaults)
	if found {
		return false, "app Id already exists"
	}
	return true, ""
}

func (k Keeper) GetWhitelistAppIDLockerRewardsCheck(ctx sdk.Context, appMappingID uint64, assetIDs []uint64) (found bool, err string) {
	lockerAssets, _ := k.locker.GetLockerProductAssetMapping(ctx, appMappingID)
	for i := range assetIDs {
		found := uint64InSlice(assetIDs[i], lockerAssets.AssetIds)
		if !found {
			return false, "assetID does not exist"
		}
	}
	return true, ""
}

func (k Keeper) GetExternalLockerRewardsCheck(ctx sdk.Context, appMappingID uint64, assetID uint64) (found bool, err string) {
	lockerAssets, found := k.GetLockerProductAssetMapping(ctx, appMappingID)
	if !found {
		return false, "locker product asset mapping doesnt exist"
	}

	found = uint64InSlice(assetID, lockerAssets.AssetIds)
	if !found {
		return false, "asset id does not exist"
	}
	return true, ""
}

func (k Keeper) GetExternalVaultRewardsCheck(ctx sdk.Context, appMappingID uint64, assetID uint64) (found bool, err string) {
	return true, ""
}

func (k Keeper) Whitelist(goCtx context.Context, msg *types.WhitelistAsset) (*types.MsgWhitelistAssetResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	klwsParams, _ := k.GetKillSwitchData(ctx, msg.AppMappingId)
	if klwsParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	esmStatus, found := k.GetESMStatus(ctx, msg.AppMappingId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}

	if err := k.WhitelistAsset(ctx, msg.AppMappingId, msg.AssetId, true); err != nil {
		return nil, err
	}
	return &types.MsgWhitelistAssetResponse{}, nil
}

func (k Keeper) WhitelistAppVault(goCtx context.Context, msg *types.WhitelistAppIdVault) (*types.MsgWhitelistAppIdVaultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	klwsParams, _ := k.GetKillSwitchData(ctx, msg.AppMappingId)
	if klwsParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	esmStatus, found := k.GetESMStatus(ctx, msg.AppMappingId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}
	if err := k.WhitelistAppIDVault(ctx, msg.AppMappingId); err != nil {
		return nil, err
	}
	return &types.MsgWhitelistAppIdVaultResponse{}, nil
}

func (k Keeper) SetLockerRewardTracker(ctx sdk.Context, rewards types.LockerRewardsTracker) {
	var (
		store = k.Store(ctx)
		key   = types.LockerRewardsTrackerKey(rewards.LockerId, rewards.AppMappingId)
		value = k.cdc.MustMarshal(&rewards)
	)

	store.Set(key, value)
}

func (k Keeper) GetLockerRewardTracker(ctx sdk.Context, id, appID uint64) (rewards types.LockerRewardsTracker, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LockerRewardsTrackerKey(id, appID)
		value = store.Get(key)
	)

	if value == nil {
		return rewards, false
	}

	k.cdc.MustUnmarshal(value, &rewards)
	return rewards, true
}

func (k Keeper) SetVaultInterestTracker(ctx sdk.Context, vault types.VaultInterestTracker) {
	var (
		store = k.Store(ctx)
		key   = types.LockerRewardsTrackerKey(vault.VaultId, vault.AppMappingId)
		value = k.cdc.MustMarshal(&vault)
	)

	store.Set(key, value)
}

func (k Keeper) GetVaultInterestTracker(ctx sdk.Context, id, appID uint64) (vault types.VaultInterestTracker, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.VaultInterestTrackerKey(id, appID)
		value = store.Get(key)
	)

	if value == nil {
		return vault, false
	}

	k.cdc.MustUnmarshal(value, &vault)
	return vault, true
}

func (k Keeper) CalculateLockerRewards(ctx sdk.Context, appID, assetID uint64, Depositor string, NetBalance sdk.Int, blockHeight int64, lockerBlockTime int64 ) (sdk.Dec, error) {

	_, found := k.GetReward(ctx, appID, assetID)
	if !found {
		return sdk.ZeroDec(), nil
	}

	collectorLookup, found := k.GetCollectorLookupTable(ctx, appID, assetID)
	if !found {
		return sdk.ZeroDec(), collectortypes.ErrorAssetDoesNotExist
	}
	collectorBTime := collectorLookup.BlockTime.Unix()
	if collectorLookup.LockerSavingRate.IsZero() {
		return sdk.ZeroDec(), nil
	} else {
		if blockHeight == 0 {
			// take bh from lsr
			newAmt, err := k.LockerRewards(ctx, NetBalance, collectorLookup.LockerSavingRate, collectorBTime)
			if err != nil {
				return sdk.ZeroDec(), err
			}
			return newAmt, nil
		} else {
			newAmt, err := k.LockerRewards(ctx, NetBalance, collectorLookup.LockerSavingRate, lockerBlockTime)
			if err != nil {
				return sdk.ZeroDec(), err
			}
			return newAmt, nil
		}
	}

	return sdk.ZeroDec(), nil
}