package amm_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidity/amm"
)

func TestBasicPool(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	for i := 0; i < 1000; i++ {
		rx, ry := sdkmath.NewInt(1+r.Int63n(100000000)), sdkmath.NewInt(1+r.Int63n(100000000))
		pool := amm.NewBasicPool(rx, ry, sdkmath.Int{})

		highest, found := pool.HighestBuyPrice()
		require.True(t, found)
		require.True(sdkmath.LegacyDecEq(t, pool.Price(), highest))
		lowest, found := pool.LowestSellPrice()
		require.True(t, found)
		require.True(sdkmath.LegacyDecEq(t, pool.Price(), lowest))
	}
}

func TestCreateBasicPool(t *testing.T) {
	for _, tc := range []struct {
		name        string
		rx, ry      sdkmath.Int
		expectedErr string
	}{
		{
			"both zero amount",
			sdkmath.NewInt(0), sdkmath.NewInt(0),
			"cannot create basic pool with zero reserve amount",
		},
		{
			"zero y amount",
			sdkmath.NewInt(1000000), sdkmath.NewInt(0),
			"cannot create basic pool with zero reserve amount",
		},
		{
			"zero x amount",
			sdkmath.NewInt(0), sdkmath.NewInt(1000000),
			"cannot create basic pool with zero reserve amount",
		},
		{
			"too low price",
			sdkmath.NewInt(1000000), sdkmath.NewIntWithDecimal(1, 26),
			"pool price is lower than min price 0.000000000000001000",
		},
		{
			"too high price",
			sdkmath.NewIntWithDecimal(1, 48), sdkmath.NewInt(1000000),
			"pool price is greater than max price 100000000000000000000.000000000000000000",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := amm.CreateBasicPool(tc.rx, tc.ry)
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestBasicPool_Price(t *testing.T) {
	for _, tc := range []struct {
		name   string
		rx, ry int64             // reserve balance
		ps     int64             // pool coin supply
		p      sdkmath.LegacyDec // expected pool price
	}{
		{
			name: "normal pool",
			ps:   10000,
			rx:   20000,
			ry:   100,
			p:    sdkmath.LegacyNewDec(200),
		},
		{
			name: "decimal rounding",
			ps:   10000,
			rx:   200,
			ry:   300,
			p:    sdkmath.LegacyMustNewDecFromStr("0.666666666666666667"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pool := amm.NewBasicPool(sdkmath.NewInt(tc.rx), sdkmath.NewInt(tc.ry), sdkmath.NewInt(tc.ps))
			require.True(sdkmath.LegacyDecEq(t, tc.p, pool.Price()))
		})
	}

	// panicking cases
	for _, tc := range []struct {
		rx, ry int64
		ps     int64
	}{
		{
			rx: 0,
			ry: 1000,
			ps: 1000,
		},
		{
			rx: 1000,
			ry: 0,
			ps: 1000,
		},
	} {
		t.Run("panics", func(t *testing.T) {
			require.Panics(t, func() {
				pool := amm.NewBasicPool(sdkmath.NewInt(tc.rx), sdkmath.NewInt(tc.ry), sdkmath.NewInt(tc.ps))
				pool.Price()
			})
		})
	}
}

func TestBasicPool_IsDepleted(t *testing.T) {
	for _, tc := range []struct {
		name       string
		rx, ry     int64 // reserve balance
		ps         int64 // pool coin supply
		isDepleted bool
	}{
		{
			name:       "empty pool",
			rx:         0,
			ry:         0,
			ps:         0,
			isDepleted: true,
		},
		{
			name:       "depleted, with some coins from outside",
			rx:         100,
			ry:         0,
			ps:         0,
			isDepleted: true,
		},
		{
			name:       "depleted, with some coins from outside #2",
			rx:         100,
			ry:         100,
			ps:         0,
			isDepleted: true,
		},
		{
			name:       "normal pool",
			rx:         10000,
			ry:         10000,
			ps:         10000,
			isDepleted: false,
		},
		{
			name:       "not depleted, but reserve coins are gone",
			rx:         0,
			ry:         10000,
			ps:         10000,
			isDepleted: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pool := amm.NewBasicPool(sdkmath.NewInt(tc.rx), sdkmath.NewInt(tc.ry), sdkmath.NewInt(tc.ps))
			require.Equal(t, tc.isDepleted, pool.IsDepleted())
		})
	}
}

func TestBasicPool_Deposit(t *testing.T) {
	for _, tc := range []struct {
		name   string
		rx, ry int64 // reserve balance
		ps     int64 // pool coin supply
		x, y   int64 // depositing coin amount
		ax, ay int64 // expected accepted coin amount
		pc     int64 // expected minted pool coin amount
	}{
		{
			name: "ideal deposit",
			rx:   2000,
			ry:   100,
			ps:   10000,
			x:    200,
			y:    10,
			ax:   200,
			ay:   10,
			pc:   1000,
		},
		{
			name: "unbalanced deposit",
			rx:   2000,
			ry:   100,
			ps:   10000,
			x:    100,
			y:    2000,
			ax:   100,
			ay:   5,
			pc:   500,
		},
		{
			name: "decimal truncation",
			rx:   222,
			ry:   333,
			ps:   333,
			x:    100,
			y:    100,
			ax:   66,
			ay:   99,
			pc:   99,
		},
		{
			name: "decimal truncation #2",
			rx:   200,
			ry:   300,
			ps:   333,
			x:    80,
			y:    80,
			ax:   53,
			ay:   80,
			pc:   88,
		},
		{
			name: "zero minting amount",
			ps:   100,
			rx:   10000,
			ry:   10000,
			x:    99,
			y:    99,
			ax:   0,
			ay:   0,
			pc:   0,
		},
		{
			name: "tiny minting amount",
			rx:   10000,
			ry:   10000,
			ps:   100,
			x:    100,
			y:    100,
			ax:   100,
			ay:   100,
			pc:   1,
		},
		{
			name: "tiny minting amount #2",
			rx:   10000,
			ry:   10000,
			ps:   100,
			x:    199,
			y:    199,
			ax:   100,
			ay:   100,
			pc:   1,
		},
		{
			name: "zero minting amount",
			rx:   10000,
			ry:   10000,
			ps:   999,
			x:    10,
			y:    10,
			ax:   0,
			ay:   0,
			pc:   0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pool := amm.NewBasicPool(sdkmath.NewInt(tc.rx), sdkmath.NewInt(tc.ry), sdkmath.NewInt(tc.ps))
			ax, ay, pc := amm.Deposit(sdkmath.NewInt(tc.rx), sdkmath.NewInt(tc.ry), sdkmath.NewInt(tc.ps), sdkmath.NewInt(tc.x), sdkmath.NewInt(tc.y))
			require.True(sdkmath.IntEq(t, sdkmath.NewInt(tc.ax), ax))
			require.True(sdkmath.IntEq(t, sdkmath.NewInt(tc.ay), ay))
			require.True(sdkmath.IntEq(t, sdkmath.NewInt(tc.pc), pc))
			// Additional assertions
			if !pool.IsDepleted() {
				require.True(t, (ax.Int64()*tc.ps) >= (pc.Int64()*tc.rx)) // (ax / rx) > (pc / ps)
				require.True(t, (ay.Int64()*tc.ps) >= (pc.Int64()*tc.ry)) // (ay / ry) > (pc / ps)
			}
		})
	}
}

func TestBasicPool_Withdraw(t *testing.T) {
	for _, tc := range []struct {
		name    string
		rx, ry  int64 // reserve balance
		ps      int64 // pool coin supply
		pc      int64 // redeeming pool coin amount
		feeRate sdkmath.LegacyDec
		x, y    int64 // withdrawn coin amount
	}{
		{
			name:    "ideal withdraw",
			rx:      2000,
			ry:      100,
			ps:      10000,
			pc:      1000,
			feeRate: sdkmath.LegacyZeroDec(),
			x:       200,
			y:       10,
		},
		{
			name:    "ideal withdraw - with fee",
			rx:      2000,
			ry:      100,
			ps:      10000,
			pc:      1000,
			feeRate: sdkmath.LegacyMustNewDecFromStr("0.003"),
			x:       199,
			y:       9,
		},
		{
			name:    "withdraw all",
			rx:      123,
			ry:      567,
			ps:      10,
			pc:      10,
			feeRate: sdkmath.LegacyMustNewDecFromStr("0.003"),
			x:       123,
			y:       567,
		},
		{
			name:    "advantageous for pool",
			rx:      100,
			ry:      100,
			ps:      10000,
			pc:      99,
			feeRate: sdkmath.LegacyZeroDec(),
			x:       0,
			y:       0,
		},
		{
			name:    "advantageous for pool",
			rx:      10000,
			ry:      100,
			ps:      10000,
			pc:      99,
			feeRate: sdkmath.LegacyZeroDec(),
			x:       99,
			y:       0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			x, y := amm.Withdraw(sdkmath.NewInt(tc.rx), sdkmath.NewInt(tc.ry), sdkmath.NewInt(tc.ps), sdkmath.NewInt(tc.pc), tc.feeRate)
			require.True(sdkmath.IntEq(t, sdkmath.NewInt(tc.x), x))
			require.True(sdkmath.IntEq(t, sdkmath.NewInt(tc.y), y))
			// Additional assertions
			require.True(t, (tc.pc*tc.rx) >= (x.Int64()*tc.ps))
			require.True(t, (tc.pc*tc.ry) >= (y.Int64()*tc.ps))
		})
	}
}

