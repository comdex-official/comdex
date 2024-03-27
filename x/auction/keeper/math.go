package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Multiply(a, b sdk.Dec) sdk.Dec {
	return a.Mul(b)
}

func (k Keeper) getOutflowTokenInitialPrice(price sdk.Int, buffer sdk.Dec) sdk.Dec {
	result := buffer.Mul(sdk.NewDecFromInt(price))
	return result
}

func (k Keeper) getOutflowTokenEndPrice(price, cusp sdk.Dec) sdk.Dec {
	result := Multiply(price, cusp)
	return result
}

func (k Keeper) getPriceFromLinearDecreaseFunction(top sdk.Dec, tau, dur sdk.Int) sdk.Dec {
	result1 := tau.Sub(dur)
	result2 := top.Mul(sdk.NewDecFromInt(result1))
	result3 := result2.Quo(sdk.NewDecFromInt(tau))
	return result3
}
