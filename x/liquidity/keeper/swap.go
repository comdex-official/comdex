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

func CalculateSwapFeeAmount(ctx sdk.Context, params types.GenericParams, calculatedOfferCoinAmt sdk.Int) sdk.Int {
	return calculatedOfferCoinAmt.ToDec().MulTruncate(params.SwapFeeRate).TruncateInt()
}

func (k Keeper) PriceLimits(ctx sdk.Context, lastPrice sdk.Dec, params types.GenericParams) (lowest, highest sdk.Dec) {
	return types.PriceLimits(lastPrice, params.MaxPriceLimitRatio, int(params.TickPrecision))
}

// ValidateMsgLimitOrder validates types.MsgLimitOrder with state and returns
// calculated offer coin and price that is fit into ticks.
func (k Keeper) ValidateMsgLimitOrder(ctx sdk.Context, msg *types.MsgLimitOrder) (offerCoin sdk.Coin, swapFeeCoin sdk.Coin, price sdk.Dec, err error) {
	_, found := k.assetKeeper.GetApp(ctx, msg.AppId)
	if !found {
		return sdk.Coin{}, sdk.Coin{}, sdk.Dec{},
			sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", msg.AppId)
	}

	params, err := k.GetGenericParams(ctx, msg.AppId)
	if err != nil {
		return sdk.Coin{}, sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrap(err, "params retreval failed")
	}

	spendable := k.bankKeeper.SpendableCoins(ctx, msg.GetOrderer())
	if spendableAmt := spendable.AmountOf(msg.OfferCoin.Denom); spendableAmt.LT(msg.OfferCoin.Amount) {
		return sdk.Coin{}, sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds, "%s is smaller than %s",
			sdk.NewCoin(msg.OfferCoin.Denom, spendableAmt), msg.OfferCoin)
	}

	if msg.OrderLifespan > params.MaxOrderLifespan {
		return sdk.Coin{}, sdk.Coin{}, sdk.Dec{},
			sdkerrors.Wrapf(types.ErrTooLongOrderLifespan, "%s is longer than %s", msg.OrderLifespan, params.MaxOrderLifespan)
	}

	pair, found := k.GetPair(ctx, msg.AppId, msg.PairId)
	if !found {
		return sdk.Coin{}, sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", msg.PairId)
	}

	var upperPriceLimit, lowerPriceLimit sdk.Dec
	if pair.LastPrice != nil {
		lowerPriceLimit, upperPriceLimit = k.PriceLimits(ctx, *pair.LastPrice, params)
	} else {
		upperPriceLimit = amm.HighestTick(int(params.TickPrecision))
		lowerPriceLimit = amm.LowestTick(int(params.TickPrecision))
	}
	switch {
	case msg.Price.GT(upperPriceLimit):
		return sdk.Coin{}, sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(types.ErrPriceOutOfRange, "%s is higher than %s", msg.Price, upperPriceLimit)
	case msg.Price.LT(lowerPriceLimit):
		return sdk.Coin{}, sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(types.ErrPriceOutOfRange, "%s is lower than %s", msg.Price, lowerPriceLimit)
	}

	switch msg.Direction {
	case types.OrderDirectionBuy:
		if msg.OfferCoin.Denom != pair.QuoteCoinDenom || msg.DemandCoinDenom != pair.BaseCoinDenom {
			return sdk.Coin{}, sdk.Coin{}, sdk.Dec{},
				sdkerrors.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)",
					msg.DemandCoinDenom, msg.OfferCoin.Denom, pair.BaseCoinDenom, pair.QuoteCoinDenom)
		}
		price = amm.PriceToDownTick(msg.Price, int(params.TickPrecision))

		offerCoin = sdk.NewCoin(msg.OfferCoin.Denom, amm.OfferCoinAmount(amm.Buy, price, msg.Amount))
		swapFeeCoin = sdk.NewCoin(msg.OfferCoin.Denom, CalculateSwapFeeAmount(ctx, params, offerCoin.Amount))

		if msg.OfferCoin.IsLT(offerCoin.Add(swapFeeCoin)) {
			return sdk.Coin{}, sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(
				types.ErrInsufficientOfferCoin, "%s is smaller than %s", msg.OfferCoin, offerCoin.Add(swapFeeCoin))
		}
	case types.OrderDirectionSell:
		if msg.OfferCoin.Denom != pair.BaseCoinDenom || msg.DemandCoinDenom != pair.QuoteCoinDenom {
			return sdk.Coin{}, sdk.Coin{}, sdk.Dec{},
				sdkerrors.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)",
					msg.OfferCoin.Denom, msg.DemandCoinDenom, pair.BaseCoinDenom, pair.QuoteCoinDenom)
		}
		price = amm.PriceToUpTick(msg.Price, int(params.TickPrecision))

		offerCoin = sdk.NewCoin(msg.OfferCoin.Denom, msg.Amount)
		swapFeeCoin = sdk.NewCoin(msg.OfferCoin.Denom, CalculateSwapFeeAmount(ctx, params, offerCoin.Amount))

		if msg.OfferCoin.Amount.LT(swapFeeCoin.Amount.Add(offerCoin.Amount)) {
			return sdk.Coin{}, sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(
				types.ErrInsufficientOfferCoin, "%s is smaller than %s", msg.OfferCoin, sdk.NewCoin(msg.OfferCoin.Denom, swapFeeCoin.Amount.Add(offerCoin.Amount)))
		}
	}
	if types.IsTooSmallOrderAmount(msg.Amount, price) {
		return sdk.Coin{}, sdk.Coin{}, sdk.Dec{}, types.ErrTooSmallOrder
	}

	return offerCoin, swapFeeCoin, price, nil
}

