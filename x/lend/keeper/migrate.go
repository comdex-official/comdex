package keeper

import (
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	migrationtypes "github.com/comdex-official/comdex/x/lend/migrations/v2/types"
	"github.com/comdex-official/comdex/x/lend/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	//  Migrate these 3 for their store: (export fix)
	// 		FundModBal
	//		FundReserveBal
	//		TotalReserveStatsByAssetID

	store := ctx.KVStore(storeKey)
	err := MigrateModBal(store, cdc)
	if err != nil {
		return err
	}

	err = MigrateResBal(store, cdc)
	if err != nil {
		return err
	}

	err = MigrateResStats(store, cdc)
	if err != nil {
		return err
	}

	return err
}

func MigrateModBal(store storetypes.KVStore, cdc codec.BinaryCodec) error {
	key := types.KeyFundModBal
	value := store.Get(key)
	var modBal types.ModBal
	cdc.MustUnmarshal(value, &modBal)

	store.Delete(key)
	SetModBal(store, cdc, modBal)

	return nil
}

func SetModBal(store storetypes.KVStore, cdc codec.BinaryCodec, modBal types.ModBal) {
	var (
		key   = types.KeyFundModBal
		value = cdc.MustMarshal(&modBal)
	)

	store.Set(key, value)
}

func MigrateResBal(store storetypes.KVStore, cdc codec.BinaryCodec) error {
	key := types.KeyFundReserveBal
	value := store.Get(key)
	var resBal types.ReserveBal
	cdc.MustUnmarshal(value, &resBal)

	store.Delete(key)
	SetResBal(store, cdc, resBal)

	return nil
}

func SetResBal(store storetypes.KVStore, cdc codec.BinaryCodec, resBal types.ReserveBal) {
	var (
		key   = types.KeyFundReserveBal
		value = cdc.MustMarshal(&resBal)
	)

	store.Set(key, value)
}

func MigrateResStats(store storetypes.KVStore, cdc codec.BinaryCodec) error {
	key1 := types.AllReserveStatsKey(1)
	value1 := store.Get(key1)
	var allReserveStats1 types.AllReserveStats
	cdc.MustUnmarshal(value1, &allReserveStats1)
	store.Delete(key1)
	SetResStats(store, cdc, allReserveStats1)

	key2 := types.AllReserveStatsKey(2)
	value2 := store.Get(key2)
	var allReserveStats2 types.AllReserveStats
	cdc.MustUnmarshal(value2, &allReserveStats2)
	store.Delete(key2)
	SetResStats(store, cdc, allReserveStats2)

	key3 := types.AllReserveStatsKey(3)
	value3 := store.Get(key3)
	var allReserveStats3 types.AllReserveStats
	cdc.MustUnmarshal(value3, &allReserveStats3)
	store.Delete(key3)
	SetResStats(store, cdc, allReserveStats3)

	return nil
}

func SetResStats(store storetypes.KVStore, cdc codec.BinaryCodec, allReserveStats types.AllReserveStats) {
	var (
		key   = types.AllReserveStatsKey(allReserveStats.AssetID)
		value = cdc.MustMarshal(&allReserveStats)
	)

	store.Set(key, value)
}

func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return MigrateStoreV2(ctx, m.keeper.storeKey, m.keeper.cdc)
}

func MigrateStoreV2(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	//  Migrate these 2 for their store:
	// 		LendPairs
	//		AssetRatesParams

	store := ctx.KVStore(storeKey)
	err := MigrateLendPairs(store, cdc)
	if err != nil {
		return err
	}

	err = MigrateAssetRatesParams(store, cdc)
	if err != nil {
		return err
	}

	return err
}

func MigrateLendPairs(store storetypes.KVStore, cdc codec.BinaryCodec) error {

	iterator := store.Iterator(types.LendPairKeyPrefix, storetypes.PrefixEndBytes(types.LendPairKeyPrefix))
	defer func(iterator storetypes.Iterator) {
		err := iterator.Close()
		if err != nil {
			return
		}
	}(iterator)

	var pair migrationtypes.Extended_Pair_Old

	for iterator.Valid() {
		key := iterator.Key()
		value := iterator.Value()

		cdc.MustUnmarshal(value, &pair)
		store.Delete(key)
		SetLendPairs(store, cdc, pair)

		iterator.Next()
	}
	return nil
}

func SetLendPairs(store storetypes.KVStore, cdc codec.BinaryCodec, pair migrationtypes.Extended_Pair_Old) {
	newPair := types.Extended_Pair{
		Id:              pair.Id,
		AssetIn:         pair.AssetIn,
		AssetOut:        pair.AssetOut,
		IsInterPool:     pair.IsInterPool,
		AssetOutPoolID:  pair.AssetOutPoolID,
		MinUsdValueLeft: pair.MinUsdValueLeft,
		IsEModeEnabled:  false,
	}

	var (
		key   = types.LendPairKey(pair.Id)
		value = cdc.MustMarshal(&newPair)
	)

	store.Set(key, value)
}

func MigrateAssetRatesParams(store storetypes.KVStore, cdc codec.BinaryCodec) error {

	iterator := store.Iterator(types.AssetRatesParamsKeyPrefix, storetypes.PrefixEndBytes(types.AssetRatesParamsKeyPrefix))
	defer func(iterator storetypes.Iterator) {
		err := iterator.Close()
		if err != nil {
			return
		}
	}(iterator)

	var assetRatesParams migrationtypes.AssetRatesParams_Old

	for iterator.Valid() {
		key := iterator.Key()
		value := iterator.Value()

		cdc.MustUnmarshal(value, &assetRatesParams)
		store.Delete(key)
		SetAssetRatesParams(store, cdc, assetRatesParams)

		iterator.Next()
	}
	return nil
}

func SetAssetRatesParams(store storetypes.KVStore, cdc codec.BinaryCodec, assetRatesParams migrationtypes.AssetRatesParams_Old) {
	newAssetRatesParams := types.AssetRatesParams{
		AssetID:               assetRatesParams.AssetID,
		UOptimal:              assetRatesParams.UOptimal,
		Base:                  assetRatesParams.Base,
		Slope1:                assetRatesParams.Slope1,
		Slope2:                assetRatesParams.Slope2,
		EnableStableBorrow:    assetRatesParams.EnableStableBorrow,
		StableBase:            assetRatesParams.StableBase,
		StableSlope1:          assetRatesParams.StableSlope1,
		StableSlope2:          assetRatesParams.StableSlope2,
		Ltv:                   assetRatesParams.Ltv,
		LiquidationThreshold:  assetRatesParams.LiquidationThreshold,
		LiquidationPenalty:    assetRatesParams.LiquidationPenalty,
		LiquidationBonus:      assetRatesParams.LiquidationBonus,
		ReserveFactor:         assetRatesParams.ReserveFactor,
		CAssetID:              assetRatesParams.CAssetID,
		IsIsolated:            false,
		ELtv:                  sdkmath.LegacyNewDec(0),
		ELiquidationThreshold: sdkmath.LegacyNewDec(0),
		ELiquidationPenalty:   sdkmath.LegacyNewDec(0),
	}

	var (
		key   = types.AssetRatesParamsKey(newAssetRatesParams.AssetID)
		value = cdc.MustMarshal(&newAssetRatesParams)
	)

	store.Set(key, value)
}
