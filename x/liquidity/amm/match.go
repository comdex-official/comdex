package amm

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
)

// PriceDirection specifies estimated price direction within this batch.
type PriceDirection int

const (
	PriceStaying PriceDirection = iota + 1
	PriceIncreasing
	PriceDecreasing
)

func (dir PriceDirection) String() string {
	switch dir {
	case PriceStaying:
		return "PriceStaying"
	case PriceIncreasing:
		return "PriceIncreasing"
	case PriceDecreasing:
		return "PriceDecreasing"
	default:
		return fmt.Sprintf("PriceDirection(%d)", dir)
	}
}

// FillOrder fills the order by given amount and price.
func FillOrder(order Order, amt sdkmath.Int, price sdkmath.LegacyDec) (quoteCoinDiff sdkmath.Int) {
	matchableAmt := MatchableAmount(order, price)
	if amt.GT(matchableAmt) {
		panic(fmt.Errorf("cannot match more than open amount; %s > %s", amt, matchableAmt))
	}
	var paid, received sdkmath.Int
	switch order.GetDirection() {
	case Buy:
		paid = price.MulInt(amt).Ceil().TruncateInt()
		received = amt
		quoteCoinDiff = paid
	case Sell:
		paid = amt
		received = price.MulInt(amt).TruncateInt()
		quoteCoinDiff = received.Neg()
	}
	order.SetPaidOfferCoinAmount(order.GetPaidOfferCoinAmount().Add(paid))
	order.SetReceivedDemandCoinAmount(order.GetReceivedDemandCoinAmount().Add(received))
	order.SetOpenAmount(order.GetOpenAmount().Sub(amt))
	return
}

// FulfillOrder fills the order by its remaining open amount at given price.
func FulfillOrder(order Order, price sdkmath.LegacyDec) (quoteCoinDiff sdkmath.Int) {
	quoteCoinDiff = sdkmath.ZeroInt()
	matchableAmt := MatchableAmount(order, price)
	if matchableAmt.IsPositive() {
		quoteCoinDiff = quoteCoinDiff.Add(FillOrder(order, matchableAmt, price))
	}
	return
}

// FulfillOrders fills multiple orders by their remaining open amount
// at given price.
func FulfillOrders(orders []Order, price sdkmath.LegacyDec) (quoteCoinDiff sdkmath.Int) {
	quoteCoinDiff = sdkmath.ZeroInt()
	for _, order := range orders {
		quoteCoinDiff = quoteCoinDiff.Add(FulfillOrder(order, price))
	}
	return
}

func FindMatchPrice(ov OrderView, tickPrec int) (matchPrice sdkmath.LegacyDec, found bool) {
	highestBuyPrice, found := ov.HighestBuyPrice()
	if !found {
		return sdkmath.LegacyDec{}, false
	}
	lowestSellPrice, found := ov.LowestSellPrice()
	if !found {
		return sdkmath.LegacyDec{}, false
	}
	if highestBuyPrice.LT(lowestSellPrice) {
		return sdkmath.LegacyDec{}, false
	}

	prec := TickPrecision(tickPrec)
	lowestTickIdx := prec.TickToIndex(prec.LowestTick())
	highestTickIdx := prec.TickToIndex(prec.HighestTick())
	var i, j int
	i, found = findFirstTrueCondition(lowestTickIdx, highestTickIdx, func(i int) bool {
		sellAmt := ov.SellAmountUnder(prec.TickFromIndex(i), true)
		return sellAmt.IsPositive() && ov.BuyAmountOver(prec.TickFromIndex(i+1), true).LTE(sellAmt)
	})
	if !found {
		return sdkmath.LegacyDec{}, false
	}
	j, found = findFirstTrueCondition(highestTickIdx, lowestTickIdx, func(i int) bool {
		buyAmt := ov.BuyAmountOver(prec.TickFromIndex(i), true)
		return buyAmt.IsPositive() && buyAmt.GTE(ov.SellAmountUnder(prec.TickFromIndex(i-1), true))
	})
	if !found {
		return sdkmath.LegacyDec{}, false
	}
	midTick := TickFromIndex(i, tickPrec).Add(TickFromIndex(j, tickPrec)).QuoInt64(2)
	return RoundPrice(midTick, tickPrec), true
}

