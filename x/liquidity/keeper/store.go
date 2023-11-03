package keeper

import (
	gogotypes "github.com/cosmos/gogoproto/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/liquidity/types"
)

// GetLastPairID returns the last pair id.
func (k Keeper) GetLastPairID(ctx sdk.Context, appID uint64) (id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLastPairIDKey(appID))
	if bz == nil {
		id = 0 // initialize the pair id
	} else {
		var val gogotypes.UInt64Value
		k.cdc.MustUnmarshal(bz, &val)
		id = val.GetValue()
	}
	return
}

// SetLastPairID stores the last pair id.
func (k Keeper) SetLastPairID(ctx sdk.Context, appID, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: id})
	store.Set(types.GetLastPairIDKey(appID), bz)
}

// GetPair returns pair object for the given pair id.
func (k Keeper) GetPair(ctx sdk.Context, appID, id uint64) (pair types.Pair, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPairKey(appID, id))
	if bz == nil {
		return
	}
	pair = types.MustUnmarshalPair(k.cdc, bz)
	return pair, true
}

// GetPairByDenoms returns a types.Pair for given denoms.
func (k Keeper) GetPairByDenoms(ctx sdk.Context, appID uint64, baseCoinDenom, quoteCoinDenom string) (pair types.Pair, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPairIndexKey(appID, baseCoinDenom, quoteCoinDenom))
	if bz == nil {
		return
	}
	var val gogotypes.UInt64Value
	k.cdc.MustUnmarshal(bz, &val)
	pair, found = k.GetPair(ctx, appID, val.Value)
	return
}

// SetPair stores the particular pair.
func (k Keeper) SetPair(ctx sdk.Context, pair types.Pair) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalPair(k.cdc, pair)
	store.Set(types.GetPairKey(pair.AppId, pair.Id), bz)
}

// SetPairIndex stores a pair index.
func (k Keeper) SetPairIndex(ctx sdk.Context, appID uint64, baseCoinDenom, quoteCoinDenom string, pairID uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: pairID})
	store.Set(types.GetPairIndexKey(appID, baseCoinDenom, quoteCoinDenom), bz)
}

// SetPairLookupIndex stores a pair lookup index for given denoms.
func (k Keeper) SetPairLookupIndex(ctx sdk.Context, appID uint64, denomA string, denomB string, pairID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetPairsByDenomsIndexKey(appID, denomA, denomB, pairID), []byte{})
}

// IterateAllPairs iterates over all the stored pairs and performs a callback function.
// Stops iteration when callback returns true.
func (k Keeper) IterateAllPairs(ctx sdk.Context, appID uint64, cb func(pair types.Pair) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAllPairsKey(appID))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		pair := types.MustUnmarshalPair(k.cdc, iter.Value())
		stop, err := cb(pair)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

// GetAllPairs returns all pairs in the store.
func (k Keeper) GetAllPairs(ctx sdk.Context, appID uint64) (pairs []types.Pair) {
	pairs = []types.Pair{}
	_ = k.IterateAllPairs(ctx, appID, func(pair types.Pair) (stop bool, err error) {
		pairs = append(pairs, pair)
		return false, nil
	})
	return pairs
}

// GetLastPoolID returns the last pool id.
func (k Keeper) GetLastPoolID(ctx sdk.Context, appID uint64) (id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLastPoolIDKey(appID))
	if bz == nil {
		id = 0 // initialize the pool id
	} else {
		var val gogotypes.UInt64Value
		k.cdc.MustUnmarshal(bz, &val)
		id = val.GetValue()
	}
	return
}

// SetLastPoolID stores the last pool id.
func (k Keeper) SetLastPoolID(ctx sdk.Context, appID, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: id})
	store.Set(types.GetLastPoolIDKey(appID), bz)
}

// GetPool returns pool object for the given pool id.
func (k Keeper) GetPool(ctx sdk.Context, appID, id uint64) (pool types.Pool, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPoolKey(appID, id))
	if bz == nil {
		return
	}
	pool = types.MustUnmarshalPool(k.cdc, bz)
	return pool, true
}

