package v5

import (
	"fmt"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	v5types "github.com/comdex-official/comdex/x/lend/migrations/v5/types"
	"github.com/comdex-official/comdex/x/lend/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	fmt.Println("MigrateStore")
	err := migrateValuesPool(store, cdc)
	if err != nil {
		return err
	}
	err = migrateValuesLend(store, cdc)
	if err != nil {
		return err
	}
	err = migrateValuesBorrow(store, cdc)
	if err != nil {
		return err
	}
	err = migrateValueAuctionParams(store, cdc)
	if err != nil {
		return err
	}
	err = migrateValueLockedBorrows(store, cdc)
	if err != nil {
		return err
	}
	return err
}

func migrateValuesPool(store sdk.KVStore, cdc codec.BinaryCodec) error {
	oldPools := GetPools(store, cdc)
	var (
		assetDataPoolOne []*types.AssetDataPoolMapping
		assetDataPoolTwo []*types.AssetDataPoolMapping
		assetData        []*types.AssetDataPoolMapping
	)
	assetDataPoolOneAssetOne := &types.AssetDataPoolMapping{
		AssetID:          1,
		AssetTransitType: 3,
		SupplyCap:        uint64(5000000000000000000),
	}
	assetDataPoolOneAssetTwo := &types.AssetDataPoolMapping{
		AssetID:          2,
		AssetTransitType: 1,
		SupplyCap:        uint64(1000000000000000000),
	}
	assetDataPoolOneAssetThree := &types.AssetDataPoolMapping{
		AssetID:          3,
		AssetTransitType: 2,
		SupplyCap:        uint64(5000000000000000000),
	}
	assetDataPoolTwoAssetFour := &types.AssetDataPoolMapping{
		AssetID:          4,
		AssetTransitType: 1,
		SupplyCap:        uint64(3000000000000000000),
	}
	assetDataPoolOne = append(assetDataPoolOne, assetDataPoolOneAssetOne, assetDataPoolOneAssetTwo, assetDataPoolOneAssetThree)
	assetDataPoolTwo = append(assetDataPoolTwo, assetDataPoolTwoAssetFour, assetDataPoolOneAssetOne, assetDataPoolOneAssetThree)

	for _, j := range oldPools {
		if j.PoolID == 1 {
			assetData = assetDataPoolOne
		} else {
			assetData = assetDataPoolTwo
		}
		newPool := types.Pool{
			PoolID:       j.PoolID,
			ModuleName:   j.ModuleName,
			CPoolName:    j.CPoolName,
			ReserveFunds: j.ReserveFunds,
			AssetData:    assetData,
		}

		for _, v := range newPool.AssetData {
			var assetStats types.PoolAssetLBMapping
			assetStats.PoolID = newPool.PoolID
			assetStats.AssetID = v.AssetID
			assetStats.TotalBorrowed = sdk.ZeroInt()
			assetStats.TotalStableBorrowed = sdk.ZeroInt()
			assetStats.TotalLend = sdk.ZeroInt()
			assetStats.TotalInterestAccumulated = sdk.ZeroInt()
			SetAssetStatsByPoolIDAndAssetID(store, cdc, assetStats)
			reserveBuybackStats, found := GetReserveBuybackAssetData(store, cdc, v.AssetID)
			if !found {
				reserveBuybackStats.AssetID = v.AssetID
				reserveBuybackStats.ReserveAmount = sdk.ZeroInt()
				reserveBuybackStats.BuybackAmount = sdk.ZeroInt()
				SetReserveBuybackAssetData(store, cdc, reserveBuybackStats)
			}
		}
		key := types.PoolKey(j.PoolID)
		store.Delete(key)
		SetPool(store, cdc, newPool)
		SetPoolID(store, cdc, newPool.PoolID)
	}
	return nil
}

