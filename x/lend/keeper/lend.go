package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k *Keeper) SetUserLendIDHistory(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendHistoryIDPrefix
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
		key   = types.LendHistoryIDPrefix
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

func (k *Keeper) SetPoolID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.PoolIDPrefix
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k *Keeper) GetPoolID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.PoolIDPrefix
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
		key   = types.PoolKey(pool.PoolID)
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
		key   = types.AssetToPairMappingKey(assetToPair.AssetID, assetToPair.PoolID)
		value = k.cdc.MustMarshal(&assetToPair)
	)

	store.Set(key, value)
}

func (k *Keeper) GetAssetToPair(ctx sdk.Context, assetID, poolID uint64) (assetToPair types.AssetToPairMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AssetToPairMappingKey(assetID, poolID)
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

func (k *Keeper) UpdateUserLendIDMapping(
	ctx sdk.Context,
	lendOwner string,
	lendID uint64,
	isInsert bool,
) error {
	userVaults, found := k.GetUserLends(ctx, lendOwner)

	if !found && isInsert {
		userVaults = types.UserLendIdMapping{
			Owner:   lendOwner,
			LendIDs: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorLendOwnerNotFound
	}

	if isInsert {
		userVaults.LendIDs = append(userVaults.LendIDs, lendID)
	} else {
		for index, id := range userVaults.LendIDs {
			if id == lendID {
				userVaults.LendIDs = append(userVaults.LendIDs[:index], userVaults.LendIDs[index+1:]...)
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
	userLendID, _ := k.GetUserLends(ctx, address)
	for _, v := range userLendID.LendIDs {
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

func (k *Keeper) UpdateLendIDByOwnerAndPoolMapping(
	ctx sdk.Context,
	lendOwner string,
	lendID uint64,
	poolID uint64,
	isInsert bool,
) error {
	userLends, found := k.GetLendIDByOwnerAndPool(ctx, lendOwner, poolID)

	if !found && isInsert {
		userLends = types.LendIdByOwnerAndPoolMapping{
			Owner:   lendOwner,
			PoolID:  poolID,
			LendIDs: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorLendOwnerNotFound
	}

	if isInsert {
		userLends.LendIDs = append(userLends.LendIDs, lendID)
	} else {
		for index, id := range userLends.LendIDs {
			if id == lendID {
				userLends.LendIDs = append(userLends.LendIDs[:index], userLends.LendIDs[index+1:]...)
				break
			}
		}
	}

	k.SetLendIDByOwnerAndPool(ctx, userLends)
	return nil
}

func (k *Keeper) GetLendIDByOwnerAndPool(ctx sdk.Context, address string, poolID uint64) (userLends types.LendIdByOwnerAndPoolMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LendByUserAndPoolKey(address, poolID)
		value = store.Get(key)
	)
	if value == nil {
		return userLends, false
	}
	k.cdc.MustUnmarshal(value, &userLends)

	return userLends, true
}

func (k *Keeper) LendIDByOwnerAndPool(ctx sdk.Context, address string, poolID uint64) (userLends []types.LendAsset, found bool) {
	userLendID, _ := k.GetLendIDByOwnerAndPool(ctx, address, poolID)
	for _, v := range userLendID.LendIDs {
		userLend, _ := k.GetLend(ctx, v)
		userLends = append(userLends, userLend)
	}
	return userLends, true
}

func (k *Keeper) SetLendIDByOwnerAndPool(ctx sdk.Context, userLends types.LendIdByOwnerAndPoolMapping) {
	var (
		store = k.Store(ctx)
		key   = types.LendByUserAndPoolKey(userLends.Owner, userLends.PoolID)
		value = k.cdc.MustMarshal(&userLends)
	)
	store.Set(key, value)
}

func (k *Keeper) SetLendIDToBorrowIDMapping(ctx sdk.Context, lendIDToBorrowIDMapping types.LendIdToBorrowIdMapping) {
	var (
		store = k.Store(ctx)
		key   = types.LendIDToBorrowIDMappingKey(lendIDToBorrowIDMapping.LendingID)
		value = k.cdc.MustMarshal(&lendIDToBorrowIDMapping)
	)

	store.Set(key, value)
}

func (k *Keeper) UpdateLendIDToBorrowIDMapping(
	ctx sdk.Context,
	lendID uint64,
	borrowID uint64,
	isInsert bool,
) error {
	lendIDToBorrowIDMapping, found := k.GetLendIDToBorrowIDMapping(ctx, lendID)

	if !found && isInsert {
		lendIDToBorrowIDMapping = types.LendIdToBorrowIdMapping{
			LendingID:   lendID,
			BorrowingID: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorLendOwnerNotFound
	}

	if isInsert {
		lendIDToBorrowIDMapping.BorrowingID = append(lendIDToBorrowIDMapping.BorrowingID, borrowID)
	} else {
		for index, id := range lendIDToBorrowIDMapping.BorrowingID {
			if id == borrowID {
				lendIDToBorrowIDMapping.BorrowingID = append(lendIDToBorrowIDMapping.BorrowingID[:index], lendIDToBorrowIDMapping.BorrowingID[index+1:]...)
				break
			}
		}
	}

	k.SetLendIDToBorrowIDMapping(ctx, lendIDToBorrowIDMapping)
	return nil
}

func (k *Keeper) GetLendIDToBorrowIDMapping(ctx sdk.Context, id uint64) (lendIDToBorrowIDMapping types.LendIdToBorrowIdMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LendIDToBorrowIDMappingKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return lendIDToBorrowIDMapping, false
	}

	k.cdc.MustUnmarshal(value, &lendIDToBorrowIDMapping)
	return lendIDToBorrowIDMapping, true
}

func (k *Keeper) DeleteLendIDToBorrowIDMapping(ctx sdk.Context, lendingID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendIDToBorrowIDMappingKey(lendingID)
	)

	store.Delete(key)
}

func (k *Keeper) SetAssetStatsByPoolIDAndAssetID(ctx sdk.Context, AssetStats types.AssetStats) {
	var (
		store = k.Store(ctx)
		key   = types.SetAssetStatsByPoolIDAndAssetID(AssetStats.AssetID, AssetStats.PoolID)
		value = k.cdc.MustMarshal(&AssetStats)
	)

	store.Set(key, value)
}

func (k *Keeper) GetAssetStatsByPoolIDAndAssetID(ctx sdk.Context, assetID, poolID uint64) (AssetStats types.AssetStats, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.SetAssetStatsByPoolIDAndAssetID(assetID, poolID)
		value = store.Get(key)
	)

	if value == nil {
		return AssetStats, false
	}

	k.cdc.MustUnmarshal(value, &AssetStats)
	return AssetStats, true
}

func (k *Keeper) AssetStatsByPoolIDAndAssetID(ctx sdk.Context, assetID, poolID uint64) (AssetStats types.AssetStats, found bool) {
	AssetStats, found = k.UpdateAPR(ctx, poolID, assetID)
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

func (k *Keeper) UpdateLendIDsMapping(
	ctx sdk.Context,
	lendID uint64,
	isInsert bool,
) error {
	userVaults, found := k.GetLends(ctx)

	if !found && isInsert {
		userVaults = types.LendMapping{
			LendIDs: nil,
		}
	} else if !found && !isInsert {
		return types.ErrorLendOwnerNotFound
	}

	if isInsert {
		userVaults.LendIDs = append(userVaults.LendIDs, lendID)
	} else {
		for index, id := range userVaults.LendIDs {
			if id == lendID {
				userVaults.LendIDs = append(userVaults.LendIDs[:index], userVaults.LendIDs[index+1:]...)
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

func (k *Keeper) GetModuleBalanceByPoolID(ctx sdk.Context, poolID uint64) (ModuleBalance types.ModuleBalance, found bool) {
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return ModuleBalance, false
	}
	for _, v := range pool.AssetData {
		asset, _ := k.GetAsset(ctx, v.AssetID)
		balance := k.ModuleBalance(ctx, pool.ModuleName, asset.Denom)
		tokenBal := sdk.NewCoin(asset.Denom, balance)
		modBalStats := types.ModuleBalanceStats{
			AssetID: asset.Id,
			Balance: tokenBal,
		}
		ModuleBalance.PoolID = poolID
		ModuleBalance.ModuleBalanceStats = append(ModuleBalance.ModuleBalanceStats, modBalStats)
	}
	return ModuleBalance, true
}

func (k *Keeper) SetUserDepositStats(ctx sdk.Context, depositStats types.DepositStats) {
	var (
		store = k.Store(ctx)
		key   = types.UserDepositStatsPrefix
		value = k.cdc.MustMarshal(&depositStats)
	)

	store.Set(key, value)
}

func (k *Keeper) GetUserDepositStats(ctx sdk.Context) (depositStats types.DepositStats, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.UserDepositStatsPrefix
		value = store.Get(key)
	)

	if value == nil {
		return depositStats, false
	}

	k.cdc.MustUnmarshal(value, &depositStats)
	return depositStats, true
}

func (k *Keeper) SetReserveDepositStats(ctx sdk.Context, depositStats types.DepositStats) {
	var (
		store = k.Store(ctx)
		key   = types.ReserveDepositStatsPrefix
		value = k.cdc.MustMarshal(&depositStats)
	)

	store.Set(key, value)
}

func (k *Keeper) GetReserveDepositStats(ctx sdk.Context) (depositStats types.DepositStats, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.ReserveDepositStatsPrefix
		value = store.Get(key)
	)

	if value == nil {
		return depositStats, false
	}

	k.cdc.MustUnmarshal(value, &depositStats)
	return depositStats, true
}

func (k *Keeper) SetBuyBackDepositStats(ctx sdk.Context, depositStats types.DepositStats) {
	var (
		store = k.Store(ctx)
		key   = types.BuyBackDepositStatsPrefix
		value = k.cdc.MustMarshal(&depositStats)
	)

	store.Set(key, value)
}

func (k *Keeper) GetBuyBackDepositStats(ctx sdk.Context) (depositStats types.DepositStats, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.BuyBackDepositStatsPrefix
		value = store.Get(key)
	)

	if value == nil {
		return depositStats, false
	}

	k.cdc.MustUnmarshal(value, &depositStats)
	return depositStats, true
}

func (k *Keeper) SetBorrowStats(ctx sdk.Context, borrowStats types.DepositStats) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowStatsPrefix
		value = k.cdc.MustMarshal(&borrowStats)
	)

	store.Set(key, value)
}

func (k *Keeper) GetBorrowStats(ctx sdk.Context) (borrowStats types.DepositStats, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.BorrowStatsPrefix
		value = store.Get(key)
	)

	if value == nil {
		return borrowStats, false
	}

	k.cdc.MustUnmarshal(value, &borrowStats)
	return borrowStats, true
}
