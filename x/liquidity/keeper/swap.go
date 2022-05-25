package keeper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/comdex-official/comdex/x/liquidity/amm"
	"github.com/comdex-official/comdex/x/liquidity/types"
)

// ValidateMsgLimitOrder validates types.MsgLimitOrder with state and returns
// calculated offer coin and price that is fit into ticks.
func (k Keeper) ValidateMsgLimitOrder(ctx sdk.Context, msg *types.MsgLimitOrder) (offerCoin sdk.Coin, price sdk.Dec, err error) {
	params := k.GetParams(ctx)

	if msg.OrderLifespan > params.MaxOrderLifespan {
		return sdk.Coin{}, sdk.Dec{},
			sdkerrors.Wrapf(types.ErrTooLongOrderLifespan, "%s is longer than %s", msg.OrderLifespan, params.MaxOrderLifespan)
	}

	pair, found := k.GetPair(ctx, msg.PairId)
	if !found {
		return sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", msg.PairId)
	}

	var upperPriceLimit, lowerPriceLimit sdk.Dec
	if pair.LastPrice != nil {
		lastPrice := *pair.LastPrice
		upperPriceLimit = lastPrice.Mul(sdk.OneDec().Add(params.MaxPriceLimitRatio))
		lowerPriceLimit = lastPrice.Mul(sdk.OneDec().Sub(params.MaxPriceLimitRatio))
	} else {
		upperPriceLimit = amm.HighestTick(int(params.TickPrecision))
		lowerPriceLimit = amm.LowestTick(int(params.TickPrecision))
	}
	switch {
	case msg.Price.GT(upperPriceLimit):
		return sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(types.ErrPriceOutOfRange, "%s is higher than %s", msg.Price, upperPriceLimit)
	case msg.Price.LT(lowerPriceLimit):
		return sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(types.ErrPriceOutOfRange, "%s is lower than %s", msg.Price, lowerPriceLimit)
	}

	switch msg.Direction {
	case types.OrderDirectionBuy:
		if msg.OfferCoin.Denom != pair.QuoteCoinDenom || msg.DemandCoinDenom != pair.BaseCoinDenom {
			return sdk.Coin{}, sdk.Dec{},
				sdkerrors.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)",
					msg.DemandCoinDenom, msg.OfferCoin.Denom, pair.BaseCoinDenom, pair.QuoteCoinDenom)
		}
		price = amm.PriceToDownTick(msg.Price, int(params.TickPrecision))
		offerCoin = sdk.NewCoin(msg.OfferCoin.Denom, amm.OfferCoinAmount(amm.Buy, price, msg.Amount))
		if msg.OfferCoin.IsLT(offerCoin) {
			return sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(
				types.ErrInsufficientOfferCoin, "%s is smaller than %s", msg.OfferCoin, offerCoin)
		}
	case types.OrderDirectionSell:
		if msg.OfferCoin.Denom != pair.BaseCoinDenom || msg.DemandCoinDenom != pair.QuoteCoinDenom {
			return sdk.Coin{}, sdk.Dec{},
				sdkerrors.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)",
					msg.OfferCoin.Denom, msg.DemandCoinDenom, pair.BaseCoinDenom, pair.QuoteCoinDenom)
		}
		price = amm.PriceToUpTick(msg.Price, int(params.TickPrecision))
		offerCoin = sdk.NewCoin(msg.OfferCoin.Denom, msg.Amount)
		if msg.OfferCoin.Amount.LT(msg.Amount) {
			return sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(
				types.ErrInsufficientOfferCoin, "%s is smaller than %s", msg.OfferCoin, sdk.NewCoin(msg.OfferCoin.Denom, msg.Amount))
		}
	}
	if types.IsTooSmallOrderAmount(msg.Amount, price) {
		return sdk.Coin{}, sdk.Dec{}, types.ErrTooSmallOrder
	}

	return offerCoin, price, nil
}