// GetPoolByReserveAddress returns pool object for the given reserve account address.
func (k Keeper) GetPoolByReserveAddress(ctx sdk.Context, appID uint64, reserveAddr sdk.AccAddress) (pool types.Pool, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetPoolByReserveAddressIndexKey(appID, reserveAddr))
	if bz == nil {
		return
	}
	var val gogotypes.UInt64Value
	k.cdc.MustUnmarshal(bz, &val)
	poolID := val.GetValue()
	return k.GetPool(ctx, appID, poolID)
}

// SetPool stores the particular pool.
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalPool(k.cdc, pool)
	store.Set(types.GetPoolKey(pool.AppId, pool.Id), bz)
}

// SetPoolByReserveIndex stores a pool by reserve account index key.
func (k Keeper) SetPoolByReserveIndex(ctx sdk.Context, pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: pool.Id})
	store.Set(types.GetPoolByReserveAddressIndexKey(pool.AppId, pool.GetReserveAddress()), bz)
}

// SetPoolsByPairIndex stores a pool by pair index key.
func (k Keeper) SetPoolsByPairIndex(ctx sdk.Context, pool types.Pool) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetPoolsByPairIndexKey(pool.AppId, pool.PairId, pool.Id), []byte{})
}

// IterateAllPools iterates over all the stored pools and performs a callback function.
// Stops iteration when callback returns true.
func (k Keeper) IterateAllPools(ctx sdk.Context, appID uint64, cb func(pool types.Pool) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAllPoolsKey(appID))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		pool := types.MustUnmarshalPool(k.cdc, iter.Value())
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

// IteratePoolsByPair iterates over all the stored pools by the pair and performs a callback function.
// Stops iteration when callback returns true.
func (k Keeper) IteratePoolsByPair(ctx sdk.Context, appID, pairID uint64, cb func(pool types.Pool) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetPoolsByPairIndexKeyPrefix(appID, pairID))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		poolID := types.ParsePoolsByPairIndexKey(iter.Key())
		pool, _ := k.GetPool(ctx, appID, poolID)
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

// GetAllPools returns all pools in the store.
func (k Keeper) GetAllPools(ctx sdk.Context, appID uint64) (pools []types.Pool) {
	pools = []types.Pool{}
	_ = k.IterateAllPools(ctx, appID, func(pool types.Pool) (stop bool, err error) {
		pools = append(pools, pool)
		return false, nil
	})
	return
}

// GetPoolsByPair returns pools within the pair.
func (k Keeper) GetPoolsByPair(ctx sdk.Context, appID, pairID uint64) (pools []types.Pool) {
	_ = k.IteratePoolsByPair(ctx, appID, pairID, func(pool types.Pool) (stop bool, err error) {
		pools = append(pools, pool)
		return false, nil
	})
	return
}

// GetDepositRequest returns the particular deposit request.
func (k Keeper) GetDepositRequest(ctx sdk.Context, appID, poolID, id uint64) (req types.DepositRequest, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetDepositRequestKey(appID, poolID, id))
	if bz == nil {
		return
	}
	req = types.MustUnmarshalDepositRequest(k.cdc, bz)
	return req, true
}

// SetDepositRequest stores deposit request for the batch execution.
func (k Keeper) SetDepositRequest(ctx sdk.Context, req types.DepositRequest) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalDepositRequest(k.cdc, req)
	store.Set(types.GetDepositRequestKey(req.AppId, req.PoolId, req.Id), bz)
}

func (k Keeper) SetDepositRequestIndex(ctx sdk.Context, req types.DepositRequest) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetDepositRequestIndexKey(req.AppId, req.GetDepositor(), req.PoolId, req.Id), []byte{})
}

// IterateAllDepositRequests iterates through all deposit requests in the store
// and call cb for each request.
func (k Keeper) IterateAllDepositRequests(ctx sdk.Context, appID uint64, cb func(req types.DepositRequest) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAllDepositRequestKey(appID))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		req := types.MustUnmarshalDepositRequest(k.cdc, iter.Value())
		stop, err := cb(req)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

