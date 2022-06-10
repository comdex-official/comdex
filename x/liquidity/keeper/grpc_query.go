package keeper

import (
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/comdex-official/comdex/x/liquidity/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper.
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// Params queries the parameters of the liquidity module.
func (k Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var params types.Params
	k.Keeper.paramSpace.GetParamSet(ctx, &params)
	return &types.QueryParamsResponse{Params: params}, nil
}

// Pools queries all pools.
func (k Querier) Pools(c context.Context, req *types.QueryPoolsRequest) (*types.QueryPoolsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id cannot be 0")
	}

	var disabled bool
	if req.Disabled != "" {
		var err error
		disabled, err = strconv.ParseBool(req.Disabled)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)

	var keyPrefix []byte
	var poolGetter func(key, value []byte) types.Pool
	var pairGetter func(id uint64) types.Pair
	pairMap := map[uint64]types.Pair{}
	switch {
	case req.PairId == 0:
		keyPrefix = types.GetAllPoolsKey(req.AppId)
		poolGetter = func(_, value []byte) types.Pool {
			return types.MustUnmarshalPool(k.cdc, value)
		}
		pairGetter = func(id uint64) types.Pair {
			pair, ok := pairMap[id]
			if !ok {
				pair, _ = k.GetPair(ctx, req.AppId, id)
				pairMap[id] = pair
			}
			return pair
		}
	default:
		keyPrefix = types.GetPoolsByPairIndexKeyPrefix(req.AppId, req.PairId)
		poolGetter = func(key, _ []byte) types.Pool {
			poolID := types.ParsePoolsByPairIndexKey(append(keyPrefix, key...))
			pool, _ := k.GetPool(ctx, req.AppId, poolID)
			return pool
		}
		pair, _ := k.GetPair(ctx, req.AppId, req.PairId)
		pairGetter = func(_ uint64) types.Pair {
			return pair
		}
	}

	poolStore := prefix.NewStore(store, keyPrefix)

	var poolsRes []types.PoolResponse
	pageRes, err := query.FilteredPaginate(poolStore, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		pool := poolGetter(key, value)
		if req.Disabled != "" {
			if pool.Disabled != disabled {
				return false, nil
			}
		}

		pair := pairGetter(pool.PairId)
		rx, ry := k.getPoolBalances(ctx, pool, pair)
		poolRes := types.PoolResponse{
			Id:                    pool.Id,
			PairId:                pool.PairId,
			ReserveAddress:        pool.ReserveAddress,
			PoolCoinDenom:         pool.PoolCoinDenom,
			Balances:              sdk.NewCoins(rx, ry),
			LastDepositRequestId:  pool.LastDepositRequestId,
			LastWithdrawRequestId: pool.LastWithdrawRequestId,
		}

		if accumulate {
			poolsRes = append(poolsRes, poolRes)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolsResponse{Pools: poolsRes, Pagination: pageRes}, nil
}

// Pool queries the specific pool.
func (k Querier) Pool(c context.Context, req *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id cannot be 0")
	}

	if req.PoolId == 0 {
		return nil, status.Error(codes.InvalidArgument, "pool id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	pool, found := k.GetPool(ctx, req.AppId, req.PoolId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "pool %d doesn't exist", req.PoolId)
	}

	rx, ry := k.GetPoolBalances(ctx, pool)
	poolRes := types.PoolResponse{
		Id:                    pool.Id,
		PairId:                pool.PairId,
		ReserveAddress:        pool.ReserveAddress,
		PoolCoinDenom:         pool.PoolCoinDenom,
		Balances:              sdk.NewCoins(rx, ry),
		LastDepositRequestId:  pool.LastDepositRequestId,
		LastWithdrawRequestId: pool.LastWithdrawRequestId,
	}

	return &types.QueryPoolResponse{Pool: poolRes}, nil
}

// PoolByReserveAddress queries the specific pool by the reserve account address.
func (k Querier) PoolByReserveAddress(c context.Context, req *types.QueryPoolByReserveAddressRequest) (*types.QueryPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id cannot be 0")
	}

	if req.ReserveAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty reserve account address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	reserveAddr, err := sdk.AccAddressFromBech32(req.ReserveAddress)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "reserve account address %s is not valid", req.ReserveAddress)
	}

	pool, found := k.GetPoolByReserveAddress(ctx, req.AppId, reserveAddr)
	if !found {
		return nil, status.Errorf(codes.NotFound, "pool by %s doesn't exist", req.ReserveAddress)
	}

	rx, ry := k.GetPoolBalances(ctx, pool)
	poolRes := types.PoolResponse{
		Id:                    pool.Id,
		PairId:                pool.PairId,
		ReserveAddress:        pool.ReserveAddress,
		PoolCoinDenom:         pool.PoolCoinDenom,
		Balances:              sdk.NewCoins(rx, ry),
		LastDepositRequestId:  pool.LastDepositRequestId,
		LastWithdrawRequestId: pool.LastWithdrawRequestId,
	}

	return &types.QueryPoolResponse{Pool: poolRes}, nil
}

