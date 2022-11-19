package amm_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	utils "github.com/petrichormoney/petri/types"
	"github.com/petrichormoney/petri/x/liquidity/amm"
)

const defTickPrec = amm.TickPrecision(3)

func TestPriceToTick(t *testing.T) {
	for _, tc := range []struct {
		price    sdk.Dec
		expected sdk.Dec
	}{
		{utils.ParseDec("0.000000000000099999"), utils.ParseDec("0.00000000000009999")},
		{utils.ParseDec("1.999999999999999999"), utils.ParseDec("1.999")},
		{utils.ParseDec("99.999999999999999999"), utils.ParseDec("99.99")},
		{utils.ParseDec("100.999999999999999999"), utils.ParseDec("100.9")},
		{utils.ParseDec("9999.999999999999999999"), utils.ParseDec("9999")},
		{utils.ParseDec("10019"), utils.ParseDec("10010")},
		{utils.ParseDec("1000100005"), utils.ParseDec("1000000000")},
	} {
		require.True(sdk.DecEq(t, tc.expected, amm.PriceToDownTick(tc.price, 3)))
	}
}

func TestTick(t *testing.T) {
	for _, tc := range []struct {
		i        int
		prec     int
		expected sdk.Dec
	}{
		{0, 3, sdk.NewDecWithPrec(1, int64(sdk.Precision-defTickPrec))},
		{1, 3, utils.ParseDec("0.000000000000001001")},
		{8999, 3, utils.ParseDec("0.000000000000009999")},
		{9000, 3, utils.ParseDec("0.000000000000010000")},
		{9001, 3, utils.ParseDec("0.000000000000010010")},
		{17999, 3, utils.ParseDec("0.000000000000099990")},
		{18000, 3, utils.ParseDec("0.000000000000100000")},
		{135000, 3, sdk.NewDec(1)},
		{135001, 3, utils.ParseDec("1.001")},
	} {
		t.Run("", func(t *testing.T) {
			res := amm.TickFromIndex(tc.i, tc.prec)
			require.True(sdk.DecEq(t, tc.expected, res))
			require.Equal(t, tc.i, amm.TickToIndex(res, tc.prec))
		})
	}
}

func TestUpTick(t *testing.T) {
	for _, tc := range []struct {
		price    sdk.Dec
		prec     int
		expected sdk.Dec
	}{
		{utils.ParseDec("1000000000000000000"), 3, utils.ParseDec("1001000000000000000")},
		{utils.ParseDec("1000"), 3, utils.ParseDec("1001")},
		{utils.ParseDec("999.9"), 3, utils.ParseDec("1000")},
		{utils.ParseDec("999.0"), 3, utils.ParseDec("999.1")},
		{utils.ParseDec("1.100"), 3, utils.ParseDec("1.101")},
		{utils.ParseDec("1.000"), 3, utils.ParseDec("1.001")},
		{utils.ParseDec("0.9999"), 3, utils.ParseDec("1.000")},
		{utils.ParseDec("0.1000"), 3, utils.ParseDec("0.1001")},
		{utils.ParseDec("0.09999"), 3, utils.ParseDec("0.1000")},
		{utils.ParseDec("0.09997"), 3, utils.ParseDec("0.09998")},
	} {
		t.Run("", func(t *testing.T) {
			require.True(sdk.DecEq(t, tc.expected, amm.UpTick(tc.price, tc.prec)))
		})
	}
}

func TestDownTick(t *testing.T) {
	for _, tc := range []struct {
		price    sdk.Dec
		prec     int
		expected sdk.Dec
	}{
		{utils.ParseDec("1000000000000000000"), 3, utils.ParseDec("999900000000000000")},
		{utils.ParseDec("10010"), 3, utils.ParseDec("10000")},
		{utils.ParseDec("100.0"), 3, utils.ParseDec("99.99")},
		{utils.ParseDec("99.99"), 3, utils.ParseDec("99.98")},
		{utils.ParseDec("1.000"), 3, utils.ParseDec("0.9999")},
		{utils.ParseDec("0.9990"), 3, utils.ParseDec("0.9989")},
		{utils.ParseDec("0.9999"), 3, utils.ParseDec("0.9998")},
		{utils.ParseDec("0.1"), 3, utils.ParseDec("0.09999")},
		{utils.ParseDec("0.00000000000001000"), 3, utils.ParseDec("0.000000000000009999")},
		{utils.ParseDec("0.000000000000001001"), 3, utils.ParseDec("0.000000000000001000")},
	} {
		t.Run("", func(t *testing.T) {
			require.True(sdk.DecEq(t, tc.expected, amm.DownTick(tc.price, tc.prec)))
		})
	}
}

