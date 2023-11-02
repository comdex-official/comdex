package amm

import (
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"

	utils "github.com/comdex-official/comdex/types"
)

var (
	_ Pool = (*BasicPool)(nil)
	_ Pool = (*RangedPool)(nil)
)

// Pool is the interface of a pool.
type Pool interface {
	Balances() (rx, ry sdkmath.Int)
	SetBalances(rx, ry sdkmath.Int, derive bool)
	PoolCoinSupply() sdkmath.Int
	Price() sdkmath.LegacyDec
	IsDepleted() bool

	HighestBuyPrice() (sdkmath.LegacyDec, bool)
	LowestSellPrice() (sdkmath.LegacyDec, bool)
	BuyAmountOver(price sdkmath.LegacyDec, inclusive bool) sdkmath.Int
	SellAmountUnder(price sdkmath.LegacyDec, inclusive bool) sdkmath.Int
	BuyAmountTo(price sdkmath.LegacyDec) sdkmath.Int
	SellAmountTo(price sdkmath.LegacyDec) sdkmath.Int

	Clone() Pool
}

// BasicPool is the basic pool type.
type BasicPool struct {
	// rx and ry are the pool's reserve balance of each x/y coin.
	// In perspective of a pair, x coin is the quote coin and
	// y coin is the base coin.
	rx, ry sdkmath.Int
	// ps is the pool's pool coin supply.
	ps sdkmath.Int
}

// NewBasicPool returns a new BasicPool.
// It is OK to pass an empty sdkmath.Int to ps when ps is not going to be used.
func NewBasicPool(rx, ry, ps sdkmath.Int) *BasicPool {
	return &BasicPool{
		rx: rx,
		ry: ry,
		ps: ps,
	}
}

func CreateBasicPool(rx, ry sdkmath.Int) (*BasicPool, error) {
	if rx.IsZero() || ry.IsZero() {
		return nil, fmt.Errorf("cannot create basic pool with zero reserve amount")
	}
	p := rx.ToLegacyDec().Quo(ry.ToLegacyDec())
	if p.LT(MinPoolPrice) {
		return nil, fmt.Errorf("pool price is lower than min price %s", MinPoolPrice)
	}
	if p.GT(MaxPoolPrice) {
		return nil, fmt.Errorf("pool price is greater than max price %s", MaxPoolPrice)
	}
	return NewBasicPool(rx, ry, InitialPoolCoinSupply(rx, ry)), nil
}

// Balances returns the balances of the pool.
func (pool *BasicPool) Balances() (rx, ry sdkmath.Int) {
	return pool.rx, pool.ry
}

func (pool *BasicPool) SetBalances(rx, ry sdkmath.Int, _ bool) {
	pool.rx = rx
	pool.ry = ry
}

// PoolCoinSupply returns the pool coin supply.
func (pool *BasicPool) PoolCoinSupply() sdkmath.Int {
	return pool.ps
}

// Price returns the pool price.
func (pool *BasicPool) Price() sdkmath.LegacyDec {
	if pool.rx.IsZero() || pool.ry.IsZero() {
		panic("pool price is not defined for a depleted pool")
	}
	return pool.rx.ToLegacyDec().Quo(pool.ry.ToLegacyDec())
}

// IsDepleted returns whether the pool is depleted or not.
func (pool *BasicPool) IsDepleted() bool {
	return pool.ps.IsZero() || pool.rx.IsZero() || pool.ry.IsZero()
}

// HighestBuyPrice returns the highest buy price of the pool.
func (pool *BasicPool) HighestBuyPrice() (price sdkmath.LegacyDec, found bool) {
	// The highest buy price is actually a bit lower than pool price,
	// but it's not important for our matching logic.
	return pool.Price(), true
}

// LowestSellPrice returns the lowest sell price of the pool.
func (pool *BasicPool) LowestSellPrice() (price sdkmath.LegacyDec, found bool) {
	// The lowest sell price is actually a bit higher than the pool price,
	// but it's not important for our matching logic.
	return pool.Price(), true
}