// LimitOrder handles types.MsgLimitOrder and stores types.Order.
func (k Keeper) LimitOrder(ctx sdk.Context, msg *types.MsgLimitOrder) (types.Order, error) {
	offerCoin, price, err := k.ValidateMsgLimitOrder(ctx, msg)
	if err != nil {
		return types.Order{}, err
	}

	refundedCoin := msg.OfferCoin.Sub(offerCoin)
	pair, _ := k.GetPair(ctx, msg.PairId)
	if err := k.bankKeeper.SendCoins(ctx, msg.GetOrderer(), pair.GetEscrowAddress(), sdk.NewCoins(offerCoin)); err != nil {
		return types.Order{}, err
	}

	requestID := k.getNextOrderIDWithUpdate(ctx, pair)
	expireAt := ctx.BlockTime().Add(msg.OrderLifespan)
	order := types.NewOrderForLimitOrder(msg, requestID, pair, offerCoin, price, expireAt, ctx.BlockHeight())
	k.SetOrder(ctx, order)
	k.SetOrderIndex(ctx, order)

	params := k.GetParams(ctx)
	ctx.GasMeter().ConsumeGas(params.OrderExtraGas, "OrderExtraGas")

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLimitOrder,
			sdk.NewAttribute(types.AttributeKeyOrderer, msg.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairID, strconv.FormatUint(msg.PairId, 10)),
			sdk.NewAttribute(types.AttributeKeyOrderDirection, msg.Direction.String()),
			sdk.NewAttribute(types.AttributeKeyOfferCoin, offerCoin.String()),
			sdk.NewAttribute(types.AttributeKeyDemandCoinDenom, msg.DemandCoinDenom),
			sdk.NewAttribute(types.AttributeKeyPrice, price.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyOrderID, strconv.FormatUint(order.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyBatchID, strconv.FormatUint(order.BatchId, 10)),
			sdk.NewAttribute(types.AttributeKeyExpireAt, order.ExpireAt.Format(time.RFC3339)),
			sdk.NewAttribute(types.AttributeKeyRefundedCoins, refundedCoin.String()),
		),
	})

	return order, nil
}

// ValidateMsgMarketOrder validates types.MsgMarketOrder with state and returns
// calculated offer coin and price.
func (k Keeper) ValidateMsgMarketOrder(ctx sdk.Context, msg *types.MsgMarketOrder) (offerCoin sdk.Coin, price sdk.Dec, err error) {
	params := k.GetParams(ctx)

	if msg.OrderLifespan > params.MaxOrderLifespan {
		return sdk.Coin{}, sdk.Dec{},
			sdkerrors.Wrapf(types.ErrTooLongOrderLifespan, "%s is longer than %s", msg.OrderLifespan, params.MaxOrderLifespan)
	}

	pair, found := k.GetPair(ctx, msg.PairId)
	if !found {
		return sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", msg.PairId)
	}

	if pair.LastPrice == nil {
		return sdk.Coin{}, sdk.Dec{}, types.ErrNoLastPrice
	}
	lastPrice := *pair.LastPrice

	switch msg.Direction {
	case types.OrderDirectionBuy:
		if msg.OfferCoin.Denom != pair.QuoteCoinDenom || msg.DemandCoinDenom != pair.BaseCoinDenom {
			return sdk.Coin{}, sdk.Dec{},
				sdkerrors.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)",
					msg.DemandCoinDenom, msg.OfferCoin.Denom, pair.BaseCoinDenom, pair.QuoteCoinDenom)
		}
		price = amm.PriceToDownTick(lastPrice.Mul(sdk.OneDec().Add(params.MaxPriceLimitRatio)), int(params.TickPrecision))
		offerCoin = sdk.NewCoin(msg.OfferCoin.Denom, amm.OfferCoinAmount(amm.Buy, price, msg.Amount))
		if msg.OfferCoin.IsLT(offerCoin) {
			return sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(
				types.ErrInsufficientOfferCoin, "%s is smaller than %s", msg.OfferCoin, offerCoin)
		}
	case types.OrderDirectionSell:
		if msg.OfferCoin.Denom != pair.BaseCoinDenom || msg.DemandCoinDenom != pair.QuoteCoinDenom {
			return sdk.Coin{}, sdk.Dec{},
				sdkerrors.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)",
					msg.OfferCoin.Denom, msg.DemandCoinDenom, pair.BaseCoinDenom, pair.QuoteCoinDenom)
		}
		price = amm.PriceToUpTick(lastPrice.Mul(sdk.OneDec().Sub(params.MaxPriceLimitRatio)), int(params.TickPrecision))
		offerCoin = sdk.NewCoin(msg.OfferCoin.Denom, msg.Amount)
		if msg.OfferCoin.Amount.LT(msg.Amount) {
			return sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(
				types.ErrInsufficientOfferCoin, "%s is smaller than %s", msg.OfferCoin, sdk.NewCoin(msg.OfferCoin.Denom, msg.Amount))
		}
	}
	if types.IsTooSmallOrderAmount(msg.Amount, price) {
		return sdk.Coin{}, sdk.Dec{}, types.ErrTooSmallOrder
	}

	return offerCoin, price, nil
}

