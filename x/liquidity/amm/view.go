package amm

import (
	"sort"

	sdkmath "cosmossdk.io/math"
)

var (
	_ OrderView = (*OrderBookView)(nil)
	_ OrderView = Pool(nil)
	_ OrderView = MultipleOrderViews(nil)
)

type OrderView interface {
	HighestBuyPrice() (sdkmath.LegacyDec, bool)
	LowestSellPrice() (sdkmath.LegacyDec, bool)
	BuyAmountOver(price sdkmath.LegacyDec, inclusive bool) sdkmath.Int
	SellAmountUnder(price sdkmath.LegacyDec, inclusive bool) sdkmath.Int
}

type OrderBookView struct {
	buyAmtAccSums, sellAmtAccSums []amtAccSum
}

func (ob *OrderBook) MakeView() *OrderBookView {
	view := &OrderBookView{
		buyAmtAccSums:  make([]amtAccSum, len(ob.buys.ticks)),
		sellAmtAccSums: make([]amtAccSum, len(ob.sells.ticks)),
	}
	for i, tick := range ob.buys.ticks {
		var prevSum sdkmath.Int
		if i == 0 {
			prevSum = sdkmath.ZeroInt()
		} else {
			prevSum = view.buyAmtAccSums[i-1].sum
		}
		view.buyAmtAccSums[i] = amtAccSum{
			price: tick.price,
			sum:   prevSum.Add(TotalMatchableAmount(tick.orders, tick.price)),
		}
	}
	for i, tick := range ob.sells.ticks {
		var prevSum sdkmath.Int
		if i == 0 {
			prevSum = sdkmath.ZeroInt()
		} else {
			prevSum = view.sellAmtAccSums[i-1].sum
		}
		view.sellAmtAccSums[i] = amtAccSum{
			price: tick.price,
			sum:   prevSum.Add(TotalMatchableAmount(tick.orders, tick.price)),
		}
	}
	return view
}

func (view *OrderBookView) Match() {
	if len(view.buyAmtAccSums) == 0 || len(view.sellAmtAccSums) == 0 {
		return
	}
	buyIdx := sort.Search(len(view.buyAmtAccSums), func(i int) bool {
		return view.BuyAmountOver(view.buyAmtAccSums[i].price, true).GT(
			view.SellAmountUnder(view.buyAmtAccSums[i].price, true))
	})
	sellIdx := sort.Search(len(view.sellAmtAccSums), func(i int) bool {
		return view.SellAmountUnder(view.sellAmtAccSums[i].price, true).GT(
			view.BuyAmountOver(view.sellAmtAccSums[i].price, true))
	})
	if buyIdx == len(view.buyAmtAccSums) && sellIdx == len(view.sellAmtAccSums) {
		return
	}
	matchAmt := sdkmath.ZeroInt()
	if buyIdx > 0 {
		matchAmt = view.buyAmtAccSums[buyIdx-1].sum
	}
	if sellIdx > 0 {
		matchAmt = sdkmath.MaxInt(matchAmt, view.sellAmtAccSums[sellIdx-1].sum)
	}
	for i, accSum := range view.buyAmtAccSums {
		if i < buyIdx {
			view.buyAmtAccSums[i].sum = zeroInt
		} else {
			view.buyAmtAccSums[i].sum = accSum.sum.Sub(matchAmt)
		}
	}
	for i, accSum := range view.sellAmtAccSums {
		if i < sellIdx {
			view.sellAmtAccSums[i].sum = zeroInt
		} else {
			view.sellAmtAccSums[i].sum = accSum.sum.Sub(matchAmt)
		}
	}
}

func (view *OrderBookView) HighestBuyPrice() (sdkmath.LegacyDec, bool) {
	if len(view.buyAmtAccSums) == 0 {
		return sdkmath.LegacyDec{}, false
	}
	i := sort.Search(len(view.buyAmtAccSums), func(i int) bool {
		return view.buyAmtAccSums[i].sum.IsPositive()
	})
	if i >= len(view.buyAmtAccSums) {
		return sdkmath.LegacyDec{}, false
	}
	return view.buyAmtAccSums[i].price, true
}