// BuyAmountOver returns the amount of buy orders for price greater than
// or equal to given price.
// amt = (X - P*Y)/P
func (pool *BasicPool) BuyAmountOver(price sdkmath.LegacyDec, _ bool) (amt sdkmath.Int) {
	origPrice := price
	if price.LT(MinPoolPrice) {
		price = MinPoolPrice
	}
	if price.GTE(pool.Price()) {
		return zeroInt
	}
	dx := pool.rx.ToLegacyDec().Sub(price.MulInt(pool.ry))
	if !dx.IsPositive() {
		return zeroInt
	}
	utils.SafeMath(func() {
		amt = dx.QuoTruncate(origPrice).TruncateInt()
		if amt.GT(MaxCoinAmount) {
			amt = MaxCoinAmount
		}
	}, func() {
		amt = MaxCoinAmount
	})
	return
}

func (pool *BasicPool) SellAmountUnder(price sdkmath.LegacyDec, _ bool) (amt sdkmath.Int) {
	if price.GT(MaxPoolPrice) {
		price = MaxPoolPrice
	}
	if price.LTE(pool.Price()) {
		return zeroInt
	}
	amt = pool.ry.ToLegacyDec().Sub(pool.rx.ToLegacyDec().QuoRoundUp(price)).TruncateInt()
	if !amt.IsPositive() {
		return zeroInt
	}
	return
}

// BuyAmountTo returns the amount of buy orders of the pool for price,
// where BuyAmountTo is used when the pool price is higher than the highest
// price of the order book.
func (pool *BasicPool) BuyAmountTo(price sdkmath.LegacyDec) (amt sdkmath.Int) {
	origPrice := price
	if price.LT(MinPoolPrice) {
		price = MinPoolPrice
	}
	if price.GTE(pool.Price()) {
		return zeroInt
	}
	sqrtRx := utils.DecApproxSqrt(pool.rx.ToLegacyDec())
	sqrtRy := utils.DecApproxSqrt(pool.ry.ToLegacyDec())
	sqrtPrice := utils.DecApproxSqrt(price)
	dx := pool.rx.ToLegacyDec().Sub(sqrtPrice.Mul(sqrtRx.Mul(sqrtRy))) // dx = rx - sqrt(P * rx * ry)
	if !dx.IsPositive() {
		return zeroInt
	}
	utils.SafeMath(func() {
		amt = dx.QuoTruncate(origPrice).TruncateInt() // dy = dx / P
		if amt.GT(MaxCoinAmount) {
			amt = MaxCoinAmount
		}
	}, func() {
		amt = MaxCoinAmount
	})
	return
}

// SellAmountTo returns the amount of sell orders of the pool for price,
// where SellAmountTo is used when the pool price is lower than the lowest
// price of the order book.
func (pool *BasicPool) SellAmountTo(price sdkmath.LegacyDec) (amt sdkmath.Int) {
	if price.GT(MaxPoolPrice) {
		price = MaxPoolPrice
	}
	if price.LTE(pool.Price()) {
		return zeroInt
	}
	sqrtRx := utils.DecApproxSqrt(pool.rx.ToLegacyDec())
	sqrtRy := utils.DecApproxSqrt(pool.ry.ToLegacyDec())
	sqrtPrice := utils.DecApproxSqrt(price)
	// dy = ry - sqrt(rx * ry / P)
	amt = pool.ry.ToLegacyDec().Sub(sqrtRx.Mul(sqrtRy).Quo(sqrtPrice)).TruncateInt()
	if !amt.IsPositive() {
		return zeroInt
	}
	return
}

func (pool *BasicPool) Clone() Pool {
	return NewBasicPool(pool.rx, pool.ry, pool.ps)
}

type RangedPool struct {
	rx, ry             sdkmath.Int
	ps                 sdkmath.Int
	minPrice, maxPrice sdkmath.LegacyDec
	transX, transY     sdkmath.LegacyDec
	xComp, yComp       sdkmath.LegacyDec
}

