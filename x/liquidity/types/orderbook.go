package types

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	"sort"

	"github.com/comdex-official/comdex/x/liquidity/amm"
)

func OrderBookBasePrice(ov amm.OrderView, tickPrec int) (sdkmath.LegacyDec, bool) {
	highestBuyPrice, foundHighestBuyPrice := ov.HighestBuyPrice()
	lowestSellPrice, foundLowestSellPrice := ov.LowestSellPrice()

	switch {
	case foundHighestBuyPrice && foundLowestSellPrice:
		return amm.RoundPrice(highestBuyPrice.Add(lowestSellPrice).QuoInt64(2), tickPrec), true
	case foundHighestBuyPrice:
		return highestBuyPrice, true
	case foundLowestSellPrice:
		return lowestSellPrice, true
	default: // not found
		return sdkmath.LegacyDec{}, false
	}
}

// OrderBookConfig defines configuration parameter for an order book response.
type OrderBookConfig struct {
	PriceUnitPower int
	MaxNumTicks    int
}

func MakeOrderBookPairResponse(pairID uint64, ov *amm.OrderBookView, lowestPrice, highestPrice sdkmath.LegacyDec, tickPrecision int, configs ...OrderBookConfig) OrderBookPairResponse {
	resp := OrderBookPairResponse{
		PairId: pairID,
	}
	basePrice, found := OrderBookBasePrice(ov, tickPrecision)
	if !found {
		return resp
	}
	resp.BasePrice = basePrice
	ammTickPrec := amm.TickPrecision(tickPrecision)

	sort.Slice(configs, func(i, j int) bool {
		return configs[i].PriceUnitPower < configs[j].PriceUnitPower
	})
	lowestPriceUnitMaxNumTicks := configs[0].MaxNumTicks

	highestBuyPrice, foundHighestBuyPrice := ov.HighestBuyPrice()
	lowestSellPrice, foundLowestSellPrice := ov.LowestSellPrice()

	var smallestPriceUnit sdkmath.LegacyDec
	if foundLowestSellPrice {
		currentPrice := lowestSellPrice
		for i := 0; i < lowestPriceUnitMaxNumTicks && currentPrice.LTE(highestPrice); {
			amtInclusive := ov.SellAmountOver(currentPrice, true)
			amtExclusive := ov.SellAmountOver(currentPrice, false)
			amt := amtInclusive.Sub(amtExclusive)
			if amt.IsPositive() {
				i++
				if i == lowestPriceUnitMaxNumTicks {
					break
				}
			}
			if !amtExclusive.IsPositive() {
				break
			}
			currentPrice = ammTickPrec.UpTick(currentPrice)
		}
		smallestPriceUnit = ammTickPrec.TickGap(currentPrice)
	} else {
		smallestPriceUnit = ammTickPrec.TickGap(highestBuyPrice)
	}

	for _, config := range configs {
		priceUnit := smallestPriceUnit
		for j := 0; j < config.PriceUnitPower; j++ {
			priceUnit = priceUnit.MulInt64(10)
		}
		ob := OrderBookResponse{
			PriceUnit: priceUnit,
			Buys:      nil,
			Sells:     nil,
		}
		if foundLowestSellPrice {
			startPrice := FitPriceToTickGap(lowestSellPrice, priceUnit, false)
			currentPrice := startPrice
			accAmt := sdkmath.ZeroInt()
			for j := 0; j < config.MaxNumTicks && currentPrice.LTE(highestPrice); {
				amt := ov.SellAmountUnder(currentPrice, true).Sub(accAmt)
				if amt.IsPositive() {
					ob.Sells = append(ob.Sells, OrderBookTickResponse{
						Price:           currentPrice,
						UserOrderAmount: amt,
						PoolOrderAmount: sdkmath.ZeroInt(),
					})
					accAmt = accAmt.Add(amt)
					j++
				}
				if !ov.SellAmountOver(currentPrice, false).IsPositive() {
					break
				}
				currentPrice = currentPrice.Add(priceUnit)
			}
			// Reverse sell ticks.
			for l, r := 0, len(ob.Sells)-1; l < r; l, r = l+1, r-1 {
				ob.Sells[l], ob.Sells[r] = ob.Sells[r], ob.Sells[l]
			}
		}
		if foundHighestBuyPrice {
			startPrice := FitPriceToTickGap(highestBuyPrice, priceUnit, true)
			currentPrice := startPrice
			accAmt := sdkmath.ZeroInt()
			for j := 0; j < config.MaxNumTicks && currentPrice.GTE(lowestPrice) && !currentPrice.IsNegative(); {
				amt := ov.BuyAmountOver(currentPrice, true).Sub(accAmt)
				if amt.IsPositive() {
					ob.Buys = append(ob.Buys, OrderBookTickResponse{
						Price:           currentPrice,
						UserOrderAmount: amt,
						PoolOrderAmount: sdkmath.ZeroInt(),
					})
					accAmt = accAmt.Add(amt)
					j++
				}
				if !ov.BuyAmountUnder(currentPrice, false).IsPositive() {
					break
				}
				currentPrice = currentPrice.Sub(priceUnit)
			}
		}
		resp.OrderBooks = append(resp.OrderBooks, ob)
	}

	return resp
}

// PrintOrderBookResponse prints out OrderBookResponse in human-readable form.
func PrintOrderBookResponse(ob OrderBookResponse, basePrice sdkmath.LegacyDec) {
	fmt.Println("+------------------------------------------------------------------------+")
	for _, tick := range ob.Sells {
		fmt.Printf("| %18s | %28s |                    |\n", tick.UserOrderAmount, tick.Price.String())
	}
	fmt.Println("|------------------------------------------------------------------------|")
	fmt.Printf("|                      %28s                      |\n", basePrice.String())
	fmt.Println("|------------------------------------------------------------------------|")
	for _, tick := range ob.Buys {
		fmt.Printf("|                    | %28s | %-18s |\n", tick.Price.String(), tick.UserOrderAmount)
	}
	fmt.Println("+------------------------------------------------------------------------+")
}

// FitPriceToTickGap fits price into given tick gap.
func FitPriceToTickGap(price, gap sdkmath.LegacyDec, down bool) sdkmath.LegacyDec {
	b := price.BigInt()
	b.Quo(b, gap.BigInt()).Mul(b, gap.BigInt())
	tick := sdkmath.LegacyNewDecFromBigIntWithPrec(b, sdkmath.LegacyPrecision)
	if !down && !tick.Equal(price) {
		tick = tick.Add(gap)
	}
	return tick
}