// MarketOrder handles types.MsgMarketOrder and stores types.Order.
func (k Keeper) MarketOrder(ctx sdk.Context, msg *types.MsgMarketOrder) (types.Order, error) {
	offerCoin, price, err := k.ValidateMsgMarketOrder(ctx, msg)
	if err != nil {
		return types.Order{}, err
	}

	refundedCoin := msg.OfferCoin.Sub(offerCoin)
	pair, _ := k.GetPair(ctx, msg.PairId)
	if err := k.bankKeeper.SendCoins(ctx, msg.GetOrderer(), pair.GetEscrowAddress(), sdk.NewCoins(offerCoin)); err != nil {
		return types.Order{}, err
	}

	requestID := k.getNextOrderIDWithUpdate(ctx, pair)
	expireAt := ctx.BlockTime().Add(msg.OrderLifespan)
	order := types.NewOrderForMarketOrder(msg, requestID, pair, offerCoin, price, expireAt, ctx.BlockHeight())
	k.SetOrder(ctx, order)
	k.SetOrderIndex(ctx, order)

	params := k.GetParams(ctx)
	ctx.GasMeter().ConsumeGas(params.OrderExtraGas, "OrderExtraGas")

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMarketOrder,
			sdk.NewAttribute(types.AttributeKeyOrderer, msg.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairID, strconv.FormatUint(msg.PairId, 10)),
			sdk.NewAttribute(types.AttributeKeyOrderDirection, msg.Direction.String()),
			sdk.NewAttribute(types.AttributeKeyOfferCoin, msg.OfferCoin.String()),
			sdk.NewAttribute(types.AttributeKeyDemandCoinDenom, msg.DemandCoinDenom),
			sdk.NewAttribute(types.AttributeKeyPrice, price.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyOrderID, strconv.FormatUint(order.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyBatchID, strconv.FormatUint(order.BatchId, 10)),
			sdk.NewAttribute(types.AttributeKeyExpireAt, order.ExpireAt.Format(time.RFC3339)),
			sdk.NewAttribute(types.AttributeKeyRefundedCoins, refundedCoin.String()),
		),
	})

	return order, nil
}

// ValidateMsgCancelOrder validates types.MsgCancelOrder and returns the order.
func (k Keeper) ValidateMsgCancelOrder(ctx sdk.Context, msg *types.MsgCancelOrder) (order types.Order, err error) {
	var found bool
	order, found = k.GetOrder(ctx, msg.PairId, msg.OrderId)
	if !found {
		return types.Order{},
			sdkerrors.Wrapf(sdkerrors.ErrNotFound, "order %d not found in pair %d", msg.OrderId, msg.PairId)
	}
	if msg.Orderer != order.Orderer {
		return types.Order{}, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "mismatching orderer")
	}
	if order.Status == types.OrderStatusCanceled {
		return types.Order{}, types.ErrAlreadyCanceled
	}
	pair, _ := k.GetPair(ctx, msg.PairId)
	if order.BatchId == pair.CurrentBatchId {
		return types.Order{}, types.ErrSameBatch
	}
	return order, nil
}

// CancelOrder handles types.MsgCancelOrder and cancels an order.
func (k Keeper) CancelOrder(ctx sdk.Context, msg *types.MsgCancelOrder) error {
	order, err := k.ValidateMsgCancelOrder(ctx, msg)
	if err != nil {
		return err
	}

	if err := k.FinishOrder(ctx, order, types.OrderStatusCanceled); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelOrder,
			sdk.NewAttribute(types.AttributeKeyOrderer, msg.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairID, strconv.FormatUint(msg.PairId, 10)),
			sdk.NewAttribute(types.AttributeKeyOrderID, strconv.FormatUint(msg.OrderId, 10)),
		),
	})

	return nil
}