func TestHighestTick(t *testing.T) {
	for _, tc := range []struct {
		prec     int
		expected string
	}{
		{3, "133400000000000000000000000000000000000000000000000000000000000000000000000000"},
		{0, "100000000000000000000000000000000000000000000000000000000000000000000000000000"},
		{1, "130000000000000000000000000000000000000000000000000000000000000000000000000000"},
	} {
		t.Run("", func(t *testing.T) {
			i, ok := new(big.Int).SetString(tc.expected, 10)
			require.True(t, ok)
			tick := amm.HighestTick(tc.prec)
			require.True(sdk.DecEq(t, sdk.NewDecFromBigInt(i), tick))
			require.Panics(t, func() {
				amm.UpTick(tick, tc.prec)
			})
		})
	}
}

func TestLowestTick(t *testing.T) {
	for _, tc := range []struct {
		prec     int
		expected sdk.Dec
	}{
		{0, sdk.NewDecWithPrec(1, 18)},
		{3, sdk.NewDecWithPrec(1, 15)},
	} {
		t.Run("", func(t *testing.T) {
			require.True(sdk.DecEq(t, tc.expected, amm.LowestTick(tc.prec)))
		})
	}
}

func TestPriceToUpTick(t *testing.T) {
	for _, tc := range []struct {
		price    sdk.Dec
		prec     int
		expected sdk.Dec
	}{
		{utils.ParseDec("1.0015"), 3, utils.ParseDec("1.002")},
		{utils.ParseDec("100"), 3, utils.ParseDec("100")},
		{utils.ParseDec("100.01"), 3, utils.ParseDec("100.1")},
		{utils.ParseDec("100.099"), 3, utils.ParseDec("100.1")},
	} {
		t.Run("", func(t *testing.T) {
			require.True(sdk.DecEq(t, tc.expected, amm.PriceToUpTick(tc.price, tc.prec)))
		})
	}
}

func TestRoundTickIndex(t *testing.T) {
	for _, tc := range []struct {
		i        int
		expected int
	}{
		{0, 0},
		{1, 2},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 6},
		{6, 6},
		{7, 8},
		{8, 8},
		{9, 10},
		{10, 10},
	} {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tc.expected, amm.RoundTickIndex(tc.i))
		})
	}
}

func TestRoundPrice(t *testing.T) {
	for _, tc := range []struct {
		price    sdk.Dec
		prec     int
		expected sdk.Dec
	}{
		{utils.ParseDec("0.000000000000001000"), 3, utils.ParseDec("0.000000000000001000")},
		{utils.ParseDec("0.000000000000010000"), 3, utils.ParseDec("0.000000000000010000")},
		{utils.ParseDec("0.000000000000010005"), 3, utils.ParseDec("0.000000000000010000")},
		{utils.ParseDec("0.000000000000010015"), 3, utils.ParseDec("0.000000000000010020")},
		{utils.ParseDec("0.000000000000010025"), 3, utils.ParseDec("0.000000000000010020")},
		{utils.ParseDec("0.000000000000010035"), 3, utils.ParseDec("0.000000000000010040")},
		{utils.ParseDec("0.000000000000010045"), 3, utils.ParseDec("0.000000000000010040")},
		{utils.ParseDec("1.0005"), 3, utils.ParseDec("1.0")},
		{utils.ParseDec("1.0015"), 3, utils.ParseDec("1.002")},
		{utils.ParseDec("1.0025"), 3, utils.ParseDec("1.002")},
		{utils.ParseDec("1.0035"), 3, utils.ParseDec("1.004")},
	} {
		t.Run("", func(t *testing.T) {
			require.True(sdk.DecEq(t, tc.expected, amm.RoundPrice(tc.price, tc.prec)))
		})
	}
}

func BenchmarkUpTick(b *testing.B) {
	b.Run("price fit in ticks", func(b *testing.B) {
		price := utils.ParseDec("0.9999")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			amm.UpTick(price, 3)
		}
	})
	b.Run("price not fit in ticks", func(b *testing.B) {
		price := utils.ParseDec("0.99995")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			amm.UpTick(price, 3)
		}
	})
}

func BenchmarkDownTick(b *testing.B) {
	b.Run("price fit in ticks", func(b *testing.B) {
		price := utils.ParseDec("0.9999")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			amm.DownTick(price, 3)
		}
	})
	b.Run("price not fit in ticks", func(b *testing.B) {
		price := utils.ParseDec("0.99995")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			amm.DownTick(price, 3)
		}
	})
	b.Run("price at edge", func(b *testing.B) {
		price := utils.ParseDec("1")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			amm.DownTick(price, 3)
		}
	})
}

