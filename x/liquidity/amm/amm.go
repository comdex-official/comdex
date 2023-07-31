package amm

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// The minimum and maximum coin amount used in the amm package.
var (
	MinCoinAmount = sdk.NewInt(100)
	MaxCoinAmount = sdkmath.NewIntWithDecimal(1, 40)
)

var (
	MinPoolPrice               = sdk.NewDecWithPrec(1, 15)            // 10^-15
	MaxPoolPrice               = sdk.NewDec(sdkmath.NewIntWithDecimal(1, 20).Int64()) // 10^20
	MinRangedPoolPriceGapRatio = sdk.NewDecWithPrec(1, 3)             // 0.001, 0.1%
)