// NewRangedPool returns a new RangedPool.
func NewRangedPool(rx, ry, ps sdkmath.Int, minPrice, maxPrice sdkmath.LegacyDec) *RangedPool {
	transX, transY := DeriveTranslation(rx, ry, minPrice, maxPrice)
	return &RangedPool{
		rx:       rx,
		ry:       ry,
		ps:       ps,
		minPrice: minPrice,
		maxPrice: maxPrice,
		transX:   transX,
		transY:   transY,
		xComp:    rx.ToLegacyDec().Add(transX),
		yComp:    ry.ToLegacyDec().Add(transY),
	}
}

// CreateRangedPool creates new RangedPool from given inputs, while validating
// the inputs and using only needed amount of x/y coins(the rest should be refunded).
func CreateRangedPool(x, y sdkmath.Int, minPrice, maxPrice, initialPrice sdkmath.LegacyDec) (pool *RangedPool, err error) {
	if !x.IsPositive() && !y.IsPositive() {
		return nil, fmt.Errorf("either x or y must be positive")
	}
	if err := ValidateRangedPoolParams(minPrice, maxPrice, initialPrice); err != nil {
		return nil, err
	}

	// P = initialPrice, M = minPrice, L = maxPrice
	var ax, ay sdkmath.Int
	switch {
	case initialPrice.Equal(minPrice): // single y asset pool
		ax = zeroInt
		ay = y
	case initialPrice.Equal(maxPrice): // single x asset pool
		ax = x
		ay = zeroInt
	default: // normal pool
		sqrt := utils.DecApproxSqrt
		xDec, yDec := x.ToLegacyDec(), y.ToLegacyDec()
		sqrtP := sqrt(initialPrice) // sqrt(P)
		sqrtM := sqrt(minPrice)     // sqrt(M)
		sqrtL := sqrt(maxPrice)     // sqrt(L)
		// Assume that we can accept all x
		ax = x
		// ay = {x / (sqrt(P)-sqrt(M))} * (1/sqrt(P) - 1/sqrt(L))
		ay = xDec.Quo(sqrtP.Sub(sqrtM)).Mul(inv(sqrtP).Sub(inv(sqrtL))).Ceil().TruncateInt()
		if ay.GT(y) {
			// Accept all y
			// ax = {y / (1/sqrt(P) - 1/sqrt(L))} * (sqrt(P) - sqrt(M))
			ax = yDec.Quo(inv(sqrtP).Sub(inv(sqrtL))).Mul(sqrtP.Sub(sqrtM)).Ceil().TruncateInt()
			ay = y
		}
	}
	return NewRangedPool(ax, ay, InitialPoolCoinSupply(ax, ay), minPrice, maxPrice), nil
}

func ValidateRangedPoolParams(minPrice, maxPrice, initialPrice sdkmath.LegacyDec) error {
	if !initialPrice.IsPositive() {
		return fmt.Errorf("initial price must be positive: %s", initialPrice)
	}
	if minPrice.LT(MinPoolPrice) {
		return fmt.Errorf("min price must not be lower than %s", MinPoolPrice)
	}
	if !maxPrice.IsPositive() {
		return fmt.Errorf("max price must be positive: %s", maxPrice)
	}
	if maxPrice.GT(MaxPoolPrice) {
		return fmt.Errorf("max price must not be higher than %s", MaxPoolPrice)
	}
	if !maxPrice.GT(minPrice) {
		return fmt.Errorf("max price must be higher than min price")
	}
	if maxPrice.Sub(minPrice).Quo(minPrice).LT(MinRangedPoolPriceGapRatio) {
		return fmt.Errorf("min price and max price are too close")
	}
	if initialPrice.LT(minPrice) {
		return fmt.Errorf("initial price must not be lower than min price")
	}
	if initialPrice.GT(maxPrice) {
		return fmt.Errorf("initial price must not be higher than max price")
	}
	return nil
}

// Balances returns the balances of the pool.
func (pool *RangedPool) Balances() (rx, ry sdkmath.Int) {
	return pool.rx, pool.ry
}

// SetBalances sets RangedPool's balances without recalculating
// transX and transY.
func (pool *RangedPool) SetBalances(rx, ry sdkmath.Int, derive bool) {
	if derive {
		pool.transX, pool.transY = DeriveTranslation(rx, ry, pool.minPrice, pool.maxPrice)
	}
	pool.rx = rx
	pool.ry = ry
	pool.xComp = pool.rx.ToLegacyDec().Add(pool.transX)
	pool.yComp = pool.ry.ToLegacyDec().Add(pool.transY)
}