// IterateDepositRequestsByDepositor iterates through deposit requests in the
// store by a depositor and call cb on each order.
func (k Keeper) IterateDepositRequestsByDepositor(
	ctx sdk.Context,
	appID uint64,

	depositor sdk.AccAddress,
	cb func(req types.DepositRequest) (stop bool, err error),
) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetDepositRequestIndexKeyPrefix(appID, depositor))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		_, poolID, reqID := types.ParseDepositRequestIndexKey(iter.Key())
		req, _ := k.GetDepositRequest(ctx, appID, poolID, reqID)
		stop, err := cb(req)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

// GetAllDepositRequests returns all deposit requests in the store.
func (k Keeper) GetAllDepositRequests(ctx sdk.Context, appID uint64) (reqs []types.DepositRequest) {
	reqs = []types.DepositRequest{}
	_ = k.IterateAllDepositRequests(ctx, appID, func(req types.DepositRequest) (stop bool, err error) {
		reqs = append(reqs, req)
		return false, nil
	})
	return
}

// GetDepositRequestsByDepositor returns deposit requests by the depositor.
func (k Keeper) GetDepositRequestsByDepositor(ctx sdk.Context, appID uint64, depositor sdk.AccAddress) (reqs []types.DepositRequest) {
	_ = k.IterateDepositRequestsByDepositor(ctx, appID, depositor, func(req types.DepositRequest) (stop bool, err error) {
		reqs = append(reqs, req)
		return false, nil
	})
	return
}

// DeleteDepositRequest deletes a deposit request.
func (k Keeper) DeleteDepositRequest(ctx sdk.Context, req types.DepositRequest) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDepositRequestKey(req.AppId, req.PoolId, req.Id))
	k.DeleteDepositRequestIndex(ctx, req)
}

func (k Keeper) DeleteDepositRequestIndex(ctx sdk.Context, req types.DepositRequest) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDepositRequestIndexKey(req.AppId, req.GetDepositor(), req.PoolId, req.Id))
}

// GetWithdrawRequest returns the particular withdraw request.
func (k Keeper) GetWithdrawRequest(ctx sdk.Context, appID, poolID, id uint64) (req types.WithdrawRequest, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetWithdrawRequestKey(appID, poolID, id))
	if bz == nil {
		return
	}
	req = types.MustUnmarshalWithdrawRequest(k.cdc, bz)
	return req, true
}

// SetWithdrawRequest stores withdraw request for the batch execution.
func (k Keeper) SetWithdrawRequest(ctx sdk.Context, req types.WithdrawRequest) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshaWithdrawRequest(k.cdc, req)
	store.Set(types.GetWithdrawRequestKey(req.AppId, req.PoolId, req.Id), bz)
}

func (k Keeper) SetWithdrawRequestIndex(ctx sdk.Context, req types.WithdrawRequest) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetWithdrawRequestIndexKey(req.AppId, req.GetWithdrawer(), req.PoolId, req.Id), []byte{})
}

// IterateAllWithdrawRequests iterates through all withdraw requests in the store
// and call cb for each request.
func (k Keeper) IterateAllWithdrawRequests(ctx sdk.Context, appID uint64, cb func(req types.WithdrawRequest) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAllWithdrawRequestKey(appID))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		req := types.MustUnmarshalWithdrawRequest(k.cdc, iter.Value())
		stop, err := cb(req)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

// IterateWithdrawRequestsByWithdrawer iterates through withdraw requests in the
// store by a withdrawer and call cb on each order.
func (k Keeper) IterateWithdrawRequestsByWithdrawer(
	ctx sdk.Context,
	appID uint64,

	withdrawer sdk.AccAddress,
	cb func(req types.WithdrawRequest) (stop bool, err error),
) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetWithdrawRequestIndexKeyPrefix(appID, withdrawer))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		_, poolID, reqID := types.ParseWithdrawRequestIndexKey(iter.Key())
		req, _ := k.GetWithdrawRequest(ctx, appID, poolID, reqID)
		stop, err := cb(req)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

// GetAllWithdrawRequests returns all withdraw requests in the store.
func (k Keeper) GetAllWithdrawRequests(ctx sdk.Context, appID uint64) (reqs []types.WithdrawRequest) {
	reqs = []types.WithdrawRequest{}
	_ = k.IterateAllWithdrawRequests(ctx, appID, func(req types.WithdrawRequest) (stop bool, err error) {
		reqs = append(reqs, req)
		return false, nil
	})
	return
}