// LimitOrder handles types.MsgLimitOrder and stores types.Order.
func (k Keeper) LimitOrder(ctx sdk.Context, msg *types.MsgLimitOrder) (types.Order, error) {
	params, err := k.GetGenericParams(ctx, msg.AppId)
	if err != nil {
		return types.Order{}, sdkerrors.Wrap(err, "params retreval failed")
	}

	offerCoin, swapFeeCoin, price, err := k.ValidateMsgLimitOrder(ctx, msg)
	if err != nil {
		return types.Order{}, err
	}

	refundedCoin := msg.OfferCoin.Sub(offerCoin.Add(swapFeeCoin))
	pair, _ := k.GetPair(ctx, msg.AppId, msg.PairId)
	if err := k.bankKeeper.SendCoins(ctx, msg.GetOrderer(), pair.GetEscrowAddress(), sdk.NewCoins(offerCoin.Add(swapFeeCoin))); err != nil {
		return types.Order{}, err
	}

	requestID := k.getNextOrderIDWithUpdate(ctx, pair)
	expireAt := ctx.BlockTime().Add(msg.OrderLifespan)
	order := types.NewOrderForLimitOrder(msg, requestID, pair, offerCoin, price, expireAt, ctx.BlockHeight())
	k.SetOrder(ctx, msg.AppId, order)
	k.SetOrderIndex(ctx, msg.AppId, order)

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
func (k Keeper) ValidateMsgMarketOrder(ctx sdk.Context, msg *types.MsgMarketOrder) (offerCoin sdk.Coin, swapFeeCoin sdk.Coin, price sdk.Dec, err error) {
	_, found := k.assetKeeper.GetApp(ctx, msg.AppId)
	if !found {
		return sdk.Coin{}, sdk.Coin{}, sdk.Dec{},
			sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", msg.AppId)
	}

	params, err := k.GetGenericParams(ctx, msg.AppId)
	if err != nil {
		return sdk.Coin{}, sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrap(err, "params retreval failed")
	}

	spendable := k.bankKeeper.SpendableCoins(ctx, msg.GetOrderer())
	if spendableAmt := spendable.AmountOf(msg.OfferCoin.Denom); spendableAmt.LT(msg.OfferCoin.Amount) {
		return sdk.Coin{}, sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds, "%s is smaller than %s",
			sdk.NewCoin(msg.OfferCoin.Denom, spendableAmt), msg.OfferCoin)
	}

	if msg.OrderLifespan > params.MaxOrderLifespan {
		return sdk.Coin{}, sdk.Coin{}, sdk.Dec{},
			sdkerrors.Wrapf(types.ErrTooLongOrderLifespan, "%s is longer than %s", msg.OrderLifespan, params.MaxOrderLifespan)
	}

	pair, found := k.GetPair(ctx, msg.AppId, msg.PairId)
	if !found {
		return sdk.Coin{}, sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", msg.PairId)
	}

	if pair.LastPrice == nil {
		return sdk.Coin{}, sdk.Coin{}, sdk.Dec{}, types.ErrNoLastPrice
	}
	lastPrice := *pair.LastPrice

	switch msg.Direction {
	case types.OrderDirectionBuy:
		if msg.OfferCoin.Denom != pair.QuoteCoinDenom || msg.DemandCoinDenom != pair.BaseCoinDenom {
			return sdk.Coin{}, sdk.Coin{}, sdk.Dec{},
				sdkerrors.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)",
					msg.DemandCoinDenom, msg.OfferCoin.Denom, pair.BaseCoinDenom, pair.QuoteCoinDenom)
		}
		price = amm.PriceToDownTick(lastPrice.Mul(sdk.OneDec().Add(params.MaxPriceLimitRatio)), int(params.TickPrecision))
		offerCoin = sdk.NewCoin(msg.OfferCoin.Denom, amm.OfferCoinAmount(amm.Buy, price, msg.Amount))
		swapFeeCoin = sdk.NewCoin(msg.OfferCoin.Denom, CalculateSwapFeeAmount(ctx, params, offerCoin.Amount))
		if msg.OfferCoin.IsLT(offerCoin.Add(swapFeeCoin)) {
			return sdk.Coin{}, sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(
				types.ErrInsufficientOfferCoin, "%s is smaller than %s", msg.OfferCoin, offerCoin.Add(swapFeeCoin))
		}
	case types.OrderDirectionSell:
		if msg.OfferCoin.Denom != pair.BaseCoinDenom || msg.DemandCoinDenom != pair.QuoteCoinDenom {
			return sdk.Coin{}, sdk.Coin{}, sdk.Dec{},
				sdkerrors.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)",
					msg.OfferCoin.Denom, msg.DemandCoinDenom, pair.BaseCoinDenom, pair.QuoteCoinDenom)
		}
		price = amm.PriceToUpTick(lastPrice.Mul(sdk.OneDec().Sub(params.MaxPriceLimitRatio)), int(params.TickPrecision))
		offerCoin = sdk.NewCoin(msg.OfferCoin.Denom, msg.Amount)
		swapFeeCoin = sdk.NewCoin(msg.OfferCoin.Denom, CalculateSwapFeeAmount(ctx, params, offerCoin.Amount))
		if msg.OfferCoin.Amount.LT(swapFeeCoin.Amount.Add(offerCoin.Amount)) {
			return sdk.Coin{}, sdk.Coin{}, sdk.Dec{}, sdkerrors.Wrapf(
				types.ErrInsufficientOfferCoin, "%s is smaller than %s", msg.OfferCoin, sdk.NewCoin(msg.OfferCoin.Denom, swapFeeCoin.Amount.Add(offerCoin.Amount)))
		}
	}
	if types.IsTooSmallOrderAmount(msg.Amount, price) {
		return sdk.Coin{}, sdk.Coin{}, sdk.Dec{}, types.ErrTooSmallOrder
	}

	return offerCoin, swapFeeCoin, price, nil
}

