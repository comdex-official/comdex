package keeper

import (
	sdkmath "cosmossdk.io/math"
)

func Multiply(a, b sdkmath.LegacyDec) sdkmath.LegacyDec {
	return a.Mul(b)
}

func (k Keeper) GetCollalteralTokenInitialPrice(price sdkmath.Int, premium sdkmath.LegacyDec) sdkmath.LegacyDec {
	result := premium.Mul(sdkmath.LegacyNewDecFromInt(price))
	return result
}

func (k Keeper) GetPriceFromLinearDecreaseFunction(CollateralTokenAuctionPrice sdkmath.LegacyDec, timeToReachZeroPrice, timeElapsed sdkmath.Int) sdkmath.LegacyDec {
	timeDifference := timeToReachZeroPrice.Sub(timeElapsed)
	resultantPrice := CollateralTokenAuctionPrice.Mul(sdkmath.LegacyNewDecFromInt(timeDifference))
	currentPrice := resultantPrice.Quo(sdkmath.LegacyNewDecFromInt(timeToReachZeroPrice))
	return currentPrice
}

func (k Keeper) GetCollateralTokenEndPrice(price, cusp sdkmath.LegacyDec) sdkmath.LegacyDec {
	result := Multiply(price, cusp)
	return result
}