// PoolCoinSupply returns the pool coin supply.
func (pool *RangedPool) PoolCoinSupply() sdkmath.Int {
	return pool.ps
}

func (pool *RangedPool) Translation() (transX, transY sdkmath.LegacyDec) {
	return pool.transX, pool.transY
}

func (pool *RangedPool) MinPrice() sdkmath.LegacyDec {
	return pool.minPrice
}

func (pool *RangedPool) MaxPrice() sdkmath.LegacyDec {
	return pool.maxPrice
}

// Price returns the pool price.
func (pool *RangedPool) Price() sdkmath.LegacyDec {
	if pool.rx.IsZero() && pool.ry.IsZero() {
		panic("pool price is not defined for a depleted pool")
	}
	return pool.xComp.Quo(pool.yComp) // (rx + transX) / (ry + transY)
}

// IsDepleted returns whether the pool is depleted or not.
func (pool *RangedPool) IsDepleted() bool {
	return pool.ps.IsZero() || (pool.rx.IsZero() && pool.ry.IsZero())
}

// HighestBuyPrice returns the highest buy price of the pool.
func (pool *RangedPool) HighestBuyPrice() (price sdkmath.LegacyDec, found bool) {
	// The highest buy price is actually a bit lower than pool price,
	// but it's not important for our matching logic.
	return pool.Price(), true
}

// LowestSellPrice returns the lowest sell price of the pool.
func (pool *RangedPool) LowestSellPrice() (price sdkmath.LegacyDec, found bool) {
	// The lowest sell price is actually a bit higher than the pool price,
	// but it's not important for our matching logic.
	return pool.Price(), true
}

// BuyAmountOver returns the amount of buy orders for price greater than
// or equal to given price.
func (pool *RangedPool) BuyAmountOver(price sdkmath.LegacyDec, _ bool) (amt sdkmath.Int) {
	origPrice := price
	if price.LT(pool.minPrice) {
		price = pool.minPrice
	}
	if price.GTE(pool.Price()) {
		return zeroInt
	}
	// dx = (rx + transX) - P * (ry + transY)
	dx := pool.xComp.Sub(price.Mul(pool.yComp))
	if !dx.IsPositive() {
		return zeroInt
	} else if dx.GT(pool.rx.ToLegacyDec()) {
		dx = pool.rx.ToLegacyDec()
	}
	utils.SafeMath(func() {
		amt = dx.QuoTruncate(origPrice).TruncateInt() // dy = dx / P
		if amt.GT(MaxCoinAmount) {
			amt = MaxCoinAmount
		}
	}, func() {
		amt = MaxCoinAmount
	})
	return
}

// SellAmountUnder returns the amount of sell orders for price less than
// or equal to given price.
func (pool *RangedPool) SellAmountUnder(price sdkmath.LegacyDec, _ bool) (amt sdkmath.Int) {
	if price.GT(pool.maxPrice) {
		price = pool.maxPrice
	}
	if price.LTE(pool.Price()) {
		return zeroInt
	}
	// dy = (ry + transY) - (rx + transX) / P
	amt = pool.yComp.Sub(pool.xComp.QuoRoundUp(price)).TruncateInt()
	if amt.GT(pool.ry) {
		amt = pool.ry
	}
	if !amt.IsPositive() {
		return zeroInt
	}
	return
}