// PoolByPoolCoinDenom queries the specific pool by the pool coin denomination.
func (k Querier) PoolByPoolCoinDenom(c context.Context, req *types.QueryPoolByPoolCoinDenomRequest) (*types.QueryPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id cannot be 0")
	}

	if req.PoolCoinDenom == "" {
		return nil, status.Error(codes.InvalidArgument, "empty pool coin denom")
	}

	ctx := sdk.UnwrapSDKContext(c)

	poolID, err := types.ParsePoolCoinDenom(req.PoolCoinDenom)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse pool coin denom: %v", err)
	}
	pool, found := k.GetPool(ctx, req.AppId, poolID)
	if !found {
		return nil, status.Errorf(codes.NotFound, "pool %d doesn't exist", poolID)
	}

	rx, ry := k.GetPoolBalances(ctx, pool)
	poolRes := types.PoolResponse{
		Id:                    pool.Id,
		PairId:                pool.PairId,
		ReserveAddress:        pool.ReserveAddress,
		PoolCoinDenom:         pool.PoolCoinDenom,
		Balances:              sdk.NewCoins(rx, ry),
		LastDepositRequestId:  pool.LastDepositRequestId,
		LastWithdrawRequestId: pool.LastWithdrawRequestId,
	}

	return &types.QueryPoolResponse{Pool: poolRes}, nil
}

