package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k *Keeper) SetUserLendIDHistory(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendHistoryIdPrefix
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k *Keeper) GetUserLendIDHistory(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LendHistoryIdPrefix
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetLend(ctx sdk.Context, lend types.LendAsset) {
	var (
		store = k.Store(ctx)
		key   = types.LendUserKey(lend.ID)
		value = k.cdc.MustMarshal(&lend)
	)

	store.Set(key, value)
}

func (k *Keeper) GetLend(ctx sdk.Context, id uint64) (lend types.LendAsset, found bool) {
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

func (k *Keeper) DeleteLend(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendUserKey(id)
	)

	store.Delete(key)
}

func (k *Keeper) SetPoolId(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.PoolIdPrefix
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k *Keeper) GetPoolId(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.PoolIdPrefix
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	var (
		store = k.Store(ctx)
		key   = types.PoolKey(pool.PoolId)
		value = k.cdc.MustMarshal(&pool)
	)

	store.Set(key, value)
}

func (k *Keeper) GetPool(ctx sdk.Context, id uint64) (pool types.Pool, found bool) {
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

func (k *Keeper) SetAssetToPair(ctx sdk.Context, assetToPair types.AssetToPairMapping) {
	var (
		store = k.Store(ctx)
		key   = types.AssetToPairMappingKey(assetToPair.AssetId, assetToPair.PoolId)
		value = k.cdc.MustMarshal(&assetToPair)
	)

	store.Set(key, value)
}

func (k *Keeper) GetAssetToPair(ctx sdk.Context, assetId, poolId uint64) (assetToPair types.AssetToPairMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AssetToPairMappingKey(assetId, poolId)
		value = store.Get(key)
	)

	if value == nil {
		return assetToPair, false
	}

	k.cdc.MustUnmarshal(value, &assetToPair)
	return assetToPair, true
}

func (k *Keeper) SetLendForAddressByAsset(ctx sdk.Context, address sdk.AccAddress, assetID, id, poolID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendForAddressByAsset(address, assetID, poolID)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) HasLendForAddressByAsset(ctx sdk.Context, address sdk.AccAddress, assetID, poolID uint64) bool {
	var (
		store = k.Store(ctx)
		key   = types.LendForAddressByAsset(address, assetID, poolID)
	)

	return store.Has(key)
}

func (k *Keeper) DeleteLendForAddressByAsset(ctx sdk.Context, address sdk.AccAddress, assetID, poolID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendForAddressByAsset(address, assetID, poolID)
	)

	store.Delete(key)
}

func (k *Keeper) UpdateUserLendIdMapping(
	ctx sdk.Context,
	lendOwner string,
	lendId uint64,
	isInsert bool,
) error {

	userVaults, found := k.GetUserLends(ctx, lendOwner)

	if !found && isInsert {
		userVaults = types.UserLendIdMapping{
			Owner:   lendOwner,
			LendIds: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorLendOwnerNotFound
	}

	if isInsert {
		userVaults.LendIds = append(userVaults.LendIds, lendId)
	} else {
		for index, id := range userVaults.LendIds {
			if id == lendId {
				userVaults.LendIds = append(userVaults.LendIds[:index], userVaults.LendIds[index+1:]...)
				break
			}
		}
	}

	k.SetUserLends(ctx, userVaults)
	return nil
}

func (k *Keeper) GetUserLends(ctx sdk.Context, address string) (userVaults types.UserLendIdMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserLendsForAddressKey(address)
		value = store.Get(key)
	)
	if value == nil {
		return userVaults, false
	}
	k.cdc.MustUnmarshal(value, &userVaults)

	return userVaults, true
}

func (k *Keeper) UserLends(ctx sdk.Context, address string) (userLends []types.LendAsset, found bool) {
	userLendId, _ := k.GetUserLends(ctx, address)
	for _, v := range userLendId.LendIds {
		userLend, _ := k.GetLend(ctx, v)
		userLends = append(userLends, userLend)
	}
	return userLends, true
}

func (k *Keeper) SetUserLends(ctx sdk.Context, userVaults types.UserLendIdMapping) {
	var (
		store = k.Store(ctx)
		key   = types.UserLendsForAddressKey(userVaults.Owner)
		value = k.cdc.MustMarshal(&userVaults)
	)
	store.Set(key, value)
}