// FindMatchableAmountAtSinglePrice returns the largest matchable amount of orders
// when matching orders at single price(batch auction).
func (ob *OrderBook) FindMatchableAmountAtSinglePrice(matchPrice sdkmath.LegacyDec) (matchableAmt sdkmath.Int, found bool) {
	type Side struct {
		ticks             []*orderBookTick
		totalMatchableAmt sdkmath.Int
		i                 int
		partialMatchAmt   sdkmath.Int
	}
	buildSide := func(ticks []*orderBookTick, priceIncreasing bool) (side *Side) {
		side = &Side{totalMatchableAmt: zeroInt}
		for i, tick := range ticks {
			if (priceIncreasing && tick.price.GT(matchPrice)) ||
				(!priceIncreasing && tick.price.LT(matchPrice)) {
				break
			}
			side.ticks = ticks[:i+1]
			side.totalMatchableAmt = side.totalMatchableAmt.Add(TotalMatchableAmount(tick.orders, matchPrice))
		}
		side.i = len(side.ticks) - 1
		return
	}
	buySide := buildSide(ob.buys.ticks, ob.buys.priceIncreasing)
	if len(buySide.ticks) == 0 {
		return sdkmath.Int{}, false
	}
	sellSide := buildSide(ob.sells.ticks, ob.sells.priceIncreasing)
	if len(sellSide.ticks) == 0 {
		return sdkmath.Int{}, false
	}
	sides := map[OrderDirection]*Side{
		Buy:  buySide,
		Sell: sellSide,
	}
	// Repeatedly check both buy/sell side to see if there is an order to drop.
	// If there is not, then the loop is finished.
	for {
		ok := true
		for _, dir := range []OrderDirection{Buy, Sell} {
			side := sides[dir]
			i := side.i
			tick := side.ticks[i]
			tickAmt := TotalMatchableAmount(tick.orders, matchPrice)
			// side.partialMatchAmt can be negative at this moment, but
			// FindMatchableAmountAtSinglePrice won't return a negative amount because
			// the if-block below would set ok = false if otherTicksAmt >= matchAmt
			// and the loop would be continued.
			matchableAmt = sdkmath.MinInt(buySide.totalMatchableAmt, sellSide.totalMatchableAmt)
			otherTicksAmt := side.totalMatchableAmt.Sub(tickAmt)
			side.partialMatchAmt = matchableAmt.Sub(otherTicksAmt)
			if otherTicksAmt.GTE(matchableAmt) ||
				(dir == Sell && matchPrice.MulInt(side.partialMatchAmt).TruncateInt().IsZero()) {
				if i == 0 { // There's no orders left, which means orders are not matchable.
					return sdkmath.Int{}, false
				}
				side.totalMatchableAmt = side.totalMatchableAmt.Sub(tickAmt)
				side.i--
				ok = false
			}
		}
		if ok {
			return matchableAmt, true
		}
	}
}

// MatchAtSinglePrice matches all matchable orders(buy orders with higher(or equal) price
// than the price and sell orders with lower(or equal) price than the price)
// at the price.
func (ob *OrderBook) MatchAtSinglePrice(matchPrice sdkmath.LegacyDec) (quoteCoinDiff sdkmath.Int, matched bool) {
	matchableAmt, found := ob.FindMatchableAmountAtSinglePrice(matchPrice)
	if !found {
		return sdkmath.Int{}, false
	}
	quoteCoinDiff = sdkmath.ZeroInt()
	distributeToTicks := func(ticks []*orderBookTick) {
		remainingAmt := matchableAmt
		for _, tick := range ticks {
			tickAmt := TotalMatchableAmount(tick.orders, matchPrice)
			if tickAmt.LTE(remainingAmt) {
				quoteCoinDiff = quoteCoinDiff.Add(FulfillOrders(tick.orders, matchPrice))
				remainingAmt = remainingAmt.Sub(tickAmt)
				if remainingAmt.IsZero() {
					break
				}
			} else {
				quoteCoinDiff = quoteCoinDiff.Add(DistributeOrderAmountToTick(tick, remainingAmt, matchPrice))
				break
			}
		}
	}
	distributeToTicks(ob.buys.ticks)
	distributeToTicks(ob.sells.ticks)
	matched = true
	return
}