// MarketOrder handles types.MsgMarketOrder and stores types.Order.
func (k Keeper) MarketOrder(ctx sdk.Context, msg *types.MsgMarketOrder) (types.Order, error) {
	params, err := k.GetGenericParams(ctx, msg.AppId)
	if err != nil {
		return types.Order{}, sdkerrors.Wrap(err, "params retreval failed")
	}

	offerCoin, swapFeeCoin, price, err := k.ValidateMsgMarketOrder(ctx, msg)
	if err != nil {
		return types.Order{}, err
	}

	refundedCoin := msg.OfferCoin.Sub(offerCoin.Add(swapFeeCoin))
	pair, _ := k.GetPair(ctx, msg.AppId, msg.PairId)
	if err := k.bankKeeper.SendCoins(ctx, msg.GetOrderer(), pair.GetEscrowAddress(), sdk.NewCoins(offerCoin.Add(swapFeeCoin))); err != nil {
		return types.Order{}, err
	}

	requestID := k.getNextOrderIDWithUpdate(ctx, pair)
	expireAt := ctx.BlockTime().Add(msg.OrderLifespan)
	order := types.NewOrderForMarketOrder(msg, requestID, pair, offerCoin, price, expireAt, ctx.BlockHeight())
	k.SetOrder(ctx, msg.AppId, order)
	k.SetOrderIndex(ctx, msg.AppId, order)

	ctx.GasMeter().ConsumeGas(params.OrderExtraGas, "OrderExtraGas")

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMarketOrder,
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

func (k Keeper) MMOrder(ctx sdk.Context, msg *types.MsgMMOrder) (orders []types.Order, err error) {
	_, found := k.assetKeeper.GetApp(ctx, msg.AppId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", msg.AppId)
	}

	params, err := k.GetGenericParams(ctx, msg.AppId)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "params retreval failed")
	}

	tickPrec := int(params.TickPrecision)

	if msg.SellAmount.IsPositive() {
		if !amm.PriceToDownTick(msg.MinSellPrice, tickPrec).Equal(msg.MinSellPrice) {
			return nil, sdkerrors.Wrapf(types.ErrPriceNotOnTicks, "min sell price is not on ticks")
		}
		if !amm.PriceToDownTick(msg.MaxSellPrice, tickPrec).Equal(msg.MaxSellPrice) {
			return nil, sdkerrors.Wrapf(types.ErrPriceNotOnTicks, "max sell price is not on ticks")
		}
	}
	if msg.BuyAmount.IsPositive() {
		if !amm.PriceToDownTick(msg.MinBuyPrice, tickPrec).Equal(msg.MinBuyPrice) {
			return nil, sdkerrors.Wrapf(types.ErrPriceNotOnTicks, "min buy price is not on ticks")
		}
		if !amm.PriceToDownTick(msg.MaxBuyPrice, tickPrec).Equal(msg.MaxBuyPrice) {
			return nil, sdkerrors.Wrapf(types.ErrPriceNotOnTicks, "max buy price is not on ticks")
		}
	}

	pair, found := k.GetPair(ctx, msg.AppId, msg.PairId)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", msg.PairId)
	}

	var lowestPrice, highestPrice sdk.Dec
	if pair.LastPrice != nil {
		lowestPrice, highestPrice = types.PriceLimits(*pair.LastPrice, params.MaxPriceLimitRatio, tickPrec)
	} else {
		lowestPrice = amm.LowestTick(tickPrec)
		highestPrice = amm.HighestTick(tickPrec)
	}

	if msg.SellAmount.IsPositive() {
		if msg.MinSellPrice.LT(lowestPrice) || msg.MinSellPrice.GT(highestPrice) {
			return nil, sdkerrors.Wrapf(types.ErrPriceOutOfRange, "min sell price is out of range [%s, %s]", lowestPrice, highestPrice)
		}
		if msg.MaxSellPrice.LT(lowestPrice) || msg.MaxSellPrice.GT(highestPrice) {
			return nil, sdkerrors.Wrapf(types.ErrPriceOutOfRange, "max sell price is out of range [%s, %s]", lowestPrice, highestPrice)
		}
	}
	if msg.BuyAmount.IsPositive() {
		if msg.MinBuyPrice.LT(lowestPrice) || msg.MinBuyPrice.GT(highestPrice) {
			return nil, sdkerrors.Wrapf(types.ErrPriceOutOfRange, "min buy price is out of range [%s, %s]", lowestPrice, highestPrice)
		}
		if msg.MaxBuyPrice.LT(lowestPrice) || msg.MaxBuyPrice.GT(highestPrice) {
			return nil, sdkerrors.Wrapf(types.ErrPriceOutOfRange, "max buy price is out of range [%s, %s]", lowestPrice, highestPrice)
		}
	}

	maxNumTicks := int(params.MaxNumMarketMakingOrderTicks)

	var buyTicks, sellTicks []types.MMOrderTick
	offerBaseCoin := sdk.NewInt64Coin(pair.BaseCoinDenom, 0)
	offerQuoteCoin := sdk.NewInt64Coin(pair.QuoteCoinDenom, 0)
	if msg.BuyAmount.IsPositive() {
		buyTicks = types.MMOrderTicks(
			types.OrderDirectionBuy, msg.MinBuyPrice, msg.MaxBuyPrice, msg.BuyAmount, maxNumTicks, tickPrec)
		for _, tick := range buyTicks {
			offerQuoteCoin = offerQuoteCoin.AddAmount(tick.OfferCoinAmount)
		}
	}
	if msg.SellAmount.IsPositive() {
		sellTicks = types.MMOrderTicks(
			types.OrderDirectionSell, msg.MinSellPrice, msg.MaxSellPrice, msg.SellAmount, maxNumTicks, tickPrec)
		for _, tick := range sellTicks {
			offerBaseCoin = offerBaseCoin.AddAmount(tick.OfferCoinAmount)
		}
	}

	orderer := msg.GetOrderer()
	spendable := k.bankKeeper.SpendableCoins(ctx, orderer)
	if spendableAmt := spendable.AmountOf(pair.BaseCoinDenom); spendableAmt.LT(offerBaseCoin.Amount) {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds, "%s is smaller than %s",
			sdk.NewCoin(pair.BaseCoinDenom, spendableAmt), offerBaseCoin)
	}
	if spendableAmt := spendable.AmountOf(pair.QuoteCoinDenom); spendableAmt.LT(offerQuoteCoin.Amount) {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds, "%s is smaller than %s",
			sdk.NewCoin(pair.QuoteCoinDenom, spendableAmt), offerQuoteCoin)
	}

	maxOrderLifespan := params.MaxOrderLifespan
	if msg.OrderLifespan > maxOrderLifespan {
		return nil, sdkerrors.Wrapf(
			types.ErrTooLongOrderLifespan, "%s is longer than %s", msg.OrderLifespan, maxOrderLifespan)
	}

	// First, cancel existing market making orders in the pair from the orderer.
	canceledOrderIds, err := k.cancelMMOrder(ctx, msg.AppId, orderer, pair, true)
	if err != nil {
		return nil, err
	}

	if err := k.bankKeeper.SendCoins(ctx, orderer, pair.GetEscrowAddress(), sdk.NewCoins(offerBaseCoin, offerQuoteCoin)); err != nil {
		return nil, err
	}

	expireAt := ctx.BlockTime().Add(msg.OrderLifespan)
	lastOrderId := pair.LastOrderId

	var orderIds []uint64
	for _, tick := range buyTicks {
		lastOrderId++
		offerCoin := sdk.NewCoin(pair.QuoteCoinDenom, tick.OfferCoinAmount)
		order := types.NewOrder(
			types.OrderTypeMM,
			lastOrderId, msg.AppId,
			pair,
			orderer,
			offerCoin,
			tick.Price,
			tick.Amount,
			expireAt,
			ctx.BlockHeight(),
		)
		k.SetOrder(ctx, order.AppId, order)
		k.SetOrderIndex(ctx, order.AppId, order)
		orders = append(orders, order)
		orderIds = append(orderIds, order.Id)
	}
	for _, tick := range sellTicks {
		lastOrderId++
		offerCoin := sdk.NewCoin(pair.BaseCoinDenom, tick.OfferCoinAmount)
		order := types.NewOrder(
			types.OrderTypeMM,
			lastOrderId, msg.AppId,
			pair,
			orderer,
			offerCoin,
			tick.Price,
			tick.Amount,
			expireAt,
			ctx.BlockHeight(),
		)
		k.SetOrder(ctx, order.AppId, order)
		k.SetOrderIndex(ctx, order.AppId, order)
		orders = append(orders, order)
		orderIds = append(orderIds, order.Id)
	}

	pair.LastOrderId = lastOrderId
	k.SetPair(ctx, pair)

	k.SetMMOrderIndex(ctx, msg.AppId, types.NewMMOrderIndex(orderer, msg.AppId, pair.Id, orderIds))

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMMOrder,
			sdk.NewAttribute(types.AttributeKeyOrderer, msg.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairID, strconv.FormatUint(msg.PairId, 10)),
			sdk.NewAttribute(types.AttributeKeyBatchID, strconv.FormatUint(pair.CurrentBatchId, 10)),
			sdk.NewAttribute(types.AttributeKeyOrderIds, types.FormatUint64s(orderIds)),
			sdk.NewAttribute(types.AttributeKeyCanceledOrderIds, types.FormatUint64s(canceledOrderIds)),
		),
	})
	return
}