func (k *Keeper) UpdateLendIdByOwnerAndPoolMapping(
	ctx sdk.Context,
	lendOwner string,
	lendId uint64,
	poolId uint64,
	isInsert bool,
) error {

	userLends, found := k.GetLendIdByOwnerAndPool(ctx, lendOwner, poolId)

	if !found && isInsert {
		userLends = types.LendIdByOwnerAndPoolMapping{
			Owner:   lendOwner,
			PoolId:  poolId,
			LendIds: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorLendOwnerNotFound
	}

	if isInsert {
		userLends.LendIds = append(userLends.LendIds, lendId)
	} else {
		for index, id := range userLends.LendIds {
			if id == lendId {
				userLends.LendIds = append(userLends.LendIds[:index], userLends.LendIds[index+1:]...)
				break
			}
		}
	}

	k.SetLendIdByOwnerAndPool(ctx, userLends)
	return nil
}

func (k *Keeper) GetLendIdByOwnerAndPool(ctx sdk.Context, address string, poolId uint64) (userLends types.LendIdByOwnerAndPoolMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LendByUserAndPoolKey(address, poolId)
		value = store.Get(key)
	)
	if value == nil {
		return userLends, false
	}
	k.cdc.MustUnmarshal(value, &userLends)

	return userLends, true
}

func (k *Keeper) LendIdByOwnerAndPool(ctx sdk.Context, address string, poolId uint64) (userLends []types.LendAsset, found bool) {
	userLendId, _ := k.GetLendIdByOwnerAndPool(ctx, address, poolId)
	for _, v := range userLendId.LendIds {
		userLend, _ := k.GetLend(ctx, v)
		userLends = append(userLends, userLend)
	}
	return userLends, true
}

func (k *Keeper) SetLendIdByOwnerAndPool(ctx sdk.Context, userLends types.LendIdByOwnerAndPoolMapping) {
	var (
		store = k.Store(ctx)
		key   = types.LendByUserAndPoolKey(userLends.Owner, userLends.PoolId)
		value = k.cdc.MustMarshal(&userLends)
	)
	store.Set(key, value)
}

func (k *Keeper) SetLendIdToBorrowIdMapping(ctx sdk.Context, LendIdToBorrowIdMapping types.LendIdToBorrowIdMapping) {
	var (
		store = k.Store(ctx)
		key   = types.LendIdToBorrowIdMappingKey(LendIdToBorrowIdMapping.LendingID)
		value = k.cdc.MustMarshal(&LendIdToBorrowIdMapping)
	)

	store.Set(key, value)
}

func (k *Keeper) UpdateLendIdToBorrowIdMapping(
	ctx sdk.Context,
	lendId uint64,
	borrowId uint64,
	isInsert bool,
) error {

	LendIdToBorrowIdMapping, found := k.GetLendIdToBorrowIdMapping(ctx, lendId)

	if !found && isInsert {
		LendIdToBorrowIdMapping = types.LendIdToBorrowIdMapping{
			LendingID:   lendId,
			BorrowingID: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorLendOwnerNotFound
	}

	if isInsert {
		LendIdToBorrowIdMapping.BorrowingID = append(LendIdToBorrowIdMapping.BorrowingID, borrowId)
	} else {
		for index, id := range LendIdToBorrowIdMapping.BorrowingID {
			if id == borrowId {
				LendIdToBorrowIdMapping.BorrowingID = append(LendIdToBorrowIdMapping.BorrowingID[:index], LendIdToBorrowIdMapping.BorrowingID[index+1:]...)
				break
			}
		}
	}

	k.SetLendIdToBorrowIdMapping(ctx, LendIdToBorrowIdMapping)
	return nil
}

func (k *Keeper) GetLendIdToBorrowIdMapping(ctx sdk.Context, id uint64) (LendIdToBorrowIdMapping types.LendIdToBorrowIdMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LendIdToBorrowIdMappingKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return LendIdToBorrowIdMapping, false
	}

	k.cdc.MustUnmarshal(value, &LendIdToBorrowIdMapping)
	return LendIdToBorrowIdMapping, true
}

