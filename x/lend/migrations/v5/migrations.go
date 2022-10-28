package v5

import (
	"fmt"
	v5types "github.com/comdex-official/comdex/x/lend/migrations/v5/types"
	"github.com/comdex-official/comdex/x/lend/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	fmt.Println("MigrateStore")
	err := migrateValuesLend(store, cdc)
	if err != nil {
		return err
	}
	err = migrateValuesBorrow(store, cdc)
	if err != nil {
		return err
	}
	return err
}

func migrateValuesLend(store sdk.KVStore, cdc codec.BinaryCodec) error {
	fmt.Println(" migrateValuesLend")
	var (
		iter = sdk.KVStorePrefixIterator(store, types.LendUserPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	var lendAssets []v5types.LendAssetOld
	for ; iter.Valid(); iter.Next() {
		var asset v5types.LendAssetOld
		cdc.MustUnmarshal(iter.Value(), &asset)
		lendAssets = append(lendAssets, asset)
	}
	counterKey := types.LendCounterIDPrefix
	for _, v := range lendAssets {
		newVal, Key := migrateValueLend(v)
		if v.AmountIn.Amount.LTE(sdk.ZeroInt()) || v.AvailableToBorrow.LT(sdk.ZeroInt()) {
			store.Delete(Key)
			continue
		}
		store.Delete(Key)
		value := cdc.MustMarshal(&newVal)
		idValue := cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: v.ID,
			},
		)
		store.Set(counterKey, idValue)
		store.Set(Key, value)
		// updating global lend stats and lending IDs
		PoolAssetLBMapping := GetAssetStatsByPoolIDAndAssetID(store, cdc, v.PoolID, v.AssetID)
		PoolAssetLBMapping.TotalLend = PoolAssetLBMapping.TotalLend.Add(v.AmountIn.Amount)
		PoolAssetLBMapping.LendIds = append(PoolAssetLBMapping.LendIds, v.ID)
		SetAssetStatsByPoolIDAndAssetID(store, cdc, PoolAssetLBMapping)
		// making UserAssetLendBorrowMapping for user
		var mappingData types.UserAssetLendBorrowMapping
		mappingData.Owner = v.Owner
		mappingData.LendId = v.ID
		mappingData.PoolId = v.PoolID
		mappingData.BorrowId = nil
		SetUserLendBorrowMapping(store, cdc, mappingData)

	}
	return nil
}

func SetAssetStatsByPoolIDAndAssetID(store sdk.KVStore, cdc codec.BinaryCodec, PoolAssetLBMapping types.PoolAssetLBMapping) {
	var (
		key   = types.SetAssetStatsByPoolIDAndAssetID(PoolAssetLBMapping.PoolID, PoolAssetLBMapping.AssetID)
		value = cdc.MustMarshal(&PoolAssetLBMapping)
	)

	store.Set(key, value)
}

func SetUserLendBorrowMapping(store sdk.KVStore, cdc codec.BinaryCodec, userMapping types.UserAssetLendBorrowMapping) {
	var (
		key   = types.UserLendBorrowMappingKey(userMapping.Owner, userMapping.LendId)
		value = cdc.MustMarshal(&userMapping)
	)

	store.Set(key, value)
}

func GetAssetStatsByPoolIDAndAssetID(store sdk.KVStore, cdc codec.BinaryCodec, poolID, assetID uint64) (PoolAssetLBMapping types.PoolAssetLBMapping) {
	var (
		key   = types.SetAssetStatsByPoolIDAndAssetID(poolID, assetID)
		value = store.Get(key)
	)

	if value == nil {
		return PoolAssetLBMapping
	}

	cdc.MustUnmarshal(value, &PoolAssetLBMapping)
	return PoolAssetLBMapping
}

func migrateValueLend(oldVal v5types.LendAssetOld) (newVal types.LendAsset, oldKey []byte) {
	fmt.Println("migrateValueLend")
	newVal = types.LendAsset{
		ID:                  oldVal.ID,
		AssetID:             oldVal.AssetID,
		PoolID:              oldVal.PoolID,
		Owner:               oldVal.Owner,
		AmountIn:            oldVal.AmountIn,
		LendingTime:         oldVal.LendingTime,
		AvailableToBorrow:   oldVal.AvailableToBorrow,
		AppID:               oldVal.AppID,
		GlobalIndex:         oldVal.GlobalIndex,
		LastInteractionTime: oldVal.LastInteractionTime,
		CPoolName:           oldVal.CPoolName,
	}
	return newVal, types.LendUserKey(newVal.ID)
}

