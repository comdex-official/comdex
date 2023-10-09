package amm

import (
	sdkmath "cosmossdk.io/math"
)

// The minimum and maximum coin amount used in the amm package.
var (
	MinCoinAmount = sdkmath.NewInt(100)
	MaxCoinAmount = sdkmath.NewIntWithDecimal(1, 40)
)

var (
	MinPoolPrice               = sdkmath.LegacyNewDecWithPrec(1, 15)                           // 10^-15
	MaxPoolPrice               = sdkmath.LegacyNewDecFromInt(sdkmath.NewIntWithDecimal(1, 20)) // 10^20
	MinRangedPoolPriceGapRatio = sdkmath.LegacyNewDecWithPrec(1, 3)                            // 0.001, 0.1%
)
