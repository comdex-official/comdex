package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"fmt"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyUnliquidatePointPercent = []byte("UnliquidatePointPercent")
)

var (
	DefaultUnliquidatePointPercent = "1.6"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		UnliquidatePointPercent: DefaultUnliquidatePointPercent,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyUnliquidatePointPercent, &p.UnliquidatePointPercent, ValidateUnliquidatePointPercent),

	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	for _, v := range []struct {
		value     interface{}
		validator func(interface{}) error
	}{
		{p.UnliquidatePointPercent, ValidateUnliquidatePointPercent},
	
	} {
		if err := v.validator(v.value); err != nil {
			return err
		}
	}
	return nil
}
func ValidateUnliquidatePointPercent(i interface{}) error {
	v,ok:=i.(string)
	if !ok{
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	   q,_:=sdk.NewDecFromStr(v)
	   u,_:=sdk.NewDecFromStr("1")
	if q.LT(u) {
		return fmt.Errorf("Unliquidate Point Percentage cannot be less than 100 percent")
	}

	return nil
}