// GetWithdrawRequestsByWithdrawer returns withdraw requests by the withdrawer.
func (k Keeper) GetWithdrawRequestsByWithdrawer(ctx sdk.Context, appID uint64, withdrawer sdk.AccAddress) (reqs []types.WithdrawRequest) {
	_ = k.IterateWithdrawRequestsByWithdrawer(ctx, appID, withdrawer, func(req types.WithdrawRequest) (stop bool, err error) {
		reqs = append(reqs, req)
		return false, nil
	})
	return
}

// DeleteWithdrawRequest deletes a withdraw request.
func (k Keeper) DeleteWithdrawRequest(ctx sdk.Context, req types.WithdrawRequest) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetWithdrawRequestKey(req.AppId, req.PoolId, req.Id))
	k.DeleteWithdrawRequestIndex(ctx, req)
}

func (k Keeper) DeleteWithdrawRequestIndex(ctx sdk.Context, req types.WithdrawRequest) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetWithdrawRequestIndexKey(req.AppId, req.GetWithdrawer(), req.PoolId, req.Id))
}

// GetOrder returns the particular order.
func (k Keeper) GetOrder(ctx sdk.Context, appID, pairID, id uint64) (order types.Order, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOrderKey(appID, pairID, id))
	if bz == nil {
		return
	}
	order = types.MustUnmarshalOrder(k.cdc, bz)
	return order, true
}

// SetOrder stores an order for the batch execution.
func (k Keeper) SetOrder(ctx sdk.Context, appID uint64, order types.Order) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshaOrder(k.cdc, order)
	store.Set(types.GetOrderKey(appID, order.PairId, order.Id), bz)
}

func (k Keeper) SetOrderIndex(ctx sdk.Context, appID uint64, order types.Order) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOrderIndexKey(appID, order.GetOrderer(), order.PairId, order.Id), []byte{})
}

// IterateAllOrders iterates through all orders in the store and all
// cb for each order.
func (k Keeper) IterateAllOrders(ctx sdk.Context, appID uint64, cb func(order types.Order) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAllOrdersKey(appID))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		order := types.MustUnmarshalOrder(k.cdc, iter.Value())
		stop, err := cb(order)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

// IterateOrdersByPair iterates through all the orders within the pair
// and call cb for each order.
func (k Keeper) IterateOrdersByPair(ctx sdk.Context, appID, pairID uint64, cb func(order types.Order) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetOrdersByPairKeyPrefix(appID, pairID))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		order := types.MustUnmarshalOrder(k.cdc, iter.Value())
		stop, err := cb(order)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

// IterateOrdersByOrderer iterates through orders in the store by an orderer
// and call cb on each order.
func (k Keeper) IterateOrdersByOrderer(
	ctx sdk.Context,
	appID uint64,

	orderer sdk.AccAddress,
	cb func(order types.Order) (stop bool, err error),
) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetOrderIndexKeyPrefix(appID, orderer))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		_, pairID, orderID := types.ParseOrderIndexKey(iter.Key())
		order, _ := k.GetOrder(ctx, appID, pairID, orderID)
		stop, err := cb(order)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

// GetAllOrders returns all orders in the store.
func (k Keeper) GetAllOrders(ctx sdk.Context, appID uint64) (orders []types.Order) {
	orders = []types.Order{}
	_ = k.IterateAllOrders(ctx, appID, func(order types.Order) (stop bool, err error) {
		orders = append(orders, order)
		return false, nil
	})
	return
}

// GetOrdersByPair returns orders within the pair.
func (k Keeper) GetOrdersByPair(ctx sdk.Context, appID, pairID uint64) (orders []types.Order) {
	_ = k.IterateOrdersByPair(ctx, appID, pairID, func(order types.Order) (stop bool, err error) {
		orders = append(orders, order)
		return false, nil
	})
	return
}

