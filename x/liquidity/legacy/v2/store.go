package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	expected "github.com/comdex-official/comdex/x/liquidity/expected"
	v1liquidity "github.com/comdex-official/comdex/x/liquidity/legacy/v1"
	"github.com/comdex-official/comdex/x/liquidity/types"
)

func MigrateGenericParams(appID uint64, store sdk.KVStore, cdc codec.BinaryCodec) error {
	var oldGenericLiquidityParams v1liquidity.GenericParams
	if err := cdc.Unmarshal(store.Get(types.GetGenericParamsKey(appID)), &oldGenericLiquidityParams); err != nil {
		return err
	}
	newGenericLiquidityParams := types.GenericParams{
		BatchSize:                    oldGenericLiquidityParams.BatchSize,
		TickPrecision:                types.DefaultTickPrecision,
		FeeCollectorAddress:          oldGenericLiquidityParams.FeeCollectorAddress,
		DustCollectorAddress:         oldGenericLiquidityParams.DustCollectorAddress,
		MinInitialPoolCoinSupply:     oldGenericLiquidityParams.MinInitialPoolCoinSupply,
		PairCreationFee:              oldGenericLiquidityParams.PairCreationFee,
		PoolCreationFee:              oldGenericLiquidityParams.PoolCreationFee,
		MinInitialDepositAmount:      oldGenericLiquidityParams.MinInitialDepositAmount,
		MaxPriceLimitRatio:           oldGenericLiquidityParams.MaxPriceLimitRatio,
		MaxOrderLifespan:             oldGenericLiquidityParams.MaxOrderLifespan,
		SwapFeeRate:                  oldGenericLiquidityParams.SwapFeeRate,
		WithdrawFeeRate:              oldGenericLiquidityParams.WithdrawFeeRate,
		DepositExtraGas:              oldGenericLiquidityParams.DepositExtraGas,
		WithdrawExtraGas:             oldGenericLiquidityParams.WithdrawExtraGas,
		OrderExtraGas:                oldGenericLiquidityParams.OrderExtraGas,
		SwapFeeDistrDenom:            oldGenericLiquidityParams.SwapFeeDistrDenom,
		SwapFeeBurnRate:              oldGenericLiquidityParams.SwapFeeBurnRate,
		AppId:                        oldGenericLiquidityParams.AppId,
		MaxNumMarketMakingOrderTicks: 10,
		MaxNumActivePoolsPerPair:     20,
	}
	bz, err := cdc.Marshal(&newGenericLiquidityParams)
	if err != nil {
		return err
	}
	store.Set(types.GetGenericParamsKey(appID), bz)

	return nil
}

func MigratePools(appID uint64, store sdk.KVStore, cdc codec.BinaryCodec) error {
	iter := sdk.KVStorePrefixIterator(store, types.GetAllPoolsKey(appID))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		var oldPool v1liquidity.Pool
		if err := cdc.Unmarshal(iter.Value(), &oldPool); err != nil {
			return err
		}
		newPool := types.Pool{
			Id:                    oldPool.Id,
			PairId:                oldPool.PairId,
			ReserveAddress:        oldPool.ReserveAddress,
			PoolCoinDenom:         oldPool.PoolCoinDenom,
			LastDepositRequestId:  oldPool.LastDepositRequestId,
			LastWithdrawRequestId: oldPool.LastWithdrawRequestId,
			Disabled:              oldPool.Disabled,
			AppId:                 oldPool.AppId,
			Type:                  types.PoolTypeBasic,
			Creator:               "",
			MinPrice:              nil,
			MaxPrice:              nil,
		}
		bz, err := cdc.Marshal(&newPool)
		if err != nil {
			return err
		}
		store.Set(iter.Key(), bz)
	}
	return nil
}

func MigrateOrders(appID uint64, store sdk.KVStore, cdc codec.BinaryCodec) error {
	iter := sdk.KVStorePrefixIterator(store, types.GetAllOrdersKey(appID))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		var oldOrder v1liquidity.Order
		if err := cdc.Unmarshal(iter.Value(), &oldOrder); err != nil {
			return err
		}
		newOrder := types.Order{
			Id:        oldOrder.Id,
			PairId:    oldOrder.PairId,
			MsgHeight: oldOrder.MsgHeight,
			Orderer:   oldOrder.Orderer,
			// Only the type has changed, not the value, so simply type-cast here
			Direction:          types.OrderDirection(oldOrder.Direction),
			OfferCoin:          oldOrder.OfferCoin,
			RemainingOfferCoin: oldOrder.RemainingOfferCoin,
			ReceivedCoin:       oldOrder.ReceivedCoin,
			Price:              oldOrder.Price,
			Amount:             oldOrder.Amount,
			OpenAmount:         oldOrder.OpenAmount,
			BatchId:            oldOrder.BatchId,
			ExpireAt:           oldOrder.ExpireAt,
			// Only the type has changed, not the value, so simply type-cast here
			Status: types.OrderStatus(oldOrder.Status),
			AppId:  appID,
			// There's no way to determine whether the order was made through
			// MsgLimitOrder or MsgMarketOrder, set the order type as OrderTypeLimit
			// as a fallback.
			Type: types.OrderTypeLimit,
		}
		bz, err := cdc.Marshal(&newOrder)
		if err != nil {
			return err
		}
		store.Set(iter.Key(), bz)
	}
	return nil
}

func MigrateStore(
	ctx sdk.Context,
	assetKeeper expected.AssetKeeper,
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
) error {
	allApps, found := assetKeeper.GetApps(ctx)
	cacheCtx, writeCache := ctx.CacheContext()
	if found {
		for _, app := range allApps {
			store := cacheCtx.KVStore(storeKey)
			if err := MigrateGenericParams(app.Id, store, cdc); err != nil {
				panic(err)
			}
			if err := MigratePools(app.Id, store, cdc); err != nil {
				panic(err)
			}
			if err := MigrateOrders(app.Id, store, cdc); err != nil {
				panic(err)
			}
		}
	}
	writeCache()
	return nil
}
