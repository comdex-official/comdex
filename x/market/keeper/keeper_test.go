package keeper

import (
	"fmt"
	rewardtypes "github.com/comdex-official/comdex/x/rewards/types"
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestFucn(t *testing.T) {
	var id uint64
	var rate uint64
	var Decimals int64
	var amt sdk.Int
	id = 5
	amt = sdk.NewInt(1230045678934223213)

	if id == 1 {
		rate = 11632845
		Decimals = 1000000
	}
	if id == 2 {
		rate = 140530
		Decimals = 1000000
	}
	if id == 3 {
		rate = 1000000
		Decimals = 1000000
	}
	if id == 4 {
		rate = 1183857
		Decimals = 1000000
	}
	if id == 5 {
		rate = 1297119384
		Decimals = 1000000000000000000
	}
	numerator := sdk.NewDecFromInt(amt).Mul(sdk.NewDecFromInt(sdk.NewIntFromUint64(rate)))
	denominator := sdk.NewDecFromInt(sdk.NewIntFromUint64(uint64(Decimals)))
	result := numerator.Quo(denominator)
	fmt.Println("result ", result)
	fmt.Println("--------------")
	newAmtwithDec := result.Quo(sdk.NewDecFromInt(sdk.NewIntFromUint64(1000000)))
	finalAMount := newAmtwithDec.Mul(sdk.NewDecFromInt(sdk.NewIntFromUint64(uint64(1000000))))
	test123 := sdk.Int(newAmtwithDec.Mul(sdk.NewDecFromInt(sdk.NewIntFromUint64(uint64(1000000)))))
	fmt.Println("test123 ", test123)
	fmt.Println("truncate ", finalAMount.TruncateInt())
}

func TestFucn1(t *testing.T) {
	x := sdk.MustNewDecFromStr("128345678.4567432")
	fmt.Println("uint64 value", x.TruncateInt().Uint64())
}

func TestFucn2(t *testing.T) {
	//arr :=[]int{12,23,34,45,56}
	//arr = arr[:0]
	// arr = []int{}
	// arr = make([]int, 0)
	lsr, _ := sdk.NewDecFromStr("0.001")
	amount := sdk.NewInt(100000000)
	yearsElapsed := sdk.NewDec(50).QuoInt64(rewardtypes.SecondsPerYear)
	perc := lsr.String()
	a, _ := sdk.NewDecFromStr("1")
	b, _ := sdk.NewDecFromStr(perc)
	factor1 := a.Add(b)

	intPerBlockFactor := math.Pow(factor1.MustFloat64(), yearsElapsed.MustFloat64())
	intAccPerBlock := intPerBlockFactor - rewardtypes.Float64One

	amtFloat := sdk.NewDecFromInt(amount).MustFloat64()
	newAmount := intAccPerBlock * amtFloat

	fmt.Println("yearsElapsed", yearsElapsed)
	fmt.Println("yearsElapsed float", yearsElapsed.MustFloat64())
	fmt.Println("factor1 ", factor1)
	fmt.Println("factor1 float", factor1.MustFloat64())
	fmt.Println("intPerBlockFactor", intPerBlockFactor)
	fmt.Println("newAmount", newAmount)
	//fmt.Println("arr new value", arr)
}

func TestFucn3(t *testing.T) {
	x := sdk.NewInt(32606975564)
	y := sdk.NewInt(2600633730)
	fmt.Println("answer", x.Sub(y))
}