// GetOrdersByOrderer returns orders by the orderer.
func (k Keeper) GetOrdersByOrderer(ctx sdk.Context, appID uint64, orderer sdk.AccAddress) (orders []types.Order) {
	_ = k.IterateOrdersByOrderer(ctx, appID, orderer, func(order types.Order) (stop bool, err error) {
		orders = append(orders, order)
		return false, nil
	})
	return
}

// DeleteOrder deletes an order.
func (k Keeper) DeleteOrder(ctx sdk.Context, appID uint64, order types.Order) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOrderKey(appID, order.PairId, order.Id))
	k.DeleteOrderIndex(ctx, appID, order)
}

func (k Keeper) DeleteOrderIndex(ctx sdk.Context, appID uint64, order types.Order) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOrderIndexKey(appID, order.GetOrderer(), order.PairId, order.Id))
}

// GetMMOrderIndex returns the market making order index.
func (k Keeper) GetMMOrderIndex(ctx sdk.Context, orderer sdk.AccAddress, appID, pairID uint64) (index types.MMOrderIndex, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMMOrderIndexKey(orderer, appID, pairID))
	if bz == nil {
		return
	}
	index = types.MustUnmarshalMMOrderIndex(k.cdc, bz)
	return index, true
}

// SetMMOrderIndex stores a market making order index.
func (k Keeper) SetMMOrderIndex(ctx sdk.Context, appID uint64, index types.MMOrderIndex) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshaMMOrderIndex(k.cdc, index)
	store.Set(types.GetMMOrderIndexKey(index.GetOrderer(), appID, index.PairId), bz)
}

// IterateAllMMOrderIndexes iterates through all market making order indexes
// in the store and call cb for each order.
func (k Keeper) IterateAllMMOrderIndexes(ctx sdk.Context, appID uint64, cb func(index types.MMOrderIndex) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAllMMOrderIndexKey(appID))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		index := types.MustUnmarshalMMOrderIndex(k.cdc, iter.Value())
		stop, err := cb(index)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return nil
}

// GetAllMMOrderIndexes returns all market making order indexes in the store.
func (k Keeper) GetAllMMOrderIndexes(ctx sdk.Context, appID uint64) (indexes []types.MMOrderIndex) {
	indexes = []types.MMOrderIndex{}
	_ = k.IterateAllMMOrderIndexes(ctx, appID, func(index types.MMOrderIndex) (stop bool, err error) {
		indexes = append(indexes, index)
		return false, nil
	})
	return
}

// DeleteMMOrderIndex deletes a market making order index.
func (k Keeper) DeleteMMOrderIndex(ctx sdk.Context, appID uint64, index types.MMOrderIndex) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetMMOrderIndexKey(index.GetOrderer(), appID, index.PairId))
}

// GetGenericLiquidityParams returns the generic liquidity params by app id.
func (k Keeper) GetGenericLiquidityParams(ctx sdk.Context, appID uint64) (genericLiquidityParams types.GenericParams, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetGenericParamsKey(appID))
	if bz == nil {
		return
	}
	genericLiquidityParams = types.MustUnmarshalGenericLiquidityParams(k.cdc, bz)
	return genericLiquidityParams, true
}

// SetGenericLiquidityParams sets the the generic liquidity params by app id.
func (k Keeper) SetGenericLiquidityParams(ctx sdk.Context, genericLiquidityParams types.GenericParams) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalGenericLiquidityParams(k.cdc, genericLiquidityParams)
	store.Set(types.GetGenericParamsKey(genericLiquidityParams.AppId), bz)
}

// SetActiveFarmer stores the active farmer.
func (k Keeper) SetActiveFarmer(ctx sdk.Context, activeFarmer types.ActiveFarmer) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalActiveFarmer(k.cdc, activeFarmer)
	store.Set(types.GetActiveFarmerKey(activeFarmer.AppId, activeFarmer.PoolId, sdk.MustAccAddressFromBech32(activeFarmer.Farmer)), bz)
}

// DeleteActiveFarmer deletes a Active farmer from store.
func (k Keeper) DeleteActiveFarmer(ctx sdk.Context, activeFarmer types.ActiveFarmer) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetActiveFarmerKey(activeFarmer.AppId, activeFarmer.PoolId, sdk.MustAccAddressFromBech32(activeFarmer.Farmer)))
}