// BuyAmountTo returns the amount of buy orders of the pool for price,
// where BuyAmountTo is used when the pool price is higher than the highest
// price of the order book.
func (pool *RangedPool) BuyAmountTo(price sdkmath.LegacyDec) (amt sdkmath.Int) {
	origPrice := price
	if price.LT(pool.minPrice) {
		price = pool.minPrice
	}
	if price.GTE(pool.Price()) {
		return zeroInt
	}
	sqrtXComp := utils.DecApproxSqrt(pool.xComp)
	sqrtYComp := utils.DecApproxSqrt(pool.yComp)
	sqrtPrice := utils.DecApproxSqrt(price)
	// dx = rx - (sqrt(P * (rx + transX) * (ry + transY)) - transX)
	dx := pool.rx.ToLegacyDec().Sub(sqrtPrice.Mul(sqrtXComp.Mul(sqrtYComp)).Sub(pool.transX))
	if !dx.IsPositive() {
		return zeroInt
	} else if dx.GT(pool.rx.ToLegacyDec()) {
		dx = pool.rx.ToLegacyDec()
	}
	utils.SafeMath(func() {
		amt = dx.QuoTruncate(origPrice).TruncateInt() // dy = dx / P
		if amt.GT(MaxCoinAmount) {
			amt = MaxCoinAmount
		}
	}, func() {
		amt = MaxCoinAmount
	})
	return
}

// SellAmountTo returns the amount of sell orders of the pool for price,
// where SellAmountTo is used when the pool price is lower than the lowest
// price of the order book.
func (pool *RangedPool) SellAmountTo(price sdkmath.LegacyDec) (amt sdkmath.Int) {
	if price.GT(pool.maxPrice) {
		price = pool.maxPrice
	}
	if price.LTE(pool.Price()) {
		return zeroInt
	}
	sqrtXComp := utils.DecApproxSqrt(pool.xComp)
	sqrtYComp := utils.DecApproxSqrt(pool.yComp)
	sqrtPrice := utils.DecApproxSqrt(price)
	// dy = ry - (sqrt((x + transX) * (y + transY) / P) - b)
	amt = pool.ry.ToLegacyDec().Sub(sqrtXComp.Mul(sqrtYComp).QuoRoundUp(sqrtPrice).Sub(pool.transY)).TruncateInt()
	if amt.GT(pool.ry) {
		amt = pool.ry
	}
	if !amt.IsPositive() {
		return zeroInt
	}
	return
}

func (pool *RangedPool) Clone() Pool {
	return &RangedPool{
		rx:       pool.rx,
		ry:       pool.ry,
		ps:       pool.ps,
		transX:   pool.transX,
		transY:   pool.transY,
		xComp:    pool.xComp,
		yComp:    pool.yComp,
		minPrice: pool.minPrice,
		maxPrice: pool.maxPrice,
	}
}

// Deposit returns accepted x and y coin amount and minted pool coin amount
// when someone deposits x and y coins.
func Deposit(rx, ry, ps, x, y sdkmath.Int) (ax, ay, pc sdkmath.Int) {
	// Calculate accepted amount and minting amount.
	// Note that we take as many coins as possible(by ceiling numbers)
	// from depositor and mint as little coins as possible.

	utils.SafeMath(func() {
		rx, ry := rx.ToLegacyDec(), ry.ToLegacyDec()
		ps := ps.ToLegacyDec()

		// pc = floor(ps * min(x / rx, y / ry))
		var ratio sdkmath.LegacyDec
		switch {
		case rx.IsZero():
			ratio = y.ToLegacyDec().QuoTruncate(ry)
		case ry.IsZero():
			ratio = x.ToLegacyDec().QuoTruncate(rx)
		default:
			ratio = sdkmath.LegacyMinDec(
				x.ToLegacyDec().QuoTruncate(rx),
				y.ToLegacyDec().QuoTruncate(ry),
			)
		}
		pc = ps.MulTruncate(ratio).TruncateInt()

		mintProportion := pc.ToLegacyDec().Quo(ps)       // pc / ps
		ax = rx.Mul(mintProportion).Ceil().TruncateInt() // ceil(rx * mintProportion)
		ay = ry.Mul(mintProportion).Ceil().TruncateInt() // ceil(ry * mintProportion)
	}, func() {
		ax, ay, pc = sdkmath.ZeroInt(), sdkmath.ZeroInt(), sdkmath.ZeroInt()
	})

	return
}

