package keeper

import (
	"time"

	"github.com/comdex-official/comdex/x/incentives/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

// EPOCHES
func (k *Keeper) SetEpochInfoByDuration(ctx sdk.Context, epochInfo types.EpochInfo) {
	var (
		store = k.Store(ctx)
		key   = types.GetEpochInfoByDurationKey(epochInfo.Duration)
		value = k.cdc.MustMarshal(&epochInfo)
	)
	store.Set(key, value)
}

func (k *Keeper) GetEpochInfoByDuration(ctx sdk.Context, duration time.Duration) (epochInfo types.EpochInfo, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.GetEpochInfoByDurationKey(duration)
		value = store.Get(key)
	)
	if value == nil {
		return epochInfo, false
	}
	k.cdc.MustUnmarshal(value, &epochInfo)
	return epochInfo, true
}

func (k *Keeper) DeleteEpochInfoByDuration(ctx sdk.Context, duration time.Duration) {
	var (
		store = k.Store(ctx)
		key   = types.GetEpochInfoByDurationKey(duration)
	)
	store.Delete(key)
}

func (k *Keeper) GetAllEpochInfos(ctx sdk.Context) (epochInfos []types.EpochInfo) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.EpochInfoByDurationKeyPrefix)
	)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var epochInfo types.EpochInfo
		k.cdc.MustUnmarshal(iter.Value(), &epochInfo)
		epochInfos = append(epochInfos, epochInfo)
	}
	return epochInfos
}

// GAUGES

func (k *Keeper) GetGaugeID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.GaugeIdKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)
	return id.GetValue()
}

func (k *Keeper) SetGaugeID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.GaugeIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k *Keeper) SetGauge(ctx sdk.Context, gauge types.Gauge) {
	var (
		store = k.Store(ctx)
		key   = types.GetGaugeKey(gauge.Id)
		value = k.cdc.MustMarshal(&gauge)
	)
	store.Set(key, value)
}

func (k *Keeper) DeleteGauge(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.GetGaugeKey(id)
	)
	store.Delete(key)
}

func (k *Keeper) GetGaugeById(ctx sdk.Context, id uint64) (gauge types.Gauge, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.GetGaugeKey(id)
		value = store.Get(key)
	)
	if value == nil {
		return gauge, false
	}
	k.cdc.MustUnmarshal(value, &gauge)
	return gauge, true
}

func (k *Keeper) GetAllGauges(ctx sdk.Context) (gauges []types.Gauge) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.GaugeKeyPrefix)
	)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var gauge types.Gauge
		k.cdc.MustUnmarshal(iter.Value(), &gauge)
		gauges = append(gauges, gauge)
	}
	return gauges
}
