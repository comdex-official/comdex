package amm_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidity/amm"
)

func newOrder(dir amm.OrderDirection, price sdkmath.LegacyDec, amt sdkmath.Int) amm.Order {
	return amm.DefaultOrderer.Order(dir, price, amt)
}

func TestFindMatchPrice(t *testing.T) {
	for _, tc := range []struct {
		name       string
		ov         amm.OrderView
		found      bool
		matchPrice sdkmath.LegacyDec
	}{
		{
			"happy case",
			amm.NewOrderBook(
				newOrder(amm.Buy, utils.ParseDec("1.1"), sdkmath.NewInt(10000)),
				newOrder(amm.Sell, utils.ParseDec("0.9"), sdkmath.NewInt(10000)),
			).MakeView(),
			true,
			utils.ParseDec("1.0"),
		},
		{
			"buy order only",
			amm.NewOrderBook(
				newOrder(amm.Buy, utils.ParseDec("1.0"), sdkmath.NewInt(10000)),
			).MakeView(),
			false,
			sdkmath.LegacyDec{},
		},
		{
			"sell order only",
			amm.NewOrderBook(
				newOrder(amm.Sell, utils.ParseDec("1.0"), sdkmath.NewInt(10000)),
			).MakeView(),
			false,
			sdkmath.LegacyDec{},
		},
		{
			"highest buy price is lower than lowest sell price",
			amm.NewOrderBook(
				newOrder(amm.Buy, utils.ParseDec("0.9"), sdkmath.NewInt(10000)),
				newOrder(amm.Sell, utils.ParseDec("1.1"), sdkmath.NewInt(10000)),
			).MakeView(),
			false,
			sdkmath.LegacyDec{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			matchPrice, found := amm.FindMatchPrice(tc.ov, int(defTickPrec))
			require.Equal(t, tc.found, found)
			if found {
				require.Equal(t, tc.matchPrice, matchPrice)
			}
		})
	}
}

func TestFindMatchPrice_Rounding(t *testing.T) {
	basePrice := utils.ParseDec("0.9990")

	for i := 0; i < 50; i++ {
		ob := amm.NewOrderBook(
			newOrder(amm.Buy, defTickPrec.UpTick(defTickPrec.UpTick(basePrice)), sdkmath.NewInt(80)),
			newOrder(amm.Sell, defTickPrec.UpTick(basePrice), sdkmath.NewInt(20)),
			newOrder(amm.Buy, basePrice, sdkmath.NewInt(10)), newOrder(amm.Sell, basePrice, sdkmath.NewInt(10)),
			newOrder(amm.Sell, defTickPrec.DownTick(basePrice), sdkmath.NewInt(70)),
		)
		matchPrice, found := amm.FindMatchPrice(ob.MakeView(), int(defTickPrec))
		require.True(t, found)
		require.True(sdkmath.LegacyDecEq(t,
			defTickPrec.RoundPrice(basePrice.Add(defTickPrec.UpTick(basePrice)).QuoInt64(2)),
			matchPrice))

		basePrice = defTickPrec.UpTick(basePrice)
	}
}