// ValidateMsgCancelOrder validates types.MsgCancelOrder and returns the order.
func (k Keeper) ValidateMsgCancelOrder(ctx sdk.Context, msg *types.MsgCancelOrder) (order types.Order, err error) {
	_, found := k.assetKeeper.GetApp(ctx, msg.AppId)
	if !found {
		return types.Order{},
			sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", msg.AppId)
	}

	order, found = k.GetOrder(ctx, msg.AppId, msg.PairId, msg.OrderId)
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
	pair, _ := k.GetPair(ctx, msg.AppId, msg.PairId)
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

	ctx.GasMeter().ConsumeGas(types.CancelOrderGas, "CancelOrderGas")

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

// ValidateMsgCancelAllOrders validates types.MsgCancelAllOrders and returns the order.
func (k Keeper) ValidateMsgCancelAllOrders(ctx sdk.Context, msg *types.MsgCancelAllOrders) error {
	_, found := k.assetKeeper.GetApp(ctx, msg.AppId)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", msg.AppId)
	}
	return nil
}

// CancelAllOrders handles types.MsgCancelAllOrders and cancels all orders.
func (k Keeper) CancelAllOrders(ctx sdk.Context, msg *types.MsgCancelAllOrders) error {
	err := k.ValidateMsgCancelAllOrders(ctx, msg)
	if err != nil {
		return err
	}

	orderPairCache := map[uint64]types.Pair{} // maps order's pair id to pair, to cache the result
	pairIdSet := map[uint64]struct{}{}        // set of pairs where to cancel orders
	var pairIds []string                      // needed to emit an event
	for _, pairId := range msg.PairIds {
		pair, found := k.GetPair(ctx, msg.AppId, pairId)
		if !found { // check if the pair exists
			return sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", pairId)
		}
		pairIdSet[pairId] = struct{}{} // add pair id to the set
		pairIds = append(pairIds, strconv.FormatUint(pairId, 10))
		orderPairCache[pairId] = pair // also cache the pair to use at below
	}

	var canceledOrderIds []string
	if err := k.IterateOrdersByOrderer(ctx, msg.AppId, msg.GetOrderer(), func(order types.Order) (stop bool, err error) {
		_, ok := pairIdSet[order.PairId] // is the pair included in the pair set?
		if len(pairIdSet) == 0 || ok {   // pair ids not specified(cancel all), or the pair is in the set
			pair, ok := orderPairCache[order.PairId]
			if !ok {
				pair, _ = k.GetPair(ctx, msg.AppId, order.PairId)
				orderPairCache[order.PairId] = pair
			}
			if order.Status != types.OrderStatusCanceled && order.BatchId < pair.CurrentBatchId {
				if err := k.FinishOrder(ctx, order, types.OrderStatusCanceled); err != nil {
					return false, err
				}
				canceledOrderIds = append(canceledOrderIds, strconv.FormatUint(order.Id, 10))
			}
		}
		return false, nil
	}); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelAllOrders,
			sdk.NewAttribute(types.AttributeKeyOrderer, msg.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairIds, strings.Join(pairIds, ",")),
			sdk.NewAttribute(types.AttributeKeyCanceledOrderIds, strings.Join(canceledOrderIds, ",")),
		),
	})

	return nil
}

