package keeper

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
)

func getBurnAmount(amount sdkmath.Int, liqPenalty sdkmath.LegacyDec) sdkmath.Int {
	liqPenalty = liqPenalty.Add(sdkmath.LegacyNewDec(1))
	result := sdkmath.LegacyNewDecFromInt(amount).Quo(liqPenalty).Ceil().TruncateInt()
	return result
}

func TestRed(t *testing.T) {
	amount := sdkmath.NewInt(100)
	liq := sdkmath.LegacyMustNewDecFromStr("0.15")
	fmt.Println(getBurnAmount(amount, liq))
}

func TestAdd(t *testing.T) {
	a := sdkmath.ZeroInt()
	b := sdkmath.NewIntFromUint64(300)
	c := sdkmath.NewIntFromUint64(100)
	d := b.Add(a.Sub(c))
	fmt.Println(d)
}
