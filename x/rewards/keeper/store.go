package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/cosmos/gogoproto/types"

	"github.com/comdex-official/comdex/x/rewards/types"
)

// EPOCHES

// SetEpochInfoByDuration sets EpochInfo with epoch duration as a key.
func (k Keeper) SetEpochInfoByDuration(ctx sdk.Context, epochInfo types.EpochInfo) {
	var (
		store = k.Store(ctx)
		key   = types.GetEpochInfoByDurationKey(epochInfo.Duration)
		value = k.cdc.MustMarshal(&epochInfo)
	)
	store.Set(key, value)
}

// GetEpochInfoByDuration gets EpochInfo by epoch duration.
func (k Keeper) GetEpochInfoByDuration(ctx sdk.Context, duration time.Duration) (epochInfo types.EpochInfo, found bool) {
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

// DeleteEpochInfoByDuration deletes the EpochInfo using epoch duration.
func (k Keeper) DeleteEpochInfoByDuration(ctx sdk.Context, duration time.Duration) {
	var (
		store = k.Store(ctx)
		key   = types.GetEpochInfoByDurationKey(duration)
	)
	store.Delete(key)
}

// GetAllEpochInfos returns all the EpochInfo.
func (k Keeper) GetAllEpochInfos(ctx sdk.Context) (epochInfos []types.EpochInfo) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.EpochInfoByDurationKeyPrefix)
	)
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		var epochInfo types.EpochInfo
		k.cdc.MustUnmarshal(iter.Value(), &epochInfo)
		epochInfos = append(epochInfos, epochInfo)
	}
	return epochInfos
}

// GAUGES

// GetGaugeID return gauge by id.
func (k Keeper) GetGaugeID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.GaugeIDKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)
	return id.GetValue()
}

// SetGaugeID sets id for the gauge.
func (k Keeper) SetGaugeID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.GaugeIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

// SetGauge sets gauge with Id as a key.
func (k Keeper) SetGauge(ctx sdk.Context, gauge types.Gauge) {
	var (
		store = k.Store(ctx)
		key   = types.GetGaugeKey(gauge.Id)
		value = k.cdc.MustMarshal(&gauge)
	)
	store.Set(key, value)
}

// DeleteGauge deletes the gauge.
func (k Keeper) DeleteGauge(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.GetGaugeKey(id)
	)
	store.Delete(key)
}

// GetGaugeByID returns gauge by id.
func (k Keeper) GetGaugeByID(ctx sdk.Context, id uint64) (gauge types.Gauge, found bool) {
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

// GetAllGauges returns all the gauges from store.
func (k Keeper) GetAllGauges(ctx sdk.Context) (gauges []types.Gauge) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.GaugeKeyPrefix)
	)
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		var gauge types.Gauge
		k.cdc.MustUnmarshal(iter.Value(), &gauge)
		gauges = append(gauges, gauge)
	}
	return gauges
}

// SetGaugeIdsByTriggerDuration sets a gauge ids by the trigger duration.
func (k Keeper) SetGaugeIdsByTriggerDuration(ctx sdk.Context, gaugesByTriggerDuration types.GaugeByTriggerDuration) {
	var (
		store = k.Store(ctx)
		key   = types.GetGaugeIdsByTriggerDurationKey(gaugesByTriggerDuration.TriggerDuration)
		value = k.cdc.MustMarshal(&gaugesByTriggerDuration)
	)
	store.Set(key, value)
}

// GetGaugeIdsByTriggerDuration returns all the gauges for the given durtion.
func (k Keeper) GetGaugeIdsByTriggerDuration(ctx sdk.Context, triggerDuration time.Duration) (gaugeIdsByTriggerDuration types.GaugeByTriggerDuration, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.GetGaugeIdsByTriggerDurationKey(triggerDuration)
		value = store.Get(key)
	)
	if value == nil {
		return gaugeIdsByTriggerDuration, false
	}
	k.cdc.MustUnmarshal(value, &gaugeIdsByTriggerDuration)
	return gaugeIdsByTriggerDuration, true
}

func (k Keeper) GetAllGaugeIdsByTriggerDuration(ctx sdk.Context) (gaugeByTriggerDuration []types.GaugeByTriggerDuration) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.GaugeIdsByTriggerDurationKeyPrefix)
	)
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		var gauge types.GaugeByTriggerDuration
		k.cdc.MustUnmarshal(iter.Value(), &gauge)
		gaugeByTriggerDuration = append(gaugeByTriggerDuration, gauge)
	}
	return gaugeByTriggerDuration
}

// GetAllGaugesByGaugeTypeID returns all the gauges with given gaugeTypeId.
func (k Keeper) GetAllGaugesByGaugeTypeID(ctx sdk.Context, gaugeTypeID uint64) (gauges []types.Gauge) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.GaugeKeyPrefix)
	)
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		var gauge types.Gauge
		k.cdc.MustUnmarshal(iter.Value(), &gauge)
		if gauge.GaugeTypeId == gaugeTypeID {
			gauges = append(gauges, gauge)
		}
	}
	return gauges
}
