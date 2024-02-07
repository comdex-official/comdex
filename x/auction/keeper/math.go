package keeper

import (
	sdkmath "cosmossdk.io/math"
)

func Multiply(a, b sdkmath.LegacyDec) sdkmath.LegacyDec {
	return a.Mul(b)
}

func (k Keeper) getOutflowTokenInitialPrice(price sdkmath.Int, buffer sdkmath.LegacyDec) sdkmath.LegacyDec {
	result := buffer.Mul(sdkmath.LegacyNewDecFromInt(price))
	return result
}

func (k Keeper) getOutflowTokenEndPrice(price, cusp sdkmath.LegacyDec) sdkmath.LegacyDec {
	result := Multiply(price, cusp)
	return result
}

func (k Keeper) getPriceFromLinearDecreaseFunction(top sdkmath.LegacyDec, tau, dur sdkmath.Int) sdkmath.LegacyDec {
	result1 := tau.Sub(dur)
	result2 := top.Mul(sdkmath.LegacyNewDecFromInt(result1))
	result3 := result2.Quo(sdkmath.LegacyNewDecFromInt(tau))
	return result3
}