func TestBasicPool_BuyAmountOver(t *testing.T) {
	pool := amm.NewBasicPool(sdkmath.NewInt(1000000), sdkmath.NewInt(1000000), sdkmath.Int{})

	for _, tc := range []struct {
		pool  *amm.BasicPool
		price sdkmath.LegacyDec
		amt   sdkmath.Int
	}{
		{pool, utils.ParseDec("1.1"), sdkmath.ZeroInt()},
		{pool, utils.ParseDec("1.0"), sdkmath.ZeroInt()},
		{pool, utils.ParseDec("0.9"), sdkmath.NewInt(111111)},
		{pool, utils.ParseDec("0.8"), sdkmath.NewInt(250000)},
	} {
		t.Run("", func(t *testing.T) {
			amt := tc.pool.BuyAmountOver(tc.price, true)
			require.True(sdkmath.IntEq(t, tc.amt, amt))
		})
	}
}

func TestBasicPool_SellAmountUnder(t *testing.T) {
	pool := amm.NewBasicPool(sdkmath.NewInt(1000000), sdkmath.NewInt(1000000), sdkmath.Int{})

	for _, tc := range []struct {
		pool  *amm.BasicPool
		price sdkmath.LegacyDec
		amt   sdkmath.Int
	}{
		{pool, utils.ParseDec("0.9"), sdkmath.ZeroInt()},
		{pool, utils.ParseDec("1.0"), sdkmath.ZeroInt()},
		{pool, utils.ParseDec("1.1"), sdkmath.NewInt(90909)},
		{pool, utils.ParseDec("1.2"), sdkmath.NewInt(166666)},
	} {
		t.Run("", func(t *testing.T) {
			amt := tc.pool.SellAmountUnder(tc.price, true)
			require.True(sdkmath.IntEq(t, tc.amt, amt))
		})
	}
}