// Withdraw returns withdrawn x and y coin amount when someone withdraws
// pc pool coin.
// Withdraw also takes care of the fee rate.
func Withdraw(rx, ry, ps, pc sdkmath.Int, feeRate sdkmath.LegacyDec) (x, y sdkmath.Int) {
	if pc.Equal(ps) {
		// Redeeming the last pool coin - give all remaining rx and ry.
		x = rx
		y = ry
		return
	}

	utils.SafeMath(func() {
		proportion := pc.ToLegacyDec().QuoTruncate(ps.ToLegacyDec())                       // pc / ps
		multiplier := sdkmath.LegacyOneDec().Sub(feeRate)                                  // 1 - feeRate
		x = rx.ToLegacyDec().MulTruncate(proportion).MulTruncate(multiplier).TruncateInt() // floor(rx * proportion * multiplier)
		y = ry.ToLegacyDec().MulTruncate(proportion).MulTruncate(multiplier).TruncateInt() // floor(ry * proportion * multiplier)
	}, func() {
		x, y = sdkmath.ZeroInt(), sdkmath.ZeroInt()
	})

	return
}

func DeriveTranslation(rx, ry sdkmath.Int, minPrice, maxPrice sdkmath.LegacyDec) (transX, transY sdkmath.LegacyDec) {
	sqrt := utils.DecApproxSqrt

	// M = minPrice, L = maxPrice
	rxDec, ryDec := rx.ToLegacyDec(), ry.ToLegacyDec()
	sqrtM := sqrt(minPrice)
	sqrtL := sqrt(maxPrice)

	var sqrtP sdkmath.LegacyDec
	switch {
	case rxDec.IsZero(): // y asset single pool
		sqrtP = sqrtM
	case ryDec.IsZero(): // x asset single pool
		sqrtP = sqrtL
	case rxDec.Quo(ryDec).IsZero(): // y asset single pool
		sqrtP = sqrtM
	case ryDec.Quo(rxDec).IsZero(): // x asset single pool
		sqrtP = sqrtL
	default: // normal pool
		// sqrtXOverY = sqrt(rx/ry)
		sqrtXOverY := sqrt(rxDec.Quo(ryDec))
		// alpha = sqrt(M)/sqrt(rx/ry) - sqrt(rx/ry)/sqrt(L)
		alpha := sqrtM.Quo(sqrtXOverY).Sub(sqrtXOverY.Quo(sqrtL))
		// sqrtP = sqrt(P) = {(alpha + sqrt(alpha^2 + 4)) / 2} * sqrt(rx/ry)
		sqrtP = alpha.Add(sqrt(alpha.Power(2).Add(fourDec))).QuoInt64(2).Mul(sqrtXOverY)
	}

	var sqrtK sdkmath.LegacyDec
	if !sqrtP.Equal(sqrtM) {
		// sqrtK = sqrt(K) = rx / (sqrt(P) - sqrt(M))
		sqrtK = rxDec.Quo(sqrtP.Sub(sqrtM))
	}
	if !sqrtP.Equal(sqrtL) {
		// sqrtK2 = sqrt(K') = ry / (1/sqrt(P) - 1/sqrt(L))
		sqrtK2 := ryDec.Quo(inv(sqrtP).Sub(inv(sqrtL)))
		if sqrtK.IsNil() { // P == M
			sqrtK = sqrtK2
		} else {
			p := sqrtP.Power(2)
			p1 := rxDec.Add(sqrtK.Mul(sqrtM)).Quo(ryDec.Add(sqrtK.Quo(sqrtL)))
			p2 := rxDec.Add(sqrtK2.Mul(sqrtM)).Quo(ryDec.Add(sqrtK2.Quo(sqrtL)))
			if p.Sub(p1).Abs().GT(p.Sub(p2).Abs()) {
				sqrtK = sqrtK2
			}
		}
	}
	transX = sqrtK.Mul(sqrtM)
	transY = sqrtK.Quo(sqrtL)

	return
}

func PoolOrders(pool Pool, orderer Orderer, lowestPrice, highestPrice sdkmath.LegacyDec, tickPrec int) []Order {
	return append(
		PoolBuyOrders(pool, orderer, lowestPrice, highestPrice, tickPrec),
		PoolSellOrders(pool, orderer, lowestPrice, highestPrice, tickPrec)...)
}

