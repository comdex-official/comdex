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
				Id:              id + 1,
				AssetIn:         msg.AssetIn,
				AssetOut:        msg.AssetOut,
				IsInterPool:     msg.IsInterPool,
				AssetOutPoolId:  msg.AssetOutPoolId,
				MinUsdValueLeft: msg.MinUsdValueLeft,
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
	depositStats, found := k.GetDepositStats(ctx)
	userDepositStats, _ := k.GetUserDepositStats(ctx)
	ReserveDepositStats, _ := k.GetReserveDepositStats(ctx)
	BuyBackDepositStats, _ := k.GetBuyBackDepositStats(ctx)
	BorrowStats, _ := k.GetBorrowStats(ctx)
	var balanceStats []types.BalanceStats
	if !found {
		for _, v := range pool.AssetData {
			balanceStat := types.BalanceStats{
				AssetId: v.AssetId,
				Amount:  sdk.ZeroInt(),
			}
			balanceStats = append(balanceStats, balanceStat)
			depositStats = types.DepositStats{BalanceStats: balanceStats}
			userDepositStats = types.DepositStats{BalanceStats: balanceStats}
			ReserveDepositStats = types.DepositStats{BalanceStats: balanceStats}
			BuyBackDepositStats = types.DepositStats{BalanceStats: balanceStats}
			BorrowStats = types.DepositStats{BalanceStats: balanceStats}
			k.SetDepositStats(ctx, depositStats)
			k.SetUserDepositStats(ctx, depositStats)
			k.SetReserveDepositStats(ctx, depositStats)
			k.SetBorrowStats(ctx, depositStats)
		}
	} else {
		balanceStat := types.BalanceStats{
			AssetId: pool.MainAssetId,
			Amount:  sdk.ZeroInt(),
		}
		balanceStats = append(depositStats.BalanceStats, balanceStat)
		depositStats = types.DepositStats{BalanceStats: balanceStats}
		k.SetDepositStats(ctx, depositStats)
		k.SetUserDepositStats(ctx, userDepositStats)
		k.SetReserveDepositStats(ctx, ReserveDepositStats)
		k.SetBuyBackDepositStats(ctx, BuyBackDepositStats)
		k.SetBorrowStats(ctx, BorrowStats)
	}

	poolID := k.GetPoolID(ctx)
	newPool := types.Pool{
		PoolId:               poolID + 1,
		ModuleName:           pool.ModuleName,
		MainAssetId:          pool.MainAssetId,
		FirstBridgedAssetId:  pool.FirstBridgedAssetId,
		SecondBridgedAssetId: pool.SecondBridgedAssetId,
		CPoolName:            pool.CPoolName,
		AssetData:            pool.AssetData,
	}
	k.SetPool(ctx, newPool)
	k.SetPoolID(ctx, newPool.PoolId)
	return nil
}

func (k Keeper) AddAssetToPair(ctx sdk.Context, assetToPair types.AssetToPairMapping) error {
	_, found := k.GetAsset(ctx, assetToPair.AssetId)
	if !found {
		return types.ErrorAssetDoesNotExist
	}
	_, found = k.GetPool(ctx, assetToPair.PoolId)
	if !found {
		return types.ErrPoolNotFound
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
				LiquidationBonus:     msg.LiquidationBonus,
				ReserveFactor:        msg.ReserveFactor,
				CAssetId:             msg.CAssetId,
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

func (k *Keeper) GetAssetRatesStats(ctx sdk.Context, assetID uint64) (assetRatesStats types.AssetRatesStats, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AssetRatesStatsKey(assetID)
		value = store.Get(key)
	)

	if value == nil {
		return assetRatesStats, false
	}

	k.cdc.MustUnmarshal(value, &assetRatesStats)
	return assetRatesStats, true
}

func (k *Keeper) SetDepositStats(ctx sdk.Context, depositStats types.DepositStats) {
	var (
		store = k.Store(ctx)
		key   = types.DepositStatsPrefix
		value = k.cdc.MustMarshal(&depositStats)
	)

	store.Set(key, value)
}

func (k *Keeper) GetDepositStats(ctx sdk.Context) (depositStats types.DepositStats, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.DepositStatsPrefix
		value = store.Get(key)
	)

	if value == nil {
		return depositStats, false
	}

	k.cdc.MustUnmarshal(value, &depositStats)
	return depositStats, true
}