func TestBasicPool_BuyAmountTo(t *testing.T) {
	pool := amm.NewBasicPool(sdkmath.NewInt(1000000), sdkmath.NewInt(1000000), sdkmath.Int{})

	for _, tc := range []struct {
		pool  *amm.BasicPool
		price sdkmath.LegacyDec
		amt   sdkmath.Int
	}{
		{pool, utils.ParseDec("1.1"), sdkmath.ZeroInt()},
		{pool, utils.ParseDec("1.0"), sdkmath.ZeroInt()},
		{pool, utils.ParseDec("0.5"), sdkmath.NewInt(585786)},
		{pool, utils.ParseDec("0.4"), sdkmath.NewInt(918861)},
	} {
		t.Run("", func(t *testing.T) {
			amt := tc.pool.BuyAmountTo(tc.price)
			require.True(sdkmath.IntEq(t, tc.amt, amt))
		})
	}
}

func TestBasicPool_SellAmountTo(t *testing.T) {
	pool := amm.NewBasicPool(sdkmath.NewInt(1000000), sdkmath.NewInt(1000000), sdkmath.Int{})

	for _, tc := range []struct {
		pool  *amm.BasicPool
		price sdkmath.LegacyDec
		amt   sdkmath.Int
	}{
		{pool, utils.ParseDec("0.9"), sdkmath.ZeroInt()},
		{pool, utils.ParseDec("1.0"), sdkmath.ZeroInt()},
		{pool, utils.ParseDec("1.4"), sdkmath.NewInt(154845)},
		{pool, utils.ParseDec("1.5"), sdkmath.NewInt(183503)},
	} {
		t.Run("", func(t *testing.T) {
			amt := tc.pool.SellAmountTo(tc.price)
			require.True(sdkmath.IntEq(t, tc.amt, amt))
		})
	}
}