func (k Keeper) cancelMMOrder(ctx sdk.Context, appID uint64, orderer sdk.AccAddress, pair types.Pair, skipIfNotFound bool) (canceledOrderIds []uint64, err error) {
	index, found := k.GetMMOrderIndex(ctx, orderer, appID, pair.Id)
	if found {
		for _, orderId := range index.OrderIds {
			order, found := k.GetOrder(ctx, pair.Id, appID, orderId)
			if !found {
				// The order has already been deleted from store.
				continue
			}
			if order.BatchId == pair.CurrentBatchId {
				return nil, sdkerrors.Wrap(types.ErrSameBatch, "couldn't cancel previously placed orders")
			}
			if order.Status.CanBeCanceled() {
				if err := k.FinishOrder(ctx, order, types.OrderStatusCanceled); err != nil {
					return nil, err
				}
				canceledOrderIds = append(canceledOrderIds, order.Id)
			}
		}
		k.DeleteMMOrderIndex(ctx, appID, index)
	} else if !skipIfNotFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrNotFound, "previous market making orders not found")
	}
	return
}

// CancelMMOrder handles types.MsgCancelMMOrder and cancels previous market making
// orders.
func (k Keeper) CancelMMOrder(ctx sdk.Context, msg *types.MsgCancelMMOrder) (canceledOrderIds []uint64, err error) {
	pair, found := k.GetPair(ctx, msg.AppId, msg.PairId)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", msg.PairId)
	}

	canceledOrderIds, err = k.cancelMMOrder(ctx, msg.AppId, msg.GetOrderer(), pair, false)
	if err != nil {
		return
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelMMOrder,
			sdk.NewAttribute(types.AttributeKeyOrderer, msg.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairID, strconv.FormatUint(pair.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyCanceledOrderIds, types.FormatUint64s(canceledOrderIds)),
		),
	})

	return canceledOrderIds, nil
}

func (k Keeper) ExecuteMatching(ctx sdk.Context, pair types.Pair) error {
	params, err := k.GetGenericParams(ctx, pair.AppId)
	if err != nil {
		return sdkerrors.Wrap(err, "params retreval failed")
	}
	ob := amm.NewOrderBook()

	if err := k.IterateOrdersByPair(ctx, pair.AppId, pair.Id, func(order types.Order) (stop bool, err error) {
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
			// TODO: add orders only when price is in the range?
			ob.AddOrder(types.NewUserOrder(order))
			if order.Status == types.OrderStatusNotExecuted {
				order.SetStatus(types.OrderStatusNotMatched)
				k.SetOrder(ctx, pair.AppId, order)
			}
		case types.OrderStatusCanceled:
		default:
			return false, fmt.Errorf("invalid order status: %s", order.Status)
		}
		return false, nil
	}); err != nil {
		return err
	}

	var pools []*types.PoolOrderer
	_ = k.IteratePoolsByPair(ctx, pair.AppId, pair.Id, func(pool types.Pool) (stop bool, err error) {
		if pool.Disabled {
			return false, nil
		}
		rx, ry := k.getPoolBalances(ctx, pool, pair)
		ps := k.GetPoolCoinSupply(ctx, pool)
		ammPool := types.NewPoolOrderer(
			pool.AMMPool(rx.Amount, ry.Amount, ps),
			pool.Id, pool.GetReserveAddress(), pair.BaseCoinDenom, pair.QuoteCoinDenom)
		if ammPool.IsDepleted() {
			k.MarkPoolAsDisabled(ctx, pool)
			return false, nil
		}
		pools = append(pools, ammPool)
		return false, nil
	})

	matchPrice, quoteCoinDiff, matched := k.Match(ctx, params, ob, pools, pair.LastPrice)
	if matched {
		orders := ob.Orders()
		if err := k.ApplyMatchResult(ctx, pair, orders, quoteCoinDiff); err != nil {
			return err
		}
		pair.LastPrice = &matchPrice
	}

	pair.CurrentBatchId++
	k.SetPair(ctx, pair)

	return nil
}