func (view *OrderBookView) LowestSellPrice() (sdkmath.LegacyDec, bool) {
	if len(view.sellAmtAccSums) == 0 {
		return sdkmath.LegacyDec{}, false
	}
	i := sort.Search(len(view.sellAmtAccSums), func(i int) bool {
		return view.sellAmtAccSums[i].sum.IsPositive()
	})
	if i >= len(view.sellAmtAccSums) {
		return sdkmath.LegacyDec{}, false
	}
	return view.sellAmtAccSums[i].price, true
}

func (view *OrderBookView) BuyAmountOver(price sdkmath.LegacyDec, inclusive bool) sdkmath.Int {
	i := sort.Search(len(view.buyAmtAccSums), func(i int) bool {
		if inclusive {
			return view.buyAmtAccSums[i].price.LT(price)
		}
		return view.buyAmtAccSums[i].price.LTE(price)
	})
	if i == 0 {
		return sdkmath.ZeroInt()
	}
	return view.buyAmtAccSums[i-1].sum
}

func (view *OrderBookView) BuyAmountUnder(price sdkmath.LegacyDec, inclusive bool) sdkmath.Int {
	i := sort.Search(len(view.buyAmtAccSums), func(i int) bool {
		if inclusive {
			return view.buyAmtAccSums[i].price.LTE(price)
		}
		return view.buyAmtAccSums[i].price.LT(price)
	})
	if i == 0 {
		return view.buyAmtAccSums[len(view.buyAmtAccSums)-1].sum
	}
	return view.buyAmtAccSums[len(view.buyAmtAccSums)-1].sum.Sub(view.buyAmtAccSums[i-1].sum)
}

func (view *OrderBookView) SellAmountUnder(price sdkmath.LegacyDec, inclusive bool) sdkmath.Int {
	i := sort.Search(len(view.sellAmtAccSums), func(i int) bool {
		if inclusive {
			return view.sellAmtAccSums[i].price.GT(price)
		}
		return view.sellAmtAccSums[i].price.GTE(price)
	})
	if i == 0 {
		return sdkmath.ZeroInt()
	}
	return view.sellAmtAccSums[i-1].sum
}

func (view *OrderBookView) SellAmountOver(price sdkmath.LegacyDec, inclusive bool) sdkmath.Int {
	i := sort.Search(len(view.sellAmtAccSums), func(i int) bool {
		if inclusive {
			return view.sellAmtAccSums[i].price.GTE(price)
		}
		return view.sellAmtAccSums[i].price.GT(price)
	})
	if i == 0 {
		return view.sellAmtAccSums[len(view.sellAmtAccSums)-1].sum
	}
	return view.sellAmtAccSums[len(view.sellAmtAccSums)-1].sum.Sub(view.sellAmtAccSums[i-1].sum)
}

type amtAccSum struct {
	price sdkmath.LegacyDec
	sum   sdkmath.Int
}

type MultipleOrderViews []OrderView

func (views MultipleOrderViews) HighestBuyPrice() (price sdkmath.LegacyDec, found bool) {
	for _, view := range views {
		p, f := view.HighestBuyPrice()
		if f && (price.IsNil() || p.GT(price)) {
			price = p
			found = true
		}
	}
	return
}

func (views MultipleOrderViews) LowestSellPrice() (price sdkmath.LegacyDec, found bool) {
	for _, view := range views {
		p, f := view.LowestSellPrice()
		if f && (price.IsNil() || p.LT(price)) {
			price = p
			found = true
		}
	}
	return
}

func (views MultipleOrderViews) BuyAmountOver(price sdkmath.LegacyDec, inclusive bool) sdkmath.Int {
	amt := sdkmath.ZeroInt()
	for _, view := range views {
		amt = amt.Add(view.BuyAmountOver(price, inclusive))
	}
	return amt
}

func (views MultipleOrderViews) SellAmountUnder(price sdkmath.LegacyDec, inclusive bool) sdkmath.Int {
	amt := sdkmath.ZeroInt()
	for _, view := range views {
		amt = amt.Add(view.SellAmountUnder(price, inclusive))
	}
	return amt
}