// GetActiveFarmer returns active farmer object for the given app id, pool id and farmer.
func (k Keeper) GetActiveFarmer(ctx sdk.Context, appID, poolID uint64, farmer sdk.AccAddress) (activeFarmer types.ActiveFarmer, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetActiveFarmerKey(appID, poolID, farmer))
	if bz == nil {
		return
	}
	activeFarmer = types.MustUnmarshalActiveFarmer(k.cdc, bz)
	return activeFarmer, true
}

// IterateAllActiveFarmers iterates over all the stored active farmers and performs a callback function.
// Stops iteration when callback returns true.
func (k Keeper) IterateAllActiveFarmers(ctx sdk.Context, appID, poolID uint64, cb func(activeFarmer types.ActiveFarmer) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAllActiveFarmersKey(appID, poolID))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		activeFarmer := types.MustUnmarshalActiveFarmer(k.cdc, iter.Value())
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

// GetAllActiveFarmers returns all active farmers in the store.
func (k Keeper) GetAllActiveFarmers(ctx sdk.Context, appID, poolID uint64) (activeFarmers []types.ActiveFarmer) {
	activeFarmers = []types.ActiveFarmer{}
	_ = k.IterateAllActiveFarmers(ctx, appID, poolID, func(activeFarmer types.ActiveFarmer) (stop bool, err error) {
		activeFarmers = append(activeFarmers, activeFarmer)
		return false, nil
	})
	return activeFarmers
}

// SetQueuedFarmer stores the queued farmer.
func (k Keeper) SetQueuedFarmer(ctx sdk.Context, queuedFarmer types.QueuedFarmer) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalQueuedFarmer(k.cdc, queuedFarmer)
	store.Set(types.GetQueuedFarmerKey(queuedFarmer.AppId, queuedFarmer.PoolId, sdk.MustAccAddressFromBech32(queuedFarmer.Farmer)), bz)
}

// GetQueuedFarmer returns queued farmer object for the given app id, pool id and farmer.
func (k Keeper) GetQueuedFarmer(ctx sdk.Context, appID, poolID uint64, farmer sdk.AccAddress) (queuedFarmer types.QueuedFarmer, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetQueuedFarmerKey(appID, poolID, farmer))
	if bz == nil {
		return
	}
	queuedFarmer = types.MustUnmarshalQueuedFarmer(k.cdc, bz)
	return queuedFarmer, true
}

// IterateAllQueuedFarmers iterates over all the stored queued farmers and performs a callback function.
// Stops iteration when callback returns true.
func (k Keeper) IterateAllQueuedFarmers(ctx sdk.Context, appID, poolID uint64, cb func(queuedFarmer types.QueuedFarmer) (stop bool, err error)) error {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetAllQueuedFarmersKey(appID, poolID))
	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)
	for ; iter.Valid(); iter.Next() {
		queuedFarmer := types.MustUnmarshalQueuedFarmer(k.cdc, iter.Value())
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

// GetAllQueuedFarmers returns all queued farmers in the store.
func (k Keeper) GetAllQueuedFarmers(ctx sdk.Context, appID, poolID uint64) (queuedFarmers []types.QueuedFarmer) {
	queuedFarmers = []types.QueuedFarmer{}
	_ = k.IterateAllQueuedFarmers(ctx, appID, poolID, func(queuedFarmer types.QueuedFarmer) (stop bool, err error) {
		queuedFarmers = append(queuedFarmers, queuedFarmer)
		return false, nil
	})
	return queuedFarmers
}

// GetAllCMSTPools returns all pools in the store.
func (k Keeper) GetAllCMSTPools(ctx sdk.Context, appID uint64) (pools []types.Pool) {
	cmstPools := []types.Pool{}
	allPools := k.GetAllPools(ctx, appID)
	for _, pool := range allPools {
		if !pool.Disabled {
			pair, found := k.GetPair(ctx, appID, pool.PairId)
			if found {
				if pair.BaseCoinDenom == "ucmst" || pair.QuoteCoinDenom == "ucmst" {
					cmstPools = append(cmstPools, pool)
				}
			}
		}
	}
	return cmstPools
}