func (k Keeper) Match(ctx sdk.Context, params types.GenericParams, ob *amm.OrderBook, pools []*types.PoolOrderer, lastPrice *sdk.Dec) (matchPrice sdk.Dec, quoteCoinDiff sdk.Int, matched bool) {
	tickPrec := int(params.TickPrecision)
	if lastPrice == nil {
		ov := amm.MultipleOrderViews{ob.MakeView()}
		for _, pool := range pools {
			ov = append(ov, pool)
		}
		var found bool
		matchPrice, found = amm.FindMatchPrice(ov, tickPrec)
		if !found {
			return sdk.Dec{}, sdk.Int{}, false
		}
		for _, pool := range pools {
			buyAmt := pool.BuyAmountOver(matchPrice, true)
			if buyAmt.IsPositive() {
				ob.AddOrder(pool.Order(amm.Buy, matchPrice, buyAmt))
			}
			sellAmt := pool.SellAmountUnder(matchPrice, true)
			if sellAmt.IsPositive() {
				ob.AddOrder(pool.Order(amm.Sell, matchPrice, sellAmt))
			}
		}
		quoteCoinDiff, matched = ob.MatchAtSinglePrice(matchPrice)
	} else {
		lowestPrice, highestPrice := k.PriceLimits(ctx, *lastPrice, params)
		for _, pool := range pools {
			poolOrders := amm.PoolOrders(pool, pool, lowestPrice, highestPrice, tickPrec)
			ob.AddOrder(poolOrders...)
		}
		matchPrice, quoteCoinDiff, matched = ob.Match(*lastPrice)
	}
	return
}

func (k Keeper) ApplyMatchResult(ctx sdk.Context, pair types.Pair, orders []amm.Order, quoteCoinDiff sdk.Int) error {
	params, err := k.GetGenericParams(ctx, pair.AppId)
	if err != nil {
		return sdkerrors.Wrap(err, "params retreval failed")
	}
	bulkOp := types.NewBulkSendCoinsOperation()
	for _, order := range orders { // TODO: need optimization to filter matched orders only
		order, ok := order.(*types.PoolOrder)
		if !ok {
			continue
		}
		if !order.IsMatched() {
			continue
		}
		paidCoin := sdk.NewCoin(order.OfferCoinDenom, order.PaidOfferCoinAmount)
		bulkOp.QueueSendCoins(order.ReserveAddress, pair.GetEscrowAddress(), sdk.NewCoins(paidCoin))
	}
	if err := bulkOp.Run(ctx, k.bankKeeper); err != nil {
		return err
	}
	bulkOp = types.NewBulkSendCoinsOperation()
	type PoolMatchResult struct {
		PoolId         uint64
		OrderDirection types.OrderDirection
		PaidCoin       sdk.Coin
		ReceivedCoin   sdk.Coin
		MatchedAmount  sdk.Int
	}
	poolMatchResultById := map[uint64]*PoolMatchResult{}
	var poolMatchResults []*PoolMatchResult
	for _, order := range orders {
		if !order.IsMatched() {
			continue
		}

		matchedAmt := order.GetAmount().Sub(order.GetOpenAmount())

		switch order := order.(type) {
		case *types.UserOrder:
			paidCoin := sdk.NewCoin(order.OfferCoinDenom, order.PaidOfferCoinAmount)
			receivedCoin := sdk.NewCoin(order.DemandCoinDenom, order.ReceivedDemandCoinAmount)

			o, _ := k.GetOrder(ctx, pair.AppId, pair.Id, order.OrderID)
			o.OpenAmount = o.OpenAmount.Sub(matchedAmt)
			o.RemainingOfferCoin = o.RemainingOfferCoin.Sub(paidCoin)
			o.ReceivedCoin = o.ReceivedCoin.Add(receivedCoin)

			if o.OpenAmount.IsZero() {
				if err := k.FinishOrder(ctx, o, types.OrderStatusCompleted); err != nil {
					return err
				}
			} else {
				o.SetStatus(types.OrderStatusPartiallyMatched)
				k.SetOrder(ctx, o.AppId, o)
			}
			bulkOp.QueueSendCoins(pair.GetEscrowAddress(), order.Orderer, sdk.NewCoins(receivedCoin))

			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventTypeUserOrderMatched,
					sdk.NewAttribute(types.AttributeKeyOrderDirection, types.OrderDirectionFromAMM(order.Direction).String()),
					sdk.NewAttribute(types.AttributeKeyOrderer, order.Orderer.String()),
					sdk.NewAttribute(types.AttributeKeyPairID, strconv.FormatUint(pair.Id, 10)),
					sdk.NewAttribute(types.AttributeKeyOrderID, strconv.FormatUint(order.OrderID, 10)),
					sdk.NewAttribute(types.AttributeKeyMatchedAmount, matchedAmt.String()),
					sdk.NewAttribute(types.AttributeKeyPaidCoin, paidCoin.String()),
					sdk.NewAttribute(types.AttributeKeyReceivedCoin, receivedCoin.String()),
				),
			})
		case *types.PoolOrder:
			paidCoin := sdk.NewCoin(order.OfferCoinDenom, order.PaidOfferCoinAmount)
			receivedCoin := sdk.NewCoin(order.DemandCoinDenom, order.ReceivedDemandCoinAmount)

			bulkOp.QueueSendCoins(pair.GetEscrowAddress(), order.ReserveAddress, sdk.NewCoins(receivedCoin))

			r, ok := poolMatchResultById[order.PoolID]
			if !ok {
				r = &PoolMatchResult{
					PoolId:         order.PoolID,
					OrderDirection: types.OrderDirectionFromAMM(order.Direction),
					PaidCoin:       sdk.NewCoin(paidCoin.Denom, sdk.ZeroInt()),
					ReceivedCoin:   sdk.NewCoin(receivedCoin.Denom, sdk.ZeroInt()),
					MatchedAmount:  sdk.ZeroInt(),
				}
				poolMatchResultById[order.PoolID] = r
				poolMatchResults = append(poolMatchResults, r)
			}
			dir := types.OrderDirectionFromAMM(order.Direction)
			if r.OrderDirection != dir {
				panic(fmt.Errorf("wrong order direction: %s != %s", dir, r.OrderDirection))
			}
			r.PaidCoin = r.PaidCoin.Add(paidCoin)
			r.ReceivedCoin = r.ReceivedCoin.Add(receivedCoin)
			r.MatchedAmount = r.MatchedAmount.Add(matchedAmt)
		default:
			panic(fmt.Errorf("invalid order type: %T", order))
		}
	}
	dustCollectorAddr, _ := sdk.AccAddressFromBech32(params.DustCollectorAddress)
	bulkOp.QueueSendCoins(pair.GetEscrowAddress(), dustCollectorAddr, sdk.NewCoins(sdk.NewCoin(pair.QuoteCoinDenom, quoteCoinDiff)))
	if err := bulkOp.Run(ctx, k.bankKeeper); err != nil {
		return err
	}
	for _, r := range poolMatchResults {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypePoolOrderMatched,
				sdk.NewAttribute(types.AttributeKeyOrderDirection, r.OrderDirection.String()),
				sdk.NewAttribute(types.AttributeKeyPairID, strconv.FormatUint(pair.Id, 10)),
				sdk.NewAttribute(types.AttributeKeyPoolID, strconv.FormatUint(r.PoolId, 10)),
				sdk.NewAttribute(types.AttributeKeyMatchedAmount, r.MatchedAmount.String()),
				sdk.NewAttribute(types.AttributeKeyPaidCoin, r.PaidCoin.String()),
				sdk.NewAttribute(types.AttributeKeyReceivedCoin, r.ReceivedCoin.String()),
			),
		})
	}
	return nil
}

