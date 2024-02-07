package types_test

import (
	sdkmath "cosmossdk.io/math"
	"testing"

	"github.com/stretchr/testify/require"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidity/types"
)

func TestMMOrderTicks(t *testing.T) {
	require.Equal(t,
		[]types.MMOrderTick{
			{OfferCoinAmount: sdkmath.NewInt(100000), Price: utils.ParseDec("105"), Amount: sdkmath.NewInt(100000)},
			{OfferCoinAmount: sdkmath.NewInt(100000), Price: utils.ParseDec("104.45"), Amount: sdkmath.NewInt(100000)},
			{OfferCoinAmount: sdkmath.NewInt(100000), Price: utils.ParseDec("103.89"), Amount: sdkmath.NewInt(100000)},
			{OfferCoinAmount: sdkmath.NewInt(100000), Price: utils.ParseDec("103.34"), Amount: sdkmath.NewInt(100000)},
			{OfferCoinAmount: sdkmath.NewInt(100000), Price: utils.ParseDec("102.78"), Amount: sdkmath.NewInt(100000)},
			{OfferCoinAmount: sdkmath.NewInt(100000), Price: utils.ParseDec("102.23"), Amount: sdkmath.NewInt(100000)},
			{OfferCoinAmount: sdkmath.NewInt(100000), Price: utils.ParseDec("101.67"), Amount: sdkmath.NewInt(100000)},
			{OfferCoinAmount: sdkmath.NewInt(100000), Price: utils.ParseDec("101.12"), Amount: sdkmath.NewInt(100000)},
			{OfferCoinAmount: sdkmath.NewInt(100000), Price: utils.ParseDec("100.56"), Amount: sdkmath.NewInt(100000)},
			{OfferCoinAmount: sdkmath.NewInt(100000), Price: utils.ParseDec("100"), Amount: sdkmath.NewInt(100000)},
		},
		types.MMOrderTicks(
			types.OrderDirectionSell, utils.ParseDec("100"), utils.ParseDec("105"),
			sdkmath.NewInt(1000000), int(types.DefaultMaxNumMarketMakingOrderTicks), 4),
	)

	require.Equal(t,
		[]types.MMOrderTick{
			{
				OfferCoinAmount: sdkmath.NewInt(5402),
				Price:           utils.ParseDec("100.02"),
				Amount:          sdkmath.NewInt(54),
			},
			{
				OfferCoinAmount: sdkmath.NewInt(5502),
				Price:           utils.ParseDec("100.03"),
				Amount:          sdkmath.NewInt(55),
			},
		},
		types.MMOrderTicks(
			types.OrderDirectionBuy, utils.ParseDec("100.02"), utils.ParseDec("100.03"),
			sdkmath.NewInt(109), int(types.DefaultMaxNumMarketMakingOrderTicks), 4),
	)
}
