package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Multiply(a, b sdk.Dec) sdk.Dec {
	return a.Mul(b)
}
func Mul(a, b sdk.Int) sdk.Int {
	return a.Mul(b)
}
func Divide(a, b sdk.Dec) sdk.Dec {
	return a.Quo(b)
}

func Add(a, b sdk.Dec) sdk.Dec {
	return a.Add(b)
}

func Sub(a, b sdk.Dec) sdk.Dec {
	return a.Sub(b)
}

func (k Keeper) getInflowTokenTargetAmount(amount, price sdk.Int) sdk.Int {
	result := amount.Mul(price)
	return result
}

func (k Keeper) getOutflowTokenInitialPrice(price sdk.Int, buffer sdk.Dec) sdk.Dec {
	result := buffer.MulInt(price)
	return result
}

func (k Keeper) getOutflowTokenEndPrice(price, cusp sdk.Dec) sdk.Dec {
	result := Multiply(price, cusp)
	return result
}

func (k Keeper) getPriceFromLinearDecreaseFunction(top sdk.Dec, tau, dur sdk.Int) sdk.Dec {
	result1 := (tau.Sub(dur)).Quo(tau)
	result2 := sdk.NewDecFromInt(result1)
	result := Multiply(top, result2)
	return result
}

func (k Keeper) getPriceFromStairStepExponentialDecreaseFunction(top, decreasePercent sdk.Dec, step, dur sdk.Int) sdk.Dec {
	cut := Sub(sdk.MustNewDecFromStr("1"), decreasePercent)
	result1 := dur.Quo(step)
	count := result1.Uint64()
	result2 := cut.Power(count)
	result := top.Mul(result2)
	return result
}

func (k Keeper) getPriceFromContinuousExponentialDecreaseFunction(top, decreasePercent sdk.Dec, dur sdk.Int) sdk.Dec {
	cut := Sub(sdk.MustNewDecFromStr("1"), decreasePercent)
	count := dur.Uint64()
	result2 := cut.Power(count)
	result := top.Mul(result2)
	return result
}

func (k Keeper) getBurnAmount(amount sdk.Int, liqPenalty sdk.Dec) sdk.Int {
	liqPenalty = liqPenalty.Add(sdk.NewDec(1))
	amount1 := sdk.NewDecFromInt(amount)
	result := amount1.Quo(liqPenalty).Ceil().TruncateInt()
	return result
}