func (k Keeper) FinishOrder(ctx sdk.Context, order types.Order, status types.OrderStatus) error {
	if order.Type == types.OrderTypeMM {
		return k.FinishMMOrder(ctx, order, status)
	}
	if order.Status == types.OrderStatusCompleted || order.Status.IsCanceledOrExpired() { // sanity check
		return nil
	}

	params, err := k.GetGenericParams(ctx, order.AppId)
	if err != nil {
		return sdkerrors.Wrap(err, "params retreval failed")
	}

	pair, _ := k.GetPair(ctx, order.AppId, order.PairId)

	accumulatedSwapFee := sdk.NewCoin(order.OfferCoin.Denom, sdk.NewInt(0))
	collectedSwapFeeAmountFromOrderer := CalculateSwapFeeAmount(ctx, params, order.OfferCoin.Amount)

	if order.RemainingOfferCoin.IsPositive() {
		refundCoin := order.RemainingOfferCoin

		if order.RemainingOfferCoin.IsEqual(order.OfferCoin) {
			// refund full swap fees back to orderer
			refundCoin.Amount = refundCoin.Amount.Add(collectedSwapFeeAmountFromOrderer)
		} else {
			// refund partial swap fees back to orderer and transfer remaining to to swap fee collector address
			swappedCoin := order.OfferCoin.Sub(order.RemainingOfferCoin)
			swapFeeAmt := CalculateSwapFeeAmount(ctx, params, swappedCoin.Amount)

			accumulatedSwapFee.Amount = accumulatedSwapFee.Amount.Add(swapFeeAmt)

			refundableSwapFeeAmt := collectedSwapFeeAmountFromOrderer.Sub(swapFeeAmt)
			refundCoin.Amount = refundCoin.Amount.Add(refundableSwapFeeAmt)
		}

		if err := k.bankKeeper.SendCoins(ctx, pair.GetEscrowAddress(), order.GetOrderer(), sdk.NewCoins(refundCoin)); err != nil {
			return err
		}
	} else {
		accumulatedSwapFee.Amount = accumulatedSwapFee.Amount.Add(collectedSwapFeeAmountFromOrderer)
	}

	if accumulatedSwapFee.IsPositive() {
		if err := k.bankKeeper.SendCoins(ctx, pair.GetEscrowAddress(), pair.GetSwapFeeCollectorAddress(), sdk.NewCoins(accumulatedSwapFee)); err != nil {
			return err
		}
	}

	order.SetStatus(status)
	k.SetOrder(ctx, order.AppId, order)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeOrderResult,
			sdk.NewAttribute(types.AttributeKeyOrderDirection, order.Direction.String()),
			sdk.NewAttribute(types.AttributeKeyOrderer, order.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairID, strconv.FormatUint(order.PairId, 10)),
			sdk.NewAttribute(types.AttributeKeyOrderID, strconv.FormatUint(order.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyAmount, order.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyOpenAmount, order.OpenAmount.String()),
			sdk.NewAttribute(types.AttributeKeyOfferCoin, order.OfferCoin.String()),
			sdk.NewAttribute(types.AttributeKeyRemainingOfferCoin, order.RemainingOfferCoin.String()),
			sdk.NewAttribute(types.AttributeKeyReceivedCoin, order.ReceivedCoin.String()),
			sdk.NewAttribute(types.AttributeKeyStatus, order.Status.String()),
		),
	})

	return nil
}

