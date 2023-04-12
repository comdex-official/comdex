package v3

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	expected "github.com/comdex-official/comdex/x/liquidity/expected"
	v2liquidity "github.com/comdex-official/comdex/x/liquidity/legacy/v2"
	"github.com/comdex-official/comdex/x/liquidity/types"
)

func GetPair(
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	appID, pairID uint64,
) (pair types.Pair, found bool) {
	bz := store.Get(types.GetPairKey(appID, pairID))
	if bz == nil {
		return
	}
	pair = types.MustUnmarshalPair(cdc, bz)
	return pair, true
}

func IterateAllPools(
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	appID uint64, cb func(pool types.Pool) (stop bool, err error),
) error {
	iter := sdk.KVStorePrefixIterator(store, types.GetAllPoolsKey(appID))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		pool := types.MustUnmarshalPool(cdc, iter.Value())
		stop, err := cb(pool)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

func GetAllPools(
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	appID uint64,
) (pools []types.Pool) {
	pools = []types.Pool{}
	_ = IterateAllPools(store, cdc, appID, func(pool types.Pool) (stop bool, err error) {
		pools = append(pools, pool)
		return false, nil
	})
	return
}

func IterateAllActiveFarmers(
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	appID, poolID uint64,
	cb func(activeFarmer types.ActiveFarmer) (stop bool, err error),
) error {
	iter := sdk.KVStorePrefixIterator(store, types.GetAllActiveFarmersKey(appID, poolID))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		activeFarmer := types.MustUnmarshalActiveFarmer(cdc, iter.Value())
		stop, err := cb(activeFarmer)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

func GetAllActiveFarmers(
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	appID, poolID uint64,
) (activeFarmers []types.ActiveFarmer) {
	activeFarmers = []types.ActiveFarmer{}
	_ = IterateAllActiveFarmers(store, cdc, appID, poolID, func(activeFarmer types.ActiveFarmer) (stop bool, err error) {
		activeFarmers = append(activeFarmers, activeFarmer)
		return false, nil
	})
	return activeFarmers
}

func IterateAllQueuedFarmers(
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	appID, poolID uint64,
	cb func(queuedFarmer types.QueuedFarmer) (stop bool, err error),
) error {
	iter := sdk.KVStorePrefixIterator(store, types.GetAllQueuedFarmersKey(appID, poolID))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		queuedFarmer := types.MustUnmarshalQueuedFarmer(cdc, iter.Value())
		stop, err := cb(queuedFarmer)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

func GetAllQueuedFarmers(
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	appID, poolID uint64,
) (queuedFarmers []types.QueuedFarmer) {
	queuedFarmers = []types.QueuedFarmer{}
	_ = IterateAllQueuedFarmers(store, cdc, appID, poolID, func(queuedFarmer types.QueuedFarmer) (stop bool, err error) {
		queuedFarmers = append(queuedFarmers, queuedFarmer)
		return false, nil
	})
	return queuedFarmers
}

func MigratePools(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	assetKeeper expected.AssetKeeper,
	appID uint64,
) error {
	iter := sdk.KVStorePrefixIterator(store, types.GetAllPoolsKey(appID))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		var oldPool v2liquidity.Pool
		if err := cdc.Unmarshal(iter.Value(), &oldPool); err != nil {
			return err
		}
		pair, found := GetPair(store, cdc, appID, oldPool.PairId)
		if !found {
			return fmt.Errorf("pair %d not found in app %d", oldPool.PairId, appID)
		}

		baseAsset, found := assetKeeper.GetAssetForDenom(ctx, pair.BaseCoinDenom)
		if !found {
			return fmt.Errorf("baseAsset %s not found", pair.BaseCoinDenom)
		}
		quoteAsset, found := assetKeeper.GetAssetForDenom(ctx, pair.QuoteCoinDenom)
		if !found {
			return fmt.Errorf("quoteAsset %s not found", pair.QuoteCoinDenom)
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
			Creator:               oldPool.Creator,
			MinPrice:              oldPool.MinPrice,
			MaxPrice:              oldPool.MaxPrice,
			FarmCoin: &types.FarmCoin{
				Denom:    types.FarmCoinDenom(appID, oldPool.Id),
				Decimals: types.FarmCoinDecimals(baseAsset.Decimals.Uint64(), quoteAsset.Decimals.Uint64()),
			},
		}
		bz, err := cdc.Marshal(&newPool)
		if err != nil {
			return err
		}
		store.Set(iter.Key(), bz)
	}
	return nil
}

func MintAndSendHelper(ctx sdk.Context, bankKeeper expected.BankKeeper, farmer sdk.AccAddress, farmCoin sdk.Coin) error {
	if err := bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(farmCoin)); err != nil {
		return err
	}
	if err := bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, farmer, sdk.NewCoins(farmCoin)); err != nil {
		return err
	}
	return nil
}

func MintFarmTokenAndTransferToAccounts(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
	bankKeeper expected.BankKeeper,
	appID uint64,
) error {
	allPools := GetAllPools(store, cdc, appID)
	for _, pool := range allPools {
		if !pool.Disabled {
			allActiveFarmers := GetAllActiveFarmers(store, cdc, appID, pool.Id)
			for _, afarmer := range allActiveFarmers {
				farmer := sdk.MustAccAddressFromBech32(afarmer.Farmer)
				farmCoin := sdk.NewCoin(pool.FarmCoin.Denom, afarmer.FarmedPoolCoin.Amount)
				if farmCoin.Amount.IsPositive() {
					if err := MintAndSendHelper(ctx, bankKeeper, farmer, farmCoin); err != nil {
						return err
					}
				}
			}

			allQueuedFarmer := GetAllQueuedFarmers(store, cdc, appID, pool.Id)
			for _, qfarmer := range allQueuedFarmer {
				farmer := sdk.MustAccAddressFromBech32(qfarmer.Farmer)
				farmCoin := sdk.NewCoin(pool.FarmCoin.Denom, sdk.ZeroInt())
				for _, qCoin := range qfarmer.QueudCoins {
					farmCoin.Amount = farmCoin.Amount.Add(qCoin.FarmedPoolCoin.Amount)
				}
				if farmCoin.Amount.IsPositive() {
					if err := MintAndSendHelper(ctx, bankKeeper, farmer, farmCoin); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func MigrateStore(
	ctx sdk.Context,
	assetKeeper expected.AssetKeeper,
	bankKeeper expected.BankKeeper,
	storeKey sdk.StoreKey,
	cdc codec.BinaryCodec,
) error {
	allApps, found := assetKeeper.GetApps(ctx)
	cacheCtx, writeCache := ctx.CacheContext()
	if found {
		for _, app := range allApps {
			store := cacheCtx.KVStore(storeKey)
			if err := MigratePools(cacheCtx, store, cdc, assetKeeper, app.Id); err != nil {
				panic(err)
			}
			if err := MintFarmTokenAndTransferToAccounts(cacheCtx, store, cdc, bankKeeper, app.Id); err != nil {
				panic(err)
			}
		}
	}
	writeCache()
	return nil
}