func TestMatchOrders(t *testing.T) {
	_, _, matched := amm.NewOrderBook().Match(utils.ParseDec("1.0"))
	require.False(t, matched)

	for _, tc := range []struct {
		name          string
		ob            *amm.OrderBook
		lastPrice     sdkmath.LegacyDec
		matched       bool
		matchPrice    sdkmath.LegacyDec
		quoteCoinDust sdkmath.Int
	}{
		{
			"happy case",
			amm.NewOrderBook(
				newOrder(amm.Buy, utils.ParseDec("1.0"), sdkmath.NewInt(10000)),
				newOrder(amm.Sell, utils.ParseDec("1.0"), sdkmath.NewInt(10000)),
			),
			utils.ParseDec("1.0"),
			true,
			utils.ParseDec("1.0"),
			sdkmath.ZeroInt(),
		},
		{
			"happy case #2",
			amm.NewOrderBook(
				newOrder(amm.Buy, utils.ParseDec("1.1"), sdkmath.NewInt(10000)),
				newOrder(amm.Sell, utils.ParseDec("0.9"), sdkmath.NewInt(10000)),
			),
			utils.ParseDec("1.0"),
			true,
			utils.ParseDec("1.0"),
			sdkmath.ZeroInt(),
		},
		{
			"positive quote coin dust",
			amm.NewOrderBook(
				newOrder(amm.Buy, utils.ParseDec("0.9999"), sdkmath.NewInt(1000)),
				newOrder(amm.Buy, utils.ParseDec("0.9999"), sdkmath.NewInt(1000)),
				newOrder(amm.Sell, utils.ParseDec("0.9999"), sdkmath.NewInt(1000)),
				newOrder(amm.Sell, utils.ParseDec("0.9999"), sdkmath.NewInt(1000)),
			),
			utils.ParseDec("0.9999"),
			true,
			utils.ParseDec("0.9999"),
			sdkmath.NewInt(2),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			matchPrice, quoteCoinDust, matched := tc.ob.Match(tc.lastPrice)
			require.Equal(t, tc.matched, matched)
			require.True(sdkmath.LegacyDecEq(t, tc.matchPrice, matchPrice))
			if matched {
				require.True(sdkmath.IntEq(t, tc.quoteCoinDust, quoteCoinDust))
				for _, order := range tc.ob.Orders() {
					if order.IsMatched() {
						paid := order.GetPaidOfferCoinAmount()
						received := order.GetReceivedDemandCoinAmount()
						var effPrice sdkmath.LegacyDec // Effective swap price
						switch order.GetDirection() {
						case amm.Buy:
							effPrice = paid.ToLegacyDec().QuoInt(received)
						case amm.Sell:
							effPrice = received.ToLegacyDec().QuoInt(paid)
						}
						require.True(t, utils.DecApproxEqual(tc.lastPrice, effPrice))
					}
				}
			}
		})
	}
}

func TestFindMatchableAmountAtSinglePrice(t *testing.T) {
	for _, tc := range []struct {
		orders       []amm.Order
		matchPrice   sdkmath.LegacyDec
		found        bool
		matchableAmt sdkmath.Int
	}{
		{
			[]amm.Order{
				newOrder(amm.Sell, utils.ParseDec("0.100"), sdkmath.NewInt(10000)),
				newOrder(amm.Sell, utils.ParseDec("0.099"), sdkmath.NewInt(9995)),
				newOrder(amm.Buy, utils.ParseDec("0.101"), sdkmath.NewInt(10000)),
			},
			utils.ParseDec("0.100"),
			true,
			sdkmath.NewInt(9995),
		},
		{
			[]amm.Order{
				newOrder(amm.Sell, utils.ParseDec("0.100"), sdkmath.NewInt(10000)),
				newOrder(amm.Sell, utils.ParseDec("0.099"), sdkmath.NewInt(9995)),
				newOrder(amm.Buy, utils.ParseDec("0.101"), sdkmath.NewInt(10000)),
				newOrder(amm.Buy, utils.ParseDec("0.100"), sdkmath.NewInt(1000)),
			},
			utils.ParseDec("0.100"),
			true,
			sdkmath.NewInt(11000),
		},
	} {
		t.Run("", func(t *testing.T) {
			ob := amm.NewOrderBook(tc.orders...)
			matchableAmt, found := ob.FindMatchableAmountAtSinglePrice(tc.matchPrice)
			require.Equal(t, tc.found, found)
			if found {
				require.True(sdkmath.IntEq(t, tc.matchableAmt, matchableAmt))
			}
		})
	}
}

func TestMatch_edgecase1(t *testing.T) {
	orders := []amm.Order{
		newOrder(amm.Sell, utils.ParseDec("0.100"), sdkmath.NewInt(10000)),
		newOrder(amm.Sell, utils.ParseDec("0.099"), sdkmath.NewInt(9995)),
		newOrder(amm.Buy, utils.ParseDec("0.101"), sdkmath.NewInt(10000)),
		newOrder(amm.Buy, utils.ParseDec("0.100"), sdkmath.NewInt(5000)),
	}
	ob := amm.NewOrderBook(orders...)
	_, _, matched := ob.Match(utils.ParseDec("0.098"))
	require.True(t, matched)
	for _, order := range orders {
		fmt.Printf(
			"%s %s (%s/%s) paid=%s, received=%s\n",
			order.GetDirection(), order.GetPrice(), order.GetOpenAmount(), order.GetAmount(),
			order.GetPaidOfferCoinAmount(), order.GetReceivedDemandCoinAmount())
	}
}
