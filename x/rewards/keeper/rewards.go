package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	"github.com/comdex-official/comdex/x/rewards/types"
)

// SetReward sets internal rewards for an asset of an app
func (k Keeper) SetReward(ctx sdk.Context, rewards types.InternalRewards) {
	var (
		store = k.Store(ctx)
		key   = types.RewardsKey(rewards.AppMappingId, rewards.AssetId)
		value = k.cdc.MustMarshal(&rewards)
	)

	store.Set(key, value)
}

// GetReward gets internal rewards for an asset of an app
func (k Keeper) GetReward(ctx sdk.Context, appID, assetID uint64) (rewards types.InternalRewards, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.RewardsKey(appID, assetID)
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

func (k Keeper) GetRewardByApp(ctx sdk.Context, appID uint64) (rewards []types.InternalRewards, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.RewardsKeyByApp(appID)
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

func (k Keeper) SetAppByAppID(ctx sdk.Context, appID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.AppIDKeyPrefix(appID)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: appID,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetAppIDByApp(ctx sdk.Context, appID uint64) (uint64, bool) {
	var (
		store = k.Store(ctx)
		key   = types.AppIDKeyPrefix(appID)
		value = store.Get(key)
	)

	if value == nil {
		return 0, false
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue(), true
}

func (k Keeper) GetAppIDs(ctx sdk.Context) (appIds []uint64) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.AppIdsVaultKeyPrefix)
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

func (k Keeper) DeleteAppIDByApp(ctx sdk.Context, appID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.AppIDKeyPrefix(appID)
	)

	store.Delete(key)
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

func (k Keeper) GetRemoveWhitelistAppIDLockerRewardsCheck(ctx sdk.Context, appMappingID uint64, assetIDs uint64) (found bool, err string) {
	_, found = k.GetRewardByApp(ctx, appMappingID)
	if !found {
		return false, "not found"
	}
	return true, ""
}

func (k Keeper) GetWhitelistAppIDVaultInterestCheck(ctx sdk.Context, appMappingID uint64) (found bool, err string) {
	_, found = k.GetAppIDByApp(ctx, appMappingID)
	if found {
		return false, "app Id already exists"
	}
	return true, ""
}

func (k Keeper) GetWhitelistAppIDLockerRewardsCheck(ctx sdk.Context, appMappingID uint64, assetID uint64) (found bool, err string) {
	_, found = k.locker.GetLockerProductAssetMapping(ctx, appMappingID, assetID)
	if !found {
		return false, "assetID does not exist"
	}
	return true, ""
}

func (k Keeper) GetExternalLockerRewardsCheck(ctx sdk.Context, appMappingID uint64, assetID uint64) (found bool, err string) {
	_, found = k.locker.GetLockerProductAssetMapping(ctx, appMappingID, assetID)
	if !found {
		return false, "asset id does not exist"
	}
	return true, ""
}

func (k Keeper) GetExternalVaultRewardsCheck(ctx sdk.Context, appMappingID uint64, assetID uint64) (found bool, err string) {
	return true, ""
}

// Whitelist an asset of an app for internal rewards
func (k Keeper) Whitelist(ctx sdk.Context, msg *types.WhitelistAsset) (*types.MsgWhitelistAssetResponse, error) {
	klwsParams, _ := k.esm.GetKillSwitchData(ctx, msg.AppMappingId)
	if klwsParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	esmStatus, found := k.esm.GetESMStatus(ctx, msg.AppMappingId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}

	if err := k.WhitelistAssetForInternalRewards(ctx, msg.AppMappingId, msg.AssetId); err != nil {
		return nil, err
	}
	return &types.MsgWhitelistAssetResponse{}, nil
}

func (k Keeper) WhitelistAppVault(ctx sdk.Context, msg *types.WhitelistAppIdVault) (*types.MsgWhitelistAppIdVaultResponse, error) {
	klwsParams, _ := k.esm.GetKillSwitchData(ctx, msg.AppMappingId)
	if klwsParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	esmStatus, found := k.esm.GetESMStatus(ctx, msg.AppMappingId)
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

func (k Keeper) DeleteLockerRewardTracker(ctx sdk.Context, rewards types.LockerRewardsTracker) {
	var (
		store = k.Store(ctx)
		key   = types.LockerRewardsTrackerKey(rewards.LockerId, rewards.AppMappingId)
	)
	store.Delete(key)
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

func (k Keeper) GetAllLockerRewardTracker(ctx sdk.Context) (Lockrewards []types.LockerRewardsTracker) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.LockerRewardsTrackerKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var rewards types.LockerRewardsTracker
		k.cdc.MustUnmarshal(iter.Value(), &rewards)
		Lockrewards = append(Lockrewards, rewards)
	}

	return Lockrewards
}

func (k Keeper) SetVaultInterestTracker(ctx sdk.Context, vault types.VaultInterestTracker) {
	var (
		store = k.Store(ctx)
		key   = types.VaultInterestTrackerKey(vault.VaultId, vault.AppMappingId)
		value = k.cdc.MustMarshal(&vault)
	)

	store.Set(key, value)
}

func (k Keeper) DeleteVaultInterestTracker(ctx sdk.Context, vault types.VaultInterestTracker) {
	var (
		store = k.Store(ctx)
		key   = types.VaultInterestTrackerKey(vault.VaultId, vault.AppMappingId)
	)
	store.Delete(key)
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

func (k Keeper) GetAllVaultInterestTracker(ctx sdk.Context) (Vaultrewards []types.VaultInterestTracker) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.VaultInterestTrackerKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var rewards types.VaultInterestTracker
		k.cdc.MustUnmarshal(iter.Value(), &rewards)
		Vaultrewards = append(Vaultrewards, rewards)
	}

	return Vaultrewards
}

func (k Keeper) CalculateLockerRewards(ctx sdk.Context, appID, assetID, lockerID uint64, Depositor string, NetBalance sdk.Int, blockHeight int64, lockerBlockTime int64) error {
	_, found := k.GetReward(ctx, appID, assetID)
	if !found {
		return nil
	}
	lockers, _ := k.locker.GetLockerLookupTable(ctx, appID, assetID)

	collectorLookup, found := k.collector.GetCollectorLookupTable(ctx, appID, assetID)
	if !found {
		return collectortypes.ErrorAssetDoesNotExist
	}
	var rewards sdk.Dec
	var err error
	collectorBTime := collectorLookup.BlockTime.Unix()
	if collectorLookup.LockerSavingRate.IsZero() {
		return nil
	}

	if blockHeight == 0 {
		// take bh from lsr
		rewards, err = k.CalculationOfRewards(ctx, NetBalance, collectorLookup.LockerSavingRate, collectorBTime)
		if err != nil {
			return err
		}
	} else {
		rewards, err = k.CalculationOfRewards(ctx, NetBalance, collectorLookup.LockerSavingRate, lockerBlockTime)
		if err != nil {
			return err
		}
	}
	lockerData, _ := k.locker.GetLocker(ctx, lockerID)
	lockerRewardsTracker, found := k.GetLockerRewardTracker(ctx, lockerData.LockerId, appID)
	if !found {
		lockerRewardsTracker = types.LockerRewardsTracker{
			LockerId:           lockerData.LockerId,
			AppMappingId:       appID,
			RewardsAccumulated: rewards,
		}
	} else {
		lockerRewardsTracker.RewardsAccumulated = lockerRewardsTracker.RewardsAccumulated.Add(rewards)
	}

	if lockerRewardsTracker.RewardsAccumulated.GTE(sdk.OneDec()) {
		// send rewards
		newReward := lockerRewardsTracker.RewardsAccumulated.TruncateInt()
		newRewardDec := sdk.NewDecFromInt(newReward)
		lockerRewardsTracker.RewardsAccumulated = lockerRewardsTracker.RewardsAccumulated.Sub(newRewardDec)
		k.SetLockerRewardTracker(ctx, lockerRewardsTracker)
		// netFeeCollectedData, found := k.collector.GetNetFeeCollectedData(ctx, appID, lockerData.AssetDepositId)
		// if !found {
		// 	return nil
		// }
		err = k.collector.DecreaseNetFeeCollectedData(ctx, appID, lockerData.AssetDepositId, newReward)
		if err != nil {
			return err
		}
		assetData, _ := k.asset.GetAsset(ctx, assetID)

		if newReward.GT(sdk.ZeroInt()) {
			err = k.bank.SendCoinsFromModuleToModule(ctx, collectortypes.ModuleName, lockertypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetData.Denom, newReward)))
			if err != nil {
				return err
			}
		}
		lockerRewardsMapping, found := k.locker.GetLockerTotalRewardsByAssetAppWise(ctx, appID, lockerData.AssetDepositId)
		if !found {
			var lockerReward lockertypes.LockerTotalRewardsByAssetAppWise
			lockerReward.AppId = appID
			lockerReward.AssetId = lockerData.AssetDepositId
			lockerReward.TotalRewards = newReward
			err = k.locker.SetLockerTotalRewardsByAssetAppWise(ctx, lockerReward)
			if err != nil {
				return err
			}
		} else {
			lockerRewardsMapping.TotalRewards = lockerRewardsMapping.TotalRewards.Add(newReward)
			err = k.locker.SetLockerTotalRewardsByAssetAppWise(ctx, lockerRewardsMapping)
			if err != nil {
				return err
			}
		}
		// updating user rewards data
		lockerData.BlockTime = ctx.BlockTime()
		lockerData.BlockHeight = ctx.BlockHeight()

		lockerData.NetBalance = lockerData.NetBalance.Add(newReward)
		lockerData.ReturnsAccumulated = lockerData.ReturnsAccumulated.Add(newReward)
		k.locker.SetLocker(ctx, lockerData)
		lockers.DepositedAmount = lockers.DepositedAmount.Add(newReward)
		k.locker.SetLockerLookupTable(ctx, lockers)
	} else {
		//	set tracker rewards
		k.SetLockerRewardTracker(ctx, lockerRewardsTracker)
		lockerData.BlockTime = ctx.BlockTime()
		lockerData.BlockHeight = ctx.BlockHeight()
		k.locker.SetLocker(ctx, lockerData)
	}

	return nil
}