// CancelAllOrders handles types.MsgCancelAllOrders and cancels all orders.
func (k Keeper) CancelAllOrders(ctx sdk.Context, msg *types.MsgCancelAllOrders) error {
	var canceledOrderIds []string
	//nolint
	cb := func(pair types.Pair, order types.Order) (stop bool, err error) {
		if order.Orderer == msg.Orderer && order.Status != types.OrderStatusCanceled && order.BatchId < pair.CurrentBatchId {
			if err := k.FinishOrder(ctx, order, types.OrderStatusCanceled); err != nil {
				return false, err
			}
			canceledOrderIds = append(canceledOrderIds, strconv.FormatUint(order.Id, 10))
		}
		return false, nil
	}

	var pairIDs []string
	if len(msg.PairIds) == 0 {
		pairMap := map[uint64]types.Pair{}
		if err := k.IterateAllOrders(ctx, func(order types.Order) (stop bool, err error) {
			pair, ok := pairMap[order.PairId]
			if !ok {
				pair, _ = k.GetPair(ctx, order.PairId)
				pairMap[order.PairId] = pair
			}
			return cb(pair, order)
		}); err != nil {
			return err
		}
	} else {
		for _, pairID := range msg.PairIds {
			pairIDs = append(pairIDs, strconv.FormatUint(pairID, 10))
			pair, found := k.GetPair(ctx, pairID)
			if !found {
				return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", pairID)
			}
			if err := k.IterateOrdersByPair(ctx, pairID, func(req types.Order) (stop bool, err error) {
				return cb(pair, req)
			}); err != nil {
				return err
			}
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelAllOrders,
			sdk.NewAttribute(types.AttributeKeyOrderer, msg.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairIds, strings.Join(pairIDs, ",")),
			sdk.NewAttribute(types.AttributeKeyCanceledOrderIds, strings.Join(canceledOrderIds, ",")),
		),
	})

	return nil
}

func (k Keeper) ExecuteMatching(ctx sdk.Context, pair types.Pair) error {
	ob := amm.NewOrderBook()
	skip := true // Whether to skip the matching since there is no orders.
	if err := k.IterateOrdersByPair(ctx, pair.Id, func(order types.Order) (stop bool, err error) {
		switch order.Status {
		case types.OrderStatusNotExecuted,
			types.OrderStatusNotMatched,
			types.OrderStatusPartiallyMatched:
			if order.Status != types.OrderStatusNotExecuted && order.ExpiredAt(ctx.BlockTime()) {
				if err := k.FinishOrder(ctx, order, types.OrderStatusExpired); err != nil {
					return false, err
				}
				return false, nil
			}
			ob.Add(types.NewUserOrder(order))
			if order.Status == types.OrderStatusNotExecuted {
				order.SetStatus(types.OrderStatusNotMatched)
				k.SetOrder(ctx, order)
			}
			skip = false
		case types.OrderStatusCanceled:
		default:
			return false, fmt.Errorf("invalid order status: %s", order.Status)
		}
		return false, nil
	}); err != nil {
		return err
	}

	if skip { //nolint TODO: update this when there are more than one pools
		return nil
	}

	var poolOrderSources []amm.OrderSource
	_ = k.IteratePoolsByPair(ctx, pair.Id, func(pool types.Pool) (stop bool, err error) {
		rx, ry := k.getPoolBalances(ctx, pool, pair)
		ps := k.GetPoolCoinSupply(ctx, pool)
		ammPool := amm.NewBasicPool(rx.Amount, ry.Amount, ps)
		if ammPool.IsDepleted() {
			k.MarkPoolAsDisabled(ctx, pool)
			return false, nil
		}
		poolOrderSource := types.NewBasicPoolOrderSource(ammPool, pool.Id, pool.GetReserveAddress(), pair.BaseCoinDenom, pair.QuoteCoinDenom)
		poolOrderSources = append(poolOrderSources, poolOrderSource)
		return false, nil
	})

	os := amm.MergeOrderSources(append(poolOrderSources, ob)...)

	params := k.GetParams(ctx)
	matchPrice, found := amm.FindMatchPrice(os, int(params.TickPrecision))
	if found {
		buyOrders := os.BuyOrdersOver(matchPrice)
		sellOrders := os.SellOrdersUnder(matchPrice)

		types.SortOrders(buyOrders, types.PriceDescending)
		types.SortOrders(sellOrders, types.PriceAscending)

		quoteCoinDust, matched := amm.MatchOrders(buyOrders, sellOrders, matchPrice)
		if matched {
			if err := k.ApplyMatchResult(ctx, pair, append(buyOrders, sellOrders...), quoteCoinDust); err != nil {
				return err
			}
			pair.LastPrice = &matchPrice
		}
	}

	pair.CurrentBatchId++
	k.SetPair(ctx, pair)

	//nolint
	// TODO: emit an event?
	return nil
}