// PriceDirection returns the estimated price direction within this batch
// considering the last price.
func (ob *OrderBook) PriceDirection(lastPrice sdkmath.LegacyDec) PriceDirection {
	// TODO: use OrderBookView
	buyAmtOverLastPrice := sdkmath.ZeroInt()
	buyAmtAtLastPrice := sdkmath.ZeroInt()
	for _, tick := range ob.buys.ticks {
		if tick.price.LT(lastPrice) {
			break
		}
		amt := TotalMatchableAmount(tick.orders, lastPrice)
		if tick.price.Equal(lastPrice) {
			buyAmtAtLastPrice = amt
			break
		}
		buyAmtOverLastPrice = buyAmtOverLastPrice.Add(amt)
	}
	sellAmtUnderLastPrice := sdkmath.ZeroInt()
	sellAmtAtLastPrice := sdkmath.ZeroInt()
	for _, tick := range ob.sells.ticks {
		if tick.price.GT(lastPrice) {
			break
		}
		amt := TotalMatchableAmount(tick.orders, lastPrice)
		if tick.price.Equal(lastPrice) {
			sellAmtAtLastPrice = amt
			break
		}
		sellAmtUnderLastPrice = sellAmtUnderLastPrice.Add(amt)
	}
	switch {
	case buyAmtOverLastPrice.GT(sellAmtAtLastPrice.Add(sellAmtUnderLastPrice)):
		return PriceIncreasing
	case sellAmtUnderLastPrice.GT(buyAmtAtLastPrice.Add(buyAmtOverLastPrice)):
		return PriceDecreasing
	default:
		return PriceStaying
	}
}

// Match matches orders sequentially, starting from buy orders with the highest price
// and sell orders with the lowest price.
// The matching continues until there's no more matchable orders.
func (ob *OrderBook) Match(lastPrice sdkmath.LegacyDec) (matchPrice sdkmath.LegacyDec, quoteCoinDiff sdkmath.Int, matched bool) {
	if len(ob.buys.ticks) == 0 || len(ob.sells.ticks) == 0 {
		return sdkmath.LegacyDec{}, sdkmath.Int{}, false
	}
	matchPrice = lastPrice
	dir := ob.PriceDirection(lastPrice)
	quoteCoinDiff, matched = ob.MatchAtSinglePrice(lastPrice)
	if dir == PriceStaying {
		return matchPrice, quoteCoinDiff, matched
	}
	if !matched {
		quoteCoinDiff = sdkmath.ZeroInt()
	}
	bi, si := 0, 0
	for bi < len(ob.buys.ticks) && si < len(ob.sells.ticks) && ob.buys.ticks[bi].price.GTE(ob.sells.ticks[si].price) {
		buyTick := ob.buys.ticks[bi]
		sellTick := ob.sells.ticks[si]
		var p sdkmath.LegacyDec
		switch dir {
		case PriceIncreasing:
			p = sellTick.price
		case PriceDecreasing:
			p = buyTick.price
		}
		buyTickOpenAmt := TotalMatchableAmount(buyTick.orders, p)
		sellTickOpenAmt := TotalMatchableAmount(sellTick.orders, p)
		if !buyTickOpenAmt.IsPositive() {
			bi++
			continue
		}
		if !sellTickOpenAmt.IsPositive() {
			si++
			continue
		}
		if buyTickOpenAmt.LTE(sellTickOpenAmt) {
			quoteCoinDiff = quoteCoinDiff.Add(DistributeOrderAmountToTick(buyTick, buyTickOpenAmt, p))
			bi++
		} else {
			quoteCoinDiff = quoteCoinDiff.Add(DistributeOrderAmountToTick(buyTick, sellTickOpenAmt, p))
		}
		if sellTickOpenAmt.LTE(buyTickOpenAmt) {
			quoteCoinDiff = quoteCoinDiff.Add(DistributeOrderAmountToTick(sellTick, sellTickOpenAmt, p))
			si++
		} else {
			quoteCoinDiff = quoteCoinDiff.Add(DistributeOrderAmountToTick(sellTick, buyTickOpenAmt, p))
		}
		matchPrice = p
		matched = true
	}
	return
}

