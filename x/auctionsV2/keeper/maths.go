package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func Multiply(a, b sdk.Dec) sdk.Dec {
	return a.Mul(b)
}

func (k Keeper) GetCollalteralTokenInitialPrice(price sdk.Int, premium sdk.Dec) sdk.Dec {
	result := premium.Mul(price.ToDec())
	return result
}

func (k Keeper) GetCollalteralTokenEndPrice(price, discount sdk.Dec) sdk.Dec {
	result := Multiply(price, discount)
	return result
}

func (k Keeper) getPriceFromLinearDecreaseFunction(top sdk.Dec, tau, dur sdk.Int) sdk.Dec {
	result1 := tau.Sub(dur)
	result2 := top.Mul(result1.ToDec())
	result3 := result2.Quo(tau.ToDec())
	return result3
}

func (k Keeper) getOutflowTokenEndPrice(price, cusp sdk.Dec) sdk.Dec {
	result := Multiply(price, cusp)
	return result
}
