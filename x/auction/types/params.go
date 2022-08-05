package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

const (
	DutchBidGas     = sdk.Gas(71183)
	DebtBidGas      = sdk.Gas(57580)
	DutchLendBidGas = sdk.Gas(71183)
	SurplusBidGas   = sdk.Gas(57580)
)

// ParamKeyTable the param key table for launch module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance.
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params.
func (p Params) Validate() error {
	for _, v := range []struct {
		value     interface{}
		validator func(interface{}) error
	}{} {
		if err := v.validator(v.value); err != nil {
			return err
		}
	}
	return nil
}