func TestValidateRangedPoolParams(t *testing.T) {
	for _, tc := range []struct {
		name               string
		minPrice, maxPrice sdkmath.LegacyDec
		initialPrice       sdkmath.LegacyDec
		expectedErr        string
	}{
		{
			"happy case",
			utils.ParseDec("0.5"), utils.ParseDec("2.0"),
			utils.ParseDec("1.0"),
			"",
		},
		{
			"single y asset pool",
			utils.ParseDec("0.5"), utils.ParseDec("2.0"),
			utils.ParseDec("0.5"),
			"",
		},
		{
			"single x asset pool",
			utils.ParseDec("0.5"), utils.ParseDec("2.0"),
			utils.ParseDec("2.0"),
			"",
		},
		{
			"too low min price",
			sdkmath.LegacyNewDecWithPrec(1, 16), utils.ParseDec("2.0"),
			utils.ParseDec("1.0"),
			"min price must not be lower than 0.000000000000001000",
		},
		{
			"too high max price",
			utils.ParseDec("0.5"), sdkmath.NewIntWithDecimal(1, 25).ToLegacyDec(),
			utils.ParseDec("1.0"),
			"max price must not be higher than 100000000000000000000.000000000000000000",
		},
		{
			"too low initial price",
			utils.ParseDec("0.5"), utils.ParseDec("2.0"),
			utils.ParseDec("0.499"),
			"initial price must not be lower than min price",
		},
		{
			"too high initial price",
			utils.ParseDec("0.5"), utils.ParseDec("2.0"),
			utils.ParseDec("2.001"),
			"initial price must not be higher than max price",
		},
		{
			"max price lower than min price",
			utils.ParseDec("2.0"), utils.ParseDec("0.5"),
			utils.ParseDec("1.0"),
			"max price must be higher than min price",
		},
		{
			"too close min price and max price",
			utils.ParseDec("0.9999"), utils.ParseDec("1.0001"),
			utils.ParseDec("1.0"),
			"min price and max price are too close",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := amm.ValidateRangedPoolParams(tc.minPrice, tc.maxPrice, tc.initialPrice)
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestCreateRangedPool(t *testing.T) {
	intApproxEq := func(exp, got sdkmath.Int) (*testing.T, bool, string, string, string) {
		c := exp.Sub(got).Abs().LTE(sdkmath.OneInt())
		if c && !exp.IsZero() {
			c = exp.ToLegacyDec().Sub(got.ToLegacyDec()).Abs().Quo(exp.ToLegacyDec()).LTE(sdkmath.LegacyNewDecWithPrec(1, 3))
		}
		return t, c, "expected:\t%v\ngot:\t\t%v", exp.String(), got.String()
	}

	for _, tc := range []struct {
		name               string
		x, y               sdkmath.Int
		minPrice, maxPrice sdkmath.LegacyDec
		initialPrice       sdkmath.LegacyDec
		expectedErr        string
		ax, ay             sdkmath.Int
	}{
		{
			"basic case",
			sdkmath.NewInt(1_000000), sdkmath.NewInt(1_000000),
			utils.ParseDec("0.5"), utils.ParseDec("2.0"),
			utils.ParseDec("1.0"),
			"",
			sdkmath.NewInt(1_000000), sdkmath.NewInt(1_000000),
		},
		{
			"basic case 2",
			sdkmath.NewInt(500000), sdkmath.NewInt(1_000000),
			utils.ParseDec("0.5"), utils.ParseDec("2.0"),
			utils.ParseDec("1.0"),
			"",
			sdkmath.NewInt(500000), sdkmath.NewInt(500000),
		},
		{
			"basic case 3",
			sdkmath.NewInt(1_000000), sdkmath.NewInt(500000),
			utils.ParseDec("0.5"), utils.ParseDec("2.0"),
			utils.ParseDec("1.0"),
			"",
			sdkmath.NewInt(500000), sdkmath.NewInt(500000),
		},
		{
			"invalid pool",
			sdkmath.ZeroInt(), sdkmath.ZeroInt(),
			utils.ParseDec("1.0"), utils.ParseDec("2.0"),
			utils.ParseDec("1.0"),
			"either x or y must be positive",
			sdkmath.Int{}, sdkmath.Int{},
		},
		{
			"single asset x pool",
			sdkmath.NewInt(1_000000), sdkmath.ZeroInt(),
			utils.ParseDec("1.0"), utils.ParseDec("2.0"),
			utils.ParseDec("2.0"),
			"",
			sdkmath.NewInt(1_000000), sdkmath.ZeroInt(),
		},
		{
			"single asset x pool - refund",
			sdkmath.NewInt(1_000000), sdkmath.NewInt(1_000000),
			utils.ParseDec("1.0"), utils.ParseDec("2.0"),
			utils.ParseDec("2.0"),
			"",
			sdkmath.NewInt(1_000000), sdkmath.ZeroInt(),
		},
		{
			"single asset y pool",
			sdkmath.ZeroInt(), sdkmath.NewInt(1_000000),
			utils.ParseDec("1.0"), utils.ParseDec("2.0"),
			utils.ParseDec("1.0"),
			"",
			sdkmath.ZeroInt(), sdkmath.NewInt(1_000000),
		},
		{
			"single asset y pool - refund",
			sdkmath.NewInt(1_000000), sdkmath.NewInt(1_000000),
			utils.ParseDec("1.0"), utils.ParseDec("2.0"),
			utils.ParseDec("1.0"),
			"",
			sdkmath.ZeroInt(), sdkmath.NewInt(1_000000),
		},
		{
			"small min price",
			sdkmath.NewInt(1_000000000000000000), sdkmath.NewInt(1_000000000000000000),
			sdkmath.LegacyNewDecWithPrec(1, 15), utils.ParseDec("2.0"),
			utils.ParseDec("1.0"),
			"",
			sdkmath.NewInt(1_000000000000000000), sdkmath.NewInt(292893228075549596),
		},
		{
			"large max price",
			sdkmath.NewInt(1_000000000000000000), sdkmath.NewInt(1_000000000000000000),
			utils.ParseDec("1.0"), sdkmath.NewIntWithDecimal(1, 20).ToLegacyDec(),
			utils.ParseDec("2.0"),
			"",
			sdkmath.NewInt(585786437709747665), sdkmath.NewInt(1_000000000000000000),
		},
		{
			"close min price and max price",
			sdkmath.NewInt(1_000000000000000000), sdkmath.NewInt(1_000000000000000000),
			utils.ParseDec("1.0"), utils.ParseDec("1.001"),
			utils.ParseDec("1.0005"),
			"",
			sdkmath.NewInt(1_000000000000000000), sdkmath.NewInt(999000936633614182),
		},
		{
			"small x asset",
			sdkmath.NewInt(9), sdkmath.NewInt(9_000000000000000000),
			utils.ParseDec("0.5"), utils.ParseDec("2.0"),
			utils.ParseDec("0.5000001"),
			"",
			sdkmath.NewInt(9), sdkmath.NewInt(89999987),
		},
		{
			"small y asset",
			sdkmath.NewInt(9_000000000000000000), sdkmath.NewInt(9),
			utils.ParseDec("0.5"), utils.ParseDec("2.0"),
			utils.ParseDec("1.9999999"),
			"",
			sdkmath.NewInt(359999969), sdkmath.NewInt(9),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pool, err := amm.CreateRangedPool(tc.x, tc.y, tc.minPrice, tc.maxPrice, tc.initialPrice)
			if tc.expectedErr == "" {
				require.NoError(t, err)
				ax, ay := pool.Balances()
				require.True(intApproxEq(tc.ax, ax))
				require.True(intApproxEq(tc.ay, ay))
				require.True(t, utils.DecApproxEqual(tc.initialPrice, pool.Price()))
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestRangedPool_Deposit(t *testing.T) {
	for _, tc := range []struct {
		name               string
		rx, ry             sdkmath.Int
		ps                 sdkmath.Int
		minPrice, maxPrice sdkmath.LegacyDec
		x, y               sdkmath.Int // depositing x and y coin amount
		ax, ay             sdkmath.Int // accepted x and y coin amount
		pc                 sdkmath.Int // expected minted pool coin amount
	}{
		{
			"ideal case",
			sdkmath.NewInt(1_000000000000000000), sdkmath.NewInt(1_000000000000000000),
			sdkmath.NewInt(1_000000000000),
			utils.ParseDec("0.5"), utils.ParseDec("2.0"),
			sdkmath.NewInt(123456789), sdkmath.NewInt(123456789),
			sdkmath.NewInt(123000000), sdkmath.NewInt(123000000),
			sdkmath.NewInt(123),
		},
		{
			"single x asset pool",
			sdkmath.NewInt(1_000000000000000000), sdkmath.NewInt(0),
			sdkmath.NewInt(1_000000000000),
			utils.ParseDec("0.5"), utils.ParseDec("2.0"),
			sdkmath.NewInt(123456789), sdkmath.NewInt(0),
			sdkmath.NewInt(123000000), sdkmath.NewInt(0),
			sdkmath.NewInt(123),
		},
		{
			"single y asset pool",
			sdkmath.NewInt(0), sdkmath.NewInt(1_000000000000000000),
			sdkmath.NewInt(1_000000000000),
			utils.ParseDec("0.5"), utils.ParseDec("2.0"),
			sdkmath.NewInt(0), sdkmath.NewInt(123456789),
			sdkmath.NewInt(0), sdkmath.NewInt(123000000),
			sdkmath.NewInt(123),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pool := amm.NewRangedPool(tc.rx, tc.ry, tc.ps, tc.minPrice, tc.maxPrice)
			ax, ay, pc := amm.Deposit(tc.rx, tc.ry, tc.ps, tc.x, tc.y)
			require.True(sdkmath.IntEq(t, tc.ax, ax))
			require.True(sdkmath.IntEq(t, tc.ay, ay))
			require.True(sdkmath.IntEq(t, tc.pc, pc))
			newPool := amm.NewRangedPool(tc.rx.Add(ax), tc.ry.Add(ay), tc.ps.Add(pc), tc.minPrice, tc.maxPrice)

			var reserveRatio sdkmath.LegacyDec
			switch {
			case tc.rx.IsZero():
				reserveRatio = ay.ToLegacyDec().Quo(tc.ry.ToLegacyDec())
			case tc.ry.IsZero():
				reserveRatio = ax.ToLegacyDec().Quo(tc.rx.ToLegacyDec())
			default:
				reserveRatio = ax.ToLegacyDec().Quo(tc.rx.ToLegacyDec())
				require.True(t, utils.DecApproxEqual(reserveRatio, ay.ToLegacyDec().Quo(tc.ry.ToLegacyDec())))
			}

			// check ax/ay == rx/ry
			if !tc.rx.IsZero() && !tc.ry.IsZero() {
				require.True(t, utils.DecApproxEqual(ax.ToLegacyDec().Quo(ay.ToLegacyDec()), tc.rx.ToLegacyDec().Quo(tc.ry.ToLegacyDec())))
			}

			// check ax/rx == ay/ry == pc/ps
			require.True(t, utils.DecApproxEqual(reserveRatio, pc.ToLegacyDec().Quo(tc.ps.ToLegacyDec())))

			// check pool price before == pool price after
			require.True(t, utils.DecApproxEqual(pool.Price(), newPool.Price()))

			transX, transY := pool.Translation()
			transXPrime, transYPrime := newPool.Translation()
			// alpha = reserveRatio
			// check transX' == transX * (1+alpha), transY' == transY * (1+alpha)
			require.True(t, utils.DecApproxEqual(reserveRatio.Add(sdkmath.LegacyOneDec()), transXPrime.Quo(transX)))
			require.True(t, utils.DecApproxEqual(reserveRatio.Add(sdkmath.LegacyOneDec()), transYPrime.Quo(transY)))
		})
	}
}

func TestRangedPool_Withdraw(t *testing.T) {
	for _, tc := range []struct {
		name               string
		rx, ry             sdkmath.Int
		ps                 sdkmath.Int
		minPrice, maxPrice sdkmath.LegacyDec
		pc                 sdkmath.Int // redeeming pool coin amount
		x, y               sdkmath.Int // withdrawn x and y coin amount
	}{
		{
			"ideal case",
			sdkmath.NewInt(1_000000000000000000), sdkmath.NewInt(1_000000000000000000),
			sdkmath.NewInt(1_000000000000),
			utils.ParseDec("0.5"), utils.ParseDec("2.0"),
			sdkmath.NewInt(123),
			sdkmath.NewInt(123000000), sdkmath.NewInt(123000000),
		},
		{
			"single x asset pool",
			sdkmath.NewInt(1_000000000000000000), sdkmath.NewInt(0),
			sdkmath.NewInt(1_000000000000),
			utils.ParseDec("0.5"), utils.ParseDec("2.0"),
			sdkmath.NewInt(123),
			sdkmath.NewInt(123000000), sdkmath.NewInt(0),
		},
		{
			"single y asset pool",
			sdkmath.NewInt(0), sdkmath.NewInt(1_000000000000000000),
			sdkmath.NewInt(1_000000000000),
			utils.ParseDec("0.5"), utils.ParseDec("2.0"),
			sdkmath.NewInt(123),
			sdkmath.NewInt(0), sdkmath.NewInt(123000000),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pool := amm.NewRangedPool(tc.rx, tc.ry, tc.ps, tc.minPrice, tc.maxPrice)
			x, y := amm.Withdraw(tc.rx, tc.ry, tc.ps, tc.pc, sdkmath.LegacyZeroDec())
			require.True(sdkmath.IntEq(t, tc.x, x))
			require.True(sdkmath.IntEq(t, tc.y, y))
			newPool := amm.NewRangedPool(tc.rx.Sub(x), tc.ry.Sub(y), tc.ps.Sub(tc.pc), tc.minPrice, tc.maxPrice)

			var reserveRatio sdkmath.LegacyDec
			switch {
			case tc.rx.IsZero():
				reserveRatio = y.ToLegacyDec().Quo(tc.ry.ToLegacyDec())
			case tc.ry.IsZero():
				reserveRatio = x.ToLegacyDec().Quo(tc.rx.ToLegacyDec())
			default:
				reserveRatio = x.ToLegacyDec().Quo(tc.rx.ToLegacyDec())
				require.True(t, utils.DecApproxEqual(reserveRatio, y.ToLegacyDec().Quo(tc.ry.ToLegacyDec())))
			}

			// check x/y == rx/ry
			if !tc.rx.IsZero() && !tc.ry.IsZero() {
				require.True(t, utils.DecApproxEqual(x.ToLegacyDec().Quo(y.ToLegacyDec()), tc.rx.ToLegacyDec().Quo(tc.ry.ToLegacyDec())))
			}

			// check x/rx == y/ry == pc/ps
			require.True(t, utils.DecApproxEqual(reserveRatio, tc.pc.ToLegacyDec().Quo(tc.ps.ToLegacyDec())))

			// check pool price before == pool price after
			require.True(t, utils.DecApproxEqual(pool.Price(), newPool.Price()))

			transX, transY := pool.Translation()
			transXPrime, transYPrime := newPool.Translation()
			// alpha = reserveRatio
			// check transX' == transX * (1+alpha), transY' == transY * (1+alpha)
			require.True(t, utils.DecApproxEqual(reserveRatio.Add(sdkmath.LegacyOneDec()), transXPrime.Quo(transX)))
			require.True(t, utils.DecApproxEqual(reserveRatio.Add(sdkmath.LegacyOneDec()), transYPrime.Quo(transY)))
		})
	}
}

func TestRangedPool_BuyAmountOver(t *testing.T) {
	pool := amm.NewRangedPool(
		sdkmath.NewInt(1000000), sdkmath.NewInt(1000000), sdkmath.Int{},
		utils.ParseDec("0.5"), utils.ParseDec("2.0"))

	for _, tc := range []struct {
		pool  *amm.RangedPool
		price sdkmath.LegacyDec
		amt   sdkmath.Int
	}{
		{pool, utils.ParseDec("1.1"), sdkmath.ZeroInt()},
		{pool, utils.ParseDec("1.0"), sdkmath.ZeroInt()},
		{pool, utils.ParseDec("0.9"), sdkmath.NewInt(379357)},
		{pool, utils.ParseDec("0.8"), sdkmath.NewInt(853553)},
	} {
		t.Run("", func(t *testing.T) {
			amt := tc.pool.BuyAmountOver(tc.price, true)
			require.True(sdkmath.IntEq(t, tc.amt, amt))
		})
	}
}

func TestRangedPool_SellAmountUnder(t *testing.T) {
	pool := amm.NewRangedPool(
		sdkmath.NewInt(1000000), sdkmath.NewInt(1000000), sdkmath.Int{},
		utils.ParseDec("0.5"), utils.ParseDec("2.0"))

	for _, tc := range []struct {
		pool  *amm.RangedPool
		price sdkmath.LegacyDec
		amt   sdkmath.Int
	}{
		{pool, utils.ParseDec("0.9"), sdkmath.ZeroInt()},
		{pool, utils.ParseDec("1.0"), sdkmath.ZeroInt()},
		{pool, utils.ParseDec("1.1"), sdkmath.NewInt(310383)},
		{pool, utils.ParseDec("1.2"), sdkmath.NewInt(569035)},
	} {
		t.Run("", func(t *testing.T) {
			amt := tc.pool.SellAmountUnder(tc.price, true)
			require.True(sdkmath.IntEq(t, tc.amt, amt))
		})
	}
}

func TestRangedPool_BuyAmountTo(t *testing.T) {
	pool := amm.NewRangedPool(
		sdkmath.NewInt(1000000), sdkmath.NewInt(1000000), sdkmath.Int{},
		utils.ParseDec("0.5"), utils.ParseDec("2.0"))

	for _, tc := range []struct {
		pool  *amm.RangedPool
		price sdkmath.LegacyDec
		amt   sdkmath.Int
	}{
		{pool, utils.ParseDec("1.1"), sdkmath.ZeroInt()},
		{pool, utils.ParseDec("1.0"), sdkmath.ZeroInt()},
		{pool, utils.ParseDec("0.8"), sdkmath.NewInt(450560)},
		{pool, utils.ParseDec("0.7"), sdkmath.NewInt(796682)},
		{
			amm.NewRangedPool(
				sdkmath.NewInt(957322), sdkmath.NewInt(3351038710333311), sdkmath.Int{},
				utils.ParseDec("0.9"), utils.ParseDec("1.1"),
			),
			utils.ParseDec("0.899580000000000000"),
			sdkmath.NewInt(1064187),
		},
	} {
		t.Run("", func(t *testing.T) {
			amt := tc.pool.BuyAmountTo(tc.price)
			require.True(sdkmath.IntEq(t, tc.amt, amt))
		})
	}
}

func TestRangedPool_SellAmountTo(t *testing.T) {
	pool := amm.NewRangedPool(
		sdkmath.NewInt(1000000), sdkmath.NewInt(1000000), sdkmath.Int{},
		utils.ParseDec("0.5"), utils.ParseDec("2.0"))

	for _, tc := range []struct {
		pool  *amm.RangedPool
		price sdkmath.LegacyDec
		amt   sdkmath.Int
	}{
		{pool, utils.ParseDec("0.9"), sdkmath.ZeroInt()},
		{pool, utils.ParseDec("1.0"), sdkmath.ZeroInt()},
		{pool, utils.ParseDec("1.4"), sdkmath.NewInt(528676)},
		{pool, utils.ParseDec("1.5"), sdkmath.NewInt(626519)},
	} {
		t.Run("", func(t *testing.T) {
			amt := tc.pool.SellAmountTo(tc.price)
			require.True(sdkmath.IntEq(t, tc.amt, amt))
		})
	}
}

func TestRangedPool_exhaust(t *testing.T) {
	for _, tc := range []struct {
		pool *amm.RangedPool
	}{
		{
			amm.NewRangedPool(
				sdkmath.NewInt(1000000), sdkmath.NewInt(1000000), sdkmath.Int{},
				utils.ParseDec("0.5"), utils.ParseDec("2.0")),
		},
		{
			amm.NewRangedPool(
				sdkmath.NewInt(1_000000000000000000), sdkmath.NewInt(9_000000000000000000), sdkmath.Int{},
				utils.ParseDec("0.1001"), utils.ParseDec("10.05")),
		},
		{
			amm.NewRangedPool(
				sdkmath.NewInt(123456789), sdkmath.NewInt(987654321), sdkmath.Int{},
				utils.ParseDec("0.05"), utils.ParseDec("20.1")),
		},
	} {
		t.Run("", func(t *testing.T) {
			rx, ry := tc.pool.Balances()
			minPrice := tc.pool.MinPrice()
			maxPrice := tc.pool.MaxPrice()
			orders := amm.PoolSellOrders(tc.pool, amm.DefaultOrderer, minPrice, maxPrice, 4)
			amt := amm.TotalAmount(orders)
			require.True(t, amt.LTE(ry))
			require.True(t, amt.GTE(ry.ToLegacyDec().Mul(utils.ParseDec("0.99")).TruncateInt()))
			orders = amm.PoolBuyOrders(tc.pool, amm.DefaultOrderer, minPrice, maxPrice, 4)
			x := sdkmath.ZeroInt()
			for _, order := range orders {
				x = x.Add(order.GetPrice().MulInt(order.GetAmount()).TruncateInt())
			}
			require.True(t, x.LTE(rx))
			require.True(t, x.GTE(rx.ToLegacyDec().Mul(utils.ParseDec("0.99")).TruncateInt()))
		})
	}
}

func TestRangedPool_SwapPriceOutOfRange(t *testing.T) {
	r := rand.New(rand.NewSource(0))

	for i := 0; i < 1000; i++ {
		rx := utils.RandomInt(r, sdkmath.NewInt(1_000000), sdkmath.NewInt(1000_00000))
		ry := utils.RandomInt(r, sdkmath.NewInt(1_000000), sdkmath.NewInt(1000_00000))
		minPrice := utils.RandomDec(r, utils.ParseDec("0.001"), utils.ParseDec("1"))
		maxPrice := utils.RandomDec(r, minPrice.Mul(utils.ParseDec("1.01")), utils.ParseDec("1000"))
		initialPrice := utils.RandomDec(r, minPrice, maxPrice)
		pool, err := amm.CreateRangedPool(rx, ry,
			minPrice, maxPrice, initialPrice)
		require.NoError(t, err)
		rx, ry = pool.Balances()

		// Price lower than min price
		p := utils.RandomDec(r, sdkmath.LegacyNewDecWithPrec(1, 5), minPrice.Mul(utils.ParseDec("0.99")))
		amt := pool.BuyAmountTo(p)
		nextRx := rx.Sub(p.MulInt(amt).Ceil().TruncateInt())
		nextRy := ry.Add(amt)
		require.True(t, nextRx.LTE(sdkmath.OneInt()))
		nextPool := amm.NewRangedPool(nextRx, nextRy, sdkmath.Int{}, minPrice, maxPrice)
		require.True(t, utils.DecApproxEqual(minPrice, nextPool.Price()))

		// Price higher than min price
		p = utils.RandomDec(r, maxPrice.Mul(utils.ParseDec("1.01")), utils.ParseDec("1000000"))
		amt = pool.SellAmountTo(p)
		nextRx = rx.Add(p.MulInt(amt).TruncateInt())
		nextRy = ry.Sub(amt)
		require.True(t, nextRy.LTE(sdkmath.OneInt()))
		nextPool = amm.NewRangedPool(nextRx, nextRy, sdkmath.Int{}, minPrice, maxPrice)
		require.True(t, utils.DecApproxEqual(maxPrice, nextPool.Price()))
	}
}

func TestInitialPoolCoinSupply(t *testing.T) {
	for _, tc := range []struct {
		x, y sdkmath.Int
		ps   sdkmath.Int
	}{
		{sdkmath.NewInt(1000000), sdkmath.NewInt(1000000), sdkmath.NewInt(10000000)},
		{sdkmath.NewInt(1000000), sdkmath.NewInt(10000000), sdkmath.NewInt(100000000)},
		{sdkmath.NewInt(1000000), sdkmath.NewInt(100000000), sdkmath.NewInt(100000000)},
		{sdkmath.NewInt(10000000), sdkmath.NewInt(100000000), sdkmath.NewInt(1000000000)},
		{sdkmath.NewInt(999999), sdkmath.NewInt(9999999), sdkmath.NewInt(10000000)},
	} {
		t.Run("", func(t *testing.T) {
			require.True(sdkmath.IntEq(t, tc.ps, amm.InitialPoolCoinSupply(tc.x, tc.y)))
		})
	}
}

func TestBasicPool_BuyAmountOverOverflow(t *testing.T) {
	n, _ := sdkmath.NewIntFromString("10000000000000000000000000000000000000000000")
	pool := amm.NewBasicPool(n, sdkmath.NewInt(1000), sdkmath.Int{})
	amt := pool.BuyAmountOver(defTickPrec.LowestTick(), true)
	require.True(sdkmath.IntEq(t, amm.MaxCoinAmount, amt))
}

func TestBasicPoolOrders(t *testing.T) {
	pool := amm.NewBasicPool(sdkmath.NewInt(862431695563), sdkmath.NewInt(37852851767), sdkmath.Int{})
	poolPrice := pool.Price()
	lowestPrice := poolPrice.Mul(sdkmath.LegacyNewDecWithPrec(9, 1))
	highestPrice := poolPrice.Mul(sdkmath.LegacyNewDecWithPrec(11, 1))
	require.Len(t, amm.PoolOrders(pool, amm.DefaultOrderer, lowestPrice, highestPrice, 4), 375)
}

func BenchmarkBasicPoolOrders(b *testing.B) {
	pool := amm.NewBasicPool(sdkmath.NewInt(862431695563), sdkmath.NewInt(37852851767), sdkmath.Int{})
	poolPrice := pool.Price()
	lowestPrice := poolPrice.Mul(sdkmath.LegacyNewDecWithPrec(9, 1))
	highestPrice := poolPrice.Mul(sdkmath.LegacyNewDecWithPrec(11, 1))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		amm.PoolOrders(pool, amm.DefaultOrderer, lowestPrice, highestPrice, 4)
	}
}
