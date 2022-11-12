package v5

import (
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	err := MigrateValueCollectorLookupTableData(store, cdc)
	if err != nil {
		return err
	}
	return nil
}

func MigrateValueCollectorLookupTableData(store sdk.KVStore, cdc codec.BinaryCodec) error {
	collectData := collectortypes.CollectorLookupTableData{
		AppId:            2,
		CollectorAssetId: 3,
		SecondaryAssetId: 9,
		SurplusThreshold: sdk.NewInt(100000000000000),
		DebtThreshold:    sdk.NewInt(1000000000),
		LockerSavingRate: sdk.MustNewDecFromStr("0.015"),
		LotSize:          sdk.NewInt(200000000),
		BidFactor:        sdk.MustNewDecFromStr("0.01"),
		DebtLotSize:      sdk.NewInt(2000000),
		BlockHeight:      0,
	}
	key := collectortypes.CollectorLookupTableMappingKey(collectData.AppId, collectData.CollectorAssetId)
	store.Delete(key)
	SetcollectData(store, cdc, collectData)

	return nil
}

func SetcollectData(store sdk.KVStore, cdc codec.BinaryCodec, collect collectortypes.CollectorLookupTableData) {
	var (
		key   = collectortypes.CollectorLookupTableMappingKey(collect.AppId, collect.CollectorAssetId)
		value = cdc.MustMarshal(&collect)
	)

	store.Set(key, value)
}