// Pairs queries all pairs.
func (k Querier) Pairs(c context.Context, req *types.QueryPairsRequest) (*types.QueryPairsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id cannot be 0")
	}

	if len(req.Denoms) > 2 {
		return nil, status.Errorf(codes.InvalidArgument, "too many denoms to query: %d", len(req.Denoms))
	}

	for _, denom := range req.Denoms {
		if err := sdk.ValidateDenom(denom); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)

	var keyPrefix []byte
	var pairGetter func(key, value []byte) types.Pair
	switch len(req.Denoms) {
	case 0:
		keyPrefix = types.GetAllPairsKey(req.AppId)
		pairGetter = func(_, value []byte) types.Pair {
			return types.MustUnmarshalPair(k.cdc, value)
		}
	case 1:
		keyPrefix = types.GetPairsByDenomIndexKeyPrefix(req.AppId, req.Denoms[0])
		pairGetter = func(key, _ []byte) types.Pair {
			_, _, pairID := types.ParsePairsByDenomsIndexKey(append(keyPrefix, key...))
			pair, _ := k.GetPair(ctx, req.AppId, pairID)
			return pair
		}
	case 2:
		keyPrefix = types.GetPairsByDenomsIndexKeyPrefix(req.AppId, req.Denoms[0], req.Denoms[1])
		pairGetter = func(key, _ []byte) types.Pair {
			_, _, pairID := types.ParsePairsByDenomsIndexKey(append(keyPrefix, key...))
			pair, _ := k.GetPair(ctx, req.AppId, pairID)
			return pair
		}
	}
	pairStore := prefix.NewStore(store, keyPrefix)
	var pairs []types.Pair
	pageRes, err := query.FilteredPaginate(pairStore, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		pair := pairGetter(key, value)

		if accumulate {
			pairs = append(pairs, pair)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPairsResponse{Pairs: pairs, Pagination: pageRes}, nil
}

// Pair queries the specific pair.
func (k Querier) Pair(c context.Context, req *types.QueryPairRequest) (*types.QueryPairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id cannot be 0")
	}

	if req.PairId == 0 {
		return nil, status.Error(codes.InvalidArgument, "pair id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	pair, found := k.GetPair(ctx, req.AppId, req.PairId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "pair %d doesn't exist", req.PairId)
	}

	return &types.QueryPairResponse{Pair: pair}, nil
}

// DepositRequests queries all deposit requests.
func (k Querier) DepositRequests(c context.Context, req *types.QueryDepositRequestsRequest) (*types.QueryDepositRequestsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id cannot be 0")
	}

	if req.PoolId == 0 {
		return nil, status.Error(codes.InvalidArgument, "pool id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	drsStore := prefix.NewStore(store, types.GetAllDepositRequestKey(req.AppId))

	var drs []types.DepositRequest
	pageRes, err := query.FilteredPaginate(drsStore, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		dr, err := types.UnmarshalDepositRequest(k.cdc, value)
		if err != nil {
			return false, err
		}

		if dr.PoolId != req.PoolId {
			return false, nil
		}

		if accumulate {
			drs = append(drs, dr)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDepositRequestsResponse{DepositRequests: drs, Pagination: pageRes}, nil
}

// DepositRequest queries the specific deposit request.
func (k Querier) DepositRequest(c context.Context, req *types.QueryDepositRequestRequest) (*types.QueryDepositRequestResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id cannot be 0")
	}

	if req.PoolId == 0 {
		return nil, status.Error(codes.InvalidArgument, "pool id cannot be 0")
	}

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	dq, found := k.GetDepositRequest(ctx, req.AppId, req.PoolId, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "deposit request of pool id %d and request id %d doesn't exist or deleted", req.PoolId, req.Id)
	}

	return &types.QueryDepositRequestResponse{DepositRequest: dq}, nil
}

// WithdrawRequests queries all withdraw requests.
func (k Querier) WithdrawRequests(c context.Context, req *types.QueryWithdrawRequestsRequest) (*types.QueryWithdrawRequestsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id cannot be 0")
	}

	if req.PoolId == 0 {
		return nil, status.Error(codes.InvalidArgument, "pool id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	drsStore := prefix.NewStore(store, types.GetAllWithdrawRequestKey(req.AppId))

	var wrs []types.WithdrawRequest
	pageRes, err := query.FilteredPaginate(drsStore, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		wr, err := types.UnmarshalWithdrawRequest(k.cdc, value)
		if err != nil {
			return false, err
		}

		if wr.PoolId != req.PoolId {
			return false, nil
		}

		if accumulate {
			wrs = append(wrs, wr)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryWithdrawRequestsResponse{WithdrawRequests: wrs, Pagination: pageRes}, nil
}

// WithdrawRequest queries the specific withdraw request.
func (k Querier) WithdrawRequest(c context.Context, req *types.QueryWithdrawRequestRequest) (*types.QueryWithdrawRequestResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id cannot be 0")
	}

	if req.PoolId == 0 {
		return nil, status.Error(codes.InvalidArgument, "pool id cannot be 0")
	}

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	wq, found := k.GetWithdrawRequest(ctx, req.AppId, req.PoolId, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "withdraw request of pool id %d and request id %d doesn't exist or deleted", req.PoolId, req.Id)
	}

	return &types.QueryWithdrawRequestResponse{WithdrawRequest: wq}, nil
}

// Orders queries all orders.
func (k Querier) Orders(c context.Context, req *types.QueryOrdersRequest) (*types.QueryOrdersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id cannot be 0")
	}

	if req.PairId == 0 {
		return nil, status.Error(codes.InvalidArgument, "pair id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	drsStore := prefix.NewStore(store, types.GetAllOrdersKey(req.AppId))

	var orders []types.Order
	pageRes, err := query.FilteredPaginate(drsStore, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		order, err := types.UnmarshalOrder(k.cdc, value)
		if err != nil {
			return false, err
		}

		if order.PairId != req.PairId {
			return false, nil
		}

		if accumulate {
			orders = append(orders, order)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOrdersResponse{Orders: orders, Pagination: pageRes}, nil
}

// Order queries the specific order.
func (k Querier) Order(c context.Context, req *types.QueryOrderRequest) (*types.QueryOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id cannot be 0")
	}

	if req.PairId == 0 {
		return nil, status.Error(codes.InvalidArgument, "pair id cannot be 0")
	}

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	order, found := k.GetOrder(ctx, req.AppId, req.PairId, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "order %d in pair %d not found", req.PairId, req.Id)
	}

	return &types.QueryOrderResponse{Order: order}, nil
}

// OrdersByOrderer returns orders made by an orderer.
func (k Querier) OrdersByOrderer(c context.Context, req *types.QueryOrdersByOrdererRequest) (*types.QueryOrdersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id cannot be 0")
	}

	orderer, err := sdk.AccAddressFromBech32(req.Orderer)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "orderer address %s is invalid", req.Orderer)
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)

	keyPrefix := types.GetOrderIndexKeyPrefix(req.AppId, orderer)
	orderStore := prefix.NewStore(store, keyPrefix)
	var orders []types.Order
	pageRes, err := query.FilteredPaginate(orderStore, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		_, pairID, orderID := types.ParseOrderIndexKey(append(keyPrefix, key...))
		if req.PairId != 0 && pairID != req.PairId {
			return false, nil
		}

		order, _ := k.GetOrder(ctx, req.AppId, pairID, orderID)

		if accumulate {
			orders = append(orders, order)
		}

		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOrdersResponse{Orders: orders, Pagination: pageRes}, nil
}

// SoftLock returns softlocks created by an depositor in specific pool.
func (k Querier) SoftLock(c context.Context, req *types.QuerySoftLockRequest) (*types.QuerySoftLockResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	depositor, err := sdk.AccAddressFromBech32(req.Depositor)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "orderer address %s is invalid", req.Depositor)
	}

	poolID := req.PoolId

	pool, found := k.GetPool(ctx, req.AppId, poolID)
	if !found {
		return nil, types.ErrInvalidPoolID
	}

	lpData, found := k.GetPoolLiquidityProvidersData(ctx, req.AppId, poolID)
	if !found {
		return &types.QuerySoftLockResponse{ActivePoolCoin: sdk.NewCoin(pool.PoolCoinDenom, sdk.NewInt(0)), QueuedPoolCoin: []types.QueuedPoolCoin{}}, nil
	}

	availableLiquidityGauges := k.rewardsKeeper.GetAllGaugesByGaugeTypeID(ctx, rewardstypes.LiquidityGaugeTypeID)
	minEpochDuration := k.GetMinimumEpochDurationFromPoolID(ctx, poolID, availableLiquidityGauges)

	queuedCoins := []types.QueuedPoolCoin{}
	for _, queuedRequest := range lpData.QueuedLiquidityProviders {
		if queuedRequest.Address == depositor.String() {
			poolCoin := sdk.Coin{}
			for _, coin := range queuedRequest.SupplyProvided {
				if coin.Denom == pool.PoolCoinDenom {
					poolCoin = *coin
					break
				}
			}
			queuedCoins = append(queuedCoins, types.QueuedPoolCoin{
				PoolCoin: poolCoin,
				DequeAt:  queuedRequest.CreatedAt.Add(minEpochDuration),
			})
		}
	}

	activePoolCoin := sdk.NewCoin(pool.PoolCoinDenom, sdk.NewInt(0))
	activeCoins, found := lpData.LiquidityProviders[depositor.String()]
	if found {
		for _, coin := range activeCoins.Coins {
			if coin.Denom == pool.PoolCoinDenom {
				activePoolCoin = coin
				break
			}
		}
	}

	return &types.QuerySoftLockResponse{ActivePoolCoin: activePoolCoin, QueuedPoolCoin: queuedCoins}, nil
}

