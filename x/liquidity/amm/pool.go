package amm

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	utils "github.com/petrichormoney/petri/types"
)

var (
	_ Pool        = (*BasicPool)(nil)
	_ OrderSource = (*MockPoolOrderSource)(nil)
)

// Pool is the interface of a pool.
// It also satisfies OrderView interface.
type Pool interface {
	OrderView
	Balances() (rx, ry sdk.Int)
	PoolCoinSupply() sdk.Int
	Price() sdk.Dec
	IsDepleted() bool
	Deposit(x, y sdk.Int) (ax, ay, pc sdk.Int)
	Withdraw(pc sdk.Int, feeRate sdk.Dec) (x, y sdk.Int)
	ProvidableXAmountOver(price sdk.Dec) sdk.Int
	ProvidableYAmountUnder(price sdk.Dec) sdk.Int
}

// BasicPool is the basic pool type.
type BasicPool struct {
	// rx and ry are the pool's reserve balance of each x/y coin.
	// In perspective of a pair, x coin is the quote coin and
	// y coin is the base coin.
	rx, ry sdk.Int
	// ps is the pool's pool coin supply.
	ps sdk.Int
}

// NewBasicPool returns a new BasicPool.
// It is OK to pass an empty sdk.Int to ps when ps is not going to be used.
func NewBasicPool(rx, ry, ps sdk.Int) *BasicPool {
	return &BasicPool{
		rx: rx,
		ry: ry,
		ps: ps,
	}
}

// Balances returns the balances of the pool.
func (pool *BasicPool) Balances() (rx, ry sdk.Int) {
	return pool.rx, pool.ry
}

// PoolCoinSupply returns the pool coin supply.
func (pool *BasicPool) PoolCoinSupply() sdk.Int {
	return pool.ps
}

// Price returns the pool price.
func (pool *BasicPool) Price() sdk.Dec {
	if pool.rx.IsZero() || pool.ry.IsZero() {
		panic("pool price is not defined for a depleted pool")
	}
	return pool.rx.ToDec().Quo(pool.ry.ToDec())
}

// IsDepleted returns whether the pool is depleted or not.
func (pool *BasicPool) IsDepleted() bool {
	return pool.ps.IsZero() || pool.rx.IsZero() || pool.ry.IsZero()
}

// Deposit returns accepted x and y coin amount and minted pool coin amount
// when someone deposits x and y coins.
func (pool *BasicPool) Deposit(x, y sdk.Int) (ax, ay, pc sdk.Int) {
	// Calculate accepted amount and minting amount.
	// Note that we take as many coins as possible(by ceiling numbers)
	// from depositor and mint as little coins as possible.

	utils.SafeMath(func() {
		rx, ry := pool.rx.ToDec(), pool.ry.ToDec()
		ps := pool.ps.ToDec()

		// pc = floor(ps * min(x / rx, y / ry))
		pc = ps.MulTruncate(sdk.MinDec(
			x.ToDec().QuoTruncate(rx),
			y.ToDec().QuoTruncate(ry),
		)).TruncateInt()

		mintProportion := pc.ToDec().Quo(ps)             // pc / ps
		ax = rx.Mul(mintProportion).Ceil().TruncateInt() // ceil(rx * mintProportion)
		ay = ry.Mul(mintProportion).Ceil().TruncateInt() // ceil(ry * mintProportion)
	}, func() {
		ax, ay, pc = sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt()
	})

	return
}

// Withdraw returns withdrawn x and y coin amount when someone withdraws
// pc pool coin.
// Withdraw also takes care of the fee rate.
func (pool *BasicPool) Withdraw(pc sdk.Int, feeRate sdk.Dec) (x, y sdk.Int) {
	if pc.Equal(pool.ps) {
		// Redeeming the last pool coin - give all remaining rx and ry.
		x = pool.rx
		y = pool.ry
		return
	}

	utils.SafeMath(func() {
		proportion := pc.ToDec().QuoTruncate(pool.ps.ToDec())                             // pc / ps
		multiplier := sdk.OneDec().Sub(feeRate)                                           // 1 - feeRate
		x = pool.rx.ToDec().MulTruncate(proportion).MulTruncate(multiplier).TruncateInt() // floor(rx * proportion * multiplier)
		y = pool.ry.ToDec().MulTruncate(proportion).MulTruncate(multiplier).TruncateInt() // floor(ry * proportion * multiplier)
	}, func() {
		x, y = sdk.ZeroInt(), sdk.ZeroInt()
	})

	return
}

// HighestBuyPrice returns the highest buy price of the pool.
func (pool *BasicPool) HighestBuyPrice() (price sdk.Dec, found bool) {
	// The highest buy price is actually a bit lower than pool price,
	// but it's not important for our matching logic.
	return pool.Price(), true
}

// LowestSellPrice returns the lowest sell price of the pool.
func (pool *BasicPool) LowestSellPrice() (price sdk.Dec, found bool) {
	// The lowest sell price is actually a bit higher than the pool price,
	// but it's not important for our matching logic.
	return pool.Price(), true
}