func (k *Keeper) DeleteLendIdToBorrowIdMapping(ctx sdk.Context, lendingId uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendIdToBorrowIdMappingKey(lendingId)
	)

	store.Delete(key)
}

func (k *Keeper) SetAssetStatsByPoolIdAndAssetId(ctx sdk.Context, AssetStats types.AssetStats) {
	var (
		store = k.Store(ctx)
		key   = types.SetAssetStatsByPoolIdAndAssetId(AssetStats.AssetId, AssetStats.PoolId)
		value = k.cdc.MustMarshal(&AssetStats)
	)

	store.Set(key, value)
}

func (k *Keeper) GetAssetStatsByPoolIdAndAssetId(ctx sdk.Context, assetID, poolId uint64) (AssetStats types.AssetStats, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.SetAssetStatsByPoolIdAndAssetId(assetID, poolId)
		value = store.Get(key)
	)

	if value == nil {
		return AssetStats, false
	}

	k.cdc.MustUnmarshal(value, &AssetStats)
	return AssetStats, true
}

func (k *Keeper) AssetStatsByPoolIdAndAssetId(ctx sdk.Context, assetID, poolId uint64) (AssetStats types.AssetStats, found bool) {
	AssetStats, found = k.UpdateAPR(ctx, poolId, assetID)
	if !found {
		return AssetStats, false
	}
	return AssetStats, true
}

func (k *Keeper) SetLastInterestTime(ctx sdk.Context, interestTime int64) error {
	store := ctx.KVStore(k.storeKey)
	timeKey := types.CreateLastInterestTimeKey()

	bz, err := k.cdc.Marshal(&protobuftypes.Int64Value{Value: interestTime})
	if err != nil {
		return err
	}

	store.Set(timeKey, bz)
	return nil
}

// GetLastInterestTime gets last time at which interest was accrued.
func (k Keeper) GetLastInterestTime(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.storeKey)
	timeKey := types.CreateLastInterestTimeKey()
	bz := store.Get(timeKey)

	val := protobuftypes.Int64Value{}

	if err := k.cdc.Unmarshal(bz, &val); err != nil {
		panic(err)
	}

	return val.Value
}

func (k *Keeper) UpdateLendIdsMapping(
	ctx sdk.Context,
	lendId uint64,
	isInsert bool,
) error {

	userVaults, found := k.GetLends(ctx)

	if !found && isInsert {
		userVaults = types.LendMapping{
			LendIds: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorLendOwnerNotFound
	}

	if isInsert {
		userVaults.LendIds = append(userVaults.LendIds, lendId)
	} else {
		for index, id := range userVaults.LendIds {
			if id == lendId {
				userVaults.LendIds = append(userVaults.LendIds[:index], userVaults.LendIds[index+1:]...)
				break
			}
		}
	}

	k.SetLends(ctx, userVaults)
	return nil
}

func (k *Keeper) GetLends(ctx sdk.Context) (userVaults types.LendMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LendsKey
		value = store.Get(key)
	)
	if value == nil {
		return userVaults, false
	}
	k.cdc.MustUnmarshal(value, &userVaults)

	return userVaults, true
}

func (k *Keeper) SetLends(ctx sdk.Context, userLends types.LendMapping) {
	var (
		store = k.Store(ctx)
		key   = types.LendsKey
		value = k.cdc.MustMarshal(&userLends)
	)
	store.Set(key, value)
}

func (k *Keeper) GetModuleBalanceByPoolId(ctx sdk.Context, poolId uint64) (ModuleBalance types.ModuleBalance, found bool) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return ModuleBalance, false
	}
	for _, v := range pool.AssetData {
		asset, _ := k.GetAsset(ctx, v.AssetId)
		balance := k.ModuleBalance(ctx, pool.ModuleName, asset.Denom)
		tokenBal := sdk.NewCoin(asset.Denom, balance)
		modBalStats := types.ModuleBalanceStats{
			AssetId: asset.Id,
			Balance: tokenBal,
		}
		ModuleBalance.PoolId = poolId
		ModuleBalance.ModuleBalanceStats = append(ModuleBalance.ModuleBalanceStats, modBalStats)
	}
	return ModuleBalance, true
}
