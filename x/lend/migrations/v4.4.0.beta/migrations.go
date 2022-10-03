package v4_4_0_beta

import (
	v4_3_0types "github.com/comdex-official/comdex/x/lend/migrations/v4.3.0/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//MigrationHandler codes goes here

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	err := migrateValuesLend(store, cdc)
	if err != nil {
		return err
	}
	err = migrateValuesBorrow(store, cdc)
	if err != nil {
		return err
	}

	return nil
}

func migrateValuesLend(store sdk.KVStore, cdc codec.BinaryCodec) error {

	lendIDKey := lendtypes.LendsKey
	lendIDValue := store.Get(lendIDKey)
	var usersLendMap lendtypes.LendMapping
	cdc.MustUnmarshal(lendIDValue, &usersLendMap)

	for _, ID := range usersLendMap.LendIDs {
		key := append(lendtypes.LendUserPrefix, sdk.Uint64ToBigEndian(ID)...)
		oldKey := store.Get(key)
		oldVal := store.Get(oldKey)

		newKey, newVal := migrateValueLend(cdc, oldKey, oldVal)
		store.Set(newKey, newVal)
		store.Delete(oldKey) // Delete old key, value
	}
	return nil
}

func migrateValueLend(cdc codec.BinaryCodec, oldKey []byte, oldVal []byte) (newKey []byte, newVal []byte) {

	// convert oldVal into lend type of previous version
	// use oldVal to create new lend of updated struct

	var lend v4_3_0types.LendAsset
	cdc.MustUnmarshal(oldVal, &lend)

	newLend := lendtypes.LendAsset{
		ID:                  lend.ID,
		AssetID:             lend.AssetID,
		PoolID:              lend.PoolID,
		Owner:               lend.Owner,
		AmountIn:            lend.AmountIn,
		LendingTime:         lend.LendingTime,
		AvailableToBorrow:   lend.AvailableToBorrow,
		AppID:               lend.AppID,
		GlobalIndex:         sdk.Dec{},
		LastInteractionTime: lend.LastInteractionTime,
		CPoolName:           lend.CPoolName,
	}

	newVal = cdc.MustMarshal(&newLend)
	return oldKey, newVal
}

func migrateValuesBorrow(store sdk.KVStore, cdc codec.BinaryCodec) error {

	borrowIDKey := lendtypes.BorrowsKey
	lendIDValue := store.Get(borrowIDKey)
	var usersBorrowMap lendtypes.BorrowMapping
	cdc.MustUnmarshal(lendIDValue, &usersBorrowMap)

	for _, ID := range usersBorrowMap.BorrowIDs {
		key := append(lendtypes.BorrowPairKeyPrefix, sdk.Uint64ToBigEndian(ID)...)
		oldKey := store.Get(key)
		oldVal := store.Get(oldKey)

		newKey, newVal := migrateValueBorrow(cdc, oldKey, oldVal)
		store.Set(newKey, newVal)
		store.Delete(oldKey) // Delete old key, value
	}
	return nil
}

func migrateValueBorrow(cdc codec.BinaryCodec, oldKey []byte, oldVal []byte) (newKey []byte, newVal []byte) {

	// convert oldVal into borrow type of previous version
	// use oldVal to create new borrow of updated struct

	var borrow v4_3_0types.BorrowAsset
	cdc.MustUnmarshal(oldVal, &borrow)

	newBorrow := lendtypes.BorrowAsset{
		ID:                  borrow.ID,
		LendingID:           borrow.LendingID,
		IsStableBorrow:      borrow.IsStableBorrow,
		PairID:              borrow.PairID,
		AmountIn:            borrow.AmountIn,
		AmountOut:           borrow.AmountOut,
		BridgedAssetAmount:  borrow.BridgedAssetAmount,
		BorrowingTime:       borrow.BorrowingTime,
		StableBorrowRate:    borrow.StableBorrowRate,
		UpdatedAmountOut:    borrow.UpdatedAmountOut,
		GlobalIndex:         borrow.GlobalIndex,
		ReserveGlobalIndex:  borrow.ReserveGlobalIndex,
		LastInteractionTime: borrow.LastInteractionTime,
		CPoolName:           borrow.CPoolName,
	}

	newVal = cdc.MustMarshal(&newBorrow)
	return oldKey, newVal
}