func TestTicks(t *testing.T) {
	dec := utils.ParseDec
	for i, tc := range []struct {
		basePrice sdk.Dec
		numTicks  int
		ticks     []sdk.Dec
	}{
		{
			dec("999.9"), 3,
			[]sdk.Dec{
				dec("1002"),
				dec("1001"),
				dec("1000"),
				dec("999.9"),
				dec("999.8"),
				dec("999.7"),
				dec("999.6"),
			},
		},
		{
			dec("2.0"), 3,
			[]sdk.Dec{
				dec("2.003"),
				dec("2.002"),
				dec("2.001"),
				dec("2.000"),
				dec("1.999"),
				dec("1.998"),
				dec("1.997"),
			},
		},
		{
			dec("1.000001"), 3,
			[]sdk.Dec{
				dec("1.003"),
				dec("1.002"),
				dec("1.001"),
				dec("1.000"),
				dec("0.9999"),
				dec("0.9998"),
			},
		},
		{
			dec("1.0"), 3,
			[]sdk.Dec{
				dec("1.003"),
				dec("1.002"),
				dec("1.001"),
				dec("1.000"),
				dec("0.9999"),
				dec("0.9998"),
				dec("0.9997"),
			},
		},
		{
			dec("0.99999"), 3,
			[]sdk.Dec{
				dec("1.002"),
				dec("1.001"),
				dec("1.000"),
				dec("0.9999"),
				dec("0.9998"),
				dec("0.9997"),
			},
		},
		{
			dec("0.9999"), 3,
			[]sdk.Dec{
				dec("1.002"),
				dec("1.001"),
				dec("1.000"),
				dec("0.9999"),
				dec("0.9998"),
				dec("0.9997"),
				dec("0.9996"),
			},
		},
		{
			dec("0.99985"), 3,
			[]sdk.Dec{
				dec("1.001"),
				dec("1.000"),
				dec("0.9999"),
				dec("0.9998"),
				dec("0.9997"),
				dec("0.9996"),
			},
		},
		{
			dec("0.9998"), 3,
			[]sdk.Dec{
				dec("1.001"),
				dec("1.000"),
				dec("0.9999"),
				dec("0.9998"),
				dec("0.9997"),
				dec("0.9996"),
				dec("0.9995"),
			},
		},
		{
			dec("0.99975"), 3,
			[]sdk.Dec{
				dec("1.000"),
				dec("0.9999"),
				dec("0.9998"),
				dec("0.9997"),
				dec("0.9996"),
				dec("0.9995"),
			},
		},
		{
			dec("0.9997"), 3,
			[]sdk.Dec{
				dec("1.000"),
				dec("0.9999"),
				dec("0.9998"),
				dec("0.9997"),
				dec("0.9996"),
				dec("0.9995"),
				dec("0.9994"),
			},
		},
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			require.Equal(t, tc.ticks, amm.Ticks(tc.basePrice, tc.numTicks, int(defTickPrec)))
		})
	}
}

func TestEvenTicks(t *testing.T) {
	dec := utils.ParseDec
	for i, tc := range []struct {
		basePrice sdk.Dec
		numTicks  int
		ticks     []sdk.Dec
	}{
		{
			dec("999.9"), 3,
			[]sdk.Dec{
				dec("1002"),
				dec("1001"),
				dec("1000"),
				dec("999.0"),
				dec("998.0"),
				dec("997.0"),
			},
		},
		{
			dec("2.0"), 3,
			[]sdk.Dec{
				dec("2.003"),
				dec("2.002"),
				dec("2.001"),
				dec("2.000"),
				dec("1.999"),
				dec("1.998"),
				dec("1.997"),
			},
		},
		{
			dec("1.000001"), 3,
			[]sdk.Dec{
				dec("1.003"),
				dec("1.002"),
				dec("1.001"),
				dec("1.000"),
				dec("0.999"),
				dec("0.998"),
			},
		},
		{
			dec("1.0"), 3,
			[]sdk.Dec{
				dec("1.003"),
				dec("1.002"),
				dec("1.001"),
				dec("1.000"),
				dec("0.999"),
				dec("0.998"),
				dec("0.997"),
			},
		},
		{
			dec("0.99999"), 3,
			[]sdk.Dec{
				dec("1.002"),
				dec("1.001"),
				dec("1.000"),
				dec("0.999"),
				dec("0.998"),
				dec("0.997"),
			},
		},
		{
			dec("0.9999"), 3,
			[]sdk.Dec{
				dec("1.002"),
				dec("1.001"),
				dec("1.000"),
				dec("0.999"),
				dec("0.998"),
				dec("0.997"),
			},
		},
		{
			dec("0.99985"), 3,
			[]sdk.Dec{
				dec("1.002"),
				dec("1.001"),
				dec("1.000"),
				dec("0.999"),
				dec("0.998"),
				dec("0.997"),
			},
		},
		{
			dec("0.9998"), 3,
			[]sdk.Dec{
				dec("1.002"),
				dec("1.001"),
				dec("1.000"),
				dec("0.999"),
				dec("0.998"),
				dec("0.997"),
			},
		},
		{
			dec("0.99975"), 3,
			[]sdk.Dec{
				dec("1.000"),
				dec("0.9999"),
				dec("0.9998"),
				dec("0.9997"),
				dec("0.9996"),
				dec("0.9995"),
			},
		},
		{
			dec("0.9997"), 3,
			[]sdk.Dec{
				dec("1.000"),
				dec("0.9999"),
				dec("0.9998"),
				dec("0.9997"),
				dec("0.9996"),
				dec("0.9995"),
				dec("0.9994"),
			},
		},
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			require.Equal(t, tc.ticks, amm.EvenTicks(tc.basePrice, tc.numTicks, int(defTickPrec)))
		})
	}
}