// DistributeOrderAmountToTick distributes the given order amount to the orders
// at the tick.
// Orders with higher priority(have lower batch id) get matched first,
// then the remaining amount is distributed to the remaining orders.
func DistributeOrderAmountToTick(tick *orderBookTick, amt sdkmath.Int, price sdkmath.LegacyDec) (quoteCoinDiff sdkmath.Int) {
	remainingAmt := amt
	quoteCoinDiff = sdkmath.ZeroInt()
	groups := GroupOrdersByBatchID(tick.orders)
	for _, group := range groups {
		openAmt := TotalMatchableAmount(group.Orders, price)
		if openAmt.IsZero() {
			continue
		}
		if remainingAmt.GTE(openAmt) {
			quoteCoinDiff = quoteCoinDiff.Add(FulfillOrders(group.Orders, price))
			remainingAmt = remainingAmt.Sub(openAmt)
		} else {
			SortOrders(group.Orders)
			quoteCoinDiff = quoteCoinDiff.Add(DistributeOrderAmountToOrders(group.Orders, remainingAmt, price))
			remainingAmt = sdkmath.ZeroInt()
		}
		if remainingAmt.IsZero() {
			break
		}
	}
	return
}

// DistributeOrderAmountToOrders distributes the given order amount to the orders
// proportional to each order's amount.
// The caller must sort orders before calling DistributeOrderAmountToOrders.
// After distributing the amount based on each order's proportion,
// remaining amount due to the decimal truncation is distributed
// to the orders again, by priority.
// This time, the proportion is not considered and each order takes up
// the amount as much as possible.
func DistributeOrderAmountToOrders(orders []Order, amt sdkmath.Int, price sdkmath.LegacyDec) (quoteCoinDiff sdkmath.Int) {
	totalAmt := TotalAmount(orders)
	totalMatchedAmt := sdkmath.ZeroInt()
	matchedAmtByOrder := map[Order]sdkmath.Int{}

	for _, order := range orders {
		matchableAmt := MatchableAmount(order, price)
		if matchableAmt.IsZero() {
			continue
		}
		orderAmt := order.GetAmount().ToLegacyDec()
		proportion := orderAmt.QuoTruncate(totalAmt.ToLegacyDec())
		matchedAmt := sdkmath.MinInt(matchableAmt, proportion.MulInt(amt).TruncateInt())
		if matchedAmt.IsPositive() {
			matchedAmtByOrder[order] = matchedAmt
			totalMatchedAmt = totalMatchedAmt.Add(matchedAmt)
		}
	}

	remainingAmt := amt.Sub(totalMatchedAmt)
	for _, order := range orders {
		if remainingAmt.IsZero() {
			break
		}
		prevMatchedAmt, ok := matchedAmtByOrder[order]
		if !ok { // TODO: is it possible?
			prevMatchedAmt = sdkmath.ZeroInt()
		}
		matchableAmt := MatchableAmount(order, price)
		matchedAmt := sdkmath.MinInt(remainingAmt, matchableAmt.Sub(prevMatchedAmt))
		matchedAmtByOrder[order] = prevMatchedAmt.Add(matchedAmt)
		remainingAmt = remainingAmt.Sub(matchedAmt)
	}

	var matchedOrders, notMatchedOrders []Order
	for _, order := range orders {
		matchedAmt, ok := matchedAmtByOrder[order]
		if !ok {
			matchedAmt = sdkmath.ZeroInt()
		}
		if !matchedAmt.IsZero() && (order.GetDirection() == Buy || price.MulInt(matchedAmt).TruncateInt().IsPositive()) {
			matchedOrders = append(matchedOrders, order)
		} else {
			notMatchedOrders = append(notMatchedOrders, order)
		}
	}

	if len(notMatchedOrders) > 0 {
		if len(matchedOrders) == 0 {
			return DistributeOrderAmountToOrders(orders[:len(orders)-1], amt, price)
		}
		return DistributeOrderAmountToOrders(matchedOrders, amt, price)
	}

	quoteCoinDiff = sdkmath.ZeroInt()
	for order, matchedAmt := range matchedAmtByOrder {
		quoteCoinDiff = quoteCoinDiff.Add(FillOrder(order, matchedAmt, price))
	}
	return
}