// for borrow function migration

func migrateValuesBorrow(store sdk.KVStore, cdc codec.BinaryCodec) error {
	fmt.Println(" migrateValuesBorrow")
	var (
		iter = sdk.KVStorePrefixIterator(store, types.BorrowPairKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	var borrowAssets []v5types.BorrowAssetOld
	for ; iter.Valid(); iter.Next() {
		var asset v5types.BorrowAssetOld
		cdc.MustUnmarshal(iter.Value(), &asset)
		borrowAssets = append(borrowAssets, asset)
	}
	counterKey := types.BorrowCounterIDPrefix
	for _, v := range borrowAssets {
		newVal, Key := migrateValueBorrow(v)
		if v.AmountIn.Amount.LTE(sdk.ZeroInt()) || v.AmountOut.Amount.LTE(sdk.ZeroInt()) {
			store.Delete(Key)
			continue
		}
		store.Delete(Key)
		value := cdc.MustMarshal(&newVal)
		idValue := cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: v.ID,
			},
		)
		store.Set(counterKey, idValue)
		store.Set(Key, value)
		// updating global borrow stats and borrowing IDs
		pair, found := GetLendPair(store, cdc, v.PairID)
		if !found {
			return types.ErrorPairNotFound
		}
		PoolAssetLBMapping := GetAssetStatsByPoolIDAndAssetID(store, cdc, pair.AssetOutPoolID, pair.AssetOut)
		PoolAssetLBMapping.TotalBorrowed = PoolAssetLBMapping.TotalBorrowed.Add(v.AmountOut.Amount)
		PoolAssetLBMapping.BorrowIds = append(PoolAssetLBMapping.BorrowIds, v.ID)
		SetAssetStatsByPoolIDAndAssetID(store, cdc, PoolAssetLBMapping)

		// updating UserAssetLendBorrowMapping for user
		lend, found := GetLend(store, cdc, v.LendingID)
		mappingData, found := GetUserLendBorrowMapping(store, cdc, lend.Owner, lend.ID)
		mappingData.BorrowId = append(mappingData.BorrowId, v.ID)
		SetUserLendBorrowMapping(store, cdc, mappingData)
	}
	return nil
}

func GetLend(store sdk.KVStore, cdc codec.BinaryCodec, id uint64) (lend types.LendAsset, found bool) {
	var (
		key   = types.LendUserKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return lend, false
	}

	cdc.MustUnmarshal(value, &lend)
	return lend, true
}

func GetLendPair(store sdk.KVStore, cdc codec.BinaryCodec, id uint64) (pair types.Extended_Pair, found bool) {
	var (
		key   = types.LendPairKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return pair, false
	}

	cdc.MustUnmarshal(value, &pair)
	return pair, true
}

func GetUserLendBorrowMapping(store sdk.KVStore, cdc codec.BinaryCodec, owner string, lendID uint64) (userMapping types.UserAssetLendBorrowMapping, found bool) {
	var (
		key   = types.UserLendBorrowMappingKey(owner, lendID)
		value = store.Get(key)
	)

	if value == nil {
		return userMapping, false
	}

	cdc.MustUnmarshal(value, &userMapping)
	return userMapping, true
}

func migrateValueBorrow(v v5types.BorrowAssetOld) (newVal types.BorrowAsset, oldKey []byte) {
	fmt.Println("migrateValueLend")
	newVal = types.BorrowAsset{
		ID:                  v.ID,
		LendingID:           v.LendingID,
		IsStableBorrow:      v.IsStableBorrow,
		PairID:              v.PairID,
		AmountIn:            v.AmountIn,
		AmountOut:           v.AmountOut,
		BridgedAssetAmount:  v.BridgedAssetAmount,
		BorrowingTime:       v.BorrowingTime,
		StableBorrowRate:    v.StableBorrowRate,
		InterestAccumulated: sdk.NewDecFromInt(v.Interest_Accumulated),
		GlobalIndex:         v.GlobalIndex,
		ReserveGlobalIndex:  v.ReserveGlobalIndex,
		LastInteractionTime: v.LastInteractionTime,
		CPoolName:           v.CPoolName,
		IsLiquidated:        false,
	}
	return newVal, types.LendUserKey(newVal.ID)
}
