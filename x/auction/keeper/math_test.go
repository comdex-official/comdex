package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
