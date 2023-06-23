package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
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

func MigrateModBal(store sdk.KVStore, cdc codec.BinaryCodec) error {
	key := types.KeyFundModBal
	value := store.Get(key)
	var modBal types.ModBal
	cdc.MustUnmarshal(value, &modBal)

	store.Delete(key)
	SetModBal(store, cdc, modBal)

	return nil
}

func SetModBal(store sdk.KVStore, cdc codec.BinaryCodec, modBal types.ModBal) {
	var (
		key   = types.KeyFundModBal
		value = cdc.MustMarshal(&modBal)
	)

	store.Set(key, value)
}

func MigrateResBal(store sdk.KVStore, cdc codec.BinaryCodec) error {
	key := types.KeyFundReserveBal
	value := store.Get(key)
	var resBal types.ReserveBal
	cdc.MustUnmarshal(value, &resBal)

	store.Delete(key)
	SetResBal(store, cdc, resBal)

	return nil
}

func SetResBal(store sdk.KVStore, cdc codec.BinaryCodec, resBal types.ReserveBal) {
	var (
		key   = types.KeyFundReserveBal
		value = cdc.MustMarshal(&resBal)
	)

	store.Set(key, value)
}

func MigrateResStats(store sdk.KVStore, cdc codec.BinaryCodec) error {
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

func SetResStats(store sdk.KVStore, cdc codec.BinaryCodec, allReserveStats types.AllReserveStats) {
	var (
		key   = types.AllReserveStatsKey(allReserveStats.AssetID)
		value = cdc.MustMarshal(&allReserveStats)
	)

	store.Set(key, value)
}