func (k Keeper) FinishMMOrder(ctx sdk.Context, order types.Order, status types.OrderStatus) error {
	if order.Status == types.OrderStatusCompleted || order.Status.IsCanceledOrExpired() { // sanity check
		return nil
	}

	if order.RemainingOfferCoin.IsPositive() {
		pair, _ := k.GetPair(ctx, order.AppId, order.PairId)
		if err := k.bankKeeper.SendCoins(ctx, pair.GetEscrowAddress(), order.GetOrderer(), sdk.NewCoins(order.RemainingOfferCoin)); err != nil {
			return err
		}
	}

	order.SetStatus(status)
	k.SetOrder(ctx, order.AppId, order)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeOrderResult,
			sdk.NewAttribute(types.AttributeKeyOrderDirection, order.Direction.String()),
			sdk.NewAttribute(types.AttributeKeyOrderer, order.Orderer),
			sdk.NewAttribute(types.AttributeKeyPairID, strconv.FormatUint(order.PairId, 10)),
			sdk.NewAttribute(types.AttributeKeyOrderID, strconv.FormatUint(order.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyAmount, order.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyOpenAmount, order.OpenAmount.String()),
			sdk.NewAttribute(types.AttributeKeyOfferCoin, order.OfferCoin.String()),
			sdk.NewAttribute(types.AttributeKeyRemainingOfferCoin, order.RemainingOfferCoin.String()),
			sdk.NewAttribute(types.AttributeKeyReceivedCoin, order.ReceivedCoin.String()),
			sdk.NewAttribute(types.AttributeKeyStatus, order.Status.String()),
		),
	})

	return nil
}

// ConvertAccumulatedSwapFeesWithSwapDistrToken swaps accumulated swap fees from -
// pair swap fee accmulator into actual distribution coin .
func (k Keeper) ConvertAccumulatedSwapFeesWithSwapDistrToken(ctx sdk.Context, appID uint64) {
	logger := k.Logger(ctx)

	params, err := k.GetGenericParams(ctx, appID)
	if err != nil {
		return
	}

	availablePools := k.GetAllPools(ctx, appID)
	const poolMapPrefix = "pool_"

	var edges [][]string
	pairPoolIDMap := make(map[string]uint64)

	for _, pool := range availablePools {
		pair, found := k.GetPair(ctx, pool.AppId, pool.PairId)
		if !found {
			continue
		}
		edges = append(edges, []string{pair.BaseCoinDenom, pair.QuoteCoinDenom})
		pairPoolIDMap[pair.BaseCoinDenom+pair.QuoteCoinDenom] = pair.Id
		pairPoolIDMap[pair.QuoteCoinDenom+pair.BaseCoinDenom] = pair.Id
		pairPoolIDMap[poolMapPrefix+pair.BaseCoinDenom+pair.QuoteCoinDenom] = pool.Id
		pairPoolIDMap[poolMapPrefix+pair.QuoteCoinDenom+pair.BaseCoinDenom] = pool.Id
	}

	undirectedGraph := types.BuildUndirectedGraph(edges)

	for _, pool := range availablePools {
		if pool.Disabled {
			continue
		}

		pair, found := k.GetPair(ctx, pool.AppId, pool.PairId)
		if !found {
			continue
		}

		availableBalances := k.bankKeeper.GetAllBalances(ctx, pair.GetSwapFeeCollectorAddress())

		for _, balance := range availableBalances {
			if balance.Denom != params.SwapFeeDistrDenom {
				shortestPath, found := types.BfsShortestPath(undirectedGraph, balance.Denom, params.SwapFeeDistrDenom)
				if found && len(shortestPath) > 1 {
					swappablePairID := pairPoolIDMap[shortestPath[0]+shortestPath[1]]
					swappablePoolID := pairPoolIDMap[poolMapPrefix+shortestPath[0]+shortestPath[1]]

					swappablePair, found := k.GetPair(ctx, appID, swappablePairID)
					if !found {
						continue
					}
					_, found = k.GetPool(ctx, appID, swappablePoolID)
					if !found {
						continue
					}

					swapFeeCoin := sdk.NewCoin(balance.Denom, CalculateSwapFeeAmount(ctx, params, balance.Amount))
					swapFeeCoin.Amount = swapFeeCoin.Amount.Mul(sdk.NewInt(3))
					// reserving extra for swap fee from the offer coin, i.e swapfee *3
					offerCoin := sdk.NewCoin(balance.Denom, balance.Amount.Sub(swapFeeCoin.Amount))

					orderDirection := types.OrderDirectionBuy
					// if balanceDenom is baseCoin in pair, order direction is sell (swap into quote coin)
					// else order direction is buy (swap into base coin)
					if balance.Denom == swappablePair.BaseCoinDenom {
						orderDirection = types.OrderDirectionSell
					}

					lastPrice := *pair.LastPrice
					var amount sdk.Int
					var demandCoinDenom string
					switch orderDirection {
					case types.OrderDirectionBuy:
						maxPrice := lastPrice.Mul(sdk.OneDec().Add(params.MaxPriceLimitRatio))
						amount = offerCoin.Amount.ToDec().Quo(maxPrice).TruncateInt()
						demandCoinDenom = swappablePair.BaseCoinDenom
					case types.OrderDirectionSell:
						amount = offerCoin.Amount
						demandCoinDenom = swappablePair.QuoteCoinDenom
					}
					offerCoin = offerCoin.Add(swapFeeCoin)

					msgMarketOrderMsg := types.NewMsgMarketOrder(
						appID,
						pair.GetSwapFeeCollectorAddress(),
						swappablePairID,
						orderDirection,
						offerCoin,
						demandCoinDenom,
						amount,
						0,
					)

					_, err := k.MarketOrder(ctx, msgMarketOrderMsg)
					if err != nil {
						logger.Info(fmt.Sprintf("warning - swap fee conversion : %v", err))
					}
				}
			}
		}
	}
}
