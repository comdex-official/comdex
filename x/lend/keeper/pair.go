package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/lend/types"
)

func (k Keeper) AddLendPairsRecords(ctx sdk.Context, records ...types.Extended_Pair) error {
	for _, msg := range records {
		_, found := k.GetLendPair(ctx, msg.Id)
		if found {
			return types.ErrorDuplicateLendPair
		}
		if msg.AssetIn == msg.AssetOut {
			return types.ErrorAssetsCanNotBeSame
		}

		var (
			id   = k.GetLendPairID(ctx)
			pair = types.Extended_Pair{
				Id:              id + 1,
				AssetIn:         msg.AssetIn,
				AssetOut:        msg.AssetOut,
				IsInterPool:     msg.IsInterPool,
				AssetOutPoolID:  msg.AssetOutPoolID,
				MinUsdValueLeft: msg.MinUsdValueLeft,
			}
		)

		k.SetLendPairID(ctx, pair.Id)
		k.SetLendPair(ctx, pair)
	}
	return nil
}

func (k Keeper) UpdateLendPairsRecords(ctx sdk.Context, msg types.Extended_Pair) error {
	pair, found := k.GetLendPair(ctx, msg.Id)
	if !found {
		return types.ErrorPairNotFound
	}

	_, found = k.GetAsset(ctx, msg.AssetIn)
	if !found {
		return types.ErrorAssetDoesNotExist
	}
	_, found = k.GetAsset(ctx, msg.AssetOut)
	if !found {
		return types.ErrorAssetDoesNotExist
	}

	if msg.AssetIn == msg.AssetOut {
		return types.ErrorAssetsCanNotBeSame
	}

	pair.AssetIn = msg.AssetIn
	pair.AssetOut = msg.AssetOut
	pair.MinUsdValueLeft = msg.MinUsdValueLeft

	k.SetLendPair(ctx, pair)
	return nil
}

func (k Keeper) AddPoolRecords(ctx sdk.Context, pool types.Pool) error {
	for _, v := range pool.AssetData {
		_, found := k.GetAsset(ctx, v.AssetID)
		if !found {
			return types.ErrorAssetDoesNotExist
		}
	}

	poolID := k.GetPoolID(ctx)
	newPool := types.Pool{
		PoolID:       poolID + 1,
		ModuleName:   pool.ModuleName,
		CPoolName:    pool.CPoolName,
		ReserveFunds: pool.ReserveFunds,
		AssetData:    pool.AssetData,
	}
	for _, v := range pool.AssetData {
		var assetStats types.PoolAssetLBMapping
		assetStats.PoolID = newPool.PoolID
		assetStats.AssetID = v.AssetID
		assetStats.TotalBorrowed = sdk.ZeroInt()
		assetStats.TotalStableBorrowed = sdk.ZeroInt()
		assetStats.TotalLend = sdk.ZeroInt()
		assetStats.TotalInterestAccumulated = sdk.ZeroInt()
		k.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
		k.UpdateAPR(ctx, newPool.PoolID, v.AssetID)
		reserveBuybackStats, found := k.GetReserveBuybackAssetData(ctx, v.AssetID)
		if !found {
			reserveBuybackStats.AssetID = v.AssetID
			reserveBuybackStats.ReserveAmount = sdk.ZeroInt()
			reserveBuybackStats.BuybackAmount = sdk.ZeroInt()
			k.SetReserveBuybackAssetData(ctx, reserveBuybackStats)
		}
	}

	k.SetPool(ctx, newPool)
	k.SetPoolID(ctx, newPool.PoolID)
	return nil
}

func (k Keeper) AddAssetToPair(ctx sdk.Context, assetToPair types.AssetToPairMapping) error {
	_, found := k.GetAsset(ctx, assetToPair.AssetID)
	if !found {
		return types.ErrorAssetDoesNotExist
	}
	_, found = k.GetPool(ctx, assetToPair.PoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	for _, v := range assetToPair.PairID {
		_, found = k.GetLendPair(ctx, v)
		if !found {
			return types.ErrorPairDoesNotExist
		}
	}
	k.SetAssetToPair(ctx, assetToPair)

	return nil
}

func (k Keeper) SetLendPairID(ctx sdk.Context, id uint64) {
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

func (k Keeper) SetLendPair(ctx sdk.Context, pair types.Extended_Pair) {
	var (
		store = k.Store(ctx)
		key   = types.LendPairKey(pair.Id)
		value = k.cdc.MustMarshal(&pair)
	)

	store.Set(key, value)
}

func (k Keeper) GetLendPair(ctx sdk.Context, id uint64) (pair types.Extended_Pair, found bool) {
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

func (k Keeper) GetLendPairs(ctx sdk.Context) (pairs []types.Extended_Pair) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.LendPairKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var pair types.Extended_Pair
		k.cdc.MustUnmarshal(iter.Value(), &pair)
		pairs = append(pairs, pair)
	}

	return pairs
}

func (k Keeper) GetLendPairID(ctx sdk.Context) uint64 {
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

func (k Keeper) AddAssetRatesParams(ctx sdk.Context, records ...types.AssetRatesParams) error {
	for _, msg := range records {
		_, found := k.GetAssetRatesParams(ctx, msg.AssetID)
		if found {
			return types.ErrorAssetRatesParamsNotFound
		}

		assetRatesParams := types.AssetRatesParams{
			AssetID:              msg.AssetID,
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
			CAssetID:             msg.CAssetID,
		}

		k.SetAssetRatesParams(ctx, assetRatesParams)
	}
	return nil
}

func (k Keeper) AddAuctionParamsData(ctx sdk.Context, param types.AuctionParams) error {
	var (
		store = k.Store(ctx)
		key   = types.AuctionParamKey(param.AppId)
		value = k.cdc.MustMarshal(&param)
	)

	store.Set(key, value)

	return nil
}

func (k Keeper) GetAddAuctionParamsData(ctx sdk.Context, appID uint64) (auctionParams types.AuctionParams, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AuctionParamKey(appID)
		value = store.Get(key)
	)

	if value == nil {
		return auctionParams, false
	}

	k.cdc.MustUnmarshal(value, &auctionParams)
	return auctionParams, true
}

func (k Keeper) GetAllAddAuctionParamsData(ctx sdk.Context) (auctionParams []types.AuctionParams) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.AuctionParamPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var asset types.AuctionParams
		k.cdc.MustUnmarshal(iter.Value(), &asset)
		auctionParams = append(auctionParams, asset)
	}
	return auctionParams
}

func (k Keeper) SetAssetRatesParams(ctx sdk.Context, assetRatesParams types.AssetRatesParams) {
	var (
		store = k.Store(ctx)
		key   = types.AssetRatesParamsKey(assetRatesParams.AssetID)
		value = k.cdc.MustMarshal(&assetRatesParams)
	)

	store.Set(key, value)
}

func (k Keeper) GetAssetRatesParams(ctx sdk.Context, assetID uint64) (assetRatesParams types.AssetRatesParams, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AssetRatesParamsKey(assetID)
		value = store.Get(key)
	)

	if value == nil {
		return assetRatesParams, false
	}

	k.cdc.MustUnmarshal(value, &assetRatesParams)
	return assetRatesParams, true
}

func (k Keeper) GetAllAssetRatesParams(ctx sdk.Context) (assetRatesParams []types.AssetRatesParams) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.AssetRatesParamsKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var asset types.AssetRatesParams
		k.cdc.MustUnmarshal(iter.Value(), &asset)
		assetRatesParams = append(assetRatesParams, asset)
	}
	return assetRatesParams
}