// DeserializePoolCoin splits poolcoin amount into actual assets provided by depositor.
func (k Querier) DeserializePoolCoin(c context.Context, req *types.QueryDeserializePoolCoinRequest) (*types.QueryDeserializePoolCoinResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	pool, pair, ammPool, err := k.GetAMMPoolInterfaceObject(ctx, req.AppId, req.PoolId)
	if err != nil {
		return nil, err
	}
	poolCoin := sdk.NewCoin(pool.PoolCoinDenom, sdk.NewInt(int64(req.PoolCoinAmount)))
	x, y, err := k.CalculateXYFromPoolCoin(ctx, ammPool, poolCoin)
	if err != nil {
		return &types.QueryDeserializePoolCoinResponse{Coins: []sdk.Coin{sdk.NewCoin(pair.QuoteCoinDenom, sdk.NewInt(0)), sdk.NewCoin(pair.BaseCoinDenom, sdk.NewInt(0))}}, nil
	}
	quoteCoin := sdk.NewCoin(pair.QuoteCoinDenom, x)
	baseCoin := sdk.NewCoin(pair.BaseCoinDenom, y)

	return &types.QueryDeserializePoolCoinResponse{Coins: []sdk.Coin{quoteCoin, baseCoin}}, nil
}

// PoolIncentives provides insights about available pool incentives.
func (k Querier) PoolIncentives(c context.Context, req *types.QueryPoolsIncentivesRequest) (*types.QueryPoolIncentivesResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	liquidityGauges := k.rewardsKeeper.GetAllGaugesByGaugeTypeID(ctx, rewardstypes.LiquidityGaugeTypeID)

	poolIncentives := []*types.PoolIncentive{}

	for _, gauge := range liquidityGauges {
		if ctx.BlockTime().Before(gauge.StartTime) || !gauge.IsActive {
			continue
		}

		// Not needed, redundant check
		// if gauge.TriggeredCount == gauge.TotalTriggers {
		// 	continue
		// }
		epochInfo, found := k.rewardsKeeper.GetEpochInfoByDuration(ctx, gauge.TriggerDuration)
		if !found {
			continue
		}
		childPoolIds := []uint64{}
		if len(gauge.GetLiquidityMetaData().ChildPoolIds) == 0 {
			pools := k.GetAllPools(ctx, req.AppId)
			for _, pool := range pools {
				if pool.Id != gauge.GetLiquidityMetaData().PoolId {
					childPoolIds = append(childPoolIds, pool.Id)
				}
			}
		} else {
			for _, poolID := range gauge.GetLiquidityMetaData().ChildPoolIds {
				if poolID != gauge.GetLiquidityMetaData().PoolId {
					childPoolIds = append(childPoolIds, poolID)
				}
			}
		}
		poolIncentives = append(poolIncentives, &types.PoolIncentive{
			PoolId:             gauge.GetLiquidityMetaData().PoolId,
			MasterPool:         gauge.GetLiquidityMetaData().GetIsMasterPool(),
			ChildPoolIds:       childPoolIds,
			TotalRewards:       gauge.DepositAmount,
			DistributedRewards: gauge.DistributedAmount,
			TotalEpochs:        gauge.TotalTriggers,
			FilledEpochs:       gauge.TriggeredCount,
			EpochDuration:      gauge.TriggerDuration,
			NextDistribution:   epochInfo.CurrentEpochStartTime.Add(epochInfo.Duration),
			IsSwapFee:          gauge.ForSwapFee,
		})
	}

	return &types.QueryPoolIncentivesResponse{PoolIncentives: poolIncentives}, nil
}

// FarmedPoolCoin returns the total pool coin in soft-lock.
func (k Querier) FarmedPoolCoin(c context.Context, req *types.QueryFarmedPoolCoinRequest) (*types.QueryFarmedPoolCoinResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.AppId == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)
	pool, found := k.GetPool(ctx, req.AppId, req.PoolId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPoolID, "pool id %d is invalid", req.PoolId)
	}
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	softLockedCoins := k.bankKeeper.GetBalance(ctx, moduleAddr, pool.PoolCoinDenom)

	return &types.QueryFarmedPoolCoinResponse{Coin: softLockedCoins}, nil
}
