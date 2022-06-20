package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k *Keeper) AddLendPairsRecords(ctx sdk.Context, records ...types.Extended_Pair) error {

	for _, msg := range records {

		_, found := k.GetLendPair(ctx, msg.Id)
		if found {
			return types.ErrorDuplicateLendPair
		}

		var (
			id   = k.GetLendPairID(ctx)
			pair = types.Extended_Pair{
				Id:             id + 1,
				AssetIn:        msg.AssetIn,
				AssetOut:       msg.AssetOut,
				IsInterPool:    msg.IsInterPool,
				AssetOutPoolId: msg.AssetOutPoolId,
			}
		)

		k.SetLendPairID(ctx, pair.Id)
		k.SetLendPair(ctx, pair)
	}
	return nil
}

func (k Keeper) AddPoolRecords(ctx sdk.Context, pool types.Pool) error {
	for _, v := range pool.AssetData {
		_, found := k.GetAsset(ctx, v.AssetId)
		if !found {
			return types.ErrorAssetDoesNotExist
		}
	}

	poolId := k.GetPoolId(ctx)
	newPool := types.Pool{
		PoolId:               poolId + 1,
		ModuleName:           pool.ModuleName,
		FirstBridgedAssetId:  pool.FirstBridgedAssetId,
		SecondBridgedAssetId: pool.SecondBridgedAssetId,
		AssetData:            pool.AssetData,
	}
	k.SetPool(ctx, newPool)
	k.SetPoolId(ctx, newPool.PoolId)
	return nil
}

func (k Keeper) AddAssetToPair(ctx sdk.Context, assetToPair types.AssetToPairMapping) error {
	_, found := k.GetAsset(ctx, assetToPair.AssetId)
	if !found {
		return types.ErrorAssetDoesNotExist
	}
	for _, v := range assetToPair.PairId {
		_, found := k.GetLendPair(ctx, v)
		if !found {
			return types.ErrorPairDoesNotExist
		}
	}
	k.SetAssetToPair(ctx, assetToPair)

	return nil
}

func (k *Keeper) UpdateLendPairRecords(ctx sdk.Context, msg types.Extended_Pair) error {
	pair, found := k.GetLendPair(ctx, msg.Id)
	if !found {
		return types.ErrorPairDoesNotExist
	}

	k.SetLendPair(ctx, pair)
	return nil
}

func (k *Keeper) SetLendPairID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendPairIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) SetLendPair(ctx sdk.Context, pair types.Extended_Pair) {
	var (
		store = k.Store(ctx)
		key   = types.LendPairKey(pair.Id)
		value = k.cdc.MustMarshal(&pair)
	)

	store.Set(key, value)
}

func (k *Keeper) GetLendPair(ctx sdk.Context, id uint64) (pair types.Extended_Pair, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LendPairKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return pair, false
	}

	k.cdc.MustUnmarshal(value, &pair)
	return pair, true
}

func (k *Keeper) GetLendPairID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LendPairIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var count protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &count)

	return count.GetValue()
}

func (k *Keeper) AddAssetRatesStats(ctx sdk.Context, records ...types.AssetRatesStats) error {

	for _, msg := range records {

		_, found := k.GetAssetRatesStats(ctx, msg.AssetId)
		if found {
			return types.ErrorDuplicateAssetRatesStats
		}

		var (
			assetRatesStats = types.AssetRatesStats{
				AssetId:              msg.AssetId,
				UOptimal:             msg.UOptimal,
				Base:                 msg.Base,
				Slope1:               msg.Slope1,
				Slope2:               msg.Slope2,
				EnableStableBorrow:   msg.EnableStableBorrow,
				StableBase:           msg.StableBase,
				StableSlope1:         msg.StableSlope1,
				StableSlope2:         msg.StableSlope2,
				Ltv:                  msg.Ltv,
				LiquidationThreshold: msg.LiquidationThreshold,
				LiquidationPenalty:   msg.LiquidationPenalty,
				ReserveFactor:        msg.ReserveFactor,
			}
		)

		k.SetAssetRatesStats(ctx, assetRatesStats)
	}
	return nil
}

func (k *Keeper) SetAssetRatesStats(ctx sdk.Context, assetRatesStats types.AssetRatesStats) {
	var (
		store = k.Store(ctx)
		key   = types.AssetRatesStatsKey(assetRatesStats.AssetId)
		value = k.cdc.MustMarshal(&assetRatesStats)
	)

	store.Set(key, value)
}

func (k *Keeper) GetAssetRatesStats(ctx sdk.Context, assetId uint64) (assetRatesStats types.AssetRatesStats, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AssetRatesStatsKey(assetId)
		value = store.Get(key)
	)

	if value == nil {
		return assetRatesStats, false
	}

	k.cdc.MustUnmarshal(value, &assetRatesStats)
	return assetRatesStats, true
}