func (k Keeper) ApplyMatchResult(ctx sdk.Context, pair types.Pair, orders []amm.Order, quoteCoinDust sdk.Int) error {
	bulkOp := types.NewBulkSendCoinsOperation()
	for _, order := range orders {
		if !order.IsMatched() {
			continue
		}
		if order, ok := order.(*types.PoolOrder); ok {
			paidCoin := order.OfferCoin.Sub(order.RemainingOfferCoin)
			bulkOp.QueueSendCoins(order.ReserveAddress, pair.GetEscrowAddress(), sdk.NewCoins(paidCoin))
		}
	}
	if err := bulkOp.Run(ctx, k.bankKeeper); err != nil {
		return err
	}
	bulkOp = types.NewBulkSendCoinsOperation()
	for _, order := range orders {
		if !order.IsMatched() {
			continue
		}
		switch order := order.(type) {
		case *types.UserOrder:
			o, _ := k.GetOrder(ctx, pair.Id, order.OrderID)
			o.OpenAmount = o.OpenAmount.Sub(order.Amount.Sub(order.OpenAmount))
			o.RemainingOfferCoin = o.RemainingOfferCoin.Sub(order.OfferCoin.Sub(order.RemainingOfferCoin))
			o.ReceivedCoin = o.ReceivedCoin.Add(order.ReceivedDemandCoin)
			if o.OpenAmount.IsZero() {
				if err := k.FinishOrder(ctx, o, types.OrderStatusCompleted); err != nil {
					return err
				}
			} else {
				o.SetStatus(types.OrderStatusPartiallyMatched)
				k.SetOrder(ctx, o)
				// nolint
				// TODO: emit an event?
			}
			bulkOp.QueueSendCoins(pair.GetEscrowAddress(), order.Orderer, sdk.NewCoins(order.ReceivedDemandCoin))
		case *types.PoolOrder:
			bulkOp.QueueSendCoins(pair.GetEscrowAddress(), order.ReserveAddress, sdk.NewCoins(order.ReceivedDemandCoin))
		}
	}
	params := k.GetParams(ctx)
	swapFeeCollectorAddr, _ := sdk.AccAddressFromBech32(params.SwapFeeCollectorAddress)
	bulkOp.QueueSendCoins(pair.GetEscrowAddress(), swapFeeCollectorAddr, sdk.NewCoins(sdk.NewCoin(pair.QuoteCoinDenom, quoteCoinDust)))
	if err := bulkOp.Run(ctx, k.bankKeeper); err != nil {
		return err
	}
	return nil
}

func (k Keeper) FinishOrder(ctx sdk.Context, order types.Order, status types.OrderStatus) error {
	if order.Status == types.OrderStatusCompleted || order.Status.IsCanceledOrExpired() { // sanity check
		return nil
	}

	if order.RemainingOfferCoin.IsPositive() {
		pair, _ := k.GetPair(ctx, order.PairId)
		if err := k.bankKeeper.SendCoins(ctx, pair.GetEscrowAddress(), order.GetOrderer(), sdk.NewCoins(order.RemainingOfferCoin)); err != nil {
			return err
		}
	}

	order.SetStatus(status)
	k.SetOrder(ctx, order)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeOrderResult,
			sdk.NewAttribute(types.AttributeKeyRequestID, strconv.FormatUint(order.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyOrderer, order.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairID, strconv.FormatUint(order.PairId, 10)),
			sdk.NewAttribute(types.AttributeKeyOrderDirection, order.Direction.String()),
			//nolint
			// TODO: include these attributes?
			//sdk.NewAttribute(types.AttributeKeyOfferCoin, order.OfferCoin.String()),
			//sdk.NewAttribute(types.AttributeKeyAmount, order.Amount.String()),
			//sdk.NewAttribute(types.AttributeKeyOpenAmount, order.OpenAmount.String()),
			//sdk.NewAttribute(types.AttributeKeyPrice, order.Price.String()),
			sdk.NewAttribute(types.AttributeKeyRemainingOfferCoin, order.RemainingOfferCoin.String()),
			sdk.NewAttribute(types.AttributeKeyReceivedCoin, order.ReceivedCoin.String()),
			sdk.NewAttribute(types.AttributeKeyStatus, order.Status.String()),
		),
	})

	return nil
}
