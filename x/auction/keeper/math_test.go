package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
)

//
//func TestMultiply(t *testing.T) {
//	a := sdk.MustNewDecFromStr("1.2")
//	b := sdk.NewInt(4)
//	c :=
//		fmt.Println(a, c)
//	require.Equal(t, sdk.NewInt(12), c)
//}

func getBurnAmount(amount sdk.Int, liqPenalty sdk.Dec) sdk.Int {
	liqPenalty = liqPenalty.Add(sdk.NewDec(1))
	result := amount.ToDec().Quo(liqPenalty).Ceil().TruncateInt()
	return result
}

func TestRed(t *testing.T) {
	amount := sdk.NewInt(100)
	liq := sdk.MustNewDecFromStr("0.15")
	fmt.Println(getBurnAmount(amount, liq))
}