func (k Keeper) CalculateVaultInterest(ctx sdk.Context, appID, extendedPairID, vaultID uint64, totalDebt sdk.Int, blockHeight int64, vaultBlockTime int64) error {
	_, found := k.GetAppIDByApp(ctx, appID)
	if !found {
		return nil
	}
	ExtPairVaultData, found := k.asset.GetPairsVault(ctx, extendedPairID)
	if !found {
		return assettypes.ErrorPairDoesNotExist
	}

	extPairVaultBTime := ExtPairVaultData.BlockTime.Unix()
	if ExtPairVaultData.StabilityFee.IsZero() || ExtPairVaultData.IsStableMintVault {
		return nil
	}

	blockTime := vaultBlockTime
	if blockHeight == 0 {
		blockTime = extPairVaultBTime
	}
	interest, err := k.CalculationOfRewards(ctx, totalDebt, ExtPairVaultData.StabilityFee, blockTime)
	if err != nil {
		return err
	}

	vaultData, _ := k.vault.GetVault(ctx, vaultID)
	vaultInterestTracker, found := k.GetVaultInterestTracker(ctx, vaultData.Id, appID)
	if !found {
		vaultInterestTracker = types.VaultInterestTracker{
			VaultId:             vaultData.Id,
			AppMappingId:        appID,
			InterestAccumulated: interest,
		}
	} else {
		vaultInterestTracker.InterestAccumulated = vaultInterestTracker.InterestAccumulated.Add(interest)
	}

	if vaultInterestTracker.InterestAccumulated.GTE(sdk.OneDec()) {
		newInterest := vaultInterestTracker.InterestAccumulated.TruncateInt()
		newInterestDec := sdk.NewDecFromInt(newInterest)
		vaultInterestTracker.InterestAccumulated = vaultInterestTracker.InterestAccumulated.Sub(newInterestDec)

		vaultData.BlockTime = ctx.BlockTime()
		vaultData.BlockHeight = ctx.BlockHeight()

		k.SetVaultInterestTracker(ctx, vaultInterestTracker)
		intAcc := vaultData.InterestAccumulated
		updatedIntAcc := (intAcc).Add(newInterest)
		vaultData.InterestAccumulated = updatedIntAcc
		k.vault.SetVault(ctx, vaultData)
	} else {
		k.SetVaultInterestTracker(ctx, vaultInterestTracker)
		vaultData.BlockTime = ctx.BlockTime()
		vaultData.BlockHeight = ctx.BlockHeight()
		k.vault.SetVault(ctx, vaultData)
	}

	return nil
}