// BuyAmountOver returns the amount of buy orders for price greater than
// or equal to given price.
func (pool *BasicPool) BuyAmountOver(price sdk.Dec) (amt sdk.Int) {
	if price.GTE(pool.Price()) {
		return sdk.ZeroInt()
	}
	utils.SafeMath(func() {
		amt = pool.rx.ToDec().QuoTruncate(price).Sub(pool.ry.ToDec()).TruncateInt()
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
func (pool *BasicPool) SellAmountUnder(price sdk.Dec) sdk.Int {
	return pool.ProvidableYAmountUnder(price)
}

// ProvidableXAmountOver returns the amount of x coin the pool would provide
// for price greater than or equal to given price.
func (pool *BasicPool) ProvidableXAmountOver(price sdk.Dec) (amt sdk.Int) {
	if price.GTE(pool.Price()) {
		return sdk.ZeroInt()
	}
	utils.SafeMath(func() {
		amt = pool.rx.ToDec().Sub(pool.ry.ToDec().Mul(price)).TruncateInt()
		if amt.GT(MaxCoinAmount) {
			amt = MaxCoinAmount
		}
	}, func() {
		amt = MaxCoinAmount
	})
	return
}

// ProvidableYAmountUnder returns the amount of y coin the pool would provide
// for price less than or equal to given price.
func (pool *BasicPool) ProvidableYAmountUnder(price sdk.Dec) sdk.Int {
	if price.LTE(pool.Price()) {
		return sdk.ZeroInt()
	}
	return pool.ry.ToDec().Sub(pool.rx.ToDec().QuoRoundUp(price)).TruncateInt()
}

// PoolsOrderBook returns an order book with orders made by pools.
// Use Ticks or EvenTicks to generate ticks where pools put orders.
// ticks should have more than 1 element.
// The orders in the order book are just mocks, so the rest fields
// other than Direction, Price and Amount of BaseOrder will have no
// special meaning.
func PoolsOrderBook(pools []Pool, ticks []sdk.Dec) *OrderBook {
	sortTicks(ticks)
	highestTick := ticks[0]
	lowestTick := ticks[len(ticks)-1]
	gap := ticks[0].Sub(ticks[1])
	ob := NewOrderBook()
	for _, pool := range pools {
		poolPrice := pool.Price()
		if poolPrice.GT(lowestTick) { // Buy orders
			accAmtX := pool.ProvidableXAmountOver(highestTick.Add(gap))
			for _, tick := range ticks {
				amtX := pool.ProvidableXAmountOver(tick).Sub(accAmtX)
				if amtX.IsPositive() {
					amt := amtX.ToDec().QuoTruncate(tick).TruncateInt()
					if amt.IsPositive() {
						ob.Add(NewBaseOrder(
							Buy, tick, amt, sdk.NewCoin("quote", OfferCoinAmount(Buy, tick, amt)), "base"))
					}
					accAmtX = accAmtX.Add(amtX)
				}
			}
		}
		if poolPrice.LT(highestTick) { // Sell orders
			accAmt := pool.SellAmountUnder(lowestTick.Sub(gap))
			for i := len(ticks) - 1; i >= 0; i-- {
				tick := ticks[i]
				amt := pool.SellAmountUnder(tick).Sub(accAmt)
				if amt.IsPositive() {
					ob.Add(NewBaseOrder(Sell, tick, amt, sdk.NewCoin("base", amt), "quote"))
					accAmt = accAmt.Add(amt)
				}
			}
		}
	}
	return ob
}

// InitialPoolCoinSupply returns ideal initial pool coin minting amount.
func InitialPoolCoinSupply(x, y sdk.Int) sdk.Int {
	cx := len(x.BigInt().Text(10)) - 1 // characteristic of x
	cy := len(y.BigInt().Text(10)) - 1 // characteristic of y
	c := ((cx + 1) + (cy + 1) + 1) / 2 // ceil(((cx + 1) + (cy + 1)) / 2)
	res := big.NewInt(10)
	res.Exp(res, big.NewInt(int64(c)), nil) // 10^c
	return sdk.NewIntFromBigInt(res)
}

// MockPoolOrderSource demonstrates how to implement a pool OrderSource.
type MockPoolOrderSource struct {
	Pool
	baseCoinDenom, quoteCoinDenom string
}

// NewMockPoolOrderSource returns a new MockPoolOrderSource for testing.
func NewMockPoolOrderSource(pool Pool, baseCoinDenom, quoteCoinDenom string) *MockPoolOrderSource {
	return &MockPoolOrderSource{
		Pool:           pool,
		baseCoinDenom:  baseCoinDenom,
		quoteCoinDenom: quoteCoinDenom,
	}
}

// BuyOrdersOver returns buy orders for price greater or equal than given price.
func (os *MockPoolOrderSource) BuyOrdersOver(price sdk.Dec) []Order {
	amt := os.BuyAmountOver(price)
	if amt.IsZero() {
		return nil
	}
	quoteCoin := sdk.NewCoin(os.quoteCoinDenom, OfferCoinAmount(Buy, price, amt))
	return []Order{NewBaseOrder(Buy, price, amt, quoteCoin, os.baseCoinDenom)}
}

// SellOrdersUnder returns sell orders for price less or equal than given price.
func (os *MockPoolOrderSource) SellOrdersUnder(price sdk.Dec) []Order {
	amt := os.SellAmountUnder(price)
	if amt.IsZero() {
		return nil
	}
	baseCoin := sdk.NewCoin(os.baseCoinDenom, amt)
	return []Order{NewBaseOrder(Sell, price, amt, baseCoin, os.quoteCoinDenom)}
}