func PoolBuyOrders(pool Pool, orderer Orderer, lowestPrice, highestPrice sdkmath.LegacyDec, tickPrec int) (orders []Order) {
	defer func() {
		if r := recover(); r != nil {
			orders = nil
		}
	}()
	poolPrice := pool.Price()
	if poolPrice.LTE(lowestPrice) {
		return nil
	}
	tmpPool := pool.Clone()
	placeOrder := func(price sdkmath.LegacyDec, amt sdkmath.Int, derive bool) {
		orders = append(orders, orderer.Order(Buy, price, amt))
		rx, ry := tmpPool.Balances()
		rx = rx.Sub(price.MulInt(amt).Ceil().TruncateInt()) // quote coin ceiling
		ry = ry.Add(amt)
		tmpPool.SetBalances(rx, ry, derive)
	}
	if poolPrice.GT(highestPrice) {
		amt := tmpPool.BuyAmountTo(highestPrice)
		if amt.GTE(MinCoinAmount) {
			placeOrder(highestPrice, amt, true)
		}
	}
	tick := PriceToDownTick(sdkmath.LegacyMinDec(highestPrice, tmpPool.Price()), tickPrec)
	for tick.GTE(lowestPrice) {
		amt := tmpPool.BuyAmountOver(tick, true)
		if amt.LT(MinCoinAmount) {
			tick = DownTick(tick, tickPrec) // TODO: check if the tick is the lowest possible tick
			continue
		}
		placeOrder(tick, amt, false)
		rx, _ := tmpPool.Balances()
		if !rx.IsPositive() {
			break
		}
		tick = PriceToDownTick(tick.Mul(oneDec.Sub(poolOrderPriceGapRatio(poolPrice, tick))), tickPrec)
	}
	return orders
}

func PoolSellOrders(pool Pool, orderer Orderer, lowestPrice, highestPrice sdkmath.LegacyDec, tickPrec int) (orders []Order) {
	defer func() {
		if r := recover(); r != nil {
			orders = nil
		}
	}()
	poolPrice := pool.Price()
	if poolPrice.GTE(highestPrice) {
		return nil
	}
	tmpPool := pool.Clone()
	placeOrder := func(price sdkmath.LegacyDec, amt sdkmath.Int, derive bool) {
		orders = append(orders, orderer.Order(Sell, price, amt))
		rx, ry := tmpPool.Balances()
		rx = rx.Add(price.MulInt(amt).TruncateInt()) // quote coin truncation
		ry = ry.Sub(amt)
		tmpPool.SetBalances(rx, ry, derive)
	}
	if poolPrice.LT(lowestPrice) {
		amt := tmpPool.SellAmountTo(lowestPrice)
		if amt.GTE(MinCoinAmount) && lowestPrice.MulInt(amt).TruncateInt().IsPositive() {
			placeOrder(lowestPrice, amt, true)
		}
	}
	tick := PriceToUpTick(sdkmath.LegacyMaxDec(lowestPrice, tmpPool.Price()), tickPrec)
	for tick.LTE(highestPrice) {
		amt := tmpPool.SellAmountUnder(tick, true)
		if amt.LT(MinCoinAmount) || tick.MulInt(amt).TruncateInt().IsZero() {
			tick = UpTick(tick, tickPrec)
			continue
		}
		placeOrder(tick, amt, false)
		_, ry := tmpPool.Balances()
		if !ry.GT(MinCoinAmount) {
			break
		}
		tick = PriceToUpTick(tick.Mul(oneDec.Add(poolOrderPriceGapRatio(poolPrice, tick))), tickPrec)
	}
	return orders
}

// InitialPoolCoinSupply returns ideal initial pool coin minting amount.
func InitialPoolCoinSupply(x, y sdkmath.Int) sdkmath.Int {
	cx := len(x.BigInt().Text(10)) - 1 // characteristic of x
	cy := len(y.BigInt().Text(10)) - 1 // characteristic of y
	c := ((cx + 1) + (cy + 1) + 1) / 2 // ceil(((cx + 1) + (cy + 1)) / 2)
	res := big.NewInt(10)
	res.Exp(res, big.NewInt(int64(c)), nil) // 10^c
	return sdkmath.NewIntFromBigInt(res)
}