func (k Keeper) GetExternalRewardLends(ctx sdk.Context) (LendExternalRewards []types.LendExternalRewards) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.ExternalRewardsLendKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var LendExternalReward types.LendExternalRewards
		k.cdc.MustUnmarshal(iter.Value(), &LendExternalReward)
		LendExternalRewards = append(LendExternalRewards, LendExternalReward)
	}

	return LendExternalRewards
}

func (k Keeper) SetExternalRewardLend(ctx sdk.Context, LendExternalRewards types.LendExternalRewards) {
	var (
		store = k.Store(ctx)
		key   = types.ExternalRewardsLendMappingKey(LendExternalRewards.Id)
		value = k.cdc.MustMarshal(&LendExternalRewards)
	)
	store.Set(key, value)
}

func (k Keeper) GetExternalRewardLend(ctx sdk.Context, id uint64) (LendExternalRewards types.LendExternalRewards) {
	var (
		store = k.Store(ctx)
		key   = types.ExternalRewardsLendMappingKey(id)
		value = store.Get(key)
	)
	if value == nil {
		return LendExternalRewards
	}
	k.cdc.MustUnmarshal(value, &LendExternalRewards)
	return LendExternalRewards
}

func (k Keeper) SetExternalRewardsLendID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.ExtRewardsLendIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetExternalRewardsLendID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.ExtRewardsLendIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) SetExternalRewardStableVault(ctx sdk.Context, VaultExternalRewards types.StableVaultExternalRewards) {
	var (
		store = k.Store(ctx)
		key   = types.ExternalRewardsStableVaultMappingKey(VaultExternalRewards.AppId)
		value = k.cdc.MustMarshal(&VaultExternalRewards)
	)
	store.Set(key, value)
}