func SetPoolID(store sdk.KVStore, cdc codec.BinaryCodec, id uint64) {
	var (
		key   = types.PoolIDPrefix
		value = cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func SetPool(store sdk.KVStore, cdc codec.BinaryCodec, pool types.Pool) {
	var (
		key   = types.PoolKey(pool.PoolID)
		value = cdc.MustMarshal(&pool)
	)

	store.Set(key, value)
}

func SetReserveBuybackAssetData(store sdk.KVStore, cdc codec.BinaryCodec, reserve types.ReserveBuybackAssetData) {
	var (
		key   = types.ReserveBuybackAssetDataKey(reserve.AssetID)
		value = cdc.MustMarshal(&reserve)
	)

	store.Set(key, value)
}

func GetReserveBuybackAssetData(store sdk.KVStore, cdc codec.BinaryCodec, id uint64) (reserve types.ReserveBuybackAssetData, found bool) {
	var (
		key   = types.ReserveBuybackAssetDataKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return reserve, false
	}

	cdc.MustUnmarshal(value, &reserve)
	return reserve, true
}

func GetPools(store sdk.KVStore, cdc codec.BinaryCodec) (pools []types.Pool) {
	var (
		iter = sdk.KVStorePrefixIterator(store, types.PoolKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var pool types.Pool
		cdc.MustUnmarshal(iter.Value(), &pool)
		pools = append(pools, pool)
	}

	return pools
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
	return newVal, types.BorrowUserKey(newVal.ID)
}

// for auction params migration

func migrateValueAuctionParams(store sdk.KVStore, cdc codec.BinaryCodec) error {
	fmt.Println("migrateValueAuctionParams")
	buffer, _ := sdk.NewDecFromStr("1.2")
	cusp, _ := sdk.NewDecFromStr("0.4")
	auctionParams := types.AuctionParams{
		AppId:                  3,
		AuctionDurationSeconds: 21600,
		Buffer:                 buffer,
		Cusp:                   cusp,
		Step:                   sdk.NewIntFromUint64(360),
		PriceFunctionType:      1,
		DutchId:                3,
		BidDurationSeconds:     10800,
	}
	var (
		oldKey = types.AuctionParamKey(1)
		key    = types.AuctionParamKey(auctionParams.AppId)
		value  = cdc.MustMarshal(&auctionParams)
	)
	store.Delete(oldKey)
	store.Set(key, value)
	return nil
}

// for locked borrow to borrow

func migrateValueLockedBorrows(store sdk.KVStore, cdc codec.BinaryCodec) error {
	fmt.Println("migrateValueAuctionParams")
	var (
		iter = sdk.KVStorePrefixIterator(store, liquidationtypes.LockedVaultKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	var lockedVaults []liquidationtypes.LockedVault
	for ; iter.Valid(); iter.Next() {
		var lockedVault liquidationtypes.LockedVault
		cdc.MustUnmarshal(iter.Value(), &lockedVault)
		if lockedVault.GetBorrowMetaData() != nil {
			if lockedVault.AmountIn.GT(sdk.ZeroInt()) && lockedVault.AmountOut.GT(sdk.ZeroInt()) {
				lockedVaults = append(lockedVaults, lockedVault)
			} else {
				key := liquidationtypes.LockedVaultKey(lockedVault.AppId, lockedVault.LockedVaultId)
				store.Delete(key)
			}
		}
	}
	counterKey := types.BorrowCounterIDPrefix

	for _, v := range lockedVaults {
		lockedVaultKey := liquidationtypes.LockedVaultKey(v.AppId, v.LockedVaultId)
		newVal, Key := migrateValueLockedBorrow(store, cdc, v)
		store.Delete(lockedVaultKey)

		store.Delete(Key)
		value := cdc.MustMarshal(&newVal)
		idValue := cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: newVal.ID,
			},
		)
		store.Set(counterKey, idValue)
		store.Set(Key, value)
		// updating global borrow stats and borrowing IDs
		pair, found := GetLendPair(store, cdc, newVal.PairID)
		if !found {
			return types.ErrorPairNotFound
		}
		PoolAssetLBMapping := GetAssetStatsByPoolIDAndAssetID(store, cdc, pair.AssetOutPoolID, pair.AssetOut)
		PoolAssetLBMapping.TotalBorrowed = PoolAssetLBMapping.TotalBorrowed.Add(newVal.AmountOut.Amount)
		PoolAssetLBMapping.BorrowIds = append(PoolAssetLBMapping.BorrowIds, newVal.ID)
		SetAssetStatsByPoolIDAndAssetID(store, cdc, PoolAssetLBMapping)

		// updating UserAssetLendBorrowMapping for user
		lend, found := GetLend(store, cdc, newVal.LendingID)
		mappingData, found := GetUserLendBorrowMapping(store, cdc, lend.Owner, lend.ID)
		mappingData.BorrowId = append(mappingData.BorrowId, newVal.ID)
		SetUserLendBorrowMapping(store, cdc, mappingData)
	}

	return nil
}

func GetAsset(store sdk.KVStore, cdc codec.BinaryCodec, id uint64) (asset assettypes.Asset, found bool) {
	var (
		key   = assettypes.AssetKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return asset, false
	}

	cdc.MustUnmarshal(value, &asset)
	return asset, true
}

func GetPool(store sdk.KVStore, cdc codec.BinaryCodec, id uint64) (pool types.Pool, found bool) {
	var (
		key   = types.PoolKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return pool, false
	}

	cdc.MustUnmarshal(value, &pool)
	return pool, true
}

func migrateValueLockedBorrow(store sdk.KVStore, cdc codec.BinaryCodec, v liquidationtypes.LockedVault) (newBorrow types.BorrowAsset, oldKey []byte) {
	pair, _ := GetLendPair(store, cdc, v.ExtendedPairId)
	assetIn, _ := GetAsset(store, cdc, pair.AssetIn)
	assetOut, _ := GetAsset(store, cdc, pair.AssetOut)
	amountIn := sdk.NewCoin(assetIn.Denom, v.AmountIn)
	amountOut := sdk.NewCoin(assetOut.Denom, v.AmountOut)
	pool, _ := GetPool(store, cdc, pair.AssetOutPoolID)

	borrowMetaData := v.GetBorrowMetaData()
	globalIndex, _ := sdk.NewDecFromStr("0.002")
	newBorrow = types.BorrowAsset{
		ID:                  v.OriginalVaultId,
		LendingID:           borrowMetaData.LendingId,
		IsStableBorrow:      borrowMetaData.IsStableBorrow,
		PairID:              v.ExtendedPairId,
		AmountIn:            amountIn,
		AmountOut:           amountOut,
		BridgedAssetAmount:  borrowMetaData.BridgedAssetAmount,
		BorrowingTime:       v.LiquidationTimestamp,
		StableBorrowRate:    borrowMetaData.StableBorrowRate,
		InterestAccumulated: sdk.NewDecFromInt(v.InterestAccumulated),
		GlobalIndex:         globalIndex,
		ReserveGlobalIndex:  sdk.OneDec(),
		LastInteractionTime: v.LiquidationTimestamp,
		CPoolName:           pool.CPoolName,
		IsLiquidated:        false,
	}
	return newBorrow, types.BorrowUserKey(newBorrow.ID)
}
